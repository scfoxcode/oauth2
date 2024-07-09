package main

import (
    "fmt"
    "os"
    "net/http"
    "sync"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

func loadConfig() error {
    deployType := os.Getenv("DEPLOY_TYPE")
    isProduction := deployType == "production"

    var envPath string;

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
    // Setup router here
    port := os.Getenv("SERVER_PORT")
    // router.Run(":" + port)
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
    go fun() {
        defer wg.Done()
        apiServer()
    }

    wg.Wait()
    select {}
}
