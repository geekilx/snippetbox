package main

import "net/http"

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.requestLog.Println(r.URL)

		next.ServeHTTP(w, r)

	})
}
