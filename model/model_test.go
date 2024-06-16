package model_test

import (
	"database/sql"
	"testing"

	"wow-book-list/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetBookByID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error(err)
	}

	queryStatement := "SELECT book_id, title, author, genre, published_year FROM book WHERE book_id = ?"

	t.Run("No book found", func(t *testing.T) {
		var id uint = 0

		mock.ExpectQuery(queryStatement).
			WithArgs(id).
			WillReturnError(sql.ErrNoRows)

		book, err := model.GetBookByID(db, id)
		if assert.Error(t, err) {
			assert.ErrorIs(t, err, sql.ErrNoRows)
			assert.Nil(t, book)
		}
	})

	t.Run("Happy path", func(t *testing.T) {
		var (
			id     uint   = 888
			title  string = "Harry Potter"
			author string = "JK Rolling in the Deep"
			genre  string = "Horror"
			year   int    = 2020
		)

		mockRow := sqlmock.NewRows([]string{"book_id", "title", "author", "genre", "published_year"}).
			AddRow(id, title, author, genre, year)

		mock.ExpectQuery(queryStatement).
			WithArgs(id).
			WillReturnRows(mockRow)

		book, err := model.GetBookByID(db, id)
		if assert.NoError(t, err) {
			assert.Equal(t, title, book.Title)
			assert.Equal(t, author, book.Author)
			assert.Equal(t, genre, book.Genre)
			assert.Equal(t, year, book.PublishedYear)
		}
	})
}
