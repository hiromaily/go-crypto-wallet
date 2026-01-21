/**
 * XRPL Client Wrapper
 *
 * Provides singleton client management, connection lifecycle handling,
 * and event subscriptions for the XRP Ledger.
 */

import { Client, type ClientOptions } from 'xrpl';
import { getConfig, type XRPLConfig } from '../config';

/**
 * Connection state for the XRPL client
 */
export type ConnectionState = 'disconnected' | 'connecting' | 'connected' | 'reconnecting';

/**
 * Event types emitted by the XRPLClientManager
 */
export interface XRPLClientEvents {
  connected: () => void;
  disconnected: (code: number) => void;
  error: (errorCode: string, errorMessage: string) => void;
  reconnecting: (attempt: number) => void;
  reconnectFailed: () => void;
}

/**
 * Logger interface for dependency injection
 */
export interface Logger {
  info: (message: string, ...args: unknown[]) => void;
  warn: (message: string, ...args: unknown[]) => void;
  error: (message: string, ...args: unknown[]) => void;
  debug: (message: string, ...args: unknown[]) => void;
}

/**
 * Default console logger
 */
const defaultLogger: Logger = {
  info: (message: string, ...args: unknown[]) => console.log(`[INFO] ${message}`, ...args),
  warn: (message: string, ...args: unknown[]) => console.warn(`[WARN] ${message}`, ...args),
  error: (message: string, ...args: unknown[]) => console.error(`[ERROR] ${message}`, ...args),
  debug: (message: string, ...args: unknown[]) => console.debug(`[DEBUG] ${message}`, ...args),
};

/**
 * XRPL Client Manager
 *
 * Manages the XRPL client lifecycle including connection, disconnection,
 * reconnection, and event handling.
 */
export class XRPLClientManager {
  private client: Client | null = null;
  private config: XRPLConfig;
  private logger: Logger;
  private connectionState: ConnectionState = 'disconnected';
  private reconnectAttempts = 0;
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  private eventListeners: Map<keyof XRPLClientEvents, Set<(...args: unknown[]) => void>> =
    new Map();

  constructor(config?: Partial<XRPLConfig>, logger?: Logger) {
    const appConfig = getConfig();
    this.config = { ...appConfig.xrpl, ...config };
    this.logger = logger || defaultLogger;
  }

  /**
   * Get the current connection state
   */
  getConnectionState(): ConnectionState {
    return this.connectionState;
  }

  /**
   * Check if the client is connected
   */
  isConnected(): boolean {
    return this.client?.isConnected() ?? false;
  }

  /**
   * Add an event listener
   */
  on<K extends keyof XRPLClientEvents>(event: K, listener: XRPLClientEvents[K]): void {
    if (!this.eventListeners.has(event)) {
      this.eventListeners.set(event, new Set());
    }
    this.eventListeners.get(event)?.add(listener as (...args: unknown[]) => void);
  }

  /**
   * Remove an event listener
   */
  off<K extends keyof XRPLClientEvents>(event: K, listener: XRPLClientEvents[K]): void {
    this.eventListeners.get(event)?.delete(listener as (...args: unknown[]) => void);
  }

  /**
   * Emit an event to all listeners
   */
  private emit<K extends keyof XRPLClientEvents>(
    event: K,
    ...args: Parameters<XRPLClientEvents[K]>
  ): void {
    const listeners = this.eventListeners.get(event);
    if (listeners) {
      for (const listener of listeners) {
        try {
          listener(...args);
        } catch (error) {
          this.logger.error(`Error in event listener for ${event}:`, error);
        }
      }
    }
  }

  /**
   * Connect to the XRPL
   */
  async connect(): Promise<Client> {
    if (this.client?.isConnected()) {
      this.logger.debug('Already connected to XRPL');
      return this.client;
    }

    this.connectionState = 'connecting';
    this.logger.info(`Connecting to XRPL at ${this.config.wsUrl}`);

    try {
      const clientOptions: ClientOptions = {
        timeout: this.config.connectionTimeout,
      };

      this.client = new Client(this.config.wsUrl, clientOptions);
      this.setupEventHandlers();

      await this.client.connect();

      this.connectionState = 'connected';
      this.reconnectAttempts = 0;
      this.logger.info('Successfully connected to XRPL');
      this.emit('connected');

      return this.client;
    } catch (error) {
      this.connectionState = 'disconnected';
      this.logger.error('Failed to connect to XRPL:', error);
      throw error;
    }
  }

  /**
   * Disconnect from the XRPL
   */
  async disconnect(): Promise<void> {
    this.clearReconnectTimer();

    if (!this.client) {
      this.logger.debug('No client to disconnect');
      return;
    }

    try {
      if (this.client.isConnected()) {
        await this.client.disconnect();
        this.logger.info('Disconnected from XRPL');
      }
    } catch (error) {
      this.logger.error('Error during disconnect:', error);
    } finally {
      this.client = null;
      this.connectionState = 'disconnected';
    }
  }

  /**
   * Get the XRPL client, connecting if necessary
   */
  async getClient(): Promise<Client> {
    if (this.client?.isConnected()) {
      return this.client;
    }
    return this.connect();
  }

  /**
   * Set up event handlers for the XRPL client
   */
  private setupEventHandlers(): void {
    if (!this.client) {
      return;
    }

    this.client.on('error', (errorCode: string, errorMessage: string) => {
      this.logger.error(`XRPL Error [${errorCode}]: ${errorMessage}`);
      this.emit('error', errorCode, errorMessage);
    });

    this.client.on('disconnected', (code: number) => {
      this.logger.warn(`XRPL disconnected with code: ${code}`);
      this.connectionState = 'disconnected';
      this.emit('disconnected', code);

      // Attempt reconnection for unexpected disconnections
      // Code 1000 indicates normal closure
      if (code !== 1000) {
        this.scheduleReconnect();
      }
    });

    this.client.on('connected', () => {
      this.logger.debug('XRPL client connected event received');
    });
  }

  /**
   * Schedule a reconnection attempt
   */
  private scheduleReconnect(): void {
    if (this.reconnectAttempts >= this.config.maxReconnectAttempts) {
      this.logger.error(`Max reconnection attempts (${this.config.maxReconnectAttempts}) reached`);
      this.emit('reconnectFailed');
      return;
    }

    this.clearReconnectTimer();
    this.connectionState = 'reconnecting';

    const delay = this.config.reconnectDelay * 2 ** this.reconnectAttempts;
    this.reconnectAttempts++;

    this.logger.info(
      `Scheduling reconnection attempt ${this.reconnectAttempts}/${this.config.maxReconnectAttempts} in ${delay}ms`,
    );
    this.emit('reconnecting', this.reconnectAttempts);

    this.reconnectTimer = setTimeout(() => {
      this.reconnect().catch((error) => {
        this.logger.error('Reconnection failed:', error);
      });
    }, delay);
  }

  /**
   * Attempt to reconnect to the XRPL
   */
  private async reconnect(): Promise<void> {
    try {
      // Clean up existing client
      if (this.client) {
        try {
          await this.client.disconnect();
        } catch {
          // Ignore disconnect errors during reconnection
        }
        this.client = null;
      }

      await this.connect();
      this.logger.info('Reconnection successful');
    } catch (error) {
      this.logger.error('Reconnection failed:', error);
      this.scheduleReconnect();
    }
  }

  /**
   * Clear any pending reconnection timer
   */
  private clearReconnectTimer(): void {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
  }
}

/**
 * Singleton instance of the XRPL client manager
 */
let clientManagerInstance: XRPLClientManager | null = null;

/**
 * Get the singleton XRPL client manager instance
 */
export function getClientManager(config?: Partial<XRPLConfig>, logger?: Logger): XRPLClientManager {
  if (!clientManagerInstance) {
    clientManagerInstance = new XRPLClientManager(config, logger);
  }
  return clientManagerInstance;
}

/**
 * Get the XRPL client, connecting if necessary
 *
 * This is a convenience function for simple use cases.
 * For more control, use getClientManager().
 */
export async function getClient(): Promise<Client> {
  return getClientManager().getClient();
}

/**
 * Disconnect the singleton client
 *
 * This is a convenience function for simple use cases.
 * For more control, use getClientManager().
 */
export async function disconnectClient(): Promise<void> {
  if (clientManagerInstance) {
    await clientManagerInstance.disconnect();
  }
}

/**
 * Reset the singleton client manager (useful for testing)
 */
export function resetClientManager(): void {
  if (clientManagerInstance) {
    clientManagerInstance.disconnect().catch(() => {
      // Ignore errors during reset
    });
  }
  clientManagerInstance = null;
}
