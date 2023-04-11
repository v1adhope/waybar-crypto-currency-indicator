package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/v1adhope/waybar-crypto-currency-indicator/structure"
)

const (
	_envAPI          = "CRYPTO_CURRENCY_API"
	_envFavoriteCoin = "CRYPTO_CURRENCY_FV"
	_envWatchList    = "CRYPTO_CURRENCY_WL"
)

type waybar struct {
	Text    string `json:"text"`
	Tooltip string `json:"tooltip"`
}

func main() {
	convertID, symbolIDs := "2781", "1,1027,6636"

	api := os.Getenv(_envAPI)
	if env := os.Getenv(_envFavoriteCoin); env != "" {
		for _, v := range env {
			if byte(v) == 44 || (byte(v) < 58 && byte(v) > 47) {
				continue
			}
			log.Fatalf("incorrect %s", _envFavoriteCoin)
		}
		convertID = env
	}
	if env := os.Getenv(_envWatchList); env != "" {
		for _, v := range env {
			if v == 44 || (v < 58 && v > 47) {
				continue
			}
			log.Fatalf("incorrect %s", _envWatchList)
		}
		symbolIDs = env
	}

	client, query := &http.Client{}, &url.Values{}

	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", api)

	query.Add("id", symbolIDs)
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
		w      waybar
		b      strings.Builder
		maxLen int
	)

	symbols := strings.Split(symbolIDs, ",")
	for _, v := range symbols {
		symLen := len(data.Data[v].Symbol)

		if maxLen < symLen {
			maxLen = symLen
		}
	}

	symFirst := symbols[0]
	if data.Data[symFirst].IsFiat == 0 {
		w.Text = strconv.FormatFloat(data.Data[symFirst].Quote[convertID].Price, 'f', 3, 64)
	} else {
		w.Text = strconv.FormatFloat(1/data.Data[symFirst].Quote[convertID].Price, 'f', 3, 64)
	}

	fmtRecord := fmt.Sprintf("%%-%ds ", maxLen)
	fmtTitle := fmt.Sprintf("<b>%%-%ds %%s</b>\n", maxLen)

	fmt.Fprintf(&b, fmtTitle, "Ticker", "Price")

	for _, v := range symbols[1:] {
		fmt.Fprintf(&b, fmtRecord, data.Data[v].Symbol)
		fmt.Fprintf(&b, "%0.3f\n", data.Data[v].Quote[convertID].Price)
	}

	w.Tooltip = strings.TrimSuffix(b.String(), "\n")

	json, err := json.Marshal(w)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(json))
}
