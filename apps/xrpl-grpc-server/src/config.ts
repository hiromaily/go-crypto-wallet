/**
 * Environment configuration for XRPL gRPC Server
 */

/**
 * XRPL network configuration
 */
export interface XRPLConfig {
  /** WebSocket URL for XRPL node connection */
  wsUrl: string;
  /** Connection timeout in milliseconds */
  connectionTimeout: number;
  /** Maximum reconnection attempts */
  maxReconnectAttempts: number;
  /** Delay between reconnection attempts in milliseconds */
  reconnectDelay: number;
}

/**
 * Server configuration
 */
export interface ServerConfig {
  /** gRPC server port */
  port: number;
  /** gRPC server host */
  host: string;
}

/**
 * Application configuration
 */
export interface AppConfig {
  xrpl: XRPLConfig;
  server: ServerConfig;
}

/**
 * Default XRPL WebSocket URLs by network
 */
const XRPL_WS_URLS = {
  mainnet: 'wss://xrplcluster.com',
  testnet: 'wss://s.altnet.rippletest.net:51233',
  devnet: 'wss://s.devnet.rippletest.net:51233',
} as const;

/**
 * Get the XRPL WebSocket URL based on environment
 */
function getXRPLWebSocketUrl(): string {
  const envUrl = process.env.XRP_WS_URL;
  if (envUrl) {
    return envUrl;
  }

  const network = process.env.XRP_NETWORK || 'testnet';
  return XRPL_WS_URLS[network as keyof typeof XRPL_WS_URLS] || XRPL_WS_URLS.testnet;
}

/**
 * Parse integer from environment variable with default
 */
function parseEnvInt(value: string | undefined, defaultValue: number): number {
  if (!value) {
    return defaultValue;
  }
  const parsed = Number.parseInt(value, 10);
  return Number.isNaN(parsed) ? defaultValue : parsed;
}

/**
 * Load configuration from environment variables
 */
export function loadConfig(): AppConfig {
  return {
    xrpl: {
      wsUrl: getXRPLWebSocketUrl(),
      connectionTimeout: parseEnvInt(process.env.XRP_CONNECTION_TIMEOUT, 20000),
      maxReconnectAttempts: parseEnvInt(process.env.XRP_MAX_RECONNECT_ATTEMPTS, 5),
      reconnectDelay: parseEnvInt(process.env.XRP_RECONNECT_DELAY, 1000),
    },
    server: {
      port: parseEnvInt(process.env.GRPC_PORT, 50051),
      host: process.env.GRPC_HOST || '0.0.0.0',
    },
  };
}

/**
 * Singleton configuration instance
 */
let configInstance: AppConfig | null = null;

/**
 * Get the application configuration (singleton)
 */
export function getConfig(): AppConfig {
  if (!configInstance) {
    configInstance = loadConfig();
  }
  return configInstance;
}

/**
 * Reset configuration (useful for testing)
 */
export function resetConfig(): void {
  configInstance = null;
}
