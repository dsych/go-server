package requests

import (
	"github.com/dsych/go-server/models/database"
)

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (src *UserRequest) ToDBModel() database.User {
	return database.User{Username: src.Username, Password: []byte(src.Password)}
}
