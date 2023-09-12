package main

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/apache/skywalking-go"
)

func main() {
	engine := gin.New()
	engine.Handle("GET", "/consumer", func(context *gin.Context) {
		resp, err := http.Get("http://localhost:8080/provider")
		if err != nil {
			log.Print(err)
			context.Status(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
			context.Status(http.StatusInternalServerError)
			return
		}
		context.String(200, string(body))
	})

	engine.Handle("GET", "/provider", func(context *gin.Context) {
		context.String(200, "success")
	})

	engine.Handle("GET", "health", func(context *gin.Context) {
		context.Status(http.StatusOK)
	})

	_ = engine.Run(":9999")
}
