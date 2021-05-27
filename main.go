package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	config "github.com/Alexiiisv/Project-Forum/v2/config"
)

var Data config.Account
var TName config.TName
var TContent config.TContent
var TopicsName config.Topics
var Dataarray config.AllAccount
var allTopics config.AllTopics
var Logged config.LoginYes
var Name, Password, TopicText, UUID, SetTopicsName, SetTopicsDescription, info, Category string
var IdTopics int
var CategoryName = []string{"Informatique", "Jeux Video", "Musique", "Design", "Communication", "Animation3D", "NSFW", "Anime", "Manga"}

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
	http.HandleFunc("/topics", AllTopics)
	http.HandleFunc("/singleTopics", singleTopics)
	http.HandleFunc("/user_account", User_Info)
	http.HandleFunc("/CreateTopicInfo", CreateTopicInfo)
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
	Dataarray.Data = append(Dataarray.Data, config.Account{Name: NameChoosen, Password: Password, Email: Email, Uuid: config.GetUUID()})
	saveUuid("accounts")
	ShowAccount(w, r)
}

//page accounts
func ShowAccount(w http.ResponseWriter, r *http.Request) {
	Dataarray.Data = readuuid("ShowAccount")
	Dataarray.Connected = Logged.Connected
	t := template.New("account-template")
	t = template.Must(t.ParseFiles("./tmpl/account.html", "./tmpl/header&footer.html", "./tmpl/content.html"))
	t.ExecuteTemplate(w, "accounts", Dataarray)
}

//page accounts
func AllTopics(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("State")
	if state == "CreateTopicInfo" {
		SetTopicsName = r.FormValue("Title")
		SetTopicsDescription = r.FormValue("Description")
		fmt.Println("Name : ", SetTopicsName, " Description : ", SetTopicsDescription)
		Category = GetCategory(r)
		SetTopicInfo(state)
	}
	allTopics.Name = readtopics()
	allTopics.Connected = Logged.Connected
	fmt.Println(allTopics)
	t := template.New("topics-template")
	t = template.Must(t.ParseFiles("./tmpl/topics.html", "./tmpl/header&footer.html", "./tmpl/content.html"))
	t.ExecuteTemplate(w, "topics", allTopics)
}

//page accounts
func CreateTopicInfo(w http.ResponseWriter, r *http.Request) {
	allTopics.Name = readtopics()
	allTopics.Connected = Logged.Connected
	fmt.Println(allTopics)
	t := template.New("topics-template")
	t = template.Must(t.ParseFiles("./tmpl/topics.html", "./tmpl/header&footer.html", "./tmpl/content.html"))
	t.ExecuteTemplate(w, "CreateTopicInfo", allTopics)
}

//page accounts
func singleTopics(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("State")
	if state == "SingleTopic" {
		IdTopics, _ = strconv.Atoi(r.FormValue("IdTopics"))
	} else if state == "PostTopic" {
		IdTopics, _ = strconv.Atoi(r.FormValue("IdTopics"))
		TopicText = r.FormValue("text")
		SetTopicText("PostTopicText")
	}
	TopicsName.Name = GetTopicsData()
	TopicsName.Account = Logged
	TopicsName.Content = GetTopicsContent()
	TopicsName.Accounts = readuuid("ShowAccount")
	t := template.New("singleTopics-template")
	t = template.Must(t.ParseFiles("./tmpl/topics.html", "./tmpl/header&footer.html", "./tmpl/content.html"))
	t.ExecuteTemplate(w, "singleTopics", TopicsName)
}

//page account information
func Info(w http.ResponseWriter, r *http.Request) {
	t := template.New("account-template")
	t = template.Must(t.ParseFiles("./tmpl/account.html", "./tmpl/header&footer.html"))
	t.ExecuteTemplate(w, "account", Logged)
}

// Display User Informations
func User_Info(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	UUID = r.FormValue("Uuid")
	Dataarray.Data = readuuid(state)
	Dataarray.Connected = Logged.Connected
	fmt.Println(r.FormValue("Uuid"))
	t := template.New("account-template")
	t = template.Must(t.ParseFiles("./tmpl/account.html", "./tmpl/header&footer.html"))
	t.ExecuteTemplate(w, "user_account", Dataarray)
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
		Logged.Account = Dataarray.Data[0]
		Logged.Connected = true
		index(w, r)
	} else {
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
func readuuid(state string) []config.Account {
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

	var result []config.Account
	for rows.Next() {
		rows.Scan(&Data.Name, &Data.Password, &Data.Email, &Data.Uuid)
		if state == "LoggedOn" && config.CheckPasswordHash(Password, Data.Password) && Data.Name == Name {
			Data.Password = Password
			result = append(result, Data)
			break
		} else if state == "ShowAccount" {
			result = append(result, Data)
		} else if state == "user_account" && UUID == Data.Uuid.String() {
			result = append(result, Data)
			break
		}
	}
	return result
}

//read database/store value from database to go code
func readtopics() []config.TName {
	db, err := sql.Open("sqlite3", "./Database/Topics.db")
	if err != nil {
		log.Fatal(err)
	}
	sql_readall := `SELECT Id, Title, Description, Category FROM Topics_Name;`

	rows, err := db.Query(sql_readall)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []config.TName
	for rows.Next() {
		rows.Scan(&TName.Id, &TName.Title, &TName.Desc, &TName.Category)
		result = append(result, TName)
	}
	return result
}

//Get from a database the information about a topic
func GetTopicsData() config.TName {
	db, err := sql.Open("sqlite3", "./Database/Topics.db")
	if err != nil {
		log.Fatal(err)
	}
	sql_readall := `SELECT Id, Title, Description, Likes FROM Topics_Name;`

	rows, err := db.Query(sql_readall)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result config.TName
	for rows.Next() {
		rows.Scan(&TName.Id, &TName.Title, &TName.Desc)
		if TName.Id == IdTopics {
			result = TName
			break
		}
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

//read database/store value from database to go code
func GetTopicsContent() []config.TContent {
	db, err := sql.Open("sqlite3", "./Database/Topics.db")
	if err != nil {
		log.Fatal(err)
	}
	sql_readall := `SELECT Id, Uuid, Text FROM Topics`

	rows, err := db.Query(sql_readall)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []config.TContent
	Compte := readuuid("ShowAccount")
	for rows.Next() {
		rows.Scan(&TContent.Id, &TContent.Uuid, &TContent.Text)
		for i := 0; i < len(Compte); i++ {
			if TContent.Uuid == Compte[i].Uuid.String() {
				TContent.Uuid = Compte[i].Name
			}
		}
		if TContent.Id == IdTopics {
			result = append(result, TContent)
		} else {
			continue
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

//Set in a database the text written in a topic
func SetTopicText(state string) {
	db, err := sql.Open("sqlite3", "./Database/Topics.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if state == "PostTopicText" {
		stmt, err := db.Prepare("insert into Topics(Id, Uuid, Text) values(?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		stmt.Exec(IdTopics, Logged.Account.Uuid.String(), TopicText)
	}
}

//Set in a database the information about a topic
func SetTopicInfo(state string) {
	db, err := sql.Open("sqlite3", "./Database/Topics.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if state == "CreateTopicInfo" {
		stmt, err := db.Prepare("insert into Topics_Name(Title, Description, Category) values(?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		stmt.Exec(SetTopicsName, SetTopicsDescription, Category)
	}
}

//
func GetCategory(r *http.Request) string {
	var result = ""
	for _, v := range CategoryName {
		a := r.FormValue(v)
		b, _ := strconv.ParseBool(a)
		if b {
			result += v
			result += ","
		}
	}
	return result[:len(result)-1]
}

