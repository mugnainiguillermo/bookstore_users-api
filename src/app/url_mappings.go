package app

import (
	"mugnainiguillermo/bookstore_users-api/src/controllers/ping"
	"mugnainiguillermo/bookstore_users-api/src/controllers/users"
)

func mapUrls() {
	router.GET("/ping", ping.Ping)

	router.GET("/users/:user_id", users.GetUser)
	router.POST("/users", users.CreateUser)
}
