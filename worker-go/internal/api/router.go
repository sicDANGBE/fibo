package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "alive"})
	})

	// Futur handler pour les tests de transfert de fichiers
	r.POST("/io-test", func(c *gin.Context) {
		// Logique de réception de fichier pour benchmark
		c.Status(http.StatusAccepted)
	})

	return r
}
