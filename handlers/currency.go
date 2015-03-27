package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/dcu/hipbot/xmpp"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

var (
	pattern = regexp.MustCompile(`(\d+(\.\d+)?) (\w+) to (\w+)`)
)

type YahooExchangeResponse struct {
	Query struct {
		Results struct {
			Exchanges []Exchange `json:"rate"`
		} `json:"results"`
	} `json:"query"`
}

type Exchange struct {
	Name string `json:"Name"`
	Rate string `json:"Rate"`
	Date string `json:"Date"`
	Time string `json:"Time"`
	Ask  string `json:"Ask"`
	Bid  string `json:"Bid"`
}

type CurrencyHandler struct {
}

func (stocks *CurrencyHandler) Matches(message *xmpp.Chat) bool {
	return pattern.MatchString(message.Text)
}

func (stocks *CurrencyHandler) Process(client *xmpp.Client, roomId string, message *xmpp.Chat) {
	matches := pattern.FindStringSubmatch(message.Text)
	if len(matches) < 5 {
		return
	}

	from := matches[3]
	to := matches[4]
	amount, _ := strconv.ParseFloat(matches[1], 64)

	query := `select * from yahoo.finance.xchange where pair in ("` + to + `", "` + from + `")`
	url := `https://query.yahooapis.com/v1/public/yql?q=` + url.QueryEscape(query) + `&format=json&env=store%3A%2F%2Fdatatables.org%2Falltableswithkeys&callback=`

	resp, err := http.Get(url)
	defer resp.Body.Close()

	if err != nil {
		return
	}

	results := &YahooExchangeResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, results)

	exchanges := results.Query.Results.Exchanges
	if len(exchanges) == 0 {
		return
	}

	fmt.Printf("%#v\n", results)
	result := exchanges[0]

	rate, _ := strconv.ParseFloat(result.Rate, 64)
	total := rate * amount
	totalStr := strconv.FormatFloat(total, 'f', 2, 64)
	formattedResult := matches[1] + ` ` + from + ` = ` + totalStr + ` ` + to
	formattedResult += ` [Ask: ` + result.Ask + ` Bid: ` + result.Bid + `]`

	client.Send(xmpp.Chat{
		Remote: roomId,
		Text:   formattedResult,
		Type:   "groupchat",
	})
}
