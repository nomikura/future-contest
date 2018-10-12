package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type AtCoder struct {
	Title     string `json:"Title"`
	Path      string `json:"Path"`
	StartTime int64  `json:"StartTime"`
	Duration  int64  `json:"Duration"`
	Rated     string `json:"Rated"`
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
		if time.Now().Unix() >= contestData.StartTime {
			continue
		}

		contest := RawContest{
			Name:        contestData.Title,
			StartTime:   contestData.StartTime,
			URL:         "https://beta.atcoder.jp" + contestData.Path,
			Duration:    contestData.Duration,
			WebSiteName: "AtCoder",
			WebSiteURL:  "https://beta.atcoder.jp/",
		}
		contests = append(contests, contest)
	}

	return contests, true
}
