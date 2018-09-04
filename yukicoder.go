package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Yukicoder struct {
	Id      int
	Name    string
	Date    time.Time
	EndDate time.Time
}

func GetYukicoder() ([]RawContest, bool) {
	// 外部からデータを取得
	res, err := http.Get("https://yukicoder.me/api/v1/contest/future")
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
	var yukicoder []Yukicoder
	json.Unmarshal(body, &yukicoder)

	var contests []RawContest
	for _, result := range yukicoder {
		name := result.Name
		url := "https://yukicoder.me/contests/" + strconv.Itoa(result.Id)
		start := result.Date.Unix()
		if now := time.Now().Unix(); now >= start {
			continue
		}
		end := result.EndDate.Unix()
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
