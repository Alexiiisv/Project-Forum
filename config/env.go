package config

import (
	"database/sql"
	"fmt"
	"log"

	uuid "github.com/gofrs/uuid"

	_ "github.com/mattn/go-sqlite3"
	hash "golang.org/x/crypto/bcrypt"
)

//const for hashing
const (
	MinCost     int = 4  // the minimum allowable cost as passed in to GenerateFromPassword
	MaxCost     int = 31 // the maximum allowable cost as passed in to GenerateFromPassword
	DefaultCost int = 10 // the cost that will actually be set if a cost below MinCost is passed into GenerateFromPassword
)

// Initiliazing the port
const LocalhostPort = ":8080"

//basic struct
type Account struct {
	Name     string
	Password string
	Email    string
	Uuid     uuid.UUID
}

type AllAccount struct {
	Connected bool
	Data      []Account
}

type LoginYes struct {
	Connected bool
	Account   Account
}

type TName struct {
	Id    int
	Title string
	Desc  string
	Category  string
	Likes int
}

type TContent struct {
	Id   int
	Uuid string
	Text string
}

//single topics
type Topics struct {
	Name     TName
	Content  []TContent
	Accounts []Account
	Account  LoginYes
}

//all topics
type AllTopics struct {
	Name      []TName
	Connected bool
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
