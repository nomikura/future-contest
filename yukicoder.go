package future

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

type Yukicoder struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	Date    time.Time `json:"date"`
	EndDate time.Time `json:"endDate"`
}

func (module *ContestModule) SetYukicoder() bool {
	var request *http.Request = module.Context.Request
	context := appengine.NewContext(request)

	client := urlfetch.Client(context)

	// GETリクエスト
	response, err := client.Get("https://yukicoder.me/api/v1/contest/future")
	if err != nil {
		log.Infof(context, "%v", err)
		return false
	}

	// データをバイナリとして取得
	byteArray, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Infof(context, "%v", err)
		return false
	}

	// バイナリをAtCoderのデータにパース
	var yukicoderContests []Yukicoder
	json.Unmarshal(byteArray, &yukicoderContests)

	var contests []Contest
	for _, contest := range yukicoderContests {
		contests = append(contests, Contest{
			ID:          strconv.Itoa(contest.ID),
			Title:       contest.Name,
			StartTime:   contest.Date.Unix(),
			Duration:    contest.EndDate.Unix() - contest.Date.Unix(),
			WebSiteName: "yukicoder",
		})
	}

	// モジュールに入れる
	module.Yukicoder = append(module.Yukicoder, contests...)

	return true
}
