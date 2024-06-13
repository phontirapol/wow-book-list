package http

import (
	"fmt"
	"log"
	"net/http"

	"wow-book-list/db"
	"wow-book-list/routes"
	"wow-book-list/templates"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

const (
	GET    = http.MethodGet
	POST   = http.MethodPost
	PUT    = http.MethodPut
	DELETE = http.MethodDelete
)

func initRouter(bookDB db.BookDBInterface, template templates.TemplateInterface) *mux.Router {
	router := mux.NewRouter()

	handler := &routes.Handler{
		BookDB:   bookDB,
		Template: template,
	}

	router.HandleFunc("/api/books", handler.GetAllBooks).Methods(GET)
	router.HandleFunc("/api/books/:{bookID}", handler.GetBookByID).Methods(GET)
	router.HandleFunc("/api/books/newbook", handler.NewBookForm).Methods(GET)
	router.HandleFunc("/api/books", handler.AddNewBook).Methods(POST)
	// router.HandleFunc("/api/books/{bookID}", handler.UpdateBook).Methods(PUT)
	// router.HandleFunc("/api/books/{bookID}", handler.DeleteBook).Methods(DELETE)

	return router
}

func StartServer() {
	fmt.Println("Starting the server")

	bookDB := db.InitDB()

	template := &templates.Template{}
	if err := template.LoadTemplates("static/*.html"); err != nil {
		log.Fatal(err)
	}

	router := initRouter(bookDB, template)
	http.Handle("/", router)

	err := http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}
}
