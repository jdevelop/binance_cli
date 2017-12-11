## [Binance](https;//binance.com) command-line interface

### Overview

Very basic command-line interface that allows to check the account status, withdraw coins and list the transaction history.


### Configuration

The application requires API keys from [Binance Account](https://www.binance.com/userCenter/myAccount.html)

1. Click on API settings
1. Create the new key (give it a name "Local" for example)
1. Make sure that the **Enable Withdrawals** checkbox is checked in the settings of the key, if you want to be able to work with the deposits/withdrawals functions!
1. Copy the API key into `public` parameter of the congiguration file (see example below)
1. Copy the Secret key into `private` parameter of the congiguration file (see example below) 

The configuration file has to be placed into `$HOME/.binance` and the content should look like

```json
{
  "binanceapi": {
    "public": "public-key-here",
    "private": "private-key-here"
  }
}

```

### Running
Available arguments:

```
Usage of binance_cli:
  -amount float
        Amount (5 digits after decimal point)
  -asset string
        Asset code (IOTA, ETH ...)
  -deposits
        List deposits
  -force
        Force deposit/withdrawal operation (DANGEROUS!)
  -interval duration
        Retry interval (5m = 5 minutes, 1h = 1 hour)
  -retries int
        Force retries on unsuccessful operation, should be used with -force
  -status
        User status
  -wallet string
        Destination wallet address
  -withdrawals
        List withdrawals
```

Sample transfer of 1 IOTA to one of my wallets ( _if you don't mind to say "thanks" in this way_  ):

```
./binance_cli -amount 1.00000 -asset IOTA -wallet MUPSEKCTPPLQFHUYLDMNQVBOXWFN9IIXOCBDHSBPQHKFMDYLVAOLEOQRMWLNZFF9N9Z9GLVKTLCCPWSTBODBBCRNHW
```