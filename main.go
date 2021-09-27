package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"cryptowidget.stuff/local/types"
	"github.com/doublehops/go-common/strings"
)

type Config struct {
	URL                  string   `json:"url"`
	Coin                 string   `json:"coin"`
	Holdings             float64  `json:"holdings"`
	AdditionalCurrencies []string `json:"additionalCurrencies"`
}

type CoinResponse struct {
	CurrentPrice  string                 `json:"currentPrice"`
	Rank          uint16                 `json:"rank"`
	CurrentValues CurrencyCurrentValue `json:"currentValues"` // [aud] = 123.00
}

type Widget struct {
	cfg *Config
	Client *http.Client
}

type CurrencyCurrentValue map[string]string

func main() {
	var configFile = flag.String("config", "config.json", "path to config file.")
	flag.Parse()
	cfg, err := getConf(*configFile)
	if err != nil {
		log.Fatalf("%s", err)
	}

	client := &http.Client{}
	w := Widget{
		cfg: cfg,
		Client: client,
	}

	coin, err := w.getCoin()
	if err != nil {
		log.Fatalf("error retrieving coin. %s", err)
	}

	output := w.prepareOutput(coin)
	j, err := json.Marshal(output)
	if err != nil {
		log.Fatalf("could not marshal response struct to JSON. %s", err)
	}

	_, err = os.Stdout.WriteString(string(j)+"\n")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to print to stderr")
	}
}

// prepareOutput will use the retrieved data and build the object for output.
func (w Widget) prepareOutput(c *types.CoinsID) *CoinResponse {
	currencyValues := make(CurrencyCurrentValue)

	for _, cur := range w.cfg.AdditionalCurrencies {
		val := c.MarketData.CurrentPrice[cur] * w.cfg.Holdings
		strValue := strings.ConvertToCurrency(val)
		currencyValues[cur] = strValue
	}

	currentUSDPrice := strings.ConvertToCurrency(c.MarketData.CurrentPrice["usd"])
	output := CoinResponse{
		CurrentPrice:  currentUSDPrice,
		Rank:          c.MarketCapRank,
		CurrentValues: currencyValues,
	}

	return &output
}

// getCoin will retrieve current data for requested coin.
func (w Widget) getCoin() (*types.CoinsID, error) {
	req, err := http.NewRequest("GET", w.cfg.URL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	resp, err := w.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var coin types.CoinsID
	err = json.NewDecoder(resp.Body).Decode(&coin)
	if err != nil {
		return &coin, fmt.Errorf("unable to parse response received. %v", err)
	}

	return &coin, nil
}

// getConf will retrieve the configuration.
func getConf(configFile string) (*Config, error) {
	var cfg Config
	cf, err := os.Open(configFile)
	if err != nil {
		return &cfg, fmt.Errorf("unable to open configuration file: %w", err)
	}
	defer cf.Close()
	jsonParser := json.NewDecoder(cf)
	err = jsonParser.Decode(&cfg)
	if err != nil {
		return &cfg, fmt.Errorf("error retrieving config. unable to unmarshal config file: %w", err)
	}

	return &cfg, nil
}
