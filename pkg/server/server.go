package server

import (
	"Desktop/golangProjects/CRUD/pkg"
	"Desktop/golangProjects/CRUD/pkg/server/database"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

type Server struct {
	rlock    *sync.RWMutex
	database database.Database
}

func New(db database.Database) Server {
	rlock := &sync.RWMutex{}
	return Server{
		rlock:    rlock,
		database: db,
	}
}

// HandleCreate serves the route /users/create
func (s *Server) HandleCreate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost, http.MethodPut:
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			log.Print("unsupported content type")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("error in reading the http body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var user pkg.HttpData
		err = json.Unmarshal(body, &user)
		if err != nil {
			log.Printf("error in unmarshalling the http body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Print("creating user")
		err = s.database.Write(user.Name, strconv.Itoa(user.Age))
		if err != nil {
			log.Printf("error in writing to the database: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// HandleUser handles the route /users/{name}
func (s *Server) HandleUsers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]
	switch r.Method {
	case http.MethodGet:
		s.rlock.RLock()
		ageAsString, err := s.database.Read(name)
		s.rlock.RUnlock()
		if err != nil {
			log.Printf("we don't have anyone by the name of %s in our database", name)
			// w.WriteHeader(http.StatusBadRequest)
			w.Header().Add("Content-Type", "text/plain; charset=us-ascii")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(fmt.Appendf(make([]byte, 0), "%s does not exist in our database", name))
			return
		}
		age, err := strconv.Atoi(ageAsString)
		if err != nil {
			log.Printf("error converting age to int: %v", err)
		}
		user := pkg.HttpData{
			Name: name,
			Age:  age,
		}
		body, err := json.Marshal(user)
		if err != nil {
			log.Printf("error marshaling the body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(body)
	case http.MethodPatch:
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			log.Print("unsupported content type")
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("error reading json body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var user pkg.HttpData
		err = json.Unmarshal(body, &user)
		if err != nil {
			log.Printf("error unmarshalling json body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		s.rlock.Lock()
		defer s.rlock.Unlock()
		s.database.Write(name, strconv.Itoa(user.Age))
		log.Printf("sucesfully updated %s's age to %d", name, user.Age)
		w.WriteHeader(http.StatusOK)
	case http.MethodDelete:
		s.rlock.Lock()
		defer s.rlock.Unlock()
		err := s.database.Delete(name)
		if err != nil {
			log.Printf("error deleting user %s: %v", name, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("succesfully deleted %s from the database \n", name)
		w.WriteHeader(http.StatusOK)
	}
	return
}
