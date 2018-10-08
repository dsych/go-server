package requests

import (
	"github.com/dsych/go-server/models/database"
)

type Employee struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Manager      string `json:"manager"`
	Gender       string `json:"gender"`
	Birthday     string `json:"DOB"`
	HealthCard   string `json:"healthCard"`
	SIN          string `json:"SIN"`
	University   string `json:"university"`
	HomeAddress  string `json:"homeAddress"`
	Email        string `json:"email"`
	EmploymentID string `json:"employmentId"`
	JobRole      string `json:"jobRole"`
	Pay          string `json:"pay"`
}

func (src *Employee) ToDBModel() database.Staff {
	return database.Staff{
		FirstName:    src.FirstName,
		EmploymentID: src.EmploymentID,
		LastName:     src.LastName,
		Manager:      src.Manager}
}

func (dest Employee) PopulateFromDB(src database.Staff) Employee {
	return Employee{
		FirstName:    src.FirstName,
		LastName:     src.LastName,
		Manager:      src.Manager,
		Gender:       src.Gender,
		Birthday:     src.Birthday,
		HealthCard:   src.HealthCard,
		SIN:          src.SIN,
		University:   src.University,
		HomeAddress:  src.HomeAddress,
		Email:        src.Email,
		EmploymentID: src.EmploymentID,
		JobRole:      src.JobRole,
		Pay:          src.Pay}
}
