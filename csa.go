package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// CSAコンテスト情報
type CSA struct {
	State struct {
		Contest []struct {
			StartTime            interface{} `json:"startTime"`
			LongName             string      `json:"longName"`
			ScoreType            int         `json:"scoreType"`
			ID                   int         `json:"id"`
			PublicSources        bool        `json:"publicSources"`
			IsVisible            bool        `json:"isVisible"`
			NumRegistered        int         `json:"numRegistered"`
			Name                 string      `json:"name"`
			ScoringID            int         `json:"scoringId"`
			Rated                bool        `json:"rated"`
			ChatID               int         `json:"chatId,omitempty"`
			Description          string      `json:"description"`
			LiveResults          bool        `json:"liveResults"`
			IsAnalysisPublic     bool        `json:"isAnalysisPublic"`
			EndTime              interface{} `json:"endTime"`
			ScoreboardType       int         `json:"scoreboardType"`
			VirtualContestID     int         `json:"virtualContestId,omitempty"`
			AnalysisDiscussionID int         `json:"analysisDiscussionId,omitempty"`
			OriginArchiveID      int         `json:"originArchiveId,omitempty"`
			DescriptionArticleID int         `json:"descriptionArticleId,omitempty"`
			AnalysisArticleID    int         `json:"analysisArticleId,omitempty"`
			BaseContestID        int         `json:"baseContestId,omitempty"`
			MaxRating            int         `json:"maxRating,omitempty"`
			LiveStats            bool        `json:"liveStats,omitempty"`
			NumSubmissions       int         `json:"numSubmissions,omitempty"`
			NumCustomRuns        int         `json:"numCustomRuns,omitempty"`
			NumCompiles          int         `json:"numCompiles,omitempty"`
			NumExampleRuns       int         `json:"numExampleRuns,omitempty"`
			OwnerID              int         `json:"ownerId,omitempty"`
			SystemGenerated      bool        `json:"systemGenerated,omitempty"`
		} `json:"Contest"`
	} `json:"state"`
}

// CSAのコンテスト情報を返す
func GetCSA() ([]RawContest, bool) {
	client := &http.Client{}

	// リクエストを生成
	req, err := http.NewRequest("GET", "https://csacademy.com/contests/", nil)
	if err != nil {
		log.Printf("Faild to create request(CSA): %v", err)
		return nil, false
	}

	// リクエストヘッダにこれをつけるとコンテスト情報をJSONで返してくれる
	req.Header.Add("x-requested-with", "XMLHttpRequest")
	// リクエストを送信
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Faild to GET reqest(CSA): %v", err)
		return nil, false
	}
	defer resp.Body.Close()

	// 取得したデータを[]byteにして返して貰う
	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Faild to read data(CSA): %v", err)
		return nil, false
	}

	// JSONを解析してCSAの構造体にデータを入れる
	var csa CSA
	json.Unmarshal(byteArray, &csa)

	// JSONから必要な情報を取り出して構造体に入れる
	var contests []RawContest
	for _, result := range csa.State.Contest {
		// データが存在しない or 正式なコンテストではない 場合に回避する
		if result.StartTime == nil || result.EndTime == nil || result.SystemGenerated {
			continue
		}

		startTime := int64(result.StartTime.(float64))
		// 過去のコンテストを回避する
		if currentTime := time.Now().Unix(); currentTime >= startTime {
			continue
		}

		endTime := int64(result.EndTime.(float64))
		contest := RawContest{
			Name:        result.LongName,
			StartTime:   startTime,
			URL:         "https://csacademy.com/contest/" + result.Name + "/",
			Duration:    endTime - startTime,
			WebSiteName: "CS Academy",
			WebSiteURL:  "https://csacademy.com/",
		}
		contests = append(contests, contest)
	}

	return contests, true
}
