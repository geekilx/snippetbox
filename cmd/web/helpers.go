package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// render is a helper method for rendering templates. It takes the
// http.ResponseWriter, the HTTP status code, the name of the template to
// render, and any dynamic data to pass to the template.
func (app *application) render(w http.ResponseWriter, status int, pageName string, data *templateData) {

	ts, ok := app.templateCache[pageName]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", pageName))
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}
	w.WriteHeader(status)

	buf.WriteTo(w)

}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}
