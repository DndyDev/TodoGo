package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"dandydev.com/todogo/pkg/models"
)

func (app *application) login(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		app.notFound(writer)
		return
	}
	app.render(writer, request, "login.page.tmpl", &templateData{})
}
func (app *application) validation(
	writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.Header().Set("Allow", http.MethodPost)
		app.clientError(writer, http.StatusMethodNotAllowed)
		return
	}

	login := request.FormValue("login")
	password := request.FormValue("password")

	user, err := app.users.Get(login, password)
	if err != nil {
		app.serverError(writer, err)
		return
	}
	id := user.ID
	http.Redirect(writer, request, fmt.Sprintf("/home?id=%d", id),
		http.StatusSeeOther)
}

func (app *application) registration(
	writer http.ResponseWriter, request *http.Request) {
	app.render(writer, request, "registration.page.tmpl", &templateData{})

}
func (app *application) adminLogin(
	writer http.ResponseWriter, request *http.Request) {
	app.render(writer, request, "admin.page.tmpl", &templateData{})
}
func (app *application) adminValid(
	writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.Header().Set("Allow", http.MethodPost)
		app.clientError(writer, http.StatusMethodNotAllowed)
		return
	}

	login := request.FormValue("login")
	password := request.FormValue("password")

	_, err := app.admins.Get(login, password)
	if err != nil {
		app.serverError(writer, err)
		return
	}
	http.Redirect(writer, request, fmt.Sprintf("/admin/users"),
		http.StatusSeeOther)
}
func (app *application) banUser(
	writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(writer)
		return
	}
	err = app.users.Ban(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(writer)
		} else {
			app.serverError(writer, err)
		}
		return
	}
	http.Redirect(writer, request, fmt.Sprintf("/admin/users"),
		http.StatusSeeOther)
}
func (app *application) unBanUser(
	writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(writer)
		return
	}
	err = app.users.UnBan(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(writer)
		} else {
			app.serverError(writer, err)
		}
		return
	}
	http.Redirect(writer, request, fmt.Sprintf("/admin/users"),
		http.StatusSeeOther)
}

func (app *application) showAdminPanel(
	writer http.ResponseWriter, request *http.Request) {

	users, err := app.users.GetAll()
	if err != nil {
		app.serverError(writer, err)
		return
	}
	app.render(writer, request, "users.page.tmpl", &templateData{
		Users: users,
	})
}

func (app *application) createUser(
	writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.Header().Set("Allow", http.MethodPost)
		app.clientError(writer, http.StatusMethodNotAllowed)
		return
	}

	login := request.FormValue("login")
	password := request.FormValue("password")
	lastName := request.FormValue("lastName")
	email := request.FormValue("email")

	err := app.users.Insert(login, lastName, email, password)
	if err != nil {
		app.serverError(writer, err)
		return
	}
	http.Redirect(writer, request, fmt.Sprintf("/"),
		http.StatusSeeOther)
}
func (app *application) home(writer http.ResponseWriter, request *http.Request) {

	if request.URL.Path != "/home" {
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

	app.render(writer, request, "note.page.tmpl", &templateData{
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
	app.render(writer, request, "create.note.page.tmpl", &templateData{
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
	app.render(writer, request, "update.note.page.tmpl", &templateData{
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
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(writer)
		return
	}
	title := request.FormValue("title")
	content := request.FormValue("content")
	expires := request.FormValue("expires")
	projectId := request.FormValue("project")
	statusId := request.FormValue("status")

	err = app.notes.Put(id, title, content, expires, projectId, statusId)
	if err != nil {
		app.serverError(writer, err)
		return
	}
	http.Redirect(writer, request, fmt.Sprintf("/note?id=%d", id),
		http.StatusSeeOther)
}

func (app *application) deleteNote(
	writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(writer)
		return
	}
	err = app.notes.Delete(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(writer)
		} else {
			app.serverError(writer, err)
		}
		return
	}
	http.Redirect(writer, request, fmt.Sprintf("/home"),
		http.StatusSeeOther)
}

func (app *application) showProjects(
	writer http.ResponseWriter, request *http.Request) {
	projects, err := app.projects.GetUserProjects(1)
	if err != nil {
		app.serverError(writer, err)
		return
	}
	app.render(writer, request, "list.project.page.tmpl", &templateData{
		Projects: projects,
	})
}
func (app *application) showProject(
	writer http.ResponseWriter, request *http.Request) {
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
	statuses, err := app.statuses.GetAllStatus()
	if err != nil {
		app.serverError(writer, err)
		return
	}

	app.render(writer, request, "project.page.tmpl", &templateData{
		Project:  project,
		Notes:    notes,
		Statuses: statuses,
	})

}
func (app *application) searchNotes(
	writer http.ResponseWriter, request *http.Request) {

	projectId, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil || projectId < 1 {
		app.notFound(writer)
		return
	}

	projects, err := app.projects.GetUserProjects(1)
	if err != nil {
		app.serverError(writer, err)
		return
	}
	statuses, err := app.statuses.GetAllStatus()
	if err != nil {
		app.serverError(writer, err)
		return
	}

	app.render(writer, request, "search.page.tmpl", &templateData{
		Projects: projects,
		Statuses: statuses,
	})
}
func (app *application) showNotesToStatus(
	writer http.ResponseWriter, request *http.Request) {

	projectId := request.FormValue("project")
	statusId := request.FormValue("status")

	notes, err := app.notes.GetProjectNotesWitchStatus(projectId, statusId)
	if err != nil {
		app.serverError(writer, err)
		return
	}

	app.render(writer, request, "list.notes.page.tmpl", &templateData{
		Notes: notes,
	})
}
func (app *application) formProject(
	writer http.ResponseWriter, request *http.Request) {

	app.render(writer, request, "create.project.page.tmpl", &templateData{})
}
func (app *application) formUpdateProject(
	writer http.ResponseWriter, request *http.Request) {
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

	app.render(writer, request, "update.project.page.tmpl", &templateData{
		Project: project,
	})
}
func (app *application) createProject(
	writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.Header().Set("Allow", http.MethodPost)
		app.clientError(writer, http.StatusMethodNotAllowed)
		return
	}

	title := request.FormValue("title")
	userId := 1
	id, err := app.projects.Insert(title, userId)
	if err != nil {
		app.serverError(writer, err)
		return
	}
	http.Redirect(writer, request, fmt.Sprintf("/project?id=%d", id),
		http.StatusSeeOther)
}

func (app *application) updateProject(
	writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.Header().Set("Allow", http.MethodPost)
		app.clientError(writer, http.StatusMethodNotAllowed)
		return
	}
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(writer)
		return
	}

	title := request.FormValue("title")
	err = app.projects.Put(id, title)
	if err != nil {
		app.serverError(writer, err)
		return
	}
	http.Redirect(writer, request, fmt.Sprintf("/project?id=%d", id),
		http.StatusSeeOther)
}

func (app *application) deleteProject(
	writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(writer)
		return
	}
	err = app.projects.Delete(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(writer)
		} else {
			app.serverError(writer, err)
		}
		return
	}
	http.Redirect(writer, request, fmt.Sprintf("/home"),
		http.StatusSeeOther)
}
