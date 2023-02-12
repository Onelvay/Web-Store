package database

import (
	"database/sql"
	"fmt"

	module "github.com/Onelvay/Web-Store/module"
	uuid "github.com/google/uuid"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Adg12332,"
	dbname   = "Store"
)

var Db *sql.DB

var id string
var name string
var email string
var balance float64
var price float64

func init() {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error
	Db, err = sql.Open("postgres", psqlconn)
	CheckError(err)
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
func CreateUser(id string, name string, email string, password int, balance float64) {

	insertDynStmt := `insert into "users"("id","name", "email","password","balance") values($1, $2,$3,$4,$5)`
	_, e := Db.Exec(insertDynStmt, id, name, email, password, 0.0)
	CheckError(e)
}

func AvgProductRate(id string) (int, float64) {
	rows, err := Db.Query(`SELECT count("product_id"),avg("user_rate") from "orders"  where "user_rate"!=0  group by "product_id" having "product_id"=$1`, id)
	CheckError(err)
	var avg float64
	var cnt int
	for rows.Next() {

		err = rows.Scan(&cnt, &avg)
		CheckError(err)
		return cnt, avg
	}
	return cnt, avg
}

func GetProducts(sort bool) []module.Product {
	var err error
	var rows *sql.Rows
	if sort {
		rows, err = Db.Query(`SELECT "id","name","price" FROM "products" order by "price" desc`)
	} else {
		rows, err = Db.Query(`SELECT "id","name","price" FROM "products" order by "price"`)
	}
	CheckError(err)

	defer rows.Close()
	products := make([]module.Product, 0, 100)
	for rows.Next() {
		err = rows.Scan(&id, &name, &price)
		cnt, avg := AvgProductRate(id)
		CheckError(err)
		product := module.Product{id, name, price, avg, cnt}
		products = append(products, product)
	}

	return products
}
func GetProduct(id string) module.Product {
	rows, err := Db.Query(`SELECT "id","name","price" FROM "products" where "id"=$1`, id)
	CheckError(err)
	var product module.Product
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id, &name, &price)
		CheckError(err)
		cnt, avg := AvgProductRate(id)
		product = module.Product{id, name, price, avg, cnt}
	}
	return product
}

func GetOrders(id string) []module.Order {
	rows, err := Db.Query(`SELECT "order_id","product_id","user_rate" from "orders" inner join "users" on "user_id"="users"."id" where "users"."id"=$1`, id)
	CheckError(err)
	defer rows.Close()
	orders := make([]module.Order, 0, 100)
	for rows.Next() {
		var user_rate int
		var order_id string
		err = rows.Scan(&order_id, &id, &user_rate)
		CheckError(err)
		order := module.Order{GetProduct(id), order_id, user_rate}
		orders = append(orders, order)
	}
	return orders
}

func SetRate(id string, rate int) {
	updateStmt := `update "orders" set "user_rate"=$1 where "order_id"=$2`
	_, e := Db.Exec(updateStmt, rate, id)
	CheckError(e)

}
func GetProductsByName(name string) []module.Product {
	rows, err := Db.Query(`SELECT "id","name","price" FROM "products" where "name"=$1`, name)
	CheckError(err)
	defer rows.Close()
	products := make([]module.Product, 0, 100)
	for rows.Next() {
		err = rows.Scan(&id, &name, &price)
		CheckError(err)
		cnt, avg := AvgProductRate(id)
		product := module.Product{id, name, price, avg, cnt}
		products = append(products, product)
	}

	return products
}

func GetUser(email string, pass int) (bool, module.User) {
	rows, err := Db.Query(`SELECT "id", "name", "email","balance" FROM "users" where "email"=$1 and "password"=$2`, email, pass)
	CheckError(err)

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id, &name, &email, &balance)
		CheckError(err)
		return true, module.User{id, name, email, balance}
	}
	return false, module.User{}
}

func IsEmailFree(email string) bool {
	rows, err := Db.Query(`SELECT "email" FROM "users" where "email"=$1`, email)
	CheckError(err)

	defer rows.Close()
	for rows.Next() {
		return false
	}
	return true
}
func BuyItem(user module.User, product module.Product) bool {
	user.Balance -= product.Price
	updateStmt := `update "users" set "balance"=$1 where "id"=$2`
	_, e := Db.Exec(updateStmt, user.Balance, user.Id)
	CheckError(e)
	order_id := uuid.New()
	insertDynStmt := `insert into "orders"("order_id","user_id", "product_id","user_rate") values($1, $2,$3,$4)`
	_, er := Db.Exec(insertDynStmt, order_id, user.Id, product.Id, 0)
	CheckError(er)

	return true
}
