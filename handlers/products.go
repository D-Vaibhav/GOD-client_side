package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vaibhav/GOD/client_side/data"
)

// it's because maybe different handler can use its own way to process entity
// say ALOO is an identity then
// some will make potato chips, some will make aallo ka paratha
type Products struct {
	// think of this as their own uniique way to use data
	l *log.Logger
}

// ABSTRACTION
// Product struct constructor, for creating Product struct
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

type keyProduct struct{}

// ------------------------------------------ GET request -----------------------------------------
func (p *Products) GetProducts(rw http.ResponseWriter, req *http.Request) {
	p.l.Println("Products Handler invoked, GetProducts() executing now...")

	productList := data.GetProducts()

	err := productList.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to Marshall(Encode) json", http.StatusInternalServerError)
	}
	p.l.Println("Done.!")
}

// -------------------------------------------- POST request ----------------------------------------
func (p *Products) AddProduct(rw http.ResponseWriter, req *http.Request) {
	p.l.Println("Products Handler invoked, AddProduct() executing now...")

	// creating new product using req.Context()
	// this Value(it'll be working as map)
	prod := req.Context().Value(keyProduct{}).(data.Product) // last (what struct is supposed to be as value )

	// adding product to productList(DB)
	data.AddProduct(&prod)
}

// ------------------------------------------- PUT request---------------------------------------------
func (p Products) UpdateProduct(rw http.ResponseWriter, req *http.Request) {
	// getting the identifier as key-value map (as string) from the input uri
	varMap := mux.Vars(req)

	// converting from string to int
	productID, err := strconv.Atoi(varMap["productID"])

	// creating a new product
	prod := req.Context().Value(keyProduct{}).(data.Product)

	// ensure that any product with productID exists
	// if exists then over-write this empty passed product
	err = data.UpdateProduct(productID, &prod)

	if err == data.ErrorProductNotFound {
		http.Error(rw, "Product not Found in our DB", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Product not Found in our DB", http.StatusInternalServerError)
		return
	}
}
