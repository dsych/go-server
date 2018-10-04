package main

type User struct {
	Username string `json:"username"`
	Password []byte `json:"password"`
	Salt     []byte `json:"salt"`
}

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
	Manager      string `json:"manager "`
}

type AccessData struct {
	EmployeeID     int    `json:"employee_id"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	ComputerNumber string `json:"computer_asset_number"`
	StaticIP       string `json:"static_ip_address"`
	MACAddress     string `json:"MAC_address"`
	AccessLevel    int    `json:"access_level"`
}

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
