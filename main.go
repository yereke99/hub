package main

import (
	"hub/room"

	"github.com/gin-gonic/gin"
)

var htmlRepo string = "index.html"

func main() {
	go room.H.Run()

	r := gin.New()
	r.LoadHTMLFiles(htmlRepo)

	r.GET("/room/:roomId", func(c *gin.Context) {
		c.HTML(200, htmlRepo, nil)
	})

	r.GET("/ws/:roomId", func(c *gin.Context) {
		roomId := c.Param("roomId")
		room.ServerWS(c.Writer, c.Request, roomId)
	})

	r.Run(":8080")
}
