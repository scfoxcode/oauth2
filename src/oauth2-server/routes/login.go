package routes;

import (
    "os"
    "fmt"
    "net/http"
    "github.com/gin-gonic/gin"
    // Okay, go is weird
    // This is the path alright, but counters is NOT a sub module
    // All modules are equal (flat) the directory structure is not relevant
    // A module is defined by it's directory path final folder
    "github.com/scfoxcode/oauth2/src/oauth2-server/models/user"
    "github.com/scfoxcode/oauth2/src/oauth2-server/models/counters"
    "github.com/scfoxcode/oauth2/src/oauth2-server/models/auth"
)

// Models as a scope sucks, why can't it be user

func loginEndpoint (invalidLoginCounter *counters.Counters) gin.HandlerFunc {
    return func (ctx *gin.Context) {
        body := user.LoginProps{}

        if err := ctx.ShouldBind(&body); err != nil {
            ctx.JSON(http.StatusBadRequest, gin.H{"errors": fmt.Sprintf("%v", err)}) 
            return
        }

        count := invalidLoginCounter.GetCount(body.Username)
        if count > 10 {
            fmt.Printf("Too many failed logins for user: %s\n", body.Username)
            ctx.String(http.StatusBadRequest, "Too many failed logins. Please come back later")
            return
        }

        token, err := user.AttemptLogin(body)
        if err != nil {
            invalidLoginCounter.Increment(body.Username)
            ctx.Header("HX-Redirect", "/login") // Incase htmx swallows native
            ctx.Redirect(http.StatusFound, "/login")
            return
        }

        secret := os.Getenv("JWT_SECRET")
        signedToken := auth.SignIDToken(&token, secret)


        domain := os.Getenv("COOKIE_DOMAIN")
        ctx.SetCookie("token", signedToken, 86400, "/", domain, false, true)
        ctx.Header("HX-Redirect", "/accounts") // Incase htmx swallows native
        ctx.Redirect(http.StatusFound, "/accounts")
    }
}

func LoginRoutes(routes *gin.RouterGroup, invalidLoginCounter *counters.Counters) {
    routes.POST("/login", loginEndpoint(invalidLoginCounter))
}
