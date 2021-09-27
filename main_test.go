package main

import (
	"testing"

	"cryptowidget.stuff/local/test/testfuncs"
	"cryptowidget.stuff/local/types"
)

func TestConfig(t *testing.T) {
	configFile := "config.json.example"

	cfg, err := getConf(configFile)
	if err != nil {
		t.Error("Error retrieving config: %w", err)
	}

	if cfg.URL != "https://api.coingecko.com/api/v3/coins/ergo?tickers=true&market_data=true&community_data=true&developer_data=true&sparkline=true" {
		t.Errorf("URL not found as expected. Received: %s", cfg.URL)
	}

	if cfg.Holdings != 10 {
		t.Errorf("Holdings not as expected. Received: %f", cfg.Holdings)
	}
}

func TestPrepareOutput(t *testing.T) {
	var currentUSDPrice float64 = 14.1
	var currentUSDPriceString string = "14.10"
	var currentAUDPrice float64 = 10.1
	var currentHoldingsValue string = "101.00"

	coin := types.CoinsID{
		MarketData: &types.MarketDataItem{
			CurrentPrice: types.AllCurrencies{
				"usd": currentUSDPrice,
				"aud": currentAUDPrice,
			},
		},
	}

	w := Widget{
		cfg: &Config{
			Holdings: 10,
			AdditionalCurrencies: []string{
				"aud",
			},
		},
	}

	output := w.prepareOutput(&coin)

	if output.CurrentPrice != currentUSDPriceString {
		t.Errorf("USD Price not as expected. Got: %s; expected: %s", output.CurrentPrice, currentUSDPriceString)
	}

	if output.CurrentValues["aud"] != currentHoldingsValue {
		t.Errorf("AUD Price not as expected. Got: %s; expected: %s", output.CurrentValues["aud"], currentHoldingsValue)
	}
}

func TestGetCoin(t *testing.T) {
	testJsonResponse := testfuncs.GetTestJsonResponse("coin_response.json")

	client := testfuncs.GetNewTestClient(testJsonResponse)

	w := Widget{
		cfg: &Config{
			URL: "someurl.com",
		},
		Client: client,
	}

	got, err := w.getCoin()
	if err != nil {
		t.Errorf("Error with func getCoin. %s", err)
	}

	expected := types.CoinsID{
		MarketCapRank: 123,
		MarketData: &types.MarketDataItem{
			CurrentPrice: types.AllCurrencies{
				"usd": 13.81,
			},
		},
	}

	if got.MarketCapRank != expected.MarketCapRank {
		t.Errorf("Marketcap rank not as expected. Expected: %d; Received: %d", got.MarketCapRank, expected.MarketCapRank)
	}

	if got.MarketData.CurrentPrice["usd"] != expected.MarketData.CurrentPrice["usd"] {
		t.Errorf("CurrentPrice not as expected. Expected: %f; Received: %f", got.MarketData.CurrentPrice["usd"], expected.MarketData.CurrentPrice["usd"])
	}
}
