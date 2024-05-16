package errors

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/lib/pq"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

var (
	AccessTokenExpired         = CreateCustomError("Access Token is Expired, Please Regenrate It.", http.StatusUnauthorized)
	UnauthorizedError          = CreateCustomError("You are Not Authorized to Perform this Action.", http.StatusUnauthorized)
	TokenNotFound              = CreateCustomError("Access Token Not Found.", http.StatusUnauthorized)
	NoUserFound                = CreateCustomError("No User Found for This Request.", http.StatusNotFound)
	PasswordMatchError         = CreateCustomError("Password is Incorrect.", http.StatusUnauthorized)
	InvalidToken               = CreateCustomError("Invalid Token", http.StatusUnauthorized)
	InternalServerError        = CreateCustomError("Internal Server Error", http.StatusInternalServerError)
	NoUserOrProjectFound       = CreateCustomError("No Rows Found Related UserId and ProjectId you entered", http.StatusNotFound)
	ProjectDoesNotExistError   = CreateCustomError("project does not exists", http.StatusNotFound)
	UserIDsRequiredError       = CreateCustomError("UserIDS are not Entered", http.StatusBadRequest)
	ScreenShotIDsRequiredError = CreateCustomError("ScreenShotIDS are not Entered", http.StatusBadRequest)
	UsersProjectsLimitError    = CreateCustomError("Limit Must between 1 to 5", http.StatusBadRequest)

	NotNullConstraintError    = CreateCustomError("Required Fields are Missing or Null", http.StatusBadRequest)
	ForeignKeyConstraintError = CreateCustomError("Your Operation violate ForeignKeyConstraint", http.StatusConflict)
	UniqueKeyConstraintError  = CreateCustomError("Data you are trying to add already exists", http.StatusConflict)
	CheckConstraintError      = CreateCustomError("Please ensure all values meet the specified criteria", http.StatusBadRequest)
	NoRowsError               = CreateCustomError("No data found matching your request", http.StatusNotFound)
)

func DatabaseErrorHandling(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23502":
			// Not-null constraint violation
			return NotNullConstraintError
		case "23503":
			// Foreign key violation
			return ForeignKeyConstraintError
		case "23505":
			// Unique constraint violation
			return UniqueKeyConstraintError
		case "23514":
			// Check constraint violation
			return CheckConstraintError
		}
	}
	if errors.Is(err, sql.ErrNoRows) {
		return NoRowsError
	}
	return InternalServerError
}

func CreateCustomError(message string, statusCode int) *gqlerror.Error {
	return &gqlerror.Error{
		Message: message,
		Extensions: map[string]interface{}{
			"StatusCode": statusCode,
		},
	}
}
