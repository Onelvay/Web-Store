package main

import (
	"log"
	"net/http"
	"text/template"

	"strconv"

	db "github.com/Onelvay/Web-Store/database"
	module "github.com/Onelvay/Web-Store/module"
	uuid "github.com/google/uuid"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var Account module.User
var Item module.Product
var Sort bool = false

func createAccount(writer http.ResponseWriter, request *http.Request) {
	login := request.FormValue("login")
	password := hashCode(request.FormValue("password"))
	name := request.FormValue("nickname")
	id := uuid.New()

	if db.IsEmailFree(login) {
		db.CreateUser(id.String(), name, login, password, 0)
		http.Redirect(writer, request, "/home", http.StatusFound)
	}
	http.Redirect(writer, request, "/error", http.StatusFound)
}
func signIn(writer http.ResponseWriter, request *http.Request) {
	login := request.FormValue("login")
	password := hashCode(request.FormValue("password"))
	found, user := db.GetUser(login, password)
	if !found {
		http.Redirect(writer, request, "/error", http.StatusFound)
	}
	Account = user
	http.Redirect(writer, request, "/home", http.StatusFound)

}

func homePage(writer http.ResponseWriter, request *http.Request) {

	html, err := template.ParseFiles("template/signedHome.html")
	check(err)
	products := db.GetProducts(Sort)
	var data = struct {
		User     module.User
		Products []module.Product
		Sort     bool
	}{
		User:     Account,
		Products: products,
		Sort:     Sort,
	}
	err = html.Execute(writer, data)
	check(err)

}

func authPage(writer http.ResponseWriter, request *http.Request) {
	html, err := template.ParseFiles("template/authorization.html")
	check(err)
	err = html.Execute(writer, nil)
	check(err)
}
func ErrorPage(writer http.ResponseWriter, request *http.Request) {
	html, err := template.ParseFiles("template/error.html")
	check(err)
	err = html.Execute(writer, nil)
	check(err)
}
func signInPage(writer http.ResponseWriter, request *http.Request) {
	html, err := template.ParseFiles("template/signin.html")
	check(err)
	err = html.Execute(writer, nil)
	check(err)
}
func PurchasePage(writer http.ResponseWriter, request *http.Request) {
	if Account.Name == "" {
		http.Redirect(writer, request, "/home/signin", http.StatusFound)
	}
	html, err := template.ParseFiles("template/purchase.html")
	check(err)
	id := request.FormValue("id")
	product := db.GetProduct(id)
	var data = struct {
		User     module.User
		Products module.Product
	}{
		User:     Account,
		Products: product,
	}
	Item = product
	err = html.Execute(writer, data)
	check(err)
}
func Purchase(writer http.ResponseWriter, request *http.Request) {
	if db.BuyItem(Account, Item) {
		Account.Balance -= Item.Price
	}
	http.Redirect(writer, request, "/home", http.StatusFound)

}
func SortPrice(writer http.ResponseWriter, request *http.Request) {
	if Sort {
		Sort = false
	} else {
		Sort = true
	}

	http.Redirect(writer, request, "/home", http.StatusFound)

}
func hashCode(s string) int {
	result := 1
	for _, c := range s {
		result += int(c)
	}
	return result
}
func Search(writer http.ResponseWriter, request *http.Request) {
	html, err := template.ParseFiles("template/search.html")
	check(err)
	search := request.FormValue("search")
	products := db.GetProductsByName(search)
	var data = struct {
		Search   string
		Products []module.Product
	}{
		Search:   search,
		Products: products,
	}
	err = html.Execute(writer, data)
	check(err)
}
func userOrders(writer http.ResponseWriter, request *http.Request) {
	html, err := template.ParseFiles("template/userOrders.html")
	check(err)
	products := db.GetOrders(Account.Id)
	var data = struct {
		Products []module.Order
	}{
		Products: products,
	}
	err = html.Execute(writer, data)
	check(err)
}

func RateOrder(writer http.ResponseWriter, request *http.Request) {
	rate := request.FormValue("quantity")
	id := request.FormValue("id")
	i, err := strconv.Atoi(rate)
	check(err)
	db.SetRate(id, i)
	http.Redirect(writer, request, "/home", http.StatusFound)

}

func main() {
	http.HandleFunc("/home/user/products", userOrders)
	http.HandleFunc("/home/authorization/create", createAccount)
	http.HandleFunc("/home", homePage)
	http.HandleFunc("/home/signin/create", signIn)
	http.HandleFunc("/home/signin", signInPage)
	http.HandleFunc("/home/authorization", authPage)
	http.HandleFunc("/error", ErrorPage)
	http.HandleFunc("/purchase", PurchasePage)
	http.HandleFunc("/buying", Purchase)
	http.HandleFunc("/sort", SortPrice)
	http.HandleFunc("/search", Search)
	http.HandleFunc("/rate", RateOrder)
	err := http.ListenAndServe("localhost:8080", nil)
	log.Fatal(err)
}
