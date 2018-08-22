package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
)

type RawContestTable struct {
	Active []RawContest
	Future []RawContest
}

type RawContest struct {
	Name        string
	URL         string
	StartTime   int64 // Unix time
	Duration    int64
	WebSiteName string
	WebSiteURL  string
}

type ContestTable struct {
	Active []Contest
	Future []Contest
}

type Contest struct {
	Name        string
	URL         string
	StartTime   string // Unix time
	Duration    string
	WebSiteName string
	WebSiteURL  string
}

const location = "Asia/Tokyo"

func init() {
	ok := SetContestData()
	if !ok {
		log.Print("Faild first SetContestData() function")
	}
	loc, err := time.LoadLocation(location)
	if err != nil {
		log.Print("Fail about load location: %v", err)
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc
}

func (c *RawContest) ToStringStartTime() string {
	startTime := time.Unix(c.StartTime, 0)
	startTimeStr := startTime.Format("2006-01-02 15:04")
	wdays := [...]string{"日", "月", "火", "水", "木", "金", "土"}
	startTimeStr += " (" + wdays[startTime.Weekday()] + ")"
	return startTimeStr
}

func (c *RawContest) ToStringDurationTime() string {
	second := c.Duration
	day, hour, minute := second/86400, (second%86400)/3600, (second%86400%3600)/60
	var duration string
	if day > 0 {
		duration += strconv.Itoa(int(day)) + ":"
	}
	if hour >= 0 && hour <= 9 {
		duration += "0" + strconv.Itoa(int(hour)) + ":"
	} else if hour > 9 {
		duration += strconv.Itoa(int(hour)) + ":"
	}
	if minute >= 0 && minute <= 9 {
		duration += "0" + strconv.Itoa(int(minute))
	} else if minute > 9 {
		duration += strconv.Itoa(int(minute))
	}
	return duration
}

type ContestToDisplay struct {
	Name        string
	URL         string
	StartTime   string
	Duration    string
	WebSiteName string
}

var AllContest []RawContest
var ContestDisplay []ContestToDisplay

// コンテスト情報をJSONファイルに保存する
func SetContestData() bool {
	// コンテスト情報をcontestTableに入れる
	var contestTable RawContestTable
	atCoder, ok := GetAtCoder()
	if !ok {
		return false
	}
	codeForces, ok := GetCodeForces()
	if !ok {
		return false
	}
	csa, ok := GetCSA()
	if !ok {
		return false
	}
	contestTable.Future = append(contestTable.Future, atCoder...)
	contestTable.Future = append(contestTable.Future, codeForces...)
	contestTable.Future = append(contestTable.Future, csa...)
	sort.Slice(contestTable.Future, func(i, j int) bool { return contestTable.Future[i].StartTime < contestTable.Future[j].StartTime })

	// コンテスト情報を書き込むJSONファイルを生成
	jsonFile, err := os.Create("contest.json")
	if err != nil {
		log.Printf("Faild to create JSON file: %v", err)
		return false
	}

	// JSONファイルのエンコード
	encoder := json.NewEncoder(jsonFile)
	err = encoder.Encode(&contestTable)
	if err != nil {
		log.Printf("Faild to encode: %v", err)
		return false
	}

	log.Print("Reload JSON file")
	return true
}

// JSONファイルからコンテストデータを取得し、閲覧用のデータに加工して返す
func GetContestData() (ContestTable, bool) {
	// JSONファイルを読み込む
	jsonFile, err := os.Open("contest.json")
	if err != nil {
		log.Printf("Faild to read JSON file: %v", err)
		return ContestTable{}, false
	}
	defer jsonFile.Close()

	var contestTable ContestTable
	decoder := json.NewDecoder(jsonFile)
	// ファイルを最後まで読み込むループ
	for {
		var rawContestTable RawContestTable
		err := decoder.Decode(&rawContestTable)
		// ファイルを最後まで読み込んだら終了
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Faild to decoding JSON: %v", err)
			return ContestTable{}, false
		}

		// 生のコンテストデータを閲覧用のコンテストデータにする
		for _, rawContest := range rawContestTable.Future {
			contest := Contest{
				Name:        rawContest.Name,
				URL:         rawContest.URL,
				StartTime:   rawContest.ToStringStartTime(),
				Duration:    rawContest.ToStringDurationTime(),
				WebSiteName: rawContest.WebSiteName,
				WebSiteURL:  rawContest.WebSiteURL,
			}
			contestTable.Future = append(contestTable.Future, contest)
		}
	}

	return contestTable, true
}

func main() {
	// 5分に1回コンテストデータを取得する
	go func() {
		cycleTime := time.NewTicker(5 * time.Minute)
		for {
			select {
			case <-cycleTime.C:
				SetContestData()
			}
		}
	}()

	http.HandleFunc("/", handle)
	http.HandleFunc("/_ah/health", healthCheckHandler)
	log.Print("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// "/"にアクセスするとこの関数が呼ばれる
func handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	t, _ := template.ParseFiles("tmpl.html")
	contestTable, _ := GetContestData()
	t.Execute(w, contestTable.Future)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}
