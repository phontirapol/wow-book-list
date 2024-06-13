package routes

import (
	"net/http"
	"strconv"

	"wow-book-list/db"
	"wow-book-list/model"
	"wow-book-list/templates"

	"github.com/gorilla/mux"
)

type Handler struct {
	BookDB   db.BookDBInterface
	Template templates.TemplateInterface
}

func (h *Handler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	db := h.BookDB.GetDB()
	books, err := model.GetAllBooks(db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = h.Template.ExecuteTemplate(w, "index.html", books)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (h *Handler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	bookIDstr := mux.Vars(r)["bookID"]
	bookID, err := strconv.ParseUint(bookIDstr, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID"))
		return
	}

	db := h.BookDB.GetDB()
	book, err := model.GetBookByID(db, uint(bookID))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("No Book with this ID"))
		return
	}

	err = h.Template.ExecuteTemplate(w, "book.html", book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (h *Handler) NewBookForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	err := h.Template.ExecuteTemplate(w, "new_book.html", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (h *Handler) AddNewBook(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	title := r.FormValue("title")
	author := r.FormValue("author")
	genre := r.FormValue("genre")
	year, _ := strconv.Atoi(r.FormValue("year"))

	db := h.BookDB.GetDB()
	err := model.AddNewBook(db, title, author, genre, year)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	http.Redirect(w, r, "/api/books", http.StatusSeeOther)
}

func (h *Handler) UpdateBook(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) DeleteBook(w http.ResponseWriter, r *http.Request) {}
