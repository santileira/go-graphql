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
var userLoaderCtxKey = &contextKey{"userLoader"}

type contextKey struct {
	name string
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userLoader := UserLoader{
				maxBatch: 100,
				wait:     1 * time.Millisecond,
				fetch: func(ids []int) ([]models.User, []error) {
					fmt.Printf("Devolviendo usuario a partir de ids en data loader %v \n", len(ids))
					usersResponse := make([]models.User, 0)

					for _, id := range ids {
						userResponse := database.Get(id)
						if userResponse != nil {
							usersResponse = append(usersResponse, *userResponse)
						}
					}

					return usersResponse, nil
				},
			}
			ctx := context.WithValue(r.Context(), userLoaderCtxKey, &userLoader)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext finds the role from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *UserLoader {
	raw, _ := ctx.Value(userLoaderCtxKey).(*UserLoader)
	return raw
}
