package main

import (
    "database/sql"
    "encoding/json"
   "fmt"
    "net/http"

    _ "github.com/go-sql-driver/mysql"
)

var db *sql.DB


func insertHandler(w http.ResponseWriter, r *http.Request) {
    var data = struct {
        Author string `json:"author"`
        Book   string `json:"book"`
    } {"JK Rowling", "Harry Potter"}
    // json.NewDecoder(r.Body).Decode(&data)

    stmt, err := db.Prepare("INSERT INTO books(author, book) VALUES(?, ?)")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    _, err = stmt.Exec(data.Author, data.Book)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        fmt.Println(err)
        return
    }
    w.Write([]byte("success"))
}

func fetchHandler(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT author, book FROM books")
    if err != nil {

        http.Error(w, err.Error(), http.StatusInternalServerError)
        fmt.Println(err)
        return
    }
    defer rows.Close()

    var books []struct {
        Author string
        Book   string
    }
    for rows.Next() {
        var book struct {
            Author string
            Book   string
        }
        rows.Scan(&book.Author, &book.Book)
        books = append(books, book)
    }

    json.NewEncoder(w).Encode(books)
}

func main() {
    var err error
    db, err = sql.Open("mysql", "user:password@tcp(localhost:3306)/SHOW DATABASE")
    if err != nil {
        fmt.Println(err)
        panic(err.Error())
    }
    defer db.Close()

    http.HandleFunc("/insert", insertHandler)
    http.HandleFunc("/fetch", fetchHandler)
    http.ListenAndServe("localhost:8080", nil)
}