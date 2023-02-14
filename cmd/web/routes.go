// routes
package main

import (
	"net/http"
	"path/filepath"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.login)
	mux.HandleFunc("/validation", app.validation)
	mux.HandleFunc("/registration", app.registration)
	mux.HandleFunc("/home", app.home)

	mux.HandleFunc("/user/create", app.createUser)

	mux.HandleFunc("/admin", app.adminLogin)
	mux.HandleFunc("/admin/validation", app.adminValid)
	mux.HandleFunc("/admin/users", app.showAdminPanel)
	mux.HandleFunc("/ban", app.banUser)

	mux.HandleFunc("/note", app.showNote)
	mux.HandleFunc("/note/create", app.createNote)
	mux.HandleFunc("/note/form", app.formNote)
	mux.HandleFunc("/note/data", app.formUpdateNote)
	mux.HandleFunc("/note/update", app.updateNote)
	mux.HandleFunc("/note/delete", app.deleteNote)

	mux.HandleFunc("/project", app.showProject)
	mux.HandleFunc("/projects", app.showProjects)
	mux.HandleFunc("/project/create", app.createProject)
	mux.HandleFunc("/project/form", app.formProject)
	mux.HandleFunc("/project/data", app.formUpdateProject)
	mux.HandleFunc("/project/update", app.updateProject)
	mux.HandleFunc("/project/delete", app.deleteProject)

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err = nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}
	return f, nil
}
