package auth

import (
	"context"
	"github.com/santileira/go-graphql"
	"net/http"
)


// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var roleCtxKey = &contextKey{"Role"}

type contextKey struct {
	name string
}

// Middleware decodes the share session cookie and packs the session into context
func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			role := r.Header.Get("Role")
			// put it in context
			if role != "" {
				var roleGoGraphQL go_graphql.Role
				roleGoGraphQL = go_graphql.Role(role)
				// Write your fancy token introspection logic here and if valid user then pass appropriate key in header
				// IMPORTANT: DO NOT HANDLE UNAUTHORISED USER HERE
				ctx = context.WithValue(ctx, roleCtxKey, &roleGoGraphQL)
			}
			// and call the next with our new context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ForContext finds the role from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *go_graphql.Role {
	raw, _ := ctx.Value(roleCtxKey).(*go_graphql.Role)
	return raw
}
