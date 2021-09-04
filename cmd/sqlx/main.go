package main

import (
	"fmt"
	"reflect"
	"strings"

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
	var record []interface{}
	var user User
	var n int
	var phs []string
	for rows.Next() {

		//rows.Scanの代わりにrows.StructScanを使う
		err := rows.StructScan(&user)
		if err != nil {
			log.Fatal(err)
		}
		rv := reflect.ValueOf(user)
		tp := rv.Type()
		for i := 0; i < tp.NumField(); i++ {
			vf := rv.Field(i)
			record = append(record, vf.Interface())
		}
		phs = append(phs, fmt.Sprintf("($%d, $%d, $%d)", n+1, n+2, n+3))
		n = 3 + n
		userlist = append(userlist, user)
	}

	//fmt.Println(userlist)
	//fmt.Println(record)
	//fmt.Println(userlist)
	/*
		query := `INSERT INTO users (id, name, age) VALUES (:id, :name, :age)`
		tx := db.MustBegin()
		_, err = tx.NamedExec(query, userlist)
		if err != nil {
			fmt.Println("エラーだよ")
			tx.Rollback()
		}
		tx.Commit()

	*/
	var list []interface{}
	list = append(list, 1)
	list = append(list, "AZ")
	list = append(list, 10)
	fmt.Println(list)

	tx := db.MustBegin()
	phsStr := strings.Join(phs, ",")
	//phsStr = `(cast($1 as int),$2,cast($3 as int))`
	query := fmt.Sprintf("insert into users select id::int, name, age::int from (values%s) as v(id, name, age)", phsStr)
	fmt.Println(query)
	fmt.Println(record)
	stmt, err := tx.Prepare(query)
	if err != nil {
		fmt.Println(err)
		return
	}
	result, err := stmt.Exec(record...)
	fmt.Println(result)
	fmt.Println(err)

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
