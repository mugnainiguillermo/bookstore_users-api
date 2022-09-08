package app

import (
	"github.com/gin-gonic/gin"
	"github.com/mugnainiguillermo/bookstore_users-api/src/logger"
)

var (
	router = gin.Default()
)

func StartApplication() {
	mapUrls()

	logger.Info("about to start the application...")

	router.SetTrustedProxies(nil)
	router.Run("localhost:9000")
}
