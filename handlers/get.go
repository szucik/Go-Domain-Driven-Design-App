package handlers

import (
	"github.com/szucik/go-simple-rest-api/data"
	"net/http"
)

// swagger:route GET /users users listUsers
// Return a list of users from the database
// responses:
//	200: usersResponse

// GetUsers returns Users from data store
func (u *Users) GetUsers(rw http.ResponseWriter, r *http.Request) {
	u.l.Println("Handle GET Users")
	rw.Header().Add("Content-Type", "application/json")
	lu := data.GetUsers()
	err := lu.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}
