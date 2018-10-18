package future

import (
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ContestModule struct {
	AtCoder      []Contest
	Codeforces   []Contest
	CSA          []Contest
	AllContest   []Contest
	ContestViews []ContestView
	Context      *gin.Context
}

type Contest struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	StartTime   int64  `json:"startTimeSeconds"`
	Duration    int64  `json:"durationSeconds"`
	WebSiteName string `json:"webSiteName"`
}

type ContestView struct {
	Name        string
	URL         string
	StartTime   string
	Duration    string
	WebSiteName string
	WebSiteURL  string
}

func SetContestData(context *gin.Context) {
	contestModule := &ContestModule{}
	contestModule.Context = context

	// 各種コンテストデータを更新
	ok := contestModule.SetAtCoder()
	if !ok {
		log.Print("Faild to Get AtCoder")
	}
	ok = contestModule.SetCodeforces()
	if !ok {
		log.Print("Faild to Get Codeforces")
	}
	ok = contestModule.SetCSA()
	if !ok {
		log.Print("Faild to Get CSA")
	}

	// 各コンテストの最新情報を全体データに統合
	var allContest []Contest
	allContest = append(allContest, contestModule.AtCoder...)
	allContest = append(allContest, contestModule.Codeforces...)
	allContest = append(allContest, contestModule.CSA...)
	contestModule.AllContest = allContest

	// 最後にソート
	sort.Slice(contestModule.AllContest, func(i, j int) bool {
		return contestModule.AllContest[i].StartTime < contestModule.AllContest[j].StartTime
	})

	// ファイルに書き込む
	contestModule.FileIO("write")

	// 出力
	context.JSON(http.StatusOK, "Finish update!!")
}

func GetContestData(context *gin.Context) {
	module := &ContestModule{}
	module.Context = context
	module.FileIO("read")
	context.JSON(http.StatusOK, module.AllContest)
}

func ToStringStartTime(startArg int64) string {
	const location = "Asia/Tokyo"
	loc, err := time.LoadLocation(location)
	if err != nil {
		// log.Print("Fail about load location: %v", err)
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc

	startTime := time.Unix(startArg, 0).Local()
	startTimeStr := startTime.Format("2006-01-02 15:04")
	wdays := [...]string{"日", "月", "火", "水", "木", "金", "土"}
	startTimeStr += " (" + wdays[startTime.Weekday()] + ")"
	return startTimeStr
}

func ToStringDurationTime(durationArg int64) string {
	second := durationArg
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
