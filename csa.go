package future

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
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

func (module ContestModule) SetCSA() bool {
	var r *http.Request = module.Context.Request
	context := appengine.NewContext(r)

	client := urlfetch.Client(context)

	// リクエストを生成
	request, err := http.NewRequest("GET", "https://csacademy.com/contests/", nil)
	if err != nil {
		log.Print(err)
		return false
	}

	// リクエストヘッダに情報を追加
	request.Header.Add("x-requested-with", "XMLHttpRequest")

	// GETリクエスト
	response, err := client.Do(request)
	if err != nil {
		log.Print(err)
		return false
	}

	// データをバイナリとして取得
	byteArray, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
		return false
	}

	// バイナリをAtCoderのデータにパース
	var csaContest CSA
	json.Unmarshal(byteArray, &csaContest)

	var contests []Contest
	for _, contest := range csaContest.State.Contest {
		// データが存在しない or 正式なコンテストではない 場合に回避する
		if contest.StartTime == nil || contest.EndTime == nil || contest.SystemGenerated {
			continue
		}

		startTime := int64(contest.StartTime.(float64))
		// 過去のコンテストを回避する
		if currentTime := time.Now().Unix(); currentTime >= startTime {
			continue
		}

		endTime := int64(contest.EndTime.(float64))
		contests = append(contests, Contest{
			ID:          contest.Name,
			Title:       contest.LongName,
			StartTime:   startTime,
			Duration:    endTime - startTime,
			WebSiteName: "csa",
		})
	}

	module.Codeforces = append(module.Codeforces, contests...)
	return true
}
