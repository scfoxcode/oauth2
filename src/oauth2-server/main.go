package main

import (
	"fmt"
	"os"
	"sync"
	"github.com/joho/godotenv"
	"github.com/scfoxcode/oauth2/src/oauth2-server/routes"
)

func loadConfig() error {
    deployType := os.Getenv("DEPLOY_TYPE")
    isProduction := deployType == "production"

    var envPath string

    if isProduction {
        envPath = "config/.live.env"
    } else {
        envPath = "config/.debug.env"
    }

    err := godotenv.Load(envPath)

    if err != nil {
        fmt.Printf("Failed to load config for deploy type %s\n", envPath)
        return err
    }

    return nil
}

func apiServer() {
    router := routes.SetupRoutes()

    port := os.Getenv("SERVER_PORT")
    err := router.Run(":" + port)

    if err != nil {
        fmt.Printf("FATAL: Failed to initialise router on port %s\n", port)
        panic(err)
    }

    fmt.Printf("Server running on port %s\n", port)
}

func main() {
    fmt.Println("Initialising Server...")

    err := loadConfig()
    if err != nil {
        fmt.Println(err)
        panic("Fatal error, failed to load configuration. Terminating")
    }

    // Start the servers
    var wg sync.WaitGroup

    wg.Add(1)
    go func() {
        defer wg.Done()
        apiServer()
    }()

    wg.Wait()
    // select {}
}
