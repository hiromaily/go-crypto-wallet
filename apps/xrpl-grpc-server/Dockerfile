# =============================================================================
# XRPL gRPC Server Dockerfile
# =============================================================================
# Multi-stage build for minimal production image
#
# Prerequisites:
#   Run `make proto-ts` from repository root to generate proto files before building.
#   The proto files use Edition 2024 which requires pre-generation.
#
# Build context: repository root (.)
# Build command:
#   make proto-ts  # Generate proto files first
#   docker compose -f compose.xrp.yaml build xrpl-grpc-server
#
# Run command:
#   docker run -p 50051:50051 -e XRP_WS_URL=ws://xrp-node:6006 xrpl-grpc-server
# =============================================================================

# -----------------------------------------------------------------------------
# Stage 1: Builder - Install dependencies and build
# -----------------------------------------------------------------------------
FROM oven/bun:1.3.6-slim AS builder

WORKDIR /app

# Copy package files for dependency installation
COPY apps/xrpl-grpc-server/package.json apps/xrpl-grpc-server/bun.lock ./

# Install all dependencies (including devDependencies for build)
RUN bun install --frozen-lockfile

# Copy source files and pre-generated proto files
# Note: Proto files must be generated locally before Docker build
#       Run `make proto-ts` from repository root
COPY apps/xrpl-grpc-server/src ./src
COPY apps/xrpl-grpc-server/tsconfig.json ./

# Verify proto files exist
RUN test -f src/gen/account_pb.ts || \
    (echo "Error: Proto files not found. Run 'make proto-ts' before building Docker image." && exit 1)

# Build the application
RUN bun build src/index.ts --outdir=dist --target=bun

# -----------------------------------------------------------------------------
# Stage 2: Production - Minimal runtime image
# -----------------------------------------------------------------------------
FROM oven/bun:1.3.6-slim AS production

# Create non-root user for security
RUN groupadd --gid 1001 xrpl && \
    useradd --uid 1001 --gid xrpl --shell /bin/false xrpl

WORKDIR /app

# Copy only production dependencies
COPY apps/xrpl-grpc-server/package.json apps/xrpl-grpc-server/bun.lock ./
RUN bun install --frozen-lockfile --production && \
    rm -rf /root/.bun/install/cache

# Copy built application and generated files
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/src/gen ./src/gen
COPY --from=builder /app/src ./src

# Change ownership to non-root user
RUN chown -R xrpl:xrpl /app

# Switch to non-root user
USER xrpl

# Expose gRPC port
EXPOSE 50051

# Environment variables with defaults
ENV GRPC_PORT=50051 \
    GRPC_HOST=0.0.0.0 \
    XRP_NETWORK=testnet \
    XRP_CONNECTION_TIMEOUT=20000

# Health check - verify the server is responding
# ConnectRPC servers return 404 for non-service paths, which indicates the server is running
HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD bun --eval "fetch('http://localhost:${GRPC_PORT}/').then(r => process.exit(r.status === 404 ? 0 : 1)).catch(() => process.exit(1))"

# Start the gRPC server
CMD ["bun", "run", "src/index.ts"]
