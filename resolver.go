package go_graphql

import (
	"context"
	"fmt"
	"github.com/santileira/go-graphql/api/database"
	"github.com/santileira/go-graphql/api/dataloaders/user"
	"github.com/santileira/go-graphql/api/models"
	"math/rand"
	"strconv"
	"time"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

var videoPublishedChannel map[string]chan *models.Video

func init() {
	videoPublishedChannel = map[string]chan *models.Video{}
}

type Resolver struct {
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Video() VideoResolver {
	return &videoResolver{r}
}

func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

// ******** MUTATION ********
type mutationResolver struct {
	*Resolver
}

func (r *mutationResolver) CreateVideo(ctx context.Context, input NewVideo) (*models.Video, error) {
	fmt.Println("Creando video")
	video := &models.Video{
		ID:          rand.Int(),
		Name:        input.Name,
		Description: input.Description,
		UserID:      input.UserID,
		URL:         input.URL,
		CreatedAt:   time.Now(),
	}

	database.AddVideo(video)

	// notify new video
	// add new video in videoPublishedChannel
	for _, observer := range videoPublishedChannel {
		observer <- video
		// this sends new video to client via socket
	}

	return video, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input NewUser) (*models.User, error) {
	fmt.Println("Creando usuario")
	id := rand.Int()
	idString := strconv.Itoa(id)

	user := &models.User{
		ID:    id,
		Name:  input.Name + "_" + idString,
		Email: input.Email + "_" + idString,
	}

	database.AddUser(user)
	return user, nil
}

// ******** QUERY ********
type queryResolver struct {
	*Resolver
}

func (r *queryResolver) Videos(ctx context.Context) ([]*models.Video, error) {
	fmt.Println("Devolviendo videos")
	return database.Videos(), nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*models.User, error) {
	fmt.Println("Devolviendo usuarios")
	return database.Users(), nil
}

// ******** VIDEO ********
type videoResolver struct {
	*Resolver
}

func (r *videoResolver) User(ctx context.Context, obj *models.Video) (*models.User, error) {
	fmt.Println("Devolviendo usuario a partir de video")
	/*var userResult *models.User
	for _, user := range database.Users() {
		if user.ID == obj.UserID {
			userResult = user
			break
		}
	}

	return userResult, nil*/
	user, err := userdataloader.ForContext(ctx).Load(obj.UserID)
	return &user, err
}

// ******** SUBSCRIPTION ********
type subscriptionResolver struct {
	*Resolver
}

func (r *subscriptionResolver) VideoPublished(ctx context.Context) (<-chan *models.Video, error) {
	fmt.Println("Subscribiendome a la creación de un video")
	id := randString(8)

	videoEvent := make(chan *models.Video, 1)
	go func() {
		<-ctx.Done()
		delete(videoPublishedChannel, id)
	}()

	videoPublishedChannel[id] = videoEvent

	return videoEvent, nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
