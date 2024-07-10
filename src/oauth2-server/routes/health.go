package routes;

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func health(ctx *gin.Context) {
    ctx.String(http.StatusOK, "Server Online")
} 

func HealthRoutes(routes *gin.RouterGroup) {
    routes.GET("/", health)
}
