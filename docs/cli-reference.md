# DeepCoin CLI Reference

## Configuration

The CLI reads credentials and base URL from environment variables.

| Purpose | Preferred Env Var | Alias | Description |
|---------|-------------------|-------|-------------|
| API key | `DEEPCOIN_API_KEY` | `DC_API_KEY` | API key for authentication |
| Secret key | `DEEPCOIN_SECRET_KEY` | `DC_SECRET_KEY` | Secret key for signing |
| Passphrase | `DEEPCOIN_PASSPHRASE` | `DC_PASSPHRASE` | API passphrase |
| Base URL | `DEEPCOIN_BASE_URL` | `DC_BASE_URL` | API base URL (default: `https://api.deepcoin.com`) |

Global command flags:

| Option | Description |
|--------|-------------|
| `--version` | Show version |
| `--help` | Show help |

Compatibility aliases are also supported: `DC_API_KEY`, `DC_SECRET_KEY`, `DC_PASSPHRASE`, and `DC_BASE_URL`.

## Agent Discovery

```bash
dcli list-tools
```

Prints a machine-readable JSON inventory of stable commands for Deepcoin skills. Agents should use this command surface instead of assembling custom API scripts.

---

## `market` — Public Market Data

No authentication required.

### `market instruments`

List tradeable instruments.

```bash
dcli market instruments --inst-type SWAP
dcli market instruments --inst-type SPOT --inst-id BTC-USDT
```

| Option | Required | Description |
|--------|----------|-------------|
| `--inst-type` | Yes | `SPOT` or `SWAP` |
| `--inst-id` | No | Filter by instrument ID |
| `--json` | No | Raw JSON output |

### `market tickers`

Get market tickers for all instruments.

```bash
dcli market tickers --inst-type SWAP
```

| Option | Required | Description |
|--------|----------|-------------|
| `--inst-type` | Yes | `SPOT` or `SWAP` |
| `--json` | No | Raw JSON output |

### `market ticker <INST_ID>`

Get ticker for a single instrument.

```bash
dcli market ticker BTC-USDT-SWAP
```

### `market orderbook <INST_ID>`

Get order book depth.

```bash
dcli market orderbook BTC-USDT-SWAP --sz 50
```

| Option | Required | Description |
|--------|----------|-------------|
| `--sz` | No | Depth levels, max 400 (default: 20) |
| `--json` | No | Raw JSON output |

### `market candles <INST_ID>`

Get K-line / candlestick data.

```bash
dcli market candles BTC-USDT-SWAP --bar 1H --limit 50
```

| Option | Required | Description |
|--------|----------|-------------|
| `--bar` | No | Interval: `1m/5m/15m/30m/1H/4H/12H/1D/1W/1M/1Y` (default: `1m`) |
| `--limit` | No | Number of candles, max 300 (default: 100) |
| `--after` | No | Pagination timestamp |
| `--json` | No | Raw JSON output |

### `market trades <INST_ID>`

Get recent trades.

```bash
dcli market trades BTC-USDT-SWAP --limit 20
```

| Option | Required | Description |
|--------|----------|-------------|
| `--limit` | No | Max 500 (default: 50) |
| `--product-group` | No | `Spot/Swap/SwapU` |
| `--json` | No | Raw JSON output |

### `market funding-rate`

Get current funding rates.

```bash
dcli market funding-rate --inst-type SwapU
dcli market funding-rate --inst-type SwapU --inst-id BTCUSDT
```

| Option | Required | Description |
|--------|----------|-------------|
| `--inst-type` | Yes | `SwapU` or `Swap` |
| `--inst-id` | No | Filter by instrument ID |

### `market funding-rate-history <INST_ID>`

Get historical funding rates.

```bash
dcli market funding-rate-history BTCUSDT --size 50
```

### `market book-spread <INST_ID>`

Get bid-ask spread information.

### `market step-margin <INST_ID>`

Get margin tier info (SWAP only).

### `market server-time`

Get server time.

### `market ping`

Check API connectivity.

---

## `trade` — Order Management

All commands require authentication.

### `trade place-order`

Place a new order.

```bash
# Limit order
dcli trade place-order \
  --inst-id BTC-USDT-SWAP \
  --td-mode isolated \
  --side buy \
  --ord-type limit \
  --sz 1 \
  --px 60000 \
  --pos-side long \
  --mrg-position merge

# Market order
dcli trade place-order \
  --inst-id BTC-USDT \
  --td-mode cash \
  --side buy \
  --ord-type market \
  --sz 0.01
```

| Option | Required | Description |
|--------|----------|-------------|
| `--inst-id` | Yes | Instrument ID (e.g. `BTC-USDT-SWAP`) |
| `--td-mode` | Yes | `isolated`, `cross`, or `cash` |
| `--side` | Yes | `buy` or `sell` |
| `--ord-type` | Yes | `market`, `limit`, `post_only`, `ioc` |
| `--sz` | Yes | Order size |
| `--px` | Conditional | Price (required for limit/post_only) |
| `--pos-side` | Conditional | `long` or `short` (required for SWAP) |
| `--mrg-position` | Conditional | `merge` or `split` (required for SWAP) |
| `--tp-trigger-px` | No | Take profit trigger price |
| `--sl-trigger-px` | No | Stop loss trigger price |
| `--cl-ord-id` | No | Custom order ID |
| `--reduce-only` | No | Reduce only flag |
| `--tgt-ccy` | No | `base_ccy` or `quote_ccy` (spot market orders) |

### `trade batch-orders`

Place multiple orders at once (max 5).

```bash
dcli trade batch-orders --orders '[{"instId":"BTC-USDT-SWAP","tdMode":"isolated","side":"buy","ordType":"limit","sz":"1","px":"60000","posSide":"long","mrgPosition":"merge"}]'
```

### `trade cancel-order`

Cancel an existing order.

```bash
dcli trade cancel-order --inst-id BTC-USDT-SWAP --ord-id 1000587866272245
```

### `trade batch-cancel`

Cancel multiple orders (max 50).

```bash
dcli trade batch-cancel --order-ids "123,456,789"
```

### `trade cancel-all`

Cancel all orders for a product group.

```bash
dcli trade cancel-all --product-group SwapU
```

### `trade amend-order`

Modify an existing order's price or volume.

```bash
dcli trade amend-order --order-id 123456 --price 61000
```

### `trade amend-order-sltp`

Modify TP/SL on an existing order.

```bash
dcli trade amend-order-sltp --order-id 123456 --tp-trigger-px 65000
```

### `trade get-order`

Get details of a specific order.

```bash
dcli trade get-order --inst-id BTC-USDT-SWAP --ord-id 123456
```

### `trade get-history-order`

Get a historical (finished) order by ID.

### `trade pending-orders`

List pending (open) orders.

```bash
dcli trade pending-orders --inst-id BTC-USDT-SWAP --limit 50
```

### `trade order-history`

Get order history.

```bash
dcli trade order-history --inst-type SWAP --state filled --limit 20
```

| Option | Required | Description |
|--------|----------|-------------|
| `--inst-type` | Yes | `SPOT` or `SWAP` |
| `--inst-id` | No | Filter by instrument |
| `--state` | No | `canceled` or `filled` |
| `--ord-type` | No | Order type filter |
| `--limit` | No | Max 100 (default: 20) |

### `trade batch-query`

Query multiple orders at once (max 5).

### `trade fills`

Get trade fill history.

```bash
dcli trade fills --inst-type SWAP --inst-id BTC-USDT-SWAP
```

### `trade trigger-order`

Place a trigger (conditional) order.

```bash
dcli trade trigger-order \
  --inst-id BTC-USDT-SWAP \
  --product-group SwapU \
  --side buy \
  --sz 1 \
  --trigger-price 58000 \
  --pos-side long
```

### `trade cancel-trigger`

Cancel a trigger order.

### `trade cancel-all-triggers`

Cancel all trigger orders.

### `trade trigger-pending`

List pending trigger orders.

### `trade trigger-history`

Get trigger order history.

### `trade set-position-sltp`

Set take-profit / stop-loss on a position.

```bash
dcli trade set-position-sltp \
  --inst-type SWAP \
  --inst-id BTC-USDT-SWAP \
  --pos-side long \
  --tp-trigger-px 70000 \
  --sl-trigger-px 55000
```

### `trade modify-position-sltp`

Modify TP/SL on a position.

### `trade cancel-position-sltp`

Cancel TP/SL on a position.

### `trade close-position`

Close positions by IDs.

```bash
dcli trade close-position --inst-id BTC-USDT-SWAP --product-group SwapU --position-ids "123,456"
```

### `trade batch-close-position`

Close all positions for an instrument.

### `trade trace-order`

Place a trace (trailing) order.

### `trade trace-orders`

List pending trace orders.

---

## `account` — Account & Portfolio

All commands require authentication.

### `account balance`

Get account balance.

```bash
dcli account balance
dcli account balance --inst-type SWAP --ccy USDT
```

| Option | Required | Description |
|--------|----------|-------------|
| `--inst-type` | No | `SPOT` or `SWAP` |
| `--ccy` | No | Currency filter (e.g. `USDT`) |

### `account all-balances`

Get unified balances across Deepcoin account types.

```bash
dcli account all-balances
dcli account all-balances --account-type funding,spot,swapU --ccy USDT,BTC --json
```

| Option | Required | Description |
|--------|----------|-------------|
| `--account-type` | No | Comma-separated account types: `funding`, `spot`, `swapU`, `swap`, `bonus`, `rebate`, `event`, `copyTrade`, `robot`, `all` |
| `--ccy` | No | Currency filter; comma-separated values are supported (e.g. `USDT,BTC`) |
| `--json` | No | Output raw JSON |

### `account positions`

Get open positions.

```bash
dcli account positions --inst-type SWAP
```

### `account bills`

Get account bill history.

```bash
dcli account bills --inst-type SWAP --limit 50
```

### `account set-leverage`

Set leverage for an instrument.

```bash
dcli account set-leverage --inst-id BTC-USDT-SWAP --lever 10 --mgn-mode isolated
```

### `account uid`

Get current account UID.

### `account sub-accounts`

List sub-accounts.

### `account sub-account-balance`

Get total balance across sub-accounts.

### `account sub-account-transfer`

Transfer between sub-accounts.

```bash
dcli account sub-account-transfer \
  --from-uid 111 --to-uid 222 \
  --from-id 1 --to-id 1 \
  --amount 100 --coin USDT
```

### `account sub-account-transfer-records`

Get sub-account transfer records.

### `account deposit-list`

Get deposit history.

### `account withdraw-list`

Get withdrawal history.

### `account transfer`

Transfer assets between accounts.

```bash
dcli account transfer --currency-id USDT --amount 100 --from-id 1 --to-id 2
```

### `account recharge-chains`

Get supported deposit chains for a currency.

### `account internal-transfer-support`

Get supported coins for internal transfer.

### `account internal-transfer`

Make an internal transfer.

### `account internal-transfer-history`

Get internal transfer history.

### `account rebate-summary`

Get rebate summary.

### `account affiliates`

Get affiliate list.

### `account trade-stats-daily`

Get daily trade statistics.

### `account trade-stats-total`

Get total trade statistics.

---

## `withdrawal` — On-Chain Withdrawals

All commands require authentication. Confirm externally before `create` or `cancel`; the CLI intentionally executes the requested command directly.

### `withdrawal config`

Get aggregated withdrawal config, optionally including whitelist addresses.

```bash
dcli withdrawal config --ccy USDT --include-addresses true
```

### `withdrawal assets`

List withdrawable assets.

```bash
dcli withdrawal assets --ccy USDT
```

### `withdrawal chains`

List withdrawal chains for a coin.

```bash
dcli withdrawal chains --ccy USDT
```

### `withdrawal addresses`

List withdrawal whitelist addresses for a coin.

```bash
dcli withdrawal addresses --ccy USDT
```

### `withdrawal records`

List withdrawal records.

```bash
dcli withdrawal records --coin USDT --page 1 --size 20
dcli withdrawal records --wd-id 123456
```

### `withdrawal status`

Get a single withdrawal status.

```bash
dcli withdrawal status --wd-id 123456 --ccy USDT
```

### `withdrawal create`

Create an on-chain withdrawal.

```bash
dcli withdrawal create \
  --ccy USDT \
  --chain USDT-TRC20 \
  --amt 10 \
  --address-id 123456 \
  --account-types funding
```

### `withdrawal cancel`

Cancel an on-chain withdrawal.

```bash
dcli withdrawal cancel --wd-id 123456 --ccy USDT
```

---

## `copytrade` — Copy Trading

All commands require authentication.

### `copytrade leader-settings`

Update leader settings.

```bash
dcli copytrade leader-settings --status 1
```

### `copytrade support-contracts`

Get supported copy trading contracts.

### `copytrade set-contracts`

Set copy trading contracts.

```bash
dcli copytrade set-contracts --contracts "BTCUSDT,ETHUSDT"
```

### `copytrade followers`

Get follower list and stats.

```bash
dcli copytrade followers --status 1
```

### `copytrade leader-positions`

Get leader's current positions.

### `copytrade position-type`

Get current position type (hedge/one-way).

### `copytrade set-position-type`

Update position type.

```bash
dcli copytrade set-position-type --type 1  # Hedge mode
```

### `copytrade estimated-profit`

Get estimated profit from followers.

### `copytrade history-profit`

Get historical profit from copy trading.

---

## `strategy` — DSL Strategy Orders

All commands require authentication.

### `strategy backtest`

Run a strategy backtest.

```bash
# From file
dcli strategy backtest \
  --symbol BTC-USDT-SWAP \
  --from-ts 2025-01-01T00:00:00Z \
  --to-ts 2025-03-01T00:00:00Z \
  --dsl @my_strategy.json

# Inline JSON
dcli strategy backtest \
  --symbol BTC-USDT-SWAP \
  --from-ts 2025-01-01T00:00:00Z \
  --to-ts 2025-03-01T00:00:00Z \
  --dsl '{"version":"1.0","indicators":[{"name":"BOLL","params":{"period":20,"std_dev":2},"conditions":[{"field":"close","op":"cross_above","ref":"upper"}]}],"then":{"entry":{"side":"buy","posSide":"long","logic":"AND"},"exit":{"side":"sell","posSide":"long","logic":"AND"}},"risk":{"stop_loss":{"percent":2},"take_profit":{"percent":5}},"execution":{"order_type":"market"}}'
```

| Option | Required | Description |
|--------|----------|-------------|
| `--symbol` | Yes | Symbol (e.g. `BTC-USDT-SWAP`) |
| `--from-ts` | Yes | Start time (ISO 8601) |
| `--to-ts` | Yes | End time (ISO 8601) |
| `--dsl` | Yes | DSL JSON string or `@filepath` |

**Supported Indicators:** `BOLL`, `MA`, `EMA`, `KDJ`, `RSI`, `WR`

**Condition Operators:** `>=`, `<=`, `>`, `<`, `==`, `cross_above`, `cross_below`

### `strategy dsl-trigger-order`

Place a live DSL-driven trigger order.

```bash
dcli strategy dsl-trigger-order \
  --symbol BTC-USDT-SWAP \
  --trade-mode isolated \
  --mrg-position merge \
  --dsl @my_strategy.json
```

| Option | Required | Description |
|--------|----------|-------------|
| `--symbol` | Yes | Symbol |
| `--trade-mode` | Yes | `isolated` or `cross` |
| `--mrg-position` | Yes | `merge` or `split` |
| `--dsl` | Yes | DSL JSON string or `@filepath` |
