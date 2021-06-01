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
	Name            string
	Password        string
	Email           string
	Uuid            uuid.UUID
	Profile_Picture string
}

type AllAccount struct {
	Connected bool
	Data      []Account
	Account   Account
}

type LoginYes struct {
	Connected bool
	Account   Account
}

type TName struct {
	Id       int
	Title    string
	Desc     string
	Category string
	Likes    int
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
	Account   Account
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

const MAX_UPLOAD_SIZE = 1024 * 1024 // 1MB

// Progress is used to track the progress of a file upload.
// It implements the io.Writer interface so it can be passed
// to an io.TeeReader()
type Progress struct {
	TotalSize int64
	BytesRead int64
}

// Write is used to satisfy the io.Writer interface.
// Instead of writing somewhere, it simply aggregates
// the total bytes on each read
func (pr *Progress) Write(p []byte) (n int, err error) {
	n, err = len(p), nil
	pr.BytesRead += int64(n)
	//pr.Print()
	return
}
