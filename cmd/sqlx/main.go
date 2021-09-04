package main

import (
	"fmt"

	//"database/sql"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"os"
	"path/filepath"
)

type config struct {
	Dsn    string `json:"dsn"`
	Suffix string `json:suffix`
}

var exePath string
var suffix string
var dsn string
var db *sqlx.DB

type User struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
	Age  int    `db:"age"`
}
type Userlist []User

func main() {
	exe, _ := os.Executable()
	exePath = filepath.Dir(exe)
	setEnv()
	DbConnection(dsn)
	var userlist Userlist
	rows, err := db.Queryx("SELECT id, name, age FROM users")
	if err != nil {
		log.Fatal(err)
	}

	var user User
	for rows.Next() {

		//rows.Scanの代わりにrows.StructScanを使う
		err := rows.StructScan(&user)
		if err != nil {
			log.Fatal(err)
		}
		userlist = append(userlist, user)
	}

	fmt.Println(userlist)
	query := `INSERT INTO users (id, name, age) VALUES (:id, :name, :age)`
	tx := db.MustBegin()
	_, err = tx.NamedExec(query, userlist)
	if err != nil {
		fmt.Println("エラーだよ")
		tx.Rollback()
	}
	tx.Commit()

}
func setEnv() {
	fname := filepath.Join(exePath, "..", "..", "..", "config", "config.json")
	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var cfg config
	err = json.NewDecoder(f).Decode(&cfg)
	dsn = cfg.Dsn
	suffix = cfg.Suffix
}

func DbConnection(dsn string) {

	d, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Println(err)
	}
	db = d
}
