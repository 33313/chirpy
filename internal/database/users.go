package database

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

func (db *DB) GetUser(id int) (User, bool) {
	data := db.loadDB()
	user, ok := data.Users[id]
	return user, ok
}

func (db *DB) GetUsers() []User {
	data := db.loadDB()
	var users []User
	for _, v := range data.Users {
		users = append(users, v)
	}
	return users
}

func (db *DB) CreateUser(email string, password string) (User, error) {
	data := db.loadDB()
	if _, ok := db.GetUserByEmail(email); ok {
		return User{}, errors.New("User already exists.")
	}
	nextID := len(data.Users) + 1
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Fatalf("Error generating password hash: %s", err)
		return User{}, err
	}
	user := User{
		ID:       nextID,
		Email:    email,
		Password: pwd,
	}
	data.Users[nextID] = user
	db.writeDB(data)
	return user, nil
}

// Returns a User struct and a bool which states whether such user exists.
func (db *DB) GetUserByEmail(email string) (User, bool) {
	data := db.loadDB()
	for _, v := range data.Users {
		if v.Email == email {
			return v, true
		}
	}
	return User{}, false
}
