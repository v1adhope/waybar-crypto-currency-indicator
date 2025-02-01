package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/v1adhope/waybar-crypto-currency-indicator/structure"
)

const (
	_envApiToken     = "CRYPTO_CURRENCY_API"
	_envFavoriteCoin = "CRYPTO_CURRENCY_FV"
	_envWatchList    = "CRYPTO_CURRENCY_WL"
)

// Read more https://github.com/Alexays/Waybar/wiki/Module:-Custom#return-type
type module struct {
	Text    string `json:"text"`
	Tooltip string `json:"tooltip,omitempty"`
}

func main() {
	convertId, rawCoinIds := "2781", "1,1027,6636"

	apiToken := os.Getenv(_envApiToken)

	if err := fillVarByEnvKey(_envFavoriteCoin, &convertId); err != nil {
		log.Fatal(err)
	}

	if err := fillVarByEnvKey(_envWatchList, &rawCoinIds); err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest", nil)
	if err != nil {
		log.Fatal(err)
	}

	query := &url.Values{}
	query.Add("id", rawCoinIds)
	query.Add("convert_id", convertId)

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", apiToken)
	req.URL.RawQuery = query.Encode()

	data := structure.Data{}
	func() {
		client := &http.Client{}
		resp, err := client.Do(req)
		if resp != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			log.Fatal(err)
		}

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Fatal(err)
		}
	}()

	module := module{}
	coinIds := strings.Split(rawCoinIds, ",")

	coinIdFirst := coinIds[0]
	if data.Data[coinIdFirst].IsFiat == 0 {
		module.Text = fmt.Sprintf("%0.3f", data.Data[coinIdFirst].Quote[convertId].Price)
	} else {
		module.Text = fmt.Sprintf("%0.3f", 1/data.Data[coinIdFirst].Quote[convertId].Price)
	}

	buf := bytes.Buffer{}

	func() {
		w := tabwriter.NewWriter(&buf, 0, 0, 1, ' ', tabwriter.AlignRight)
		defer w.Flush()

		fmt.Fprintf(w, "%-6s\t%s\t%s\t\n", "Ticker", "Price", "24h%")

		for _, coinId := range coinIds[1:] {
			fmt.Fprintf(w,
				"%-6s\t%0.3f\t%0.2f\t\n",
				data.Data[coinId].Symbol,
				data.Data[coinId].Quote[convertId].Price,
				data.Data[coinId].Quote[convertId].PercentChange24h,
			)
		}
	}()

	t, err := time.Parse(time.RFC3339Nano, data.Status.Timestamp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(&buf, "\n <b>Timestamp</b>\n "+t.Local().Format(time.DateTime))

	module.Tooltip = buf.String()
	module.Tooltip = strings.Replace(module.Tooltip, "Ticker", "<b>Ticker</b>", 1)
	module.Tooltip = strings.Replace(module.Tooltip, "Price", "<b>Price</b>", 1)
	module.Tooltip = strings.Replace(module.Tooltip, "24h%", "<b>24h%</b>", 1)

	json, err := json.Marshal(module)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(json))
}

func fillVarByEnvKey(key string, placeholder *string) error {
	if env := os.Getenv(key); env != "" {
		for _, v := range env {
			if byte(v) == 44 || (byte(v) < 58 && byte(v) > 47) {
				continue
			}

			return errors.New("incorrect by key: " + key)
		}
		*placeholder = env
	}

	return nil
}
