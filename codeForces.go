package future

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
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

func (module *ContestModule) SetCodeforces() bool {
	var request *http.Request = module.Context.Request
	context := appengine.NewContext(request)

	client := urlfetch.Client(context)

	// GETリクエスト
	response, err := client.Get("https://codeforces.com/api/contest.list")
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
	var codeforcesContest CodeForces
	json.Unmarshal(byteArray, &codeforcesContest)

	var contests []Contest
	for _, contest := range codeforcesContest.Result {
		if contest.Phase != "BEFORE" {
			continue
		}

		contests = append(contests, Contest{
			ID:          strconv.Itoa(contest.ID),
			Title:       contest.Name,
			StartTime:   int64(contest.StartTimeSeconds),
			Duration:    int64(contest.DurationSeconds),
			WebSiteName: "codeforces",
		})
	}

	module.Codeforces = append(module.Codeforces, contests...)
	return true
}
