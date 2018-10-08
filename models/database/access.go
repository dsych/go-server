package database

type AccessData struct {
	EmployeeID     int    `json:"employee_id"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	ComputerNumber string `json:"computer_asset_number"`
	StaticIP       string `json:"static_ip_address"`
	MACAddress     string `json:"MAC_address"`
	AccessLevel    int    `json:"access_level"`
}
