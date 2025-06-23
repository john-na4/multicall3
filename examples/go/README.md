# Go Multicall3 Example

This example demonstrates how to use Multicall3 with Go and the `go-ethereum` library.

## Prerequisites

- Go 1.21 or later
- A valid Ethereum mainnet RPC URL

## Setup

1. Set your Ethereum mainnet RPC URL as an environment variable:
   ```bash
   export MAINNET_RPC_URL="https://your-rpc-endpoint"
   ```

   Or create a `.env` file in this directory:
   ```
   MAINNET_RPC_URL=https://your-rpc-endpoint
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Running the Example

```bash
go run main.go
```

## What it does

This example demonstrates batching multiple Ethereum calls using Multicall3:

1. Fetches DAI token symbol
2. Fetches DAI token decimals
3. Fetches Vitalik's DAI balance

The example shows how to:
- Pack multiple contract calls into a single multicall
- Execute the multicall using the Multicall3 contract
- Unpack and format the results
- Handle both contract calls and direct RPC calls

## Output

The example will output something like:
```
Block Number: 12345678
DAI Symbol: DAI
DAI Decimals: 18
Vitalik's DAI balance: 1234.567890123456789000
```

## Key Differences from Other Examples

Unlike the Rust example which uses `ethers-rs` with built-in Multicall3 support, this Go example manually constructs the multicall by:
- Packing individual function calls
- Creating the multicall payload
- Unpacking the results manually

This approach gives you more control and understanding of how Multicall3 works under the hood. 