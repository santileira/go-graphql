package database

import "github.com/santileira/go-graphql/api/models"

// users collection
var users []*models.User

// init creates storage in memory
func init() {
	users = make([]*models.User, 0)
}

func Users() []*models.User {
	return users
}

func Get(id int) *models.User {

	var userResponse *models.User

	for _, user := range users {
		if user.ID == id {
			userResponse = user
			break
		}
	}

	return userResponse
}

func AddUser(user *models.User) {
	users = append(users, user)
}
