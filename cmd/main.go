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

		dataBySymIndx := make([]struct{ ticker, price, change24 string }, len(symIDs))
		maxColsLen := make([]int, len(titleParts))

		for colIndx, col := range titleParts {
			maxColsLen[colIndx] = len(col)

			for simIndx, symID := range symIDs[1:] {
				var symLen int

				switch colIndx {
				case one:
					dataBySymIndx[simIndx].ticker = data.Data[symID].Symbol
					symLen = len(dataBySymIndx[simIndx].ticker)
				case two:
					dataBySymIndx[simIndx].price = fmt.Sprintf("%0.3f", data.Data[symID].Quote[convertID].Price)
					symLen = len(dataBySymIndx[simIndx].price)
				case three:
					dataBySymIndx[simIndx].change24 = fmt.Sprintf("%0.2f", data.Data[symID].Quote[convertID].PercentChange24h)
					symLen = len(dataBySymIndx[simIndx].change24)
				}

				if maxColsLen[colIndx] < symLen {
					maxColsLen[colIndx] = symLen
				}
			}

			if colIndx != one {
				maxColsLen[colIndx]++ // Add an indent between columns
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
			fmt.Fprintf(&b, fmtPattern(one), dataBySymIndx[i].ticker)
			fmt.Fprintf(&b, fmtPattern(two), dataBySymIndx[i].price)
			fmt.Fprintf(&b, fmtPattern(three), dataBySymIndx[i].change24)
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
