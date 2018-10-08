package database

type Staff struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Gender       string `json:"gender"`
	Birthday     string `json:"date_of_birth"`
	HealthCard   string `json:"health_card_number"`
	SIN          string `json:"SIN"`
	University   string `json:"university"`
	HomeAddress  string `json:"home_address"`
	Email        string `json:"email"`
	EmploymentID string `json:"employment_id"`
	JobRole      string `json:"job_role"`
	Pay          string `json:"pay"`
	Manager      string `json:"manager"`
}
