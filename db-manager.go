package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"database/sql"
	"errors"
	"log"

	databaseModel "github.com/dsych/go-server/models/database"
)

const database = "a1"
const userTable = "users"
const staffTable = "staff_data"
const accessTable = "system_access_data"

type DBManager struct {
	Username string
	Password string
	Host     string
	Database string

	db          *sql.DB
	isConnected bool
}

func (m *DBManager) Connect() error {
	if len(m.Username) == 0 || len(m.Password) == 0 || len(m.Host) == 0 || len(m.Database) == 0 {
		panic("Connection credentials are not provided. GO_USERNAME: '" + m.Username + "', GO_PASSWORD: '" + m.Password + "', GO_HOST: '" + m.Host + "', GO_DATABASE: '" + m.Database + "'.")
	}
	db, err := sql.Open("mysql", m.Username+":"+m.Password+"@tcp("+m.Host+")/"+m.Database)

	if err == nil {
		m.db = db
	}

	return err
}

func (m *DBManager) searchAccess(access databaseModel.AccessData) ([]databaseModel.AccessData, error) {
	query := "select " +
		"employee_id, " +
		"username, " +
		"password, " +
		"computer_asset_number, " +
		"static_ip_address, " +
		"MAC_address, " +
		"access_level" +
		" from " + accessTable + " where employee_id = ? or username = ? or access_level = ?"
	rows, err := m.db.Query(
		query,
		access.EmployeeID, access.Username, access.AccessLevel)

	defer rows.Close()

	if err != nil {
		log.Println(err)
		return nil, errors.New("Unable to retrieve records")
	}

	rc := make([]databaseModel.AccessData, 0)

	for rows.Next() {
		var row databaseModel.AccessData
		if err := rows.Scan(&row.EmployeeID, &row.Username, &row.Password, &row.ComputerNumber, &row.StaticIP, &row.MACAddress, &row.AccessLevel); err != nil {
			log.Println(err)
		}
		rc = append(rc, row)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, errors.New("Failed reading records")
	}

	return rc, nil
}

func (m *DBManager) searchStaff(access databaseModel.Staff) ([]databaseModel.Staff, error) {

	query := "select " +
		"first_name, " +
		"last_name, " +
		"gender, " +
		"date_of_birth, " +
		"health_card_number, " +
		"SIN, " +
		"university, " +
		"home_address, " +
		"email, " +
		"employment_id, " +
		"job_role, " +
		"pay, " +
		"manager " +
		" from " + staffTable + " where first_name = ? or last_name = ? or employment_id = ? or manager like ?"
	rows, err := m.db.Query(
		query,
		access.FirstName, access.LastName, access.EmploymentID, access.Manager)

	if err != nil {
		log.Println(err)
		return nil, errors.New("Unable to retrieve records")
	}
	defer rows.Close()

	rc := make([]databaseModel.Staff, 0)

	for rows.Next() {
		var row databaseModel.Staff
		if err := rows.Scan(&row.FirstName,
			&row.LastName,
			&row.Gender,
			&row.Birthday,
			&row.HealthCard,
			&row.SIN,
			&row.University,
			&row.HomeAddress,
			&row.Email,
			&row.EmploymentID,
			&row.JobRole,
			&row.Pay,
			&row.Manager); err != nil {
			log.Println(err)
		}
		rc = append(rc, row)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, errors.New("Failed reading records")
	}

	return rc, nil
}

func (m *DBManager) Authenticate(user databaseModel.User) error {
	invalidError := errors.New("Invalid credentials provided")

	dbUser := user

	row := m.db.QueryRow("select * from "+userTable+" where username = ?", user.Username)
	if err := row.Scan(&dbUser.Username, &dbUser.Password, &dbUser.Salt); err != nil {
		return invalidError
	}

	user.Salt = dbUser.Salt

	err := m.generateSaltedPassword(&user)

	if err != nil {
		return err
	} else if !bytes.Equal(user.Password, dbUser.Password) {
		return invalidError
	} else {
		return nil
	}
}

func (m *DBManager) Register(user databaseModel.User) error {

	// make sure that salt is empty
	user.Salt = nil
	err := m.generateSaltedPassword(&user)

	if err != nil {
		return err
	}

	res, err := m.db.Exec("insert into "+userTable+" values(?,?,?)", user.Username, user.Password, user.Salt)

	if err != nil {
		return err
	}

	if affected, err := res.RowsAffected(); err != nil {
		return err
	} else if affected <= 0 {
		return errors.New("Failed to insert records")
	} else {
		return nil
	}
}

// returns error or nil
func (m *DBManager) generateSaltedPassword(user *databaseModel.User) error {
	// generate salt
	tmp := make([]byte, 10)
	if _, err := rand.Read(tmp); err != nil {
		return errors.New("Unable to generate salt")

	}

	var salt []byte

	// if salt is not present, generate it.
	// if present, just use it
	if user.Salt == nil || len(user.Salt) == 0 {
		a := sha256.Sum256(tmp)
		salt = a[:]
	} else {
		salt = user.Salt
	}

	hasher := sha512.New()
	hasher.Write(append(user.Password, salt...))
	hashedPassword := hasher.Sum(nil)

	user.Password = hashedPassword
	user.Salt = salt

	return nil
}

func (m *DBManager) CloseConnection() error {
	return m.db.Close()
}
