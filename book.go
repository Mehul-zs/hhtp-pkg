package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Book struct {
	id            int
	Title         string  `json:"title"`
	Author        *Author `json:"-"`
	Publication   string  `json:"publication"`
	PublishedDate string  `json:"published_date"`
}

type Author struct {
	id        int
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Dob       string `json:"dob"`
	PenName   string `json:"pen_name"`
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "HelloMehul1@"
	dbName := "test"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(127.0.0.1:3306)/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func getBook(response http.ResponseWriter, request *http.Request) {
	db := dbConn()
	defer dbConn().Close()
	title := request.URL.Query().Get("Title")
	Author := request.URL.Query().Get("Author")
	var rows *sql.Rows
	var err error
	if title == "" {
		rows, err = db.Query("select * from Books;")
	} else {
		rows, err = db.Query("select * from Books where title=?;", title)
	}
	if err != nil {
		log.Print(err)
	}
	books := []Book{}
	for rows.Next() {
		book := Book{}
		err = rows.Scan(&book.id, &book.Title, &book.Publication, &book.PublishedDate, &book.Author.id)
		if err != nil {
			log.Print(err)
		}
		if Author == "true" {
			row := db.QueryRow("select * from Authors where id=?", book.Author.id)
			row.Scan(&book.Author.id, &book.Author.FirstName, &book.Author.LastName, &book.Author.Dob, &book.Author.PenName)
		}
		books = append(books, book)
	}
	json.NewEncoder(response).Encode(books)
	//json.NewEncoder(response).Encode([]Book{{1, "Mehul", nil, "Penguin", "10/08/2000"}, {2, "Mehul", nil, "Penguin", "10/08/2000"}})
}

func getBookById(res http.ResponseWriter, req *http.Request) {
	//json.NewEncoder(res).Encode(Book{1, "Charlie", nil, "", "10/08/2000"})

	id := mux.Vars(req)["id"]
	db := dbConn()
	defer db.Close()
	bookrow := db.QueryRow("select * from Books where id=?;", id)
	book := Book{}
	author := Author{}
	author_id := 0
	err := bookrow.Scan(&book.id, &book.Title, &book.Publication, &book.PublishedDate, &author_id)
	if err != nil {
		log.Print(err)
	}
	authorrow := db.QueryRow("select * from Authors where id=?;", author_id)
	err = authorrow.Scan(&author.id, &author.FirstName, &author.LastName, &author.Dob, &author.PenName)
	if err != nil {
		log.Print(err)
	}
	book.Author = Author
	json.NewEncoder(res).Encode(book)
}

func PostByBook(res http.ResponseWriter, req *http.Request) {

}

func main() {

	fmt.Println("Go MySQL Tutorial")

	db, _ := sql.Open("mysql", "root:HelloMehul1@@tcp(127.0.0.1:3306)/bookstore")
	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/book", getBook).Methods(http.MethodGet)
	r.HandleFunc("/book/{id}", getBookById).Methods(http.MethodGet)
	r.HandleFunc("/book", PostByBook).Methods(http.MethodPost)

	Server := http.Server{
		Addr:    ":8000",
		Handler: r,
	}

	fmt.Println("Server started at 8000")
	Server.ListenAndServe()
}
