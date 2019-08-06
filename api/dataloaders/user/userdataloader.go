package userdataloader

import (
	"context"
	"fmt"
	"github.com/santileira/go-graphql/api/database"
	"github.com/santileira/go-graphql/api/models"
	"net/http"
	"time"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var userLoaderCtxKey = &contextKey{"UserLoader"}

// contextKey structure to store value of UserLoader
type contextKey struct {
	name string
}

// Middleware creates de UserLoader and put value in the context with key "UserLoader".
// Generates function fetch to returns slice of users by slice of ids.
func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userLoader := UserLoader{
				maxBatch: 10,
				wait:     1 * time.Millisecond,
				fetch: func(ids []int) ([]models.User, []error) {
					fmt.Println("Fetching users in user data loader")
					usersResponse := make([]models.User, 0)

					for _, id := range ids {
						fmt.Println(fmt.Printf("Fetching user %v in user data loader", id))
						userResponse := database.Get(id)
						if userResponse != nil {
							usersResponse = append(usersResponse, *userResponse)
						}
					}

					return usersResponse, nil
				},
			}

			// put userLoader in the context
			ctx := context.WithValue(r.Context(), userLoaderCtxKey, &userLoader)
			r = r.WithContext(ctx)
			// and call the next with our new context
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext finds the UserLoader from the context.
func ForContext(ctx context.Context) *UserLoader {
	raw, _ := ctx.Value(userLoaderCtxKey).(*UserLoader)
	return raw
}
