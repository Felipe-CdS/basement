package main

import (
	"net/http"

	layouts_view "nugu.dev/basement/views/layouts"
	secrets_view "nugu.dev/basement/views/secrets"
	"nugu.dev/basement/views/static_views"
)

func (a *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/{$}", a.Landing)

	mux.HandleFunc("/activities", func(w http.ResponseWriter, r *http.Request) {

		if err := a.AuthMiddleware(w, r); err != nil {
			return
		}

		a.Activities(w, r)
	})

	mux.HandleFunc("/gallery", func(w http.ResponseWriter, r *http.Request) {

		if err := a.AuthMiddleware(w, r); err != nil {
			return
		}

		a.Gallery(w, r)
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

		// not logged in, render login page
		if err := a.AuthMiddleware(w, r); err != nil {
			a.Login(w, r)
			return
		}

		// logged in trying to see /login page, redirect to parameter or /home
		var redirectURL string
		if redirectURL = r.URL.Query().Get("redirect"); redirectURL == "" {
			redirectURL = "/"
		}

		http.Redirect(w, r, redirectURL, http.StatusFound)
	})

	mux.HandleFunc("/log/create", func(w http.ResponseWriter, r *http.Request) {
		if err := a.AuthMiddleware(w, r); err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		a.CreateDailyLog(w, r)
	})

	mux.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		a.GetDailyLog(w, r)
	})

	mux.HandleFunc("/log/{date}", func(w http.ResponseWriter, r *http.Request) {
		a.GetDailyLog(w, r)
	})

	mux.HandleFunc("/log/{id}/edit", func(w http.ResponseWriter, r *http.Request) {
		a.EditDailyLog(w, r)
	})

	mux.HandleFunc("/auth-token", func(w http.ResponseWriter, r *http.Request) {
		a.GetAuthTokens(w, r)
	})

	mux.HandleFunc("/reads", func(w http.ResponseWriter, r *http.Request) {
		component := static_views.Reads(
			layouts_view.StaticHomeData{
				LoggedUser:   a.isLoggedUser(r),
				BuildVersion: a.GitTag,
				BuildTime:    a.BuildTime,
				GitSHA:       a.GitShortSHA,
			})
		component.Render(r.Context(), w)
	})

	mux.HandleFunc("/bookmarks", func(w http.ResponseWriter, r *http.Request) {
		component := static_views.Bookmarks(
			layouts_view.StaticHomeData{
				LoggedUser:   a.isLoggedUser(r),
				BuildVersion: a.GitTag,
				BuildTime:    a.BuildTime,
				GitSHA:       a.GitShortSHA,
			})
		component.Render(r.Context(), w)
	})

	mux.HandleFunc("/birthdays", func(w http.ResponseWriter, r *http.Request) {
		component := static_views.Birthdays(
			layouts_view.StaticHomeData{
				LoggedUser:   a.isLoggedUser(r),
				BuildVersion: a.GitTag,
				BuildTime:    a.BuildTime,
				GitSHA:       a.GitShortSHA,
			})
		component.Render(r.Context(), w)
	})

	mux.HandleFunc("/secrets", func(w http.ResponseWriter, r *http.Request) {
		component := secrets_view.Secrets(
			layouts_view.StaticHomeData{
				LoggedUser:   a.isLoggedUser(r),
				BuildVersion: a.GitTag,
				BuildTime:    a.BuildTime,
				GitSHA:       a.GitShortSHA,
			})
		component.Render(r.Context(), w)
	})

	fs := http.FileServer(http.Dir("./assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	return mux
}

func (a *application) isLoggedUser(r *http.Request) bool {
	loggedUser := true
	authToken, err := r.Cookie("t")

	if err != nil || authToken.Value != a.AuthToken {
		loggedUser = false
	}

	return loggedUser
}
