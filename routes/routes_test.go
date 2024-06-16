package routes_test

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"wow-book-list/routes"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const (
	GET = http.MethodGet
)

type mockBookDB struct {
	Database *sql.DB
}

func (db *mockBookDB) GetDB() *sql.DB {
	return db.Database
}

type mockTemplate struct {
	err error
}

func (t *mockTemplate) ExecuteTemplate(w http.ResponseWriter, template string, data interface{}) error {
	return t.err
}

func TestGetBookByID(t *testing.T) {
	baseUrl := "/api/books/:"

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error not expected while opening mock db, %v", err)
	}

	queryStatement := "SELECT book_id, title, author, genre, published_year FROM book WHERE book_id = ?"

	t.Run("Invalid ID", func(t *testing.T) {
		testSuite := []string{
			"-1",
			"1.2",
			"abc",
			"1/1",
		}

		for _, test := range testSuite {
			urlVar := map[string]string{"bookID": test}

			r, err := http.NewRequest(GET, baseUrl+urlVar["bookID"], nil)
			if err != nil {
				t.Error(err)
			}

			r = mux.SetURLVars(r, urlVar)

			mockDB := &mockBookDB{Database: db}

			mockHandler := &routes.Handler{
				BookDB: mockDB,
			}

			w := httptest.NewRecorder()
			handler := http.HandlerFunc(mockHandler.GetBookByID)
			handler.ServeHTTP(w, r)

			expectedCode := http.StatusBadRequest
			expectedError := "Invalid ID"

			if assert.Equal(t, expectedCode, w.Code) {
				assert.Equal(t, expectedError, w.Body.String())
			}
		}
	})

	t.Run("No book with this ID", func(t *testing.T) {
		urlVar := map[string]string{"bookID": "888"}

		r, err := http.NewRequest(GET, baseUrl+urlVar["bookID"], nil)
		if err != nil {
			t.Error(err)
		}

		r = mux.SetURLVars(r, urlVar)

		mock.ExpectQuery(queryStatement).
			WillReturnError(sql.ErrNoRows)

		mockDB := &mockBookDB{Database: db}

		mockHandler := &routes.Handler{
			BookDB: mockDB,
		}

		w := httptest.NewRecorder()
		handler := http.HandlerFunc(mockHandler.GetBookByID)
		handler.ServeHTTP(w, r)

		expectedCode := http.StatusUnprocessableEntity
		expectedError := "No Book with this ID"

		if assert.Equal(t, expectedCode, w.Code) {
			assert.Equal(t, expectedError, w.Body.String())
		}
	})

	t.Run("Template error", func(t *testing.T) {
		var (
			id     string = "888"
			title  string = "Harry Potter"
			author string = "JK Rolling in the Deep"
			genre  string = "Horror"
			year   int    = 2020
		)

		urlVar := map[string]string{"bookID": id}

		r, err := http.NewRequest(GET, baseUrl+urlVar["bookID"], nil)
		if err != nil {
			t.Error(err)
		}

		r = mux.SetURLVars(r, urlVar)

		mockRow := sqlmock.NewRows([]string{"book_id", "title", "author", "genre", "published_year"}).
			AddRow(id, title, author, genre, year)

		mock.ExpectQuery(queryStatement).
			WithArgs(id).
			WillReturnRows(mockRow)

		mockDB := &mockBookDB{Database: db}

		expectedError := "Random template error"
		mockTmpl := &mockTemplate{err: errors.New(expectedError)}

		mockHandler := &routes.Handler{
			BookDB:   mockDB,
			Template: mockTmpl,
		}

		w := httptest.NewRecorder()
		handler := http.HandlerFunc(mockHandler.GetBookByID)
		handler.ServeHTTP(w, r)

		expectedCode := http.StatusInternalServerError

		if assert.Equal(t, expectedCode, w.Code) {
			assert.Equal(t, expectedError, w.Body.String())
		}
	})

	t.Run("Happy path", func(t *testing.T) {
		var (
			id     string = "888"
			title  string = "Harry Potter"
			author string = "JK Rolling in the Deep"
			genre  string = "Horror"
			year   int    = 2020
		)

		urlVar := map[string]string{"bookID": id}

		r, err := http.NewRequest(GET, baseUrl+urlVar["bookID"], nil)
		if err != nil {
			t.Error(err)
		}

		r = mux.SetURLVars(r, urlVar)

		mockRow := sqlmock.NewRows([]string{"book_id", "title", "author", "genre", "published_year"}).
			AddRow(id, title, author, genre, year)

		mock.ExpectQuery(queryStatement).
			WithArgs(id).
			WillReturnRows(mockRow)

		mockDB := &mockBookDB{Database: db}

		mockTmpl := &mockTemplate{err: nil}

		mockHandler := &routes.Handler{
			BookDB:   mockDB,
			Template: mockTmpl,
		}

		w := httptest.NewRecorder()
		handler := http.HandlerFunc(mockHandler.GetBookByID)
		handler.ServeHTTP(w, r)

		expectedCode := http.StatusOK

		assert.Equal(t, expectedCode, w.Code)
	})
}
