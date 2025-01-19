package main

import (
	"crypto/md5"
	"encoding/json"
	"io"
	"log"
	"net/http"
    "fmt"
    "io"
    "time"
    "crypto/sha256"
    "encoding/hex"
	"github.com/gorilla/mux"
)

type Block struct{
    Pos int
    Data BookCheckout
    TimeStamp string
    Hash string
    PrevHash string 
}

type Book struct {
    ID          string `json:"id"`
    Title       string `json:"title"`
    Author      string `json:"author"`
    PublishDate string `json:"publish_date"`
    ISBN        string `json:"isbn"`
}


type BookCheckout struct{
    BookID string `json:"book_id"`
    User string `json:"user"`
    CheckoutDate string `json:"checkout_date"`
    IsGenesis bool `json:"is_genesis"`
}

type Blockchain struct{
    blocks []*Block
    
}

var Blockchain *Blockchain

func writeBlock(w http.ResponseWriter, r *http.Request){
    var checkoutitem BookCheckout
    if err := json.NewDecoder(r.Body).Decode(&checkoutitem); err!=nil{
        r.WriteHeader(http.StatusInternalServerError)
        log.Printf("could not write block:%v", err)
        w.Write([]byte("could not write block"))
    }


    
}


func newBook(w http.ResponseWriter, r *http.Request){
    var book Book 
    if err := json.NewDecoder(r.Body).Decode(&book); err!=nil{
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("could not create:%v", err)
        w.Write([]byte("could not create new book"))
        return
    }
    h := md5.New()
    io.WriteString(h, book.ISBN+book.PublishDate)
    book.ID = fmt.Sprintf("%x", h.Sum(nil))
    resp, err := json.MarshalIndent(book, "", " ")
    if err != nil{
        w.WriteHeader(http.StatusInternalServerError)
        log.Printf("could not marshal payload: %v", err)
        w.Write([]byte("could not save book data"))
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Write(resp)
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/", getBlockchain).Methods("GET")
    r.HandleFunc("/", writeBlock).Methods("POST")
    r.HandleFunc("/new", newBook).Methods("POST")

    log.Print("Listening on port 3000")
    log.Fatal(http.ListenAndServe(":3000", r))
}
