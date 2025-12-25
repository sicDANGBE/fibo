package api

import (
	"fibo-orchestrateur/internal/orchestrator"

	"github.com/gin-gonic/gin"
)

func SetupRouter(orch *orchestrator.Engine, hub *Hub) *gin.Engine {
	r := gin.Default()

	// Chemins relatifs par rapport à la racine du projet où air est lancé
	r.Static("/js", "./web/js")
	r.LoadHTMLFiles("./web/index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/ws", func(c *gin.Context) {
		ServeWs(hub, c.Writer, c.Request)
	})

	r.POST("/run", func(c *gin.Context) {
		var req struct {
			Handler string                 `json:"handler"`
			Params  map[string]interface{} `json:"params"`
		}
		if err := c.BindJSON(&req); err == nil {
			orch.StartTask(req.Handler, req.Params)
			c.Status(200)
		}
	})

	return r
}
