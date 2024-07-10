package user 

import (
    "fmt"
    "errors"
    "golang.org/x/crypto/bcrypt"
    "github.com/google/uuid"
    "github.com/scfoxcode/oauth2/src/oauth2-server/models/auth"
)

type LoginProps struct {
    Username string `json:"username" form:"username"`
    Password string `json:"password" form:"password"` 
}

func LeakNothingError(hiddenError string) error {
    fmt.Println(hiddenError);
    return errors.New("Invalid Login")
}

func HashPassword(password string) (string, error) {
    hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        println(err)
        return "", LeakNothingError(fmt.Sprintf("Failed to hash password"))
    }
    return string(hashedPw), nil
}

func ComparePasswordAndStoredHash(hash string, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// This is in the wrong module
func AttemptLogin(props LoginProps) (auth.IDToken, error) {
    var token auth.IDToken
    // Because we don't db for this yet. Not worth the dev cost
    users := make(map[string]string)

    // This is dogshit
    users["stephen"] = "$2a$10$kalTzHyCUZjNuwCLdtMaU.I1K6AQYl4ushivZVK1qNwrnm4yNSH06"
    users["heritage"] = "$2a$10$XsgDlv/GO0hXP7zY3neyRelVFXwoA.eEakcGnzozH9woy3IQTd32i"

    storedHash, ok := users[props.Username]
    if !ok {
        return token, LeakNothingError(fmt.Sprintf("User \"%s\" not found", props.Username))
    }

    // We need to actually hash the prrovided passowrd here and test
    err := ComparePasswordAndStoredHash(storedHash, props.Password)
    if err != nil {
        return token, LeakNothingError(fmt.Sprintf("Password doesn't match for User \"%s\"", props.Username))
    }

    creation := auth.IDTokenCreationData {
        Uid: uuid.New().String(), 
        DisplayName: props.Username,
        Email: props.Username,
    }

    token = auth.CreateIDToken(&creation)

    return token, nil;
}

