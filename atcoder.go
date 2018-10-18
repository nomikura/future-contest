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

type AtCoder struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	StartTime int64  `json:"startTimeSeconds"`
	Duration  int64  `json:"durationSeconds"`
	Rated     string `json:"ratedRange"`
}

func (module *ContestModule) SetAtCoder() bool {
	var request *http.Request = module.Context.Request
	context := appengine.NewContext(request)

	client := urlfetch.Client(context)

	// GETリクエスト
	response, err := client.Get("https://atcoder-api.appspot.com/contests")
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
	var atcoderContests []AtCoder
	json.Unmarshal(byteArray, &atcoderContests)

	var contests []Contest
	for _, contest := range atcoderContests {
		if time.Now().Unix() >= contest.StartTime {
			continue
		}

		contests = append(contests, Contest{
			ID:          contest.ID,
			Title:       contest.Title,
			StartTime:   contest.StartTime,
			Duration:    contest.Duration,
			WebSiteName: "atcoder",
		})
	}

	// モジュールに入れる
	module.AtCoder = append(module.AtCoder, contests...)

	return true
}
