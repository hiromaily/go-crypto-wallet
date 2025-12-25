// Package http provides HTTP/REST API handlers for the wallet services.
//
// This package implements the interface adapter layer for HTTP requests,
// following Clean Architecture principles. It translates HTTP requests
// into use case calls and formats responses back to HTTP.
//
// The HTTP adapter allows CLI and HTTP interfaces to share the same
// use cases, following the Ports and Adapters (Hexagonal Architecture) pattern.
package http
