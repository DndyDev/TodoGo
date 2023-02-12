//Парсинг настроек конфигурации среды выполнения для приложения;
//Установление зависимостей для обработчиков;
//Запуск HTTP-сервера.

package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"dandydev.com/todogo/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	notes         *mysql.NoteModel
	projects      *mysql.ProjectModel
	statuses      *mysql.StatusModel
	users         *mysql.UserModel
	admins        *mysql.AdminModel
	templateCache map[string]*template.Template
}

func main() {

	addres := flag.String("addr", ":4000", "Web addres for HTTP")
	dsn := flag.String("dsn", "web:localhost@/todogo?parseTime=true",
		"Название MySQL источника данных")

	flag.Parse()

	f, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCashe("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		notes:         &mysql.NoteModel{DB: db},
		projects:      &mysql.ProjectModel{DB: db},
		statuses:      &mysql.StatusModel{DB: db},
		users:         &mysql.UserModel{DB: db},
		admins:        &mysql.AdminModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addres,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Server start on localhost%s", *addres)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
