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

}

func (app *application) formNote(
	writer http.ResponseWriter, request *http.Request) {

	statuses, err := app.statuses.GetAllStatus()
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(writer)
		} else {
			app.serverError(writer, err)
		}
		return
	}
	projects, err := app.projects.GetUserProjects(1)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(writer)
		} else {
			app.serverError(writer, err)
		}
		return
	}
	app.render(writer, request, "create_note.page.tmpl", &templateData{
		Statuses: statuses,
		Projects: projects,
	})
}

func (app *application) formUpdateNote(
	writer http.ResponseWriter, request *http.Request) {
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

	// project, err := app.projects.Get(note.ProjectID)
	// if err != nil {
	// 	if errors.Is(err, models.ErrNoRecord) {
	// 		app.notFound(writer)
	// 	} else {
	// 		app.serverError(writer, err)
	// 	}
	// 	return
	// }

	statuses, err := app.statuses.GetAllStatus()
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(writer)
		} else {
			app.serverError(writer, err)
		}
		return
	}
	projects, err := app.projects.GetUserProjects(1)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(writer)
		} else {
			app.serverError(writer, err)
		}
		return
	}
	app.render(writer, request, "update_note.page.tmpl", &templateData{
		Statuses: statuses,
		Projects: projects,
		Note:     note,
	})
}

func (app *application) createNote(
	writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.Header().Set("Allow", http.MethodPost)
		app.clientError(writer, http.StatusMethodNotAllowed)
		return
	}

	title := request.FormValue("title")
	content := request.FormValue("content")
	expires := request.FormValue("expires")
	projectId := request.FormValue("project")
	statusId := request.FormValue("status")

	id, err := app.notes.Insert(title, content, expires, projectId, statusId)
	if err != nil {
		app.serverError(writer, err)
		return
	}
	http.Redirect(writer, request, fmt.Sprintf("/note?id=%d", id),
		http.StatusSeeOther)
}

func (app *application) updateNote(
	writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.Header().Set("Allow", http.MethodPost)
		app.clientError(writer, http.StatusMethodNotAllowed)
		return
	}

	title := request.FormValue("title")
	content := request.FormValue("content")
	expires := request.FormValue("expires")
	projectId := request.FormValue("project")
	statusId := request.FormValue("status")

	id, err := app.notes.Put(title, content, expires, projectId, statusId)
	if err != nil {
		app.serverError(writer, err)
		return
	}
	http.Redirect(writer, request, fmt.Sprintf("/note?id=%d", id),
		http.StatusSeeOther)
}

func (app *application) showProject(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(writer)
		return
	}

	project, err := app.projects.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(writer)
		} else {
			app.serverError(writer, err)
		}
		return
	}
	notes, err := app.notes.GetProjectNotes(id)
	if err != nil {
		app.serverError(writer, err)
		return
	}

	app.render(writer, request, "table.page.tmpl", &templateData{
		Project: project,
		Notes:   notes,
	})

}
