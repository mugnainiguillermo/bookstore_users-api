package users

import (
	"errors"
	"fmt"
	"github.com/mugnainiguillermo/bookstore_users-api/src/datasources/mysql/users_db"
	"github.com/mugnainiguillermo/bookstore_users-api/src/logger"
	"github.com/mugnainiguillermo/bookstore_utils-go/date_utils"
	"github.com/mugnainiguillermo/bookstore_utils-go/mysql_utils"
	"github.com/mugnainiguillermo/bookstore_utils-go/rest_errors"
	"strings"
)

const (
	queryGetUser              = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id=?;"
	queryInsertUser           = "INSERT INTO users(first_name, last_name, email, date_created, password, status) VALUES(?, ?, ?, ?, ?, ?);"
	queryUpdateUser           = "UPDATE users SET first_name=?, last_name=?, email=? WHERE id=?;"
	queryDeleteUser           = "DELETE FROM users WHERE id=?;"
	queryFindUserByStatus     = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?;"
	queryFindEmailAndPassword = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE email=? AND password=? AND status=?;"
)

func (user *User) Get() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error while trying to prepare get user statement", err)
		return rest_errors.NewInternalServerError("error while trying to get user", errors.New("database error"))
	}
	defer stmt.Close()

	user.DateCreated = date_utils.GetNowString()

	result := stmt.QueryRow(user.Id)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		logger.Error("error while trying to get user by id", getErr)
		return rest_errors.NewInternalServerError("error while trying to get user", errors.New("database error"))
	}

	return nil
}

func (user *User) Save() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error while trying to prepare insert user statement", err)
		return rest_errors.NewInternalServerError("error while trying to save user", errors.New("database error"))
	}
	defer stmt.Close()

	result, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Password, user.Status)
	if saveErr != nil {
		logger.Error("error while trying to insert user", saveErr)
		return rest_errors.NewInternalServerError("error while trying to save user", errors.New("database error"))
	}

	userId, err := result.LastInsertId()
	if err != nil {
		logger.Error("error while trying to get last insert id", err)
		return rest_errors.NewInternalServerError("error while trying to save user", errors.New("database error"))
	}

	user.Id = userId

	return nil
}

func (user *User) Update() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error while trying to prepare update user statement", err)
		return rest_errors.NewInternalServerError("error while trying to update user", errors.New("database error"))
	}
	defer stmt.Close()

	_, updateErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if updateErr != nil {
		logger.Error("error while trying to update user", err)
		return rest_errors.NewInternalServerError("error while trying to update user", errors.New("database error"))
	}

	return nil
}

func (user *User) Delete() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error while trying to prepare delete user statement", err)
		return rest_errors.NewInternalServerError("error while trying to delete user", errors.New("database error"))
	}
	defer stmt.Close()

	if _, deleteErr := stmt.Exec(user.Id); deleteErr != nil {
		logger.Error("error while trying to delete user", err)
		return rest_errors.NewInternalServerError("error while trying to delete user", errors.New("database error"))
	}

	return nil
}

func (user *User) FindByStatus(status string) ([]User, *rest_errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		logger.Error("error while trying to prepare find users statement", err)
		return nil, rest_errors.NewInternalServerError("error while trying to find user", errors.New("database error"))
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error while trying to find users", err)
		return nil, rest_errors.NewInternalServerError("error while trying to find user", errors.New("database error"))
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("error while trying to find users", err)
			return nil, rest_errors.NewInternalServerError("error while trying to find user", errors.New("database error"))
		}
		results = append(results, user)
	}

	if len(results) == 0 {
		return nil, rest_errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}

	return results, nil
}

func (user *User) FindByEmailAndPassword() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryFindEmailAndPassword)
	if err != nil {
		logger.Error("error while trying to prepare get user by email and password statement", err)
		rest_errors.NewInternalServerError("error while trying to get user by email and password", errors.New("database error"))
	}
	defer stmt.Close()

	user.DateCreated = date_utils.GetNowString()

	result := stmt.QueryRow(user.Email, user.Password, StatusActive)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		if strings.Contains(getErr.Error(), mysql_utils.ErrorNoRows) {
			return rest_errors.NewNotFoundError("invalid user credentials")
		}
		logger.Error("error while trying to get user by email and password", getErr)
		return rest_errors.NewInternalServerError("error while trying to get user by email and password", errors.New("database error"))
	}

	return nil
}
