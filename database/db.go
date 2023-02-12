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
		var name string
		var price float64
		var id string
		err = rows.Scan(&id, &name, &price)
		CheckError(err)
		product := module.Product{id, name, price}
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
		var name string
		var price float64
		var id string
		err = rows.Scan(&id, &name, &price)
		CheckError(err)
		product = module.Product{id, name, price}
	}
	return product
}
func GetProductsByName(name string) []module.Product {
	rows, err := Db.Query(`SELECT "id","name","price" FROM "products" where "name"=$1`, name)
	CheckError(err)
	defer rows.Close()
	products := make([]module.Product, 0, 100)
	for rows.Next() {
		var name string
		var price float64
		var id string
		err = rows.Scan(&id, &name, &price)
		CheckError(err)
		product := module.Product{id, name, price}
		products = append(products, product)
	}

	return products
}

func GetUser(email string, pass int) (bool, module.User) {
	rows, err := Db.Query(`SELECT "id", "name", "email","balance" FROM "users" where "email"=$1 and "password"=$2`, email, pass)
	CheckError(err)

	defer rows.Close()
	for rows.Next() {
		var id string
		var name string
		var email string
		var balance float64
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
	insertDynStmt := `insert into "orders"("order_id","user_id", "product_id") values($1, $2,$3)`
	_, er := Db.Exec(insertDynStmt, order_id, user.Id, product.Id)
	CheckError(er)

	return true
}
