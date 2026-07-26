package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	"nugu.dev/basement/pkg/models"
	"nugu.dev/basement/views/activity_views"
)

func (app *application) Activities(w http.ResponseWriter, r *http.Request) {

	// Create New Activity
	if r.Method == "POST" {
		app.startActivity(w, r)
		return
	}

	// Finish Open Activity
	if r.Method == "PATCH" {
		app.finishActivity(w, r)
		return
	}

	last, queryErr := app.activitiesRepository.GetLastActivity()

	if queryErr != nil && queryErr != models.ErrNotFound {
		http.Error(w, queryErr.Error(), http.StatusInternalServerError)
		return
	}

	listDone, queryErr := app.activitiesRepository.GetDailyLog(time.Now())

	if queryErr != nil && queryErr != models.ErrNotFound {
		http.Error(w, queryErr.Error(), http.StatusInternalServerError)
		return
	}

	var component templ.Component

	// if last activity doesnt exists or already ended, then delete cookie
	if queryErr == models.ErrNotFound || !last.EndTime.IsZero() {
		component = activity_views.ActivityIndex(false, listDone) // start button
	} else {
		component = activity_views.ActivityIndex(true, listDone) // stop button
	}

	setCookie(w, last.StartTime, last.EndTime)
	component.Render(r.Context(), w)
}

func (app *application) startActivity(w http.ResponseWriter, r *http.Request) {

	_, err := app.activitiesRepository.StartActivity()

	if err != nil {
		switch err {
		case models.ErrNotFinished:
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	last, queryErr := app.activitiesRepository.GetLastActivity()

	if queryErr != nil && queryErr != models.ErrNotFound {
		http.Error(w, queryErr.Error(), http.StatusInternalServerError)
		return
	}

	setCookie(w, last.StartTime, last.EndTime)
	component := activity_views.StopButton()
	component.Render(r.Context(), w)
}

func (app *application) finishActivity(w http.ResponseWriter, r *http.Request) {

	err := app.activitiesRepository.EndActivity()
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	last, queryErr := app.activitiesRepository.GetLastActivity()

	if queryErr != nil && queryErr != models.ErrNotFound {
		http.Error(w, queryErr.Error(), http.StatusInternalServerError)
		return
	}

	setCookie(w, last.StartTime, last.EndTime)

	component := activity_views.StartButton()
	component.Render(r.Context(), w)
}

/* ========================================================================== */

func (app *application) GetDailyLog(w http.ResponseWriter, r *http.Request) {

	loggedUser := true
	authToken, err := r.Cookie("t")

	if err != nil || authToken.Value != app.AuthToken {
		loggedUser = false
	}

	partialLog := app.ShowDailyLog(r.PathValue("date"), loggedUser)
	reqType := r.URL.Query().Get("partial")
	if reqType == "true" {
		partialLog.Render(r.Context(), w)
		return
	}

	tabSelected := r.URL.Query().Get("tab")

	var calendarLog []models.ActivityDayOverview

	calendarStats := models.StatsOverview{
		DailyAverage:  0,
		LongestStreak: 0,
		CurrentStreak: 0,
	}

	from, to := r.URL.Query().Get("from"), r.URL.Query().Get("to")

	if from == "" || to == "" {
		today := time.Now()
		startDate := time.Date(today.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(1, 0, 0).AddDate(0, 0, -1)
		calendarHolder, err := app.activitiesRepository.GetIntervalLog(startDate, endDate, tabSelected)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

	Outer:
		for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {

			for _, x := range calendarHolder {
				if x.Date.Equal(d) {
					calendarLog = append(calendarLog, x)

					calendarStats.DailyAverage += x.TotalSec

					if x.TotalSec >= 3600 {
						calendarStats.CurrentStreak++
					} else {
						if calendarStats.LongestStreak < calendarStats.CurrentStreak {
							calendarStats.LongestStreak = calendarStats.CurrentStreak
						}

						calendarStats.CurrentStreak = 0
					}

					continue Outer
				}
			}

			// Didnt found the date record on database
			if calendarStats.LongestStreak < calendarStats.CurrentStreak {
				calendarStats.LongestStreak = calendarStats.CurrentStreak
			}

			x := models.ActivityDayOverview{Date: d, TotalSec: 0}

			calendarLog = append(calendarLog, x)
		}

		if len(calendarHolder) > 0 {
			calendarStats.DailyAverage = (calendarStats.DailyAverage / len(calendarHolder)) / 3600
		}
	}

	page := activity_views.Log(calendarLog, partialLog, calendarStats, loggedUser)
	page.Render(r.Context(), w)
}

/* ========================================================================== */

// GET Show single daily log partial
func (app *application) ShowDailyLog(selected string, loggedUser bool) templ.Component {

	if selected == "" {
		return activity_views.NoLogSelected()
	}

	dateReq, err := time.Parse(time.DateOnly, selected)

	if err != nil {
		return activity_views.NoLogSelected()
	}

	dailyLog, err := app.activitiesRepository.GetDailyLog(dateReq)

	if err != nil {
		return activity_views.NoLogSelected()
	}

	tags, err := app.tagsRepository.GetActivityTags()

	if err != nil {
		return activity_views.NoLogSelected()
	}

	return activity_views.DetailedLog(dateReq, dailyLog, tags, loggedUser)
}

/* ========================================================================== */

// POST form to create new log
func (app *application) CreateDailyLog(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	var a models.Activity
	var err error

	a.Title = r.FormValue("title")
	a.Description = r.FormValue("description")
	a.StartTime, err = From24ClockToTime(r.FormValue("date"), r.FormValue("start"))

	for k, vs := range r.Form {
		for _, v := range vs {
			if strings.Contains(k, "check-") {
				var t models.Tag
				t.ID, _ = strconv.Atoi(strings.Split(k, "-")[1])
				t.Name = v
				a.Tags = append(a.Tags, t)
			}
		}
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	a.EndTime, err = From24ClockToTime(r.FormValue("date"), r.FormValue("end"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if a.StartTime.After(a.EndTime) {
		http.Error(w, fmt.Errorf("start > finish").Error(), http.StatusBadRequest)
		return
	}

	_, err = app.activitiesRepository.NewCompleteActivity(a)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/log/%s", r.FormValue("date")), http.StatusSeeOther)
}

/* ========================================================================== */

func (app *application) EditDailyLog(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		log, err := app.activitiesRepository.GetSingleDetailedLogById(r.PathValue("id"))

		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		tags, err := app.tagsRepository.GetActivityTags()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		component := activity_views.EditDailyLogModalInternal(log, tags)
		component.Render(r.Context(), w)
		return
	}
}

/* ========================================================================== */

func setCookie(w http.ResponseWriter, startTime time.Time, endTime time.Time) {

	startTimerCookie := http.Cookie{Name: "start"}

	// endTime != null, last activity already finished. Kill the cookie
	if !endTime.IsZero() {
		startTimerCookie.Value = ""
		startTimerCookie.Expires = time.Now()
	} else {
		startTimerCookie.Value = strconv.FormatInt(startTime.Unix(), 10)
		startTimerCookie.Expires = time.Now().Add(time.Hour * 24)
	}

	http.SetCookie(w, &startTimerCookie)
}
