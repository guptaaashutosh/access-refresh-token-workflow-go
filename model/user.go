package model

// type User struct {
// 	Id int64
// 	Name string
// 	Age int64
// 	Address string
// }

// type Login struct {
// 	Id int64
// 	Name string
// }

type User struct {
	Id       int64
	Email    string
	Password string
	Name     string
	Age      int64
	Address  string
	Sid      []int64
}
type GetUser struct {
	Id      int64
	Email   string
	Name    string    `json:",omitempty"` //Go's encoding/json package to control how fields are encoded in JSON
	Age     int64    `json:",omitempty"`
	Address string    `json:",omitempty"`
	Service []Service `json:",omitempty"`
}

type Login struct {
	Email    string
	Password string
}

type Service struct {
	Sid     int64
	Service string
}

