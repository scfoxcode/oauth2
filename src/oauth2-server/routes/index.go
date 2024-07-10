package routes

import (
    "fmt"
    "time"
    "github.com/gin-gonic/gin"
    "github.com/scfoxcode/oauth2/src/oauth2-server/models/counters"
)

func SetupRoutes() *gin.Engine {
    fmt.Println("Initialising Endpoints")

    // Track login attempts
    var loginCounter counters.Counters
    loginCounter.Init()
    invalidLoginTicker := time.NewTicker(5 * time.Minute)
    go func() {
        for {
            select {
                case <- invalidLoginTicker.C:
                    loginCounter.Init()
            }
        }
    }()


    router := gin.Default()

    // Unauthed routes at the top
    HealthRoutes(router.Group("/health"))
    LoginRoutes(router.Group("/user"), &loginCounter) // Questionable path here

    // TODO Re use the auth middle ware and rate limitting safety I created for my accounting software

    // TODO Authed routes below here

    return router
}
