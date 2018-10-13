package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type AtCoder struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	StartTimeSeconds int64  `json:"startTimeSeconds"`
	DurationSeconds  int64  `json:"durationSeconds"`
	RatedRange       string `json:"ratedRange"`
}

func GetAtCoder() ([]RawContest, bool) {
	res, err := http.Get("https://atcoder-api.appspot.com/contests")
	if err != nil {
		log.Printf("Faild to GET reqest(AtCoder): %v", err)
		return nil, false
	}
	defer res.Body.Close()

	byteArray, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Faild to read data(AtCoder): %v", err)
		return nil, false
	}

	var atCoderContests []AtCoder
	json.Unmarshal(byteArray, &atCoderContests)

	var contests []RawContest
	for _, contestData := range atCoderContests {
		if time.Now().Unix() >= contestData.StartTimeSeconds {
			continue
		}

		contest := RawContest{
			Name:        contestData.Title,
			StartTime:   contestData.StartTimeSeconds,
			URL:         "https://beta.atcoder.jp/contests/" + contestData.ID,
			Duration:    contestData.DurationSeconds,
			WebSiteName: "AtCoder",
			WebSiteURL:  "https://beta.atcoder.jp/",
		}
		contests = append(contests, contest)
	}

	return contests, true
}
