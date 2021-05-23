package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	config "github.com/Alexiiisv/Project-Forum/v2/config"
)

var Data config.DataSend
var Dataarray config.AllAccount
var Logged config.LoginYes
var Name, Password string

func main() {
	fmt.Println("Please connect to\u001b[31m localhost", config.LocalhostPort, "\u001b[0m")
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets")))) // Join Assets Directory to the server
	http.HandleFunc("/", index)
	http.HandleFunc("/accounts", ShowAccount)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/logout", Logout)
	http.HandleFunc("/information", Info)
	http.HandleFunc("/register", Register)
	http.HandleFunc("/createAcc", CreateAccount)
	http.HandleFunc("/connect", LoggedOn)
	err := http.ListenAndServe(config.LocalhostPort, nil) // Set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//page index
func index(w http.ResponseWriter, r *http.Request) {
	Dataarray.Connected = Logged.Connected
	t := template.New("index-template")
	t = template.Must(t.ParseFiles("index.html", "./tmpl/header&footer.html"))
	t.ExecuteTemplate(w, "index", Dataarray)
}

//create account
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	NameChoosen := r.FormValue("Name")
	Password := r.FormValue("Password")
	Email := r.FormValue("Email")
	Dataarray.Data = append(Dataarray.Data, config.DataSend{Name: NameChoosen, Password: Password, Email: Email, Uuid: config.GetUUID()})
	saveUuid("accounts")
	ShowAccount(w, r)
}

//page accounts
func ShowAccount(w http.ResponseWriter, r *http.Request) {
	Dataarray.Data = readuuid("ShowAccount")
	Dataarray.Connected = Logged.Connected
	t := template.New("account-template")
	t = template.Must(t.ParseFiles("./tmpl/account.html", "./tmpl/header&footer.html"))
	t.ExecuteTemplate(w, "accounts", Dataarray)
}

//page account information
func Info(w http.ResponseWriter, r *http.Request) {
	t := template.New("account-template")
	t = template.Must(t.ParseFiles("./tmpl/account.html", "./tmpl/header&footer.html"))
	t.ExecuteTemplate(w, "account", Logged)
}

//page login
func Login(w http.ResponseWriter, r *http.Request) {
	t := template.New("account-template")
	t = template.Must(t.ParseFiles("./tmpl/login&register.html", "./tmpl/header&footer.html"))
	t.ExecuteTemplate(w, "login", Logged)
}

//page register
func Register(w http.ResponseWriter, r *http.Request) {
	Dataarray.Connected = Logged.Connected
	t := template.New("account-template")
	t = template.Must(t.ParseFiles("./tmpl/login&register.html", "./tmpl/header&footer.html"))
	t.ExecuteTemplate(w, "register", Dataarray)
}

func LoggedOn(w http.ResponseWriter, r *http.Request) {
	Name = r.FormValue("Name")
	Password = r.FormValue("Password")
	Dataarray.Data = readuuid("LoggedOn")
	if len(Dataarray.Data) == 1 {
		Logged.Account.Data = Dataarray.Data[0]
		Logged.Connected = true
		index(w, r)
	}else {
		Login(w, r)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	Logged.Account = config.Account{}
	Logged.Connected = false
	Login(w, r)
}

/*
// //get the length of a table
func GetCount(schemadottablename string, db *sql.DB) int {
	var cnt int
	_ = db.QueryRow(`select count(*) from ` + schemadottablename).Scan(&cnt)
	return cnt
}
*/

//read database/store value from database to go code
func readuuid(state string) []config.DataSend {
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

	var result []config.DataSend
	for rows.Next() {
		rows.Scan(&Data.Name, &Data.Password, &Data.Email, &Data.Uuid)
		if state == "LoggedOn" && config.CheckPasswordHash(Password, Data.Password) && Data.Name == Name {
			Data.Password = Password
			result = append(result, Data)
			break
		} else if state == "ShowAccount" {
			result = append(result, Data)
		}
	}
	return result
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
		for index := range Dataarray.Data {
			if config.UserExists(db, Dataarray.Data[index].Uuid.String()) {
				continue
			}
			stmt.Exec(Dataarray.Data[index].Name, string(config.HashPassword(Dataarray.Data[len(Dataarray.Data)-1].Password)), Dataarray.Data[index].Email, Dataarray.Data[index].Uuid.String())
		}
	}
}
