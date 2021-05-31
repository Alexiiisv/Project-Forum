package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"text/template"
	"time"

	config "github.com/Alexiiisv/Project-Forum/v2/config"
)

var Data config.Account
var TName config.TName
var TContent config.TContent
var TopicsName config.Topics
var Dataarray config.AllAccount
var allTopics config.AllTopics
var Logged config.LoginYes
var Name, Password, TopicText, UUID, SetTopicsName, SetTopicsDescription, info, Likes, Category string
var IdTopics int
var CategoryName = []string{"Informatique", "Jeux Video", "Musique", "Design", "Communication", "Animation3D", "NSFW", "Anime", "Manga"}
var pp_name string

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
	http.HandleFunc("/like", Like)
	// http.HandleFunc("/Like", Like)
	http.HandleFunc("/upload_pp", uploadHandler)
	err := http.ListenAndServe(config.LocalhostPort, nil) // Set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//page account information
func Info(w http.ResponseWriter, r *http.Request) {
	t := template.New("account-template")
	t = template.Must(t.ParseFiles("./tmpl/account.html", "./tmpl/header&footer.html"))
	t.ExecuteTemplate(w, "account", Logged)
}

//page index
func index(w http.ResponseWriter, r *http.Request) {
	t := template.New("index-template")
	t = template.Must(t.ParseFiles("index.html", "./tmpl/header&footer.html"))
	t.ExecuteTemplate(w, "index", Logged)
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
	Dataarray.Account = Logged.Account
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
		// fmt.Println("Name : ", SetTopicsName, " Description : ", SetTopicsDescription)
		Category = GetCategory(r)
		SetTopicInfo(state)
	}
	allTopics.Name = readtopics()
	allTopics.Connected = Logged.Connected
	// fmt.Println(allTopics.Name)
	t := template.New("topics-template")
	t = template.Must(t.ParseFiles("./tmpl/topics.html", "./tmpl/header&footer.html", "./tmpl/content.html"))
	t.ExecuteTemplate(w, "topics", allTopics)
}

//page accounts
func CreateTopicInfo(w http.ResponseWriter, r *http.Request) {
	allTopics.Name = readtopics()
	allTopics.Connected = Logged.Connected
	t := template.New("topics-template")
	t = template.Must(t.ParseFiles("./tmpl/topics.html", "./tmpl/header&footer.html", "./tmpl/content.html"))
	t.ExecuteTemplate(w, "CreateTopicInfo", allTopics)
}

//page accounts
func singleTopics(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("State")
	IdTopics, _ = strconv.Atoi(r.FormValue("IdTopics"))
	if state == "PostTopic" {
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

// Display User Informations
func User_Info(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	UUID = r.FormValue("Uuid")
	Dataarray.Data = readuuid(state)
	Dataarray.Connected = Logged.Connected
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

// Connection to account
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

// TODO: PUSH into UUID Cringe ?
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 32 MB is the default used by FormFile
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get a reference to the fileHeaders
	files := r.MultipartForm.File["AddPP"]

	for _, fileHeader := range files {
		if fileHeader.Size > config.MAX_UPLOAD_SIZE {
			http.Error(w, fmt.Sprintf("The uploaded image is too big: %s. Please use an image less than 1MB in size", fileHeader.Filename), http.StatusBadRequest)
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer file.Close()

		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		filetype := http.DetectContentType(buff)
		if filetype != "image/jpeg" && filetype != "image/png" {
			http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG image", http.StatusBadRequest)
			return
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = os.MkdirAll("./assets/image/Account_pp", os.ModePerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pp_name = strconv.FormatInt(time.Now().UnixNano(), 10)
		f, err := os.Create(fmt.Sprintf("./assets/image/Account_pp/%s%s", pp_name, filepath.Ext(fileHeader.Filename)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer f.Close()

		pr := &config.Progress{
			TotalSize: fileHeader.Size,
		}

		_, err = io.Copy(f, io.TeeReader(file, pr))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	state := "pp"
	saveUuid(state)
	Info(w, r)
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
	sql_readall := `SELECT Name, Password, Email, Uuid, Profile_Picture FROM Accounts`

	rows, err := db.Query(sql_readall)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []config.Account
	for rows.Next() {
		rows.Scan(&Data.Name, &Data.Password, &Data.Email, &Data.Uuid, &Data.Profile_Picture)
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
	sql_readall := `SELECT Id, Title, Description, Category, Like FROM Topics_Name;`

	rows, err := db.Query(sql_readall)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result []config.TName
	for rows.Next() {
		rows.Scan(&TName.Id, &TName.Title, &TName.Desc, &TName.Category, &TName.Likes)
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
	sql_readall := `SELECT Id, Title, Description, Category, Like FROM Topics_Name;`

	rows, err := db.Query(sql_readall)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var result config.TName
	for rows.Next() {
		rows.Scan(&TName.Id, &TName.Title, &TName.Desc, &TName.Category, &TName.Likes)
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

//TODO: Why don't push into DB ?
//write in a database
func saveUuid(state string) {
	db, err := sql.Open("sqlite3", "./Database/User.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if state == "accounts" {
		stmt, err := db.Prepare("insert into Accounts(Name, Password, Email, Uuid, Profile_Picture) values(?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		for index := range Dataarray.Data {
			if config.UserExists(db, Dataarray.Data[index].Uuid.String()) {
				continue
			}
			stmt.Exec(Dataarray.Data[index].Name, string(config.HashPassword(Dataarray.Data[len(Dataarray.Data)-1].Password)), Dataarray.Data[index].Email, Dataarray.Data[index].Uuid.String(), "Standard_Pic.png")
		}
	} else if state == "pp" {
		if err != nil {
			panic(err)
		}
		stmt, err := db.Prepare("update Accounts set Profile_Picture = ? where Uuid = ?")
		if err != nil {
			log.Fatal(err)
		}
		if Logged.Account.Profile_Picture != "Standard_Pic.png" {
			toremove := "./assets/image/Account_pp/"
			toremove += Logged.Account.Profile_Picture
			os.Remove(toremove)
		}

		link := pp_name
		link += ".png"
		Logged.Account.Profile_Picture = link
		stmt.Exec(link, Logged.Account.Uuid.String())
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

func Like(w http.ResponseWriter, r *http.Request) {

	Likes = r.FormValue("Like")
	fmt.Println(Likes)
	t := template.New("like-template")
	t = template.Must(t.ParseFiles("./tmpl/topics.html", "./tmpl/header&footer.html", "./tmpl/content.html"))
	t.ExecuteTemplate(w, "singleTopics", TopicsName)
}
