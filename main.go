package main //Main code used to run the server
import (
	"fmt"
	"log"
	"net/http"
	"text/template"

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
	fmt.Fprintln(w, Dataa.Name)
}
