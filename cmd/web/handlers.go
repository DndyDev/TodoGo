package main

import (
	"errors"
	"fmt"

	"net/http"
	"strconv"

	"dandydev.com/todogo/pkg/models"
)

func (app *application) home(writer http.ResponseWriter, request *http.Request) {

	if request.URL.Path != "/" {
		app.notFound(writer)
		return
	}
	lastsNotes, err := app.notes.Latest()
	if err != nil {
		app.serverError(writer, err)
		return
	}

	app.render(writer, request, "home.page.tmpl", &templateData{
		Notes: lastsNotes,
	})
	// data := &templateData{Notes: lastsNotes}
	// templateFiles := []string{
	// 	"./ui/html/home.page.tmpl",
	// 	"./ui/html/base.layout.tmpl",
	// 	"./ui/html/footer.partical.tmpl",
	// }
	// templates, err := template.ParseFiles(templateFiles...)
	// if err != nil {
	// 	app.serverError(writer, err)
	// 	return
	// }
	// err = templates.Execute(writer, data)
	// if err != nil {
	// 	app.serverError(writer, err)
	// }

}

func (app *application) showNote(writer http.ResponseWriter, request *http.Request) {
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

	app.render(writer, request, "show.page.tmpl", &templateData{
		Note: note,
	})

	// data := &templateData{Note: note}

	// fileTemplates := []string{
	// 	"./ui/html/show.page.tmpl",
	// 	"./ui/html/base.layout.tmpl",
	// 	"./ui/html/footer.partical.tmpl",
	// }

	// templates, err := template.ParseFiles(fileTemplates...)
	// if err != nil {
	// 	app.serverError(writer, err)
	// 	return
	// }

	// err = templates.Execute(writer, data)
	// if err != nil {
	// 	app.serverError(writer, err)
	// }
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
