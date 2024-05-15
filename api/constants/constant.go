package constants

type contextKey string

var (
	UserIDCtxKey    = contextKey("userID")
	AuthTokenCtxKey = contextKey("authToken")
)

const (
	INVALID_TOKEN  = "Invalid Token Found"
	INVALID_CLAIMS = "Invalid Cliams"
)
