package model

import (
	"database/sql"
)

type Book struct {
	ID            uint   `json:"book_id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	Genre         string `json:"genre"`
	PublishedYear int    `json:"published_year"`
}

func GetAllBooks(db *sql.DB) ([]*Book, error) {
	var books []*Book

	statement := "SELECT book_id, title, author, genre, published_year FROM book"
	rows, err := db.Query(statement)
	if err == sql.ErrNoRows {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		book := Book{}
		_ = rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Genre,
			&book.PublishedYear,
		)

		books = append(books, &book)
	}

	return books, nil
}

func GetBookByID(db *sql.DB, bookID uint) (*Book, error) {
	book := Book{}

	statement := "SELECT book_id, title, author, genre, published_year FROM book WHERE book_id = ?"
	row := db.QueryRow(statement, bookID)
	err := row.Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.Genre,
		&book.PublishedYear,
	)
	if err == sql.ErrNoRows {
		return nil, err
	}

	return &book, nil
}
