package database

import "log"

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (db *DB) GetUser(id int) (User, bool) {
	data, err := db.loadDB()
	if err != nil {
		log.Fatal(err)
	}

	chirp, ok := data.Users[id]
	return chirp, ok
}

func (db *DB) GetUsers() []User {
	data, err := db.loadDB()
	if err != nil {
		log.Fatal(err)
	}

	var users []User
	for _, v := range data.Users {
		users = append(users, v)
	}

	return users
}

func (db *DB) CreateUser(body string) (User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		log.Fatalf("Error creating user: %s", err)
		return User{}, err
	}
	nextID := len(dbStruct.Users) + 1
	user := User{
		ID:    nextID,
		Email: body,
	}
	dbStruct.Users[nextID] = user
	db.writeDB(dbStruct)
	return user, nil
}
