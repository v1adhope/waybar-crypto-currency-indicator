package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/v1adhope/waybar-crypto-currency-indicator/structure"
)

const (
	_envAPI          = "CRYPTO_CURRENCY_API"
	_envFavoriteCoin = "CRYPTO_CURRENCY_FV"
	_envWatchList    = "CRYPTO_CURRENCY_WL"
)

// NOTE: Output structure. Read more
// https://github.com/Alexays/Waybar/wiki/Module:-Custom#return-type
type waybar struct {
	Text    string `json:"text"`
	Tooltip string `json:"tooltip,omitempty"`
}

func main() {
	convertID, rawSymIDs := "2781", "1,1027,6636"

	api := os.Getenv(_envAPI)

	err := getEnv(_envFavoriteCoin, &convertID)
	if err != nil {
		log.Fatal(err)
	}

	err = getEnv(_envWatchList, &rawSymIDs)
	if err != nil {
		log.Fatal(err)
	}

	client, query := &http.Client{}, &url.Values{}

	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", api)

	query.Add("id", rawSymIDs)
	query.Add("convert_id", convertID)

	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	var data structure.Data

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	var (
		w waybar
		b strings.Builder
	)

	symIDs := strings.Split(rawSymIDs, ",")

	symFirst := symIDs[0]
	if data.Data[symFirst].IsFiat == 0 {
		w.Text = strconv.FormatFloat(data.Data[symFirst].Quote[convertID].Price, 'f', 3, 64)
	} else {
		w.Text = strconv.FormatFloat(1/data.Data[symFirst].Quote[convertID].Price, 'f', 3, 64)
	}

	if len(symIDs) > 1 {
		titleParts := []string{
			"Ticker",
			"Price",
			"24h%",
		}

		const ( // Column numbers
			one = iota
			two
			three
		)

		dataBySymID := make([]struct{ ticker, price, change24 string }, len(symIDs))
		maxColsLen := make([]int, len(titleParts))

		for i, col := range titleParts {
			maxColsLen[i] = len(col)

			for i2, symID := range symIDs[1:] {
				var symLen int

				switch i {
				case one:
					dataBySymID[i2].ticker = data.Data[symID].Symbol
					symLen = len(dataBySymID[one].ticker)
				case two:
					dataBySymID[i2].price = fmt.Sprintf("%0.3f", data.Data[symID].Quote[convertID].Price)
					symLen = len(dataBySymID[one].price)
				case three:
					dataBySymID[i2].change24 = fmt.Sprintf("%0.2f", data.Data[symID].Quote[convertID].PercentChange24h)
					symLen = len(dataBySymID[i2].change24)
				}

				if maxColsLen[i] < symLen {
					maxColsLen[i] = symLen
				}
			}

			if i != one {
				maxColsLen[i]++ // Add an indent between columns
			}
		}

		fmtPattern := func(colNum int) string {
			if colNum == one {
				return fmt.Sprintf("%%-%ds", maxColsLen[colNum])
			} else {
				return fmt.Sprintf("%%%ds", maxColsLen[colNum])
			}
		}

		fmt.Fprint(&b, "<b>")
		for colNum, title := range titleParts {
			fmt.Fprintf(&b, fmtPattern(colNum), title)
		}
		fmt.Fprint(&b, "</b>\n")

		for i := range symIDs[1:] {
			fmt.Fprintf(&b, fmtPattern(one), dataBySymID[i].ticker)
			fmt.Fprintf(&b, fmtPattern(two), dataBySymID[i].price)
			fmt.Fprintf(&b, fmtPattern(three), dataBySymID[i].change24)
			fmt.Fprintln(&b)
		}
	}

	t, err := time.Parse(time.RFC3339Nano, data.Status.Timestamp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(&b, "\n<b>Timestamp</b>\n%s", t.Local().Format(time.DateTime))

	w.Tooltip = b.String()

	json, err := json.Marshal(w)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(json))
}

func getEnv(key string, placeholder *string) error {
	if env := os.Getenv(key); env != "" {
		for _, v := range env {
			if byte(v) == 44 || (byte(v) < 58 && byte(v) > 47) {
				continue
			}

			return errors.New(fmt.Sprintf("incorrect %s", _envFavoriteCoin))
		}
		*placeholder = env
	}

	return nil
}
