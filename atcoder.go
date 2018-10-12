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

/*
// AtCoderのコンテスト情報を返す
func GetAtCoder() ([]RawContest, bool) {
	// var tmp []RawContest
	// return tmp, true

	client := &http.Client{}

	// GETリクエストを送ってHTMLを返して貰う
	req, err := http.NewRequest("GET", "https://beta.atcoder.jp/contests/", nil)
	if err != nil {
		log.Printf("Faild to create request(AtCoder): %v", err)
		return nil, false
	}

	// GETリクエストにつけるヘッダ。実験した結果、言語を日本語にすれば完全なHTMLが返ってくる
	req.Header.Add("cookie", "language=ja;")
	// リクエストを送信
	resp, _ := client.Do(req)
	if err != nil {
		log.Printf("Faild to GET reqest(AtCoder): %v", err)
		return nil, false
	}
	defer resp.Body.Close()

	// HTMLをstring型にしてcontentに代入
	byteArray, _ := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Faild to read data(AtCoder): %v", err)
		return nil, false
	}
	content := string(byteArray)

	// HTMLをgoqueryが使える形式に置き換える:
	reader := strings.NewReader(content)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Printf("Faild to create Document ojbect(AtCoder): %v", err)
		return nil, false
	}

	// goqueryでスクレイピング
	var contests []RawContest
	// 予定されたコンテストを見つけるためにh3を探す
	doc.Find("h3").Each(func(i int, s *goquery.Selection) {
		if h3 := s.Text(); h3 == "予定されたコンテスト" {
			// h3が予定されたコンテストの場合、次の要素はそのコンテストのテーブルとなる
			t := s.Next()
			t.Find("div > table > tbody > tr").Each(func(j int, u *goquery.Selection) {
				var before []string
				var href string
				u.Find("td").Each(func(k int, v *goquery.Selection) {
					before = append(before, v.Text())
					if k == 1 {
						href, _ = v.Find("a").Attr("href") // コンテストのURLを取得
					}
				})

				// Durationを求める
				str := strings.Replace(before[2], ":", "h", 1) + "m"
				tim, _ := time.ParseDuration(str)
				duration := int64(tim.Seconds())

				// StartTimeを求める
				start := before[0]
				atoi := func(str string) int {
					ret, _ := strconv.Atoi(str)
					return ret
				}
				// [2018-09-22 21:00:00+0900]の形式で抜き出した時間を無理矢理Timeオブジェクトにする
				year, month, day, hour, minute := atoi(start[:4]), atoi(start[5:7]), atoi(start[8:10]), atoi(start[11:13]), atoi(start[14:])
				// 取得する時間はJSTなので、日本時間をTimeオブジェクトにするように処理する
				jst, _ := time.LoadLocation("Asia/Tokyo")
				startTime := time.Date(year, time.Month(month), day, hour, minute, 0, 0, jst)
				unix := startTime.Unix()

				contest := RawContest{
					Name:        before[1],
					StartTime:   unix,
					Duration:    duration,
					URL:         "https://beta.atcoder.jp" + href,
					WebSiteName: "AtCoder",
					WebSiteURL:  "https://beta.atcoder.jp/",
				}
				contests = append(contests, contest)
			})
		}
	})

	return contests, true

}

*/
