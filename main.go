package main //Main code used to run the server
import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	uuid "github.com/satori/go.uuid"

	_ "github.com/mattn/go-sqlite3"

	config "./config"
)

type Data struct {
	Name string
	Uuid uuid.UUID
}

var Dataa Data
var firstinit bool

func main() {

	fmt.Println("Please connect to\u001b[31m localhost", config.LocalhostPort, "\u001b[0m")
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets")))) // Join Assets Directory to the server
	http.HandleFunc("/", troll)
	err := http.ListenAndServe(config.LocalhostPort, nil) // Set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func troll(w http.ResponseWriter, r *http.Request) {
	Dataa = Data{"Alexis", GetUUID()}
	t := template.New("index-template")
	t = template.Must(t.ParseFiles("index.html"))
	t.ExecuteTemplate(w, "index", Dataa)
	WriteInDatabase()
}

func GetUUID() uuid.UUID {
	// Creating UUID Version 4
	// panic on error
	var u1 uuid.UUID
	if Dataa.Uuid == uuid.Nil {
		u1 = uuid.Must(uuid.NewV4())
	} else {
		u1 = Dataa.Uuid
	}
	fmt.Printf("UUIDv4: %s\n", u1)

	// Parsing UUID from string input
	u2, err := uuid.FromString(u1.String())
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
	}
	fmt.Printf("Successfully parsed: %s\n", u2)
	return u2
}

func WriteInDatabase() {

	db, err := sql.Open("sqlite3", "./Database/User.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var sqlStmt string

	if firstinit {
		sqlStmt = `
			drop table foo;
			create table foo (id integer not null primary key, name text, uuid text);
			delete from foo;
		`
	} else {
		sqlStmt = `
			delete from foo;
		`
	}
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
	stmt, _ := db.Prepare("insert into foo(id, name, uuid) values(?, ?, ?)")
	_, _ = stmt.Exec(1, Dataa.Name, Dataa.Uuid.String())
}
