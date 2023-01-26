package main

import (
	"errors"
	"fmt"

	// "html/template"
	"net/http"
	"strconv"

	"dandydev.com/todogo/pkg/models"
)

func (app *application) home (writer http.ResponseWriter, request *http.Request) {
		
	if request.URL.Path != "/" {
		app.NotFound(writer)
		return
	}
	lastsNotes, err := app.notes.Latest()
	if err != nil {
		app.serverError(writer, err)
		return
	
	for _, note := range lastsNotes {
		fmt.Fprintf(writer, "%v\n", note)
	}

	// templatePaths := []string{
	// 	"./ui/html/home.page.tmpl",
	// 	"./ui/html/base.layout.tmpl",
	// 	"./ui/html/footer.partical.tmpl",
	// }
	// templates, err := template.ParseFiles(templatePaths...)
	// if err != nil {
	// 	app.serverError(writer, err)
	// 	return
	// }
	// err = templates.Execute(writer, nil)
	// if err != nil {
	// 	app.serverError(writer, err)
	// 	http.Error(writer, "Internal Server Error", 500)
	// }

}

func (app *application) showNote( writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(writer)
		return
	}

	note, err := app.notes.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(writer)
		} else {
			app.serverError(writer, err)
		}
		return
	}
	fmt.Fprintf(writer, "%v", note)
}

func (app *application) createNote(
	writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.Header().Set("Allow", http.MethodPost)
		app.clientError(writer, http.StatusMethodNotAllowed)
		return
	}

	title := "История про улитку"
	content := "Улитка выползла из раковины," +
		"\nвытянула рожки,\nи опять подобрала их."
	expires := "7"

	id, err := app.notes.Insert(title, content, expires)
	if err != nil {
		app.serverError(writer, err)
		return
	}
	http.Redirect(writer, request, fmt.Sprintf("/note?id=%d", id),
		http.StatusSeeOther)
}
