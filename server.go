package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"

	"github.com/gorilla/sessions"
)

var (
	keyPass    = "./keys/server.key"
	store      *sessions.CookieStore
	cookieName = "auth"
	authValue  = "authValue"
	users      = map[string]string{}
)

func helloServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Example response from https server!"))
}

func main() {

	users["user1"] = "1234"
	users["abc"] = "1234"
	users["123"] = "1234"

	store = initSessionStore()

	router := mux.NewRouter()
	router.Schemes("https")
	router.HandleFunc("/hello", helloServer)
	router.HandleFunc("/login", login)
	subrouter := router.PathPrefix("/content").Subrouter()

	subrouter.HandleFunc("/logout", logout).Methods("GET")
	subrouter.HandleFunc("/hello", helloServer)
	subrouter.Use(authMiddleware)

	err := http.ListenAndServeTLS("localhost:1443", "./keys/server.crt", "./keys/server.key", context.ClearHandler(router))

	if err != nil {
		log.Fatal("Unable to serve: ", err)
	}
}

func initSessionStore() *sessions.CookieStore {
	key, err := ioutil.ReadFile(keyPass)

	if err != nil {
		panic(err)
	}

	return sessions.NewCookieStore(key)
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if auth, ok := isAuthenticated(r); auth != nil && !auth.(bool) {
			http.Error(w, "Forbidden", http.StatusForbidden)
		} else if !ok {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func isAuthenticated(r *http.Request) (interface{}, bool) {
	session, err := store.Get(r, cookieName)
	if err != nil {
		log.Fatal("Session is not available")
		return false, false
	}
	return session.Values[authValue], true
}

func login(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, cookieName)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	user, okU := r.URL.Query()["user"]
	password, okP := r.URL.Query()["password"]

	if (!okU || len(user[0]) < 1) || (!okP || len(password[0]) < 1) {
		log.Println("Url Param 'user' is missing")
		return
	}

	if pass, found := users[user[0]]; !found || password[0] != string(pass) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	session.Values[authValue] = true
	session.Save(r, w)
	log.Println("Logged in")

}

func logout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, cookieName)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	session.Values[authValue] = false
	session.Save(r, w)
	log.Println("Logged out")
}
