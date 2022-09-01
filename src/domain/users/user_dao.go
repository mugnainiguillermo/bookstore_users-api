package users

import (
	"fmt"
	"mugnainiguillermo/bookstore_users-api/src/datasources/mysql/users_db"
	"mugnainiguillermo/bookstore_users-api/src/utils/errors"
)

const (
	queryInsertUser = "INSERT INTO users(first_name, last_name, email, date_created) VALUES(?, ?, ?, ?);"
)

var (
	usersDB = make(map[int64]*User)
)

func (user *User) Get() *errors.RestErr {
	if err := users_db.Client.Ping(); err != nil {
		// TODO: Proper error handling
		panic(err)
	}

	result := usersDB[user.Id]

	if result == nil {
		return errors.NewNotFoundError(fmt.Sprintf("user %d not found", user.Id))
	}

	user.Id = result.Id
	user.FirstName = result.FirstName
	user.LastName = result.LastName
	user.Email = result.Email
	user.DateCreated = result.DateCreated

	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare("SELECT * FROM users;")
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close() // Defer Closing statement as soon as we create it to avoid an open statement.

	insertResult, err := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated)
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user: %s", err.Error()))
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to save user: %s", err.Error()))
	}
	user.Id = userId
	/*
		current := usersDB[user.Id]

		if current != nil {
			if current.Email == user.Email {
				return errors.NewBadRequestError(fmt.Sprintf("email %s already exists", user.Email))
			}
			return errors.NewBadRequestError(fmt.Sprintf("user %d already exists", user.Id))
		}

		user.DateCreated = time.Now().UTC().Format(time.RFC3339)
		usersDB[user.Id] = user
	*/
	return nil
}
