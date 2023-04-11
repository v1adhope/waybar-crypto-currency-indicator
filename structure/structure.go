package structure

type Data struct {
	Status struct {
		Timestamp    string `json:"timestamp"`
		ErrorCode    int    `json:"error_code"`
		ErrorMessage string `json:"error_message"` //NOTE: may be null
		Elapsed      int    `json:"elapsed"`
		CreditCount  int    `json:"credit_count"`
		Notice       string `json:"notice"` //NOTE: may be null
	} `json:"status"`
	Data map[string]struct {
		ID             int    `json:"id"`
		Name           string `json:"name"`
		Symbol         string `json:"symbol"`
		Slug           string `json:"slug"`
		NumMarketPairs int    `json:"num_market_pairs"`
		DateAdded      string `json:"date_added"`
		Tags           []struct {
			Slug     string `json:"slug"`
			Name     string `json:"name"`
			Category string `json:"category"`
		} `json:"tags"`
		MaxSupply         int     `json:"max_supply"`
		CirculatingSupply float64 `json:"circulating_supply"`
		TotalSupply       float64 `json:"total_supply"`
		IsActive          int     `json:"is_active"`
		InfiniteSupply    bool    `json:"infinite_supply"`
		Platform          struct {
			ID           int    `json:"id"`
			Name         string `json:"name"`
			Symbol       string `json:"symbol"`
			Slug         string `json:"slug"`
			TokenAddress string `json:"token_address"`
		} `json:"platform"`
		CMCRank                       int     `json:"cmc_rank"`
		IsFiat                        int     `json:"is_fiat"`
		SelfReportedCirculatingSupply float64 `json:"self_reported_circulating_supply"` //NOTE: may be null
		SelfReportedMarketCap         float64 `json:"self_reported_market_cap"`         //NOTE: may be null
		TVLRatio                      string  `json:"tvl_ratio"`                        //NOTE: may be null
		LastUpdated                   string  `json:"last_updated"`
		Quote                         map[string]struct {
			Price                 float64 `json:"price"`
			Volume24h             float64 `json:"volume_24h"`
			VolumeChange24h       float64 `json:"volume_change_24h"`
			PercentChange1h       float64 `json:"percent_change_1h"`
			PercentChange24h      float64 `json:"percent_change_24h"`
			PercentChange7d       float64 `json:"percent_change_7d"`
			PercentChange30d      float64 `json:"percent_change_30d"`
			PercentChange60d      float64 `json:"percent_change_60d"`
			PercentChange90d      float64 `json:"percent_change_90d"`
			MarketCap             float64 `json:"market_cap"`
			MarketcapDominance    float64 `json:"market_cap_dominance"`
			FullyDilutedMarketCap float64 `json:"fully_diluted_market_cap"`
			TVL                   string  `json:"tvl"` //NOTE: may be null
			LastUpdated           string  `json:"last_updated"`
		} `json:"quote"`
	} `json:"data"`
}
