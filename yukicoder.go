package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Yukicoder struct {
	Kind             string        `json:"kind"`
	Etag             string        `json:"etag"`
	Summary          string        `json:"summary"`
	Description      string        `json:"description"`
	Updated          time.Time     `json:"updated"`
	TimeZone         string        `json:"timeZone"`
	AccessRole       string        `json:"accessRole"`
	DefaultReminders []interface{} `json:"defaultReminders"`
	NextSyncToken    string        `json:"nextSyncToken"`
	Items            []struct {
		Kind        string    `json:"kind"`
		Etag        string    `json:"etag"`
		ID          string    `json:"id"`
		Status      string    `json:"status"`
		HTMLLink    string    `json:"htmlLink"`
		Created     time.Time `json:"created"`
		Updated     time.Time `json:"updated"`
		Summary     string    `json:"summary"`
		Description string    `json:"description"`
		Location    string    `json:"location"`
		Creator     struct {
			Email       string `json:"email"`
			DisplayName string `json:"displayName"`
			Self        bool   `json:"self"`
		} `json:"creator"`
		Organizer struct {
			Email       string `json:"email"`
			DisplayName string `json:"displayName"`
			Self        bool   `json:"self"`
		} `json:"organizer"`
		Start struct {
			DateTime time.Time `json:"dateTime"`
		} `json:"start"`
		End struct {
			DateTime time.Time `json:"dateTime"`
		} `json:"end"`
		ICalUID  string `json:"iCalUID"`
		Sequence int    `json:"sequence"`
	} `json:"items"`
}

// yukicoderのJSONファイルのURLを返す
func GetYukicoderURL() string {
	now := time.Now()
	itoa := func(t int) string {
		return strconv.Itoa(t)
	}
	year, month, day := now.Year(), itoa(int(now.Month())), itoa(now.Day())
	// 1文字なら戦闘に0を付け加える
	if len(month) == 1 {
		month = "0" + month
	}
	if len(day) == 1 {
		day = "0" + day
	}

	// テキストに書いてるURL情報を読み込む
	urlByteArray, err := ioutil.ReadFile("yukicoder_url.txt")
	if err != nil {
		log.Printf("Faild to open yukicoder_url.txt: %v", err)
	}
	urlSlice := strings.Split(string(urlByteArray), ",")

	// 今日からプラスマイナス1年分のデータを抜き出すURLを生成
	url := urlSlice[0]
	url += itoa(year-1) + "-" + month + "-" + day
	url += urlSlice[1]
	url += itoa(year+1) + "-" + month + "-" + day
	url += urlSlice[2]

	return url
}

func GetYukicoder() ([]RawContest, bool) {
	// 外部からデータを取得
	res, err := http.Get(GetYukicoderURL())
	if err != nil {
		log.Printf("Faild to GET reqest(yukicoder): %v", err)
		return nil, false
	}
	defer res.Body.Close()

	// 取得したデータを読み込み[]byteを返してもらう
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Faild to read data(yukicoder): %v", err)
		return nil, false
	}

	// jsonを解析してCodeForcesの構造体にデータを入れる
	var yukicoder Yukicoder
	json.Unmarshal(body, &yukicoder)

	var contests []RawContest
	for _, result := range yukicoder.Items {
		name := result.Summary
		url := "https://yukicoder.me/contests/" + result.ICalUID
		start := result.Start.DateTime.Unix()
		if now := time.Now().Unix(); now >= start {
			continue
		}
		end := result.End.DateTime.Unix()
		contest := RawContest{
			Name:        name,
			URL:         url,
			StartTime:   start,
			Duration:    end - start,
			WebSiteName: "yukicoder",
			WebSiteURL:  "https://yukicoder.me/",
		}
		contests = append(contests, contest)
	}

	return contests, true
}
