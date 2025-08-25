package main

import (
	"net/http"

	"nugu.dev/basement/views/static_views"
)

func (a *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", a.Landing)

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

	mux.HandleFunc("/reads", func(w http.ResponseWriter, r *http.Request) {
		component := static_views.Reads()
		component.Render(r.Context(), w)
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

	mux.HandleFunc("/bookmarks", func(w http.ResponseWriter, r *http.Request) {
		component := static_views.Bookmarks()
		component.Render(r.Context(), w)
	})

	mux.HandleFunc("/birthdays", func(w http.ResponseWriter, r *http.Request) {
		component := static_views.Birthdays()
		component.Render(r.Context(), w)
	})

	fs := http.FileServer(http.Dir("./assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))

	return mux
}
