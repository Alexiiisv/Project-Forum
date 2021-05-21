package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	uuid "github.com/gofrs/uuid"

	hash "golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"

	config "github.com/Alexiiisv/Project-Forum/v2/config"
)

//basic struct
type DataSend struct {
	Name     string
	Password string
	Email    string
	Uuid     uuid.UUID
}

//const for hashing
const (
	MinCost     int = 4  // the minimum allowable cost as passed in to GenerateFromPassword
	MaxCost     int = 31 // the maximum allowable cost as passed in to GenerateFromPassword
	DefaultCost int = 10 // the cost that will actually be set if a cost below MinCost is passed into GenerateFromPassword
)

var Data DataSend
var Dataarray []DataSend
var firstinit bool

func main() {
	fmt.Println("Please connect to\u001b[31m localhost", config.LocalhostPort, "\u001b[0m")
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets")))) // Join Assets Directory to the server
	http.HandleFunc("/", index)
	http.HandleFunc("/accounts", ShowAccount)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/register", Register)
	http.HandleFunc("/log", CreateAccount)
	err := http.ListenAndServe(config.LocalhostPort, nil) // Set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//page index (main page)
func index(w http.ResponseWriter, r *http.Request) {
	t := template.New("index-template")
	t = template.Must(t.ParseFiles("index.html", "./tmpl/header&footer.html"))
	t.ExecuteTemplate(w, "index", Dataarray)
	// saveUuid()
}

//create account with value from page
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	NameChoosen := r.FormValue("Name")
	Password := r.FormValue("Password")
	Email := r.FormValue("Email")
	Dataarray = readUuid()
	Dataarray = append(Dataarray, DataSend{NameChoosen, Password, Email, GetUUID()})
	fmt.Println(string(HashPassword(Dataarray[len(Dataarray)-1].Password)))
	saveUuid("accounts")
	t := template.New("account-template")
	t = template.Must(t.ParseFiles("./tmpl/account.html", "./tmpl/header&footer.html"))
	t.ExecuteTemplate(w, "account", Dataarray)
}

//page to show all accounts existing
func ShowAccount(w http.ResponseWriter, r *http.Request) {
	Dataarray = readUuid()
	t := template.New("account-template")
	t = template.Must(t.ParseFiles("./tmpl/account.html", "./tmpl/header&footer.html"))
	t.ExecuteTemplate(w, "account", Dataarray)
}

//page login
func Login(w http.ResponseWriter, r *http.Request) {
	t := template.New("account-template")
	t = template.Must(t.ParseFiles("./tmpl/login&register.html", "./tmpl/header&footer.html"))
	t.ExecuteTemplate(w, "login", Dataarray)
}

//page login
func Register(w http.ResponseWriter, r *http.Request) {
	t := template.New("account-template")
	t = template.Must(t.ParseFiles("./tmpl/login&register.html", "./tmpl/header&footer.html"))
	t.ExecuteTemplate(w, "register", Dataarray)
}

//give a unique uuid to a user
func GetUUID() uuid.UUID {
	// Creating UUID Version 4
	// panic on error
	//var u1 uuid.UUID
	// if Dataa[len(Dataa)].Uuid == uuid.Nil {
	var u1 uuid.UUID = uuid.Must(uuid.NewV4())
	// } else {
	// 	u1 = Dataa[len(Dataa)].Uuid
	// }

	u2, err := uuid.FromString(u1.String())
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
	}
	return u2
}

//write in a database
func saveUuid(state string) {
	db, err := sql.Open("sqlite3", "./Database/User.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if state == "accounts" {
		stmt, err := db.Prepare("insert into Accounts(Name, Password, Email, Uuid) values(?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println((Dataarray))
		for index := range Dataarray {
			if UserExists(db, Dataarray[index].Uuid.String()) {
				continue
			}
			result, _ := stmt.Exec(Dataarray[index].Name, Dataarray[index].Password, Dataarray[index].Email, Dataarray[index].Uuid.String())
			fmt.Println("resultat ", result)
		}
	}
}

//read database/store value from database to go code
func readUuid() []DataSend {
	db, err := sql.Open("sqlite3", "./Database/User.db")
	if err != nil {
		log.Fatal(err)
	}
	sql_readall := `SELECT Name, Password, Email, Uuid FROM Accounts`

	rows, err := db.Query(sql_readall)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []DataSend
	for rows.Next() {
		rows.Scan(&Data.Name, &Data.Password, &Data.Email, &Data.Uuid)
		result = append(result, Data)
	}
	return result
}

/*
// //get the length of a table
func GetCount(schemadottablename string, db *sql.DB) int {
	var cnt int
	_ = db.QueryRow(`select count(*) from ` + schemadottablename).Scan(&cnt)
	return cnt
}
*/

//check if an account exist
func UserExists(db *sql.DB, uuid string) bool {
	sqlStmt := `SELECT Uuid FROM Accounts WHERE Uuid = ?`
	err := db.QueryRow(sqlStmt, uuid).Scan(&uuid)
	if err != nil {
		if err != sql.ErrNoRows {
			// a real error happened! you should change your function return
			// to "(bool, error)" and return "false, err" here
			log.Print(err)
		}
		fmt.Println("faux ", uuid)
		return false
	}
	fmt.Println("vrai ", uuid)
	return true
}

//crypt password
func HashPassword(passwd string) []byte {
	has, _ := hash.GenerateFromPassword([]byte(passwd), DefaultCost)
	return has
}

func CheckPasswordHash(password string, hashpass string) bool {
	var err error = hash.CompareHashAndPassword([]byte(hashpass), []byte(password))
	return err == nil
}
