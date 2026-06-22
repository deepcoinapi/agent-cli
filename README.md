# DeepCoin Agent CLI

Command-line tool for interacting with the [DeepCoin](https://api.deepcoin.com) exchange. Built for AI agents and human traders alike.

## Features

- **Market Data** — tickers, orderbook, candles, trades, funding rates
- **Trading** — place/cancel/amend orders, trigger orders, TP/SL, batch operations
- **Account** — balances, positions, leverage, bills, sub-accounts, asset transfers
- **Withdrawal** — on-chain withdrawal config, whitelist, create, cancel, status, records
- **Copy Trading** — leader/follower management, position tracking, profit stats
- **Strategy** — DSL-based automated trading with backtesting

## Agent-Skill Contract

`dcli` is the stable execution layer for Deepcoin agent skills.

- Use `dcli list-tools` for a machine-readable command inventory.
- Agent skills should call CLI commands instead of generating temporary API clients or signing scripts.
- If a Deepcoin API workflow is not represented here, add a stable CLI command first, then reference it from the skill.

## Installation

```bash
go install github.com/deepcoinapi/agent-cli/cmd/dcli@latest
```

Or build from source:

```bash
git clone https://github.com/deepcoinapi/agent-cli.git
cd agent-cli
go build -o dcli .
```

## Configuration

Set your API credentials via environment variables or a `.env` file:

```bash
export DEEPCOIN_API_KEY=your-api-key
export DEEPCOIN_SECRET_KEY=your-secret-key
export DEEPCOIN_PASSPHRASE=your-passphrase
export DEEPCOIN_BASE_URL=https://api.deepcoin.com  # optional
```

The CLI also accepts `DC_API_KEY`, `DC_SECRET_KEY`, `DC_PASSPHRASE`, and `DC_BASE_URL` as compatibility aliases.

Or copy `.env.example` to `.env` and fill in your credentials.

## Quick Start

```bash
# Check connectivity
dcli market ping

# Get BTC ticker
dcli market ticker BTC-USDT-SWAP

# Get account balance
dcli account balance --inst-type SWAP

# Place a limit order
dcli trade place-order \
  --inst-id BTC-USDT-SWAP \
  --td-mode isolated \
  --side buy \
  --ord-type limit \
  --sz 1 \
  --px 60000 \
  --pos-side long \
  --mrg-position merge

# View open positions
dcli account positions --inst-type SWAP

# Run a strategy backtest
dcli strategy backtest \
  --symbol BTC-USDT-SWAP \
  --from-ts 2025-01-01T00:00:00Z \
  --to-ts 2025-03-01T00:00:00Z \
  --dsl @my_strategy.json
```

## Command Groups

| Group | Description |
|-------|-------------|
| `market` | Public market data (no auth required) |
| `trade` | Order management (auth required) |
| `account` | Account, positions, assets, sub-accounts (auth required) |
| `withdrawal` | On-chain withdrawal operations (auth required) |
| `copytrade` | Copy trading management (auth required) |
| `strategy` | DSL strategy orders & backtesting (auth required) |
| `list-tools` | Machine-readable stable command inventory for agents |

Run `dcli <group> --help` to see all commands in a group.

## JSON Output

All commands support `--json` flag for raw JSON output, useful for piping to other tools:

```bash
dcli market tickers --inst-type SWAP --json | jq '.data[0]'
```

## Full CLI Reference

See [docs/cli-reference.md](docs/cli-reference.md) for the complete command reference.

## Authentication

All private endpoints use HMAC-SHA256 signature authentication:

- `DC-ACCESS-KEY` — Your API Key
- `DC-ACCESS-SIGN` — `Base64(HMAC-SHA256(timestamp + method + requestPath + body, secretKey))`
- `DC-ACCESS-TIMESTAMP` — ISO 8601 timestamp
- `DC-ACCESS-PASSPHRASE` — Your passphrase

## License

MIT
