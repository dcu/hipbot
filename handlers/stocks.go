package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/dcu/hipbot/xmpp"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type YahooResponse struct {
	Query struct {
		Results struct {
			Quote Quote `json:"quote"`
		} `json:"results"`
	} `json:"query"`
}

type Quote struct {
	AverageDailyVolume   string `json:"AverageDailyVolume"`
	Change               string `json:"Change"`
	DaysLow              string `json:"DaysLow"`
	DaysHigh             string `json:"DaysHigh"`
	YearLow              string `json:"YearLow"`
	YearHigh             string `json:"YearHigh"`
	MarketCapitalization string `json:"MarketCapitalization"`
	LastTradePriceOnly   string `json:"LastTradePriceOnly"`
	DaysRange            string `json:"DaysRange"`
	Name                 string `json:"Name"`
	Symbol               string `json:"Symbol"`
	Volume               string `json:"Volume"`
	StockExchange        string `json:"StockExchange"`
}

type StocksHandler struct {
}

func (stocks *StocksHandler) Matches(message *xmpp.Chat) bool {
	return strings.HasPrefix(message.Text, "stock:")
}

func (stocks *StocksHandler) Process(client *xmpp.Client, roomId string, message *xmpp.Chat) {
	symbol := strings.Replace(message.Text, "stock:", "", 1)

	query := url.QueryEscape(`select * from yahoo.finance.quote where symbol = "` + symbol + `"`)
	url := `https://query.yahooapis.com/v1/public/yql?q=` + query + `&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys&callback=`

	resp, err := http.Get(url)
	defer resp.Body.Close()

	if err != nil {
		return
	}

	results := &YahooResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, results)

	quote := results.Query.Results.Quote

	formattedResult := quote.Symbol + ` price: ` + quote.LastTradePriceOnly + ` change: ` + quote.Change + ` market cap: ` + quote.MarketCapitalization + ` volume: ` + quote.Volume + `/` + quote.AverageDailyVolume

	fmt.Printf("%#v\n", results)

	client.Send(xmpp.Chat{
		Remote: roomId,
		Text:   formattedResult,
		Type:   "groupchat",
	})
}
