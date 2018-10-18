package future

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	r := gin.New()
	r.GET("/update", UpdateHandler)
	r.GET("/contests", ContestHandler)
	r.GET("/", HomeHandler)
	http.Handle("/", r)
}

func HomeHandler(context *gin.Context) {
	contestModule := &ContestModule{}
	contestModule.Context = context
	contestModule.FileIO("read")

	for _, contest := range contestModule.AllContest {
		contestURL := ""
		webSiteURL := ""
		webSiteName := ""
		switch contest.WebSiteName {
		case "atcoder":
			contestURL = "https://beta.atcoder.jp/contests/"
			webSiteURL = "https://beta.atcoder.jp/"
			webSiteName = "AtCoder"
		case "csa":
			contestURL = "https://csacademy.com/contest/"
			webSiteURL = "https://csacademy.com/"
			webSiteName = "CS Academy"
		case "codeforces":
			contestURL = "http://codeforces.com/contests/"
			webSiteURL = "http://codeforces.com/"
			webSiteName = "Codeforces"
		}
		contestModule.ContestViews = append(contestModule.ContestViews, ContestView{
			Name:        contest.Title,
			URL:         contestURL + contest.ID,
			StartTime:   ToStringStartTime(contest.StartTime),
			Duration:    ToStringDurationTime(contest.Duration),
			WebSiteName: webSiteName,
			WebSiteURL:  webSiteURL,
		})
	}

	// context.JSON(http.StatusOK, contestModule.ContestViews)
	var w http.ResponseWriter = context.Writer
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// w.Header().Set("cache-control", "public, max-age=3600")
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, contestModule.ContestViews)

	// t := template.Must(template.ParseFiles("templates/index.html"))
	// if err := t.ExecuteTemplate(w, "index.html", contestModule.ContestViews); err != nil {
	// 	log.Infof(context, "template error: %v", err)
	// }

	// context.HTML(http.StatusOK, "index.html", contestModule.ContestViews)
}

func UpdateHandler(context *gin.Context) {
	SetContestData(context)
}

func ContestHandler(context *gin.Context) {
	GetContestData(context)
}
