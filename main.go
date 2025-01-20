package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"time"
)

type Block struct {
	Pos       int
	Data      BookCheckout
	TimeStamp string
	Hash      string
	PrevHash  string
}

type Book struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	PublishDate string `json:"publish_date"`
	ISBN        string `json:"isbn"`
}

type BookCheckout struct {
	BookID       string `json:"book_id"`
	User         string `json:"user"`
	CheckoutDate string `json:"checkout_date"`
	IsGenesis    bool   `json:"is_genesis"`
}

type Blockchain struct {
	blocks []*Block
}

func (bc *Blockchain) AddBlock(data BookCheckout){
    prevBlock := bc.blocks[len(bc.blocks)-1]
    block := CreateBlock(prevBlock, data)
    if validBlock(block, prevBlock){
        bc.blocks = append(bc.blocks, block)
    }
}

var Blockchain *Blockchain

func (b *Block) generateHash(){
    bytes, _ := json.Marshal(b.Data)
    data := string(b.Pos) + b.TimeStamp + string(bytes) + b.PrevHash
    hash := sha256.New()
    b.Hash = hex.EncodeToString(hash.Sum(nil))
    hash.Write([]byte(data))
}


func CreateBlock(prevBlock *Block, checkoutitem BookCheckout) *Block{
    block := &Block{}
    block.Pos = prevBlock.Pos +1
    block.PrevHash = prevBlock.Hash
    block.TimeStamp = time.Now().String()
    block.generateHash()
    return block
}

func writeBlock(w http.ResponseWriter, r *http.Request) {
	var checkoutitem BookCheckout
	if err := json.NewDecoder(r.Body).Decode(&checkoutitem); err != nil {
		r.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not write block:%v", err)
		w.Write([]byte("could not write block"))
	}
	Blockchain.AddBlock(checkoutitem)

}

func newBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not create:%v", err)
		w.Write([]byte("could not create new book"))
		return
	}
	h := md5.New()
	io.WriteString(h, book.ISBN+book.PublishDate)
	book.ID = fmt.Sprintf("%x", h.Sum(nil))
	resp, err := json.MarshalIndent(book, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not marshal payload: %v", err)
		w.Write([]byte("could not save book data"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func GenesisBlock() *Block{
    return CreateBlock(&Block{}, BookCheckout{IsGenesis:true})
}
func NewBlockChain() *Blockchain{
    return &Blockchain{[]*Block{GenesisBlock()}}
}


func main() {
    Blockchain = NewBlockChain()
	r := mux.NewRouter()
	r.HandleFunc("/", getBlockchain).Methods("GET")
	r.HandleFunc("/", writeBlock).Methods("POST")
	r.HandleFunc("/new", newBook).Methods("POST")

	log.Print("Listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
