package requests

import (
	"github.com/dsych/go-server/models/database"
)

type Access struct {
	EmployeeID     int    `json:"employeeId"`
	Username       string `json:"username"`
	AccessLevel    int    `json:"accessLvl"`
	Password       string `json:"password"`
	ComputerNumber string `json:"computerAccess"`
	StaticIP       string `json:"IP"`
	MACAddress     string `json:"MAC"`
}

func (a *Access) ToDBModel() database.AccessData {
	return database.AccessData{Username: a.Username, AccessLevel: a.AccessLevel, EmployeeID: a.EmployeeID}
}

func (dest Access) PopulateFromDB(src database.AccessData) Access {
	return Access{
		EmployeeID:     src.EmployeeID,
		Username:       src.Username,
		AccessLevel:    src.AccessLevel,
		Password:       src.Password,
		ComputerNumber: src.ComputerNumber,
		StaticIP:       src.StaticIP,
		MACAddress:     src.MACAddress}
}
