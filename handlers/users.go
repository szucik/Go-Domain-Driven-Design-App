package handlers

import (
	"github.com/zlyjoker102/simple-rest-api/data"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Users struct {
	l *log.Logger
}

func (u *Users) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		u.getUsers(rw, r)
		return
	}
	if r.Method == http.MethodPost {
		u.addUser(rw, r)
		return
	}
	if r.Method == http.MethodPut {
		rgx := regexp.MustCompile(`/([0-9]+)`)
		g := rgx.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		if len(g[0]) != 2 {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			u.l.Println("Invalid URI unable to convert to number", idString)
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		u.updateUser(id, rw, r)
		return
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (u *Users) updateUser(id int,rw http.ResponseWriter, r * http.Request) {
	u.l.Println("Handle PUT Users")
	user := &data.User{}
	err := user.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusInternalServerError)
	}

	err = data.UpdateUser(id,user)
	if err == data.ErrUserNotFound {
		http.Error(rw, "User not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "User not found", http.StatusInternalServerError)
		return
	}
}

func (u *Users) addUser(rw http.ResponseWriter, r * http.Request) {
	u.l.Println("Handle POST Users")
	usr := &data.User{}
	err := usr.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusInternalServerError)
	}
	data.AddUser(usr)
}

func (u *Users) getUsers(rw http.ResponseWriter, r * http.Request) {
	u.l.Println("Handle GET Users")
	lu := data.GetUsers()
	err := lu.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func NewUser(l *log.Logger) *Users {
	return &Users{l}
}
