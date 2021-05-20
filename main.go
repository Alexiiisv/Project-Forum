package main //Main code used to run the server
import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	uuid "github.com/satori/go.uuid"

	config "./config"
)

type Data struct {
	Name string
}

var Dataa = Data{"alexis"}

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
	t := template.New("index-template")
	t = template.Must(t.ParseFiles("index.html"))
	t.ExecuteTemplate(w, "index", Dataa)
	uuid := GetUUID()
	fmt.Fprintln(w, uuid)
}

func GetUUID() uuid.UUID {
	// Creating UUID Version 4
	// panic on error
	u1 := uuid.Must(uuid.NewV4())
	fmt.Printf("UUIDv4: %s\n", u1)

	// or error handling
	u2, err := uuid.NewV4()
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
	}
	fmt.Printf("UUIDv4: %s\n", u2)

	// Parsing UUID from string input
	u2, err = uuid.FromString(u1.String())
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
	}
	fmt.Printf("Successfully parsed: %s\n", u2)
	return u2
}
