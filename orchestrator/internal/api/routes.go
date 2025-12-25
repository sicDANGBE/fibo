package api

import (
	"embed"
	"fibo-orchestrateur/internal/orchestrator"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed all:web
var webAssets embed.FS

func SetupRouter(orch *orchestrator.Engine, hub *Hub) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	// Endpoint utilisé par le healthcheck du Compose
	r.GET("/health", func(c *gin.Context) {
		orch.Mu.Lock()
		// Vérification de l'état du canal RabbitMQ
		isRMQConnected := orch.Channel != nil && !orch.Conn.IsClosed()
		orch.Mu.Unlock()

		if isRMQConnected {
			c.JSON(http.StatusOK, gin.H{"status": "UP"})
		} else {
			// Retourne 503 pour que curl -f échoue
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "DOWN"})
		}
	})

	webRoot, _ := fs.Sub(webAssets, "web")
	r.GET("/", func(c *gin.Context) {
		index, _ := fs.ReadFile(webRoot, "index.html")
		c.Data(http.StatusOK, "text/html; charset=utf-8", index)
	})
	r.StaticFS("/static", http.FS(webRoot))
	r.GET("/ws", func(c *gin.Context) {
		ServeWs(hub, c.Writer, c.Request)

		// On envoie immédiatement la liste des workers déjà connus
		orch.Mu.Lock()
		for _, w := range orch.Workers {
			hub.BroadcastMessage(map[string]interface{}{
				"type": "WORKER_JOIN",
				"data": w,
			})
		}
		orch.Mu.Unlock()
	})
	r.POST("/run", func(c *gin.Context) {
		var req struct {
			Handler string                 `json:"handler"`
			Params  map[string]interface{} `json:"params"`
		}
		if err := c.BindJSON(&req); err == nil {
			orch.StartTask(req.Handler, req.Params)
			c.Status(http.StatusAccepted)
		}
	})
	return r
}
