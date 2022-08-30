package app

import (
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	mapUrls()
	router.SetTrustedProxies(nil)
	router.Run(":8080")
}
