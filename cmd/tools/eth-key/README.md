# eth-key

A simple command-line tool to extract Ethereum addresses from keystore files.

## Overview

`eth-key` reads Ethereum keystore files (JSON format) and extracts the corresponding Ethereum address from the file.

## Installation

### Building

```bash
cd cmd/tools/eth-key
go build -o eth-key main.go
```

Or using Makefile:

```bash
make build
```

## Usage

### Command-line Arguments

| Argument | Required | Default | Description |
|----------|----------|---------|-------------|
| `-keydir` | Yes | `./tmp` | Path to the keystore directory (preferably `tmp` directory) |
| `-keyfile` | Yes | - | Full path to the keystore file |
| `-password` | Yes | - | Password used to decrypt the keystore file |

### Examples

#### Basic Usage

```bash
./eth-key -keydir ./tmp -keyfile ./keys/key-sample.json -password password
```

#### Using Makefile

```bash
# Build and execute
make exec

# Or run directly
make run
```

### Output

The tool outputs the Ethereum address extracted from the keystore file to stdout:

```
0x71678cd07cfac46c2dc427f999abf46aae115925
```

## Notes

- All three arguments (`-keydir`, `-keyfile`, `-password`) are required
- The keystore file must be in standard Ethereum keystore format (JSON)
- An error will occur if the password is incorrect
- Currently, unmarshal errors may occur with some keystore formats (see FIXME comment in code)

## Troubleshooting

### Error: "args `keyfile` must not be empty"

The `-keyfile` argument is not specified. Please provide the full path to the keystore file.

### Error: "args `password` must not be empty"

The `-password` argument is not specified. Please provide the password used to decrypt the keystore file.

### Error: "cannot unmarshal object into Go struct field CryptoJSON.crypto.kdf of type string"

This error may occur with some keystore formats (especially Parity). It is recommended to use Geth-format keystore files.

## Related Files

- `main.go`: Main implementation file
- `Makefile`: Makefile for building and running
- `keys/`: Contains sample keystore files
