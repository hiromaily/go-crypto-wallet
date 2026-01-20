/**
 * XRPL gRPC Server Entry Point
 *
 * Main entry point for the XRPL gRPC server.
 * Initializes the XRPL client, starts the ConnectRPC server, and sets up graceful shutdown.
 */

import { getConfig } from './config';
import { createServerFetchHandler } from './server';
import { disconnectClient, getClientManager } from './xrpl';

/** Server instance for graceful shutdown */
let server: ReturnType<typeof Bun.serve> | null = null;

/**
 * Initialize the XRPL client and set up event handlers
 */
async function initializeXRPLClient(): Promise<void> {
  const clientManager = getClientManager();

  // Set up event handlers for logging
  clientManager.on('connected', () => {
    console.log('XRPL client connected successfully');
  });

  clientManager.on('disconnected', (code) => {
    console.log(`XRPL client disconnected with code: ${code}`);
  });

  clientManager.on('error', (errorCode, errorMessage) => {
    console.error(`XRPL error [${errorCode}]: ${errorMessage}`);
  });

  clientManager.on('reconnecting', (attempt) => {
    console.log(`XRPL reconnection attempt ${attempt}`);
  });

  clientManager.on('reconnectFailed', () => {
    console.error('XRPL reconnection failed after max attempts');
  });

  // Connect to XRPL
  await clientManager.connect();
}

/**
 * Start the ConnectRPC server using Bun.serve()
 *
 * @param host - Server host address
 * @param port - Server port number
 * @returns Bun server instance
 */
function startServer(host: string, port: number): ReturnType<typeof Bun.serve> {
  const fetchHandler = createServerFetchHandler();

  return Bun.serve({
    hostname: host,
    port,
    fetch: fetchHandler,
  });
}

/**
 * Set up graceful shutdown handlers
 */
function setupGracefulShutdown(): void {
  const shutdown = async (signal: string) => {
    console.log(`\nReceived ${signal}, shutting down gracefully...`);

    let exitCode = 0;
    try {
      // Stop the gRPC server first
      if (server) {
        server.stop();
        console.log('gRPC server stopped');
      }

      // Then disconnect from XRPL
      await disconnectClient();
      console.log('XRPL client disconnected');
    } catch (error) {
      console.error('Error during shutdown:', error);
      exitCode = 1;
    }

    process.exit(exitCode);
  };

  process.on('SIGINT', () => shutdown('SIGINT'));
  process.on('SIGTERM', () => shutdown('SIGTERM'));
}

/**
 * Main entry point
 */
async function main(): Promise<void> {
  const config = getConfig();

  console.log('XRPL gRPC Server starting...');
  console.log(`  XRPL WebSocket URL: ${config.xrpl.wsUrl}`);
  console.log(`  gRPC Server: ${config.server.host}:${config.server.port}`);

  setupGracefulShutdown();

  try {
    // Initialize XRPL client
    await initializeXRPLClient();
    console.log('XRPL client initialized successfully');

    // Start ConnectRPC server
    server = startServer(config.server.host, config.server.port);
    console.log(`gRPC server listening on ${config.server.host}:${config.server.port}`);

    console.log('XRPL gRPC Server is ready to accept connections');
  } catch (error) {
    console.error('Failed to initialize XRPL gRPC Server:', error);
    process.exit(1);
  }
}

// Run the main function
main().catch((error) => {
  console.error('Unhandled error:', error);
  process.exit(1);
});

export { getConfig } from './config';
export { createRouter, createServerFetchHandler } from './server';
export { accountService, addressService } from './services';
export { disconnectClient, getClient, getClientManager } from './xrpl';
