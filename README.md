# Usage

The indicator uses [CoinMarketCap API](https://coinmarketcap.com/api/)

Each run consumes 1 credit, for a free account 10,000 per month is given.

Three environmental variables are required for full work:
- `CRYPTO_CURRENCY_API`
- `CRYPTO_CURRENCY_FV`
- `CRYPTO_CURRENCY_WL`

`API` - after registration. `FV` - main coin ID, by default `2781`(USD). `WL` - list of coin IDs that you want to track are separated by commas, the first coin will be displayed on the panel, by default `1,1027,6636`(BTC,ETH,DOT).

```
curl -H "X-CMC_PRO_API_KEY: <YOUR_API>" -H "Accept: application/json" -d "symbol=<TICKER_WITH_A_COMMA_SEPARATING>" -G https://pro-api.coinmarketcap.com/v1/cryptocurrency/map | jq .
```
To search for coin IDs, consumes 1 credit `*required curl and jq`

## Sample

Updated every 30 minutes.

```
"custom/crypto-currency": {
    "max-length": 15,
    "return-type": "json",
    "format": "{}",
    "exec": "$HOME/.config/waybar/scripts/crypto-currency",
    "interval": 1800
}
```

# Preview

![preview](/assets/preview.png)
