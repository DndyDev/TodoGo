//Парсинг настроек конфигурации среды выполнения для приложения;
//Установление зависимостей для обработчиков;
//Запуск HTTP-сервера.

package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {

	addres := flag.String("addr", ":4000", "Web addres for HTTP")

	flag.Parse()

	f, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:     *addres,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Server start on %s", *addres)
	srv.ListenAndServe()
	errorLog.Fatal(err)
}
