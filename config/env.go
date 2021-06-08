package config

import uuid "github.com/gofrs/uuid"

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
	CreationDate string
	Creator string
	Pic bool //to do
	Like     int
	Liker    string
	Disliker    string
	Liked    int
}

type TContent struct {
	Id   int
	Uuid string
	Text string
	Picture string
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

type Erreur struct {
	Connected bool
	Miss bool
}