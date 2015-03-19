package handlers

import (
	"encoding/json"
	"github.com/daneharrigan/hipchat"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

type GiphyImageData struct {
	URL    string
	Width  string
	Height string
	Size   string
	Frames string
}

type GiphyGif struct {
	Type               string
	Id                 string
	URL                string
	Tags               string
	BitlyGifURL        string `json:"bitly_gif_url"`
	BitlyFullscreenURL string `json:"bitly_fullscreen_url"`
	BitlyTiledURL      string `json:"bitly_tiled_url"`
	Images             struct {
		Original               GiphyImageData
		FixedHeight            GiphyImageData `json:"fixed_height"`
		FixedHeightStill       GiphyImageData `json:"fixed_height_still"`
		FixedHeightDownsampled GiphyImageData `json:"fixed_height_downsampled"`
		FixedWidth             GiphyImageData `json:"fixed_width"`
		FixedwidthStill        GiphyImageData `json:"fixed_width_still"`
		FixedwidthDownsampled  GiphyImageData `json:"fixed_width_downsampled"`
	}
}

type GiphyResults struct {
	Data []GiphyGif
}

type GiphyHandler struct {
}

func (giphy *GiphyHandler) Matches(message *hipchat.Message) bool {
	return strings.HasPrefix(message.Body, "giphy:")
}

func (giphy *GiphyHandler) Process(client *hipchat.Client, roomId string, message *hipchat.Message) {
	search := strings.Replace(message.Body, "giphy:", "", 1)

	url := `http://api.giphy.com/v1/gifs/search?q=` + url.QueryEscape(search) + `&limit=10&api_key=dc6zaTOxFJmzC`
	resp, err := http.Get(url)
	defer resp.Body.Close()

	if err != nil {
		return
	}

	results := &GiphyResults{}
	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, results)
	println(string(body))

	resultsCount := len(results.Data)

	if resultsCount > 0 {
		image := results.Data[rand.Intn(resultsCount)]
		client.Say(roomId, "Giphy", image.Images.Original.URL)
	} else {
		client.Say(roomId, "Giphy", "Couldn't find any result.")
	}
}
