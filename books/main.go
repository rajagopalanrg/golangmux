package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/subosito/gotenv"
)

var db *sql.DB

type Book struct {
	ID     int    `json:id`
	Title  string `json:title`
	Author string `json:author`
}

func init() {
	gotenv.Load()
}
func logfatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	pgurl, err := pq.ParseURL(os.Getenv("postgres_url"))

	if err != nil {
		logfatal(err)
	}
	//pginfo := fmt.Sprint("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err = sql.Open("postgres", pgurl)
	router := mux.NewRouter()
	router.HandleFunc("/", getBooks).Methods("GET")
	router.HandleFunc("/{ID}", getBook).Methods("GET")
	router.HandleFunc("/", addBook).Methods("POST")
	router.HandleFunc("/", updateBooks).Methods("PUT")
	router.HandleFunc("/", deleteBook).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":80", router))
	//logfatal(err)

}

func getBooks(w http.ResponseWriter, r *http.Request) {

	log.Printf("get all books")

	var books []Book
	var book Book

	rows, err := db.Query("select * from books")
	logfatal(err)

	defer rows.Close()

	for rows.Next() {
		//book = new(Book)
		err := rows.Scan(&book.ID, &book.Title, &book.Author)
		logfatal(err)
		books = append(books, book)
	}
	//fmt.Printf("%+v", books)
	//db.Close()
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {

	log.Printf("get a book")

	var books []Book
	var book Book

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["ID"])
	logfatal(err)

	rows, err := db.Query("select *from books where id = $1", id)
	logfatal(err)

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&book.ID, &book.Title, &book.Author)
		logfatal(err)
		books = append(books, book)
	}
	//db.Close()
	json.NewEncoder(w).Encode(books)
}
func addBook(w http.ResponseWriter, r *http.Request) {
	log.Println("add a book")

	var book Book
	var bookid int
	json.NewDecoder(r.Body).Decode(&book)

	err := db.QueryRow("insert into books(title,author)values($1,$2) returning id;", book.Title, book.Author).Scan(&bookid)
	logfatal(err)

	json.NewEncoder(w).Encode(bookid)
}
func updateBooks(w http.ResponseWriter, r *http.Request) {

	log.Println("update a book")

	var book Book
	//var booktitle string
	json.NewDecoder(r.Body).Decode(&book)

	result, err := db.Exec("update books set title = $1, author = $2 where id = $3 returning id;", book.Title, book.Author, book.ID)
	//err := db.QueryRow("update books set title=$1, author=$2 where id=$3 returning title;", book.Title, book.Author, book.ID).Scan(&booktitle)
	logfatal(err)
	rowsUpdated, err := result.RowsAffected()
	logfatal(err)

	json.NewEncoder(w).Encode(rowsUpdated)

}
func deleteBook(w http.ResponseWriter, r *http.Request) {
	log.Println("delete a book")

	var book Book
	json.NewDecoder(r.Body).Decode(&book)
	log.Println(book.ID)

	result, err := db.Exec("delete from books where id = $1 ;", book.ID)

	logfatal(err)

	rowsDeleted, err := result.RowsAffected()

	logfatal(err)

	json.NewEncoder(w).Encode(rowsDeleted)
}
