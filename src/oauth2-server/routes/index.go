package routes

import (
    "fmt"
    "github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
    fmt.Println("Initialising Endpoints")

    router := gin.Default()

    // Unauthed routes at the top
    Health(router.Group("/health"))

    // TODO Re use the auth middle ware and rate limitting safety I created for my accounting software

    // TODO Authed routes below here

    return router
}
