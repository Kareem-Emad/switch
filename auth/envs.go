package auth

import "os"

var jwtSecret = os.Getenv("SWITCH_JWT_SECRET")
