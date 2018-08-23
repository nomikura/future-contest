package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// Cedeforcesコンテスト情報
type CodeForces struct {
	Status string `json:"status"`
	Result []struct {
		ID                  int    `json:"id"`
		Name                string `json:"name"`
		Type                string `json:"type"`
		Phase               string `json:"phase"`
		Frozen              bool   `json:"frozen"`
		DurationSeconds     int    `json:"durationSeconds"`
		StartTimeSeconds    int    `json:"startTimeSeconds"`
		RelativeTimeSeconds int    `json:"relativeTimeSeconds"`
	} `json:"result"`
}

// Codeforcesのコンテスト情報を返す
func GetCodeForces() ([]RawContest, bool) {
	// 外部からデータを取得
	res, err := http.Get("https://codeforces.com/api/contest.list")
	if err != nil {
		log.Printf("Faild to GET reqest(Codeforces): %v", err)
		return nil, false
	}
	defer res.Body.Close()

	// 取得したデータを読み込み[]byteを返してもらう
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Faild to read data(Codeforces): %v", err)
		return nil, false
	}

	// jsonを解析してCodeForcesの構造体にデータを入れる
	var codeForces CodeForces
	json.Unmarshal(body, &codeForces)

	// JSONから必要な情報をとりだして構造体に入れる
	var contests []RawContest
	for _, result := range codeForces.Result {
		if result.Phase != "BEFORE" {
			continue
		}

		contest := RawContest{
			Name:        result.Name,
			StartTime:   int64(result.StartTimeSeconds),
			URL:         "http://codeforces.com/contests/" + strconv.Itoa(result.ID),
			Duration:    int64(result.DurationSeconds),
			WebSiteName: "Codeforces",
			WebSiteURL:  "http://codeforces.com/",
		}
		contests = append(contests, contest)
	}

	return contests, true
}
