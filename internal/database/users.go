package database

import (
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"password,omitempty"`
	Refresh  string `json:"refresh_token,omitempty"`
	Red      bool   `json:"is_chirpy_red"`
}

var ErrUserNotFound = errors.New("User not found")

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
	pwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error generating password hash: %s", err)
		return User{}, err
	}
	user := User{
		ID:       nextID,
		Email:    email,
		Password: pwd,
		Red:      false,
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

func (db *DB) GetUserByToken(refreshToken string) (User, bool) {
	data := db.loadDB()
	for _, v := range data.Users {
		if v.Refresh == "" {
			continue
		}
		if v.Refresh == refreshToken {
			return v, true
		}
	}
	return User{}, false
}

func (db *DB) UpdateUser(id int, newUser User) (User, error) {
	data := db.loadDB()
	_, ok := data.Users[id]
	if !ok {
		return User{}, ErrUserNotFound
	}
	data.Users[id] = newUser
	db.writeDB(data)
	return newUser, nil
}

func (db *DB) UpgradeUser(id int) error {
	data := db.loadDB()
	user, ok := data.Users[id]
	if !ok {
		return ErrUserNotFound
	}
	data.Users[id] = User{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
		Refresh:  user.Refresh,
		Red:      true,
	}
    db.writeDB(data)
	return nil
}
