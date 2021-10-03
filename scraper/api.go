package lolapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"golang.org/x/time/rate"
)

const (
	queueType = "RANKED_SOLO_5x5"
)

type Rank int
type Division int

const (
	RankChall Rank = iota
	RankGM
	RankM
	RankD
	RankP
	RankG
	RankS
	RankB
	RankI

	Div1 Division = iota
	Div2
	Div3
	Div4
)

var Ranks = map[Rank]string{
	RankChall: "CHALLENGER",
	RankGM:    "GRANDMASTER",
	RankM:     "MASTER",
	RankD:     "DIAMOND",
	RankP:     "PLATINUM",
	RankG:     "GOLD",
	RankS:     "SILVER",
	RankB:     "BRONZE",
	RankI:     "IRON",
}

var Divisions = map[Division]string{
	Div1: "I",
	Div2: "II",
	Div3: "III",
	Div4: "IV",
}

type PlayerId struct {
	SummonerId   string `json:"summonerId"`
	SummonerName string `json:"summonerName"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
	GamesPlayed  int    `json:"gamesPlayed"`
}

type ApiClient struct {
	apiKey  string
	limiter *rate.Limiter
	region  string
}

func NewApiClient(apiKey string) ApiClient {
	// We only want one api call every 1.25 seconds.
	limiter := rate.NewLimiter(0.75, 1)
	return ApiClient{
		apiKey:  apiKey,
		limiter: limiter,
		region:  "https://na1.api.riotgames.com",
	}
}

func (c *ApiClient) makeRequest(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.region+endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Riot-Token", c.apiKey)
	// Now is when we want to see if we should block the api request and wait
	if err = c.limiter.Wait(context.Background()); err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	log.Printf("Limits: %v %v\n", res.Header.Get("X-App-Rate-Limit-Count"),
		res.Header.Get("X-Method-Rate-Limit-Count"))
	if res.StatusCode != 200 {
		return nil, errors.New(res.Status)
	}
	return res, err
}

func (c *ApiClient) GetLeagueEntries(tier Rank, div Division, page int) ([]PlayerId, error) {
	// Make the endpoint string
	if page < 1 {
		return nil, errors.New("page number cannot be than 1")
	}

	if (tier == RankChall || tier == RankGM || tier == RankM) && div != Div1 {
		return nil, errors.New("apex tiers only use division 1")
	}
	endpoint := fmt.Sprintf("/lol/league-exp/v4/entries/%v/%v/%v?page=%v", queueType, Ranks[tier],
		Divisions[div], page)

	res, err := c.makeRequest(endpoint)
	if err != nil {
		return nil, err
	}

	buf, _ := io.ReadAll(res.Body)

	var players []PlayerId

	if err := json.Unmarshal(buf, &players); err != nil {
		return nil, err
	}

	for _, p := range players {
		p.GamesPlayed = p.Wins + p.Losses
	}

	return players, nil
}
