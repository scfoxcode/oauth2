package routes;

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func health(ctx *gin.Context) {
    ctx.String(http.StatusOK, "Server Online")
} 

func Health(routes *gin.RouterGroup) {
    routes.GET("/", health)
}
