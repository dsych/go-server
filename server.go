package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"github.com/dsych/go-server/models/requests"
	_ "github.com/go-sql-driver/mysql"
)

var (
	keyPass          = "./keys/session.key"
	store            *sessions.CookieStore
	cookieName       = "auth"
	authValue        = "authValue"
	users            = map[string]string{}
	expirationPeriod = 5
	db               = DBManager{Username: os.Getenv("GO_USERNAME"), Password: os.Getenv("GO_PASSWORD"), Host: os.Getenv("GO_HOST"), Database: os.Getenv("GO_DATABASE")}
)

func main() {
	if err := db.Connect(); err != nil {
		log.Fatal("Unable to connect to database")
	}

	defer db.CloseConnection()

	store = initSessionStore()

	router := mux.NewRouter()
	router.Schemes("https")
	router.HandleFunc("/api/login", login).Methods("POST")
	router.HandleFunc("/api/content/searchAccess", searchAccess).Methods("POST")
	router.HandleFunc("/api/content/searchStaff", searchStaff).Methods("POST")
	router.HandleFunc("/api/register", register)
	router.HandleFunc("/api/logout", logout).Methods("GET")
	fs := FileSystem{fs: http.Dir("./public"), readDirBatchSize: 2}
	router.PathPrefix("/").Handler(noCacheMiddleware(http.FileServer(fs)))
	router.Use(authMiddleware)

	err := http.ListenAndServe("localhost:1444", context.ClearHandler(router))

	if err != nil {
		log.Fatal("Unable to serve: ", err)
	}
}

func initSessionStore() *sessions.CookieStore {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	key, err := ioutil.ReadFile(path.Join(dir, keyPass))

	if err != nil {
		panic(err)
	}

	return sessions.NewCookieStore(key)
}

func noCacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
		w.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
		w.Header().Set("Expires", "0")
		next.ServeHTTP(w, r) // Proxies.
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//check if this path needs to authenticated
		pathsEnforced := []string{"/content"}
		shouldValidate := false
		for _, path := range pathsEnforced {
			if strings.Contains(r.URL.EscapedPath(), path) {
				shouldValidate = true
				break
			}
		}

		// no need for validation
		if !shouldValidate {
			next.ServeHTTP(w, r)
			return
		}

		if auth, ok := isAuthenticated(r); auth == nil || auth.IsNew || (auth.Values[authValue] != nil && !auth.Values[authValue].(bool)) {
			http.Error(w, "Forbidden", http.StatusForbidden)
		} else if !ok {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		} else {
			setSession(auth, expirationPeriod)
			auth.Save(r, w)

			next.ServeHTTP(w, r)
		}
	})
}

func isAuthenticated(r *http.Request) (*sessions.Session, bool) {
	session, err := store.Get(r, cookieName)
	if err != nil {
		log.Println(err, "Client's cookie is out of date")
		return nil, false
	}

	return session, true
}

func login(w http.ResponseWriter, r *http.Request) {
	//expire previous session
	deleteSession(w, r)

	session, err := store.Get(r, cookieName)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	decoder := json.NewDecoder(r.Body)
	var body requests.UserRequest
	err = decoder.Decode(&body)

	if err != nil || len(body.Username) < 1 || len(body.Password) < 1 {
		res := "Missing username or password"
		log.Println(res)
		http.Error(w, res, http.StatusInternalServerError)
		return
	}

	if err := db.Authenticate(body.ToDBModel()); err != nil {
		log.Println(err)
		http.Error(w, "Invalid Credentials", http.StatusForbidden)
		return
	}

	setSession(session, expirationPeriod)

	session.Values[authValue] = true
	session.Save(r, w)
	log.Println("Logged in")

}

func deleteSession(w http.ResponseWriter, r *http.Request) error {
	session, err := store.Get(r, cookieName)

	// encryption keys were updated, so we want to re-create the session
	if err != nil && !session.IsNew {
		log.Println(err)
		return errors.New("Internal Server Error")
	}

	setSession(session, -1)

	session.Values[authValue] = false
	err = session.Save(r, w)
	return nil
}

func setSession(s *sessions.Session, minutes int) {
	s.Options.HttpOnly = true
	s.Options.MaxAge = minutes * 60
}

func logout(w http.ResponseWriter, r *http.Request) {
	if err := deleteSession(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	log.Println("Logged out")
}

func register(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var body requests.UserRequest
	err := decoder.Decode(&body)

	if err != nil || len(body.Username) < 1 || len(body.Password) < 1 {
		res := "Missing username or password"
		log.Println(res)
		http.Error(w, res, http.StatusBadRequest)
		return
	}

	if err := db.Register(body.ToDBModel()); err != nil {
		log.Println(err)
		http.Error(w, "Unable to register", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Registered"))
}

func searchAccess(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var body requests.Access
	err := decoder.Decode(&body)

	if err != nil || (body.AccessLevel == 0 && len(body.Username) < 1 && body.EmployeeID == 0) {
		res := "Empty request received"
		log.Println(res)
		http.Error(w, res, http.StatusBadRequest)
		return
	}

	ac, err := db.searchAccess(body.ToDBModel())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rc := make([]requests.Access, 0)

	for _, a := range ac {
		rc = append(rc, requests.Access.PopulateFromDB(requests.Access{}, a))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rc)
}

func searchStaff(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var body requests.Employee
	err := decoder.Decode(&body)

	if err != nil || (len(body.FirstName) < 1 && len(body.LastName) < 1 && len(body.EmploymentID) < 1 && len(body.Manager) < 1) {
		res := "Empty request received"
		log.Println(res)
		http.Error(w, res, http.StatusBadRequest)
		return
	}

	ac, err := db.searchStaff(body.ToDBModel())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rc := make([]requests.Employee, 0)

	for _, a := range ac {
		rc = append(rc, requests.Employee.PopulateFromDB(requests.Employee{}, a))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rc)
}
