package middlewares

import (
	"context"
	"net/http"
	"project_management/api/constants"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := r.Header.Get(string(constants.AuthTokenCtxKey))
		// fmt.Println(authToken)
		var ctx = context.WithValue(r.Context(), constants.AuthTokenCtxKey, authToken)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
