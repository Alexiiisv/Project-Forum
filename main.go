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
var Dataarray []Data
var firstinit bool

func main() {

	fmt.Println("Please connect to\u001b[31m localhost", config.LocalhostPort, "\u001b[0m")
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets")))) // Join Assets Directory to the server
	http.HandleFunc("/", troll1)
	http.HandleFunc("/lel", troll2)
	err := http.ListenAndServe(config.LocalhostPort, nil) // Set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func troll1(w http.ResponseWriter, r *http.Request) {
	t := template.New("index-template")
	t = template.Must(t.ParseFiles("index.html"))
	t.ExecuteTemplate(w, "index", Dataarray)
	// saveUuid()
}


func troll2(w http.ResponseWriter, r *http.Request) {
	NameChoosen := r.FormValue("Name")
	Dataarray = append(Dataarray, Data{NameChoosen, GetUUID()})
	t := template.New("account-template")
	t = template.Must(t.ParseFiles("./tmpl/account.html"))
	t.ExecuteTemplate(w, "account", Dataarray)
	saveUuid("Account")
}

//give a unique uuid to a user
func GetUUID() uuid.UUID {
	// Creating UUID Version 4
	// panic on error
	var u1 uuid.UUID
	// if Dataa[len(Dataa)].Uuid == uuid.Nil {
		u1 = uuid.Must(uuid.NewV4())
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

	if state == "Account" {
		stmt, _ := db.Prepare("insert into Account(id, name, uuid) values(?, ?, ?)")
		for index := range Dataarray {
			if UserExists(db, Dataarray[index].Uuid.String()) {
				continue
			}
			_, _ = stmt.Exec(GetCount("Account", db), Dataarray[index].Name, Dataarray[index].Uuid.String())
		}
	}
}

func GetCount(schemadottablename string, db *sql.DB) int {
    var cnt int
    _ = db.QueryRow(`select count(*) from ` + schemadottablename).Scan(&cnt)
    return cnt 
}

func UserExists(db * sql.DB, uuid string) bool {
    sqlStmt := `SELECT uuid FROM Account WHERE uuid = ?`
    err := db.QueryRow(sqlStmt, uuid).Scan(&uuid)
    if err != nil {
        if err != sql.ErrNoRows {
            // a real error happened! you should change your function return
            // to "(bool, error)" and return "false, err" here
            log.Print(err)
        }

        return false
    }

    return true
}