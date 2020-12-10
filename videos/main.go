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

type Video struct {
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
	router.HandleFunc("/", getvideos).Methods("GET")
	router.HandleFunc("/{ID}", getvideo).Methods("GET")
	router.HandleFunc("/", addvideo).Methods("POST")
	router.HandleFunc("/", updatevideos).Methods("PUT")
	router.HandleFunc("/", deletevideo).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":80", router))
	//logfatal(err)

}

func getvideos(w http.ResponseWriter, r *http.Request) {

	log.Printf("get all videos")

	var videos []Video
	var video Video

	rows, err := db.Query("select * from videos")
	logfatal(err)

	defer rows.Close()

	for rows.Next() {
		//video = new(video)
		err := rows.Scan(&video.ID, &video.Title, &video.Author)
		logfatal(err)
		videos = append(videos, video)
	}
	//fmt.Printf("%+v", videos)
	//db.Close()
	json.NewEncoder(w).Encode(videos)
}

func getvideo(w http.ResponseWriter, r *http.Request) {

	log.Printf("get a video")

	var videos []Video
	var video Video

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["ID"])
	logfatal(err)

	rows, err := db.Query("select *from videos where id = $1", id)
	logfatal(err)

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&video.ID, &video.Title, &video.Author)
		logfatal(err)
		videos = append(videos, video)
	}
	//db.Close()
	json.NewEncoder(w).Encode(videos)
}
func addvideo(w http.ResponseWriter, r *http.Request) {
	log.Println("add a video")

	var video Video
	var videoid int
	json.NewDecoder(r.Body).Decode(&video)

	err := db.QueryRow("insert into videos(title,author)values($1,$2) returning id;", video.Title, video.Author).Scan(&videoid)
	logfatal(err)

	json.NewEncoder(w).Encode(videoid)
}
func updatevideos(w http.ResponseWriter, r *http.Request) {

	log.Println("update a video")

	var video Video
	//var videotitle string
	json.NewDecoder(r.Body).Decode(&video)

	result, err := db.Exec("update videos set title = $1, author = $2 where id = $3 returning id;", video.Title, video.Author, video.ID)
	//err := db.QueryRow("update videos set title=$1, author=$2 where id=$3 returning title;", video.Title, video.Author, video.ID).Scan(&videotitle)
	logfatal(err)
	rowsUpdated, err := result.RowsAffected()
	logfatal(err)

	json.NewEncoder(w).Encode(rowsUpdated)

}
func deletevideo(w http.ResponseWriter, r *http.Request) {
	log.Println("delete a video")

	var video Video
	json.NewDecoder(r.Body).Decode(&video)
	log.Println(video.ID)

	result, err := db.Exec("delete from videos where id = $1 ;", video.ID)

	logfatal(err)

	rowsDeleted, err := result.RowsAffected()

	logfatal(err)

	json.NewEncoder(w).Encode(rowsDeleted)
}
