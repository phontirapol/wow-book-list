package model

import (
	"database/sql"
	"encoding/json"
)

type Book struct {
	ID            uint   `json:"book_id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	Genre         string `json:"genre"`
	PublishedYear int    `json:"published_year,string"`
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

func AddNewBook(db *sql.DB, title, author, genre string, year int) error {
	ps, err := db.Begin()
	if err != nil {
		return err
	}

	insertStatement := `
		INSERT INTO book(title, author, genre, published_year) VALUES(?,?,?,?)
	`

	statement, err := db.Prepare(insertStatement)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(title, author, genre, year)
	if err != nil {
		return err
	}

	err = ps.Commit()
	if err != nil {
		return err
	}

	return nil
}

func UpdateBook(db *sql.DB, bookID uint, bookData []byte) error {
	book := Book{ID: bookID}
	err := json.Unmarshal(bookData, &book)
	if err != nil {
		return err
	}

	ps, err := db.Begin()
	if err != nil {
		return err
	}

	updateStatement := `
		UPDATE book 
		SET title = ?, author = ?, genre = ?, published_year = ?
		WHERE book_id = ?
	`
	statement, err := db.Prepare(updateStatement)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(book.Title, book.Author, book.Genre, book.PublishedYear, book.ID)
	if err != nil {
		return err
	}

	err = ps.Commit()
	if err != nil {
		return err
	}

	return nil
}
