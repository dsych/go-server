package database

type User struct {
	Username string `json:"username"`
	Password []byte `json:"password"`
	Salt     []byte `json:"salt"`
}
