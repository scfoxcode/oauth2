package auth

import (
    "crypto/sha256"
    "encoding/base64"
    "encoding/json"
    "regexp"
    "fmt"
    "time"
)

// Note IDToken not access token. Different purposes in oauth2 standard
type IDTokenHeader struct {
    Alg string `json:"alg"`
    Typ string `json:"type"`
}

type IDTokenPayload struct {
    Iss string `json:"iss"`
    Sub string `json:"sub"` // this can be user uid
    Aud string `json:"aud"`
    Nonce string `json:"nonce"`
    DisplayName string `json:"display_name"`
    Email string `json:"Email"`
    AuthTime time.Time `json:"auth_time"`
    Iat time.Time `json:"iat"`
    Exp time.Time `json:"exp"`
}

type IDToken struct {
    Header IDTokenHeader
    Payload IDTokenPayload
}

type IDTokenCreationData struct {
    Uid string
    DisplayName string
    Email string
}

func PayloadFromSignedToken(token string) (IDTokenPayload, error) {
    var payload IDTokenPayload
    tokenShape := `^([^\.]+)\.([^\.]+)\.([^\.]+)$`
    re := regexp.MustCompile(tokenShape)
    matches := re.FindStringSubmatch(token)
    match := matches[2]
    dehashedPayload, err := base64.StdEncoding.DecodeString(match)
    if err != nil {
        println("Failed to decode payload")
        return payload, err
    }
    err = json.Unmarshal(dehashedPayload, &payload)
    if err != nil {
        println("Failed to unmarshal payload")
        return payload, err
    }

    return payload, nil
}

func CreateIDToken(data *IDTokenCreationData) IDToken {
    token := IDToken {
        Header: IDTokenHeader {
            Alg: "HS256",
            Typ: "JWT",
        },
        Payload: IDTokenPayload {
            Iss: "Auth Server",
            Sub: data.Uid,
            Aud: data.Uid,
            Nonce: "somerandomcrap",
            DisplayName: data.DisplayName,
            Email: data.Email,
            AuthTime: time.Now(),
            Iat: time.Now(),
            Exp: time.Now().Add(time.Hour * 2),
        },
    }
    return token
}

// The key idea here is the payload and secret are combined
// in the hash. The client doesn't know the secret
// The hash therfore can only be constructed by the server
// Tampering with the message without knowing the secret, would
// make forging a hash very difficult
func SignIDToken(unsignedToken *IDToken, secret string) string {
    jsonHeader, err := json.Marshal(unsignedToken.Header)
    if err != nil {
        fmt.Println("Error encoding token header as JSON", err)
        return ""
    }

    jsonPayload, err := json.Marshal(unsignedToken.Payload)
    if err != nil {
        fmt.Println("Error encoding token payload as JSON", err)
        return ""
    }

    header := base64.StdEncoding.EncodeToString(jsonHeader)
    payload := base64.StdEncoding.EncodeToString(jsonPayload)
    toHash := []byte(header + "." + payload + secret)
    hash := sha256.Sum256(toHash)

    finalToken := header + "." + payload + "." + base64.StdEncoding.EncodeToString(hash[:])
    return finalToken
}

func ValidateToken(token string, secret string) bool {
    tokenShape := `^([^\.]+)\.([^\.]+)\.([^\.]+)$`
    re := regexp.MustCompile(tokenShape)
    matches := re.FindStringSubmatch(token)

    // index 0 is the entire matched string
    if len(matches) != 4 {
        fmt.Println("Invalid token shape")
        return false 
    }

    header := matches[1]
    payload := matches[2]
    toHash := []byte(header + "." + payload + secret)
    hash := sha256.Sum256(toHash)
    encodedHash := base64.StdEncoding.EncodeToString(hash[:])

    return encodedHash == matches[3]
}


