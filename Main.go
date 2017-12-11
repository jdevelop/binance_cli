package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"time"

	"github.com/jdevelop/binance-transfer/transfer"
	"github.com/spf13/viper"
)

type Config struct {
	BinanceAPI struct {
		PubKey  string `mapstructure:"public,omitempty"`
		PrivKey string `mapstructure:"private,omitempty"`
	} `mapstructure:"binanceapi,omitempty"`
}

func main() {

	asset := flag.String("asset", "", "Asset code (IOTA, ETH ...)")
	dst := flag.String("wallet", "", "Destination wallet address")
	amount := flag.Float64("amount", -1, "Amount (5 digits after decimal point)")
	depList := flag.Bool("deposits", false, "List deposits")
	wthList := flag.Bool("withdrawals", false, "List withdrawals")
	status := flag.Bool("status", false, "User status")
	force := flag.Bool("force", false, "Force deposit/withdrawal operation (DANGEROUS!)")

	flag.Parse()

	user, err := user.Current()

	if err != nil {
		log.Fatal(err)
	}

	viper.SetConfigFile(user.HomeDir + "/.binance")
	viper.SetConfigType("json")
	err = viper.ReadInConfig()

	if err != nil {
		log.Fatal(err)
	}

	var c Config

	err = viper.Unmarshal(&c)

	if err != nil {
		log.Fatal(err)
	}

	if c.BinanceAPI.PrivKey == "" || c.BinanceAPI.PubKey == "" {
		log.Fatal("No binance public/private keys provided")
	}

	b := transfer.MakeBinance(&http.Client{
		Timeout: 20 * time.Second,
	}, c.BinanceAPI.PrivKey, c.BinanceAPI.PubKey)

	if *depList {
		fmt.Println("DEPOSIT")
		recs, err := b.DepositHistory()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(recs)
	} else if *wthList {
		fmt.Println("WITHDRAW")
		recs, err := b.WithdrawHistory()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(recs)
	} else if *asset != "" && *dst != "" && *amount > 0 {
		withdraw := func() {
			err := b.Withdraw(*asset, *dst, *amount)
			if err != nil {
				log.Fatal(err)
			}
		}
		if *force {
			withdraw()
		} else {
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Printf("Are you sure you want to withdraw %.5f of %s to %s? (y/N) \n", *amount, *asset, *dst)
			scanner.Scan()
			switch scanner.Text() {
			case "Y":
				fallthrough
			case "y":
				withdraw()
			default:
				fmt.Println("Transaction cancelled")
			}
		}
	} else if *status {
		s, err := b.Status()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(s)
	} else {
		flag.Usage()
	}

}
