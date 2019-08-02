package go_graphql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/santileira/go-graphql/api/database"
	"github.com/santileira/go-graphql/api/dataloaders/user"
	"github.com/santileira/go-graphql/api/models"
	"math/rand"
	"strconv"
	"time"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

// createdVideosObservers contains all the observers with your specific id.
var createdVideosObservers map[string]chan *models.Video

// init initializes created videos observers variable.
func init() {
	createdVideosObservers = map[string]chan *models.Video{}
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

// CreateVideo adds video in database generating random id.
// Verifies if user id exists, if isn't exist returns error.
func (r *mutationResolver) CreateVideo(ctx context.Context, input NewVideo) (*models.Video, error) {

	fmt.Println("Handle request to create video")

	user := database.Get(input.UserID)
	if user == nil {
		fmt.Println("User not exists")
		return nil, errors.New("user not exists")
	}

	video := &models.Video{
		ID:          rand.Int(),
		Name:        input.Name,
		Description: input.Description,
		UserID:      input.UserID,
		URL:         input.URL,
		CreatedAt:   time.Now().UTC(),
	}

	database.AddVideo(video)

	videoJSON, err := json.Marshal(video)
	if err != nil {
		fmt.Println(fmt.Printf("Error marshalling video %s", err.Error()))
		return nil, err
	}

	fmt.Println(fmt.Printf("Create video %s", string(videoJSON)))


	// notify new video
	// add new video in createdVideosObservers
	for _, observer := range createdVideosObservers {
		observer <- video
		// this sends new video to client via socket
	}

	return video, nil
}

// CreateUser adds user in database generating random id.
func (r *mutationResolver) CreateUser(ctx context.Context, input NewUser) (*models.User, error) {

	fmt.Println("Handle request to create user")

	id := rand.Int()
	idString := strconv.Itoa(id)

	user := &models.User{
		ID:    id,
		Name:  input.Name + "_" + idString,
		Email: input.Email + "_" + idString,
	}

	database.AddUser(user)

	userJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Println(fmt.Printf("Error marshalling user %s", err.Error()))
		return nil, err
	}

	fmt.Println(fmt.Printf("Create user %s", string(userJSON)))

	return user, nil
}

// ******** QUERY ********
type queryResolver struct {
	*Resolver
}

// Videos returns videos and error.
func (r *queryResolver) Videos(ctx context.Context) ([]*models.Video, error) {
	fmt.Println("Returning videos")
	return database.Videos(), nil
}

// Users returns users and error.
func (r *queryResolver) Users(ctx context.Context) ([]*models.User, error) {
	fmt.Println("Returning users")
	return database.Users(), nil
}

// ******** VIDEO ********
type videoResolver struct {
	*Resolver
}

// User returns user by user id in video.
func (r *videoResolver) User(ctx context.Context, video *models.Video) (*models.User, error) {

	fmt.Println("Returning user by video")

	user, err := userdataloader.ForContext(ctx).Load(video.UserID)
	if err != nil {
		fmt.Println(fmt.Printf("Error searching user %s", err.Error()))
		return nil, err
	}

	return &user, err
}

// ******** SUBSCRIPTION ********
type subscriptionResolver struct {
	*Resolver
}

// VideoCreated creates video channel and returns it.
// When context is done, deletes the channel.
func (r *subscriptionResolver) VideoCreated(ctx context.Context) (<-chan *models.Video, error) {
	fmt.Println("Subscribing me to creation of videos")

	id := rand.Int()
	idStr := strconv.Itoa(id)
	videoCreatedChannel := make(chan *models.Video, 1)

	go func() {
		<-ctx.Done()
		fmt.Println("Closing subscription")
		delete(createdVideosObservers, idStr)
	}()

	createdVideosObservers[idStr] = videoCreatedChannel

	return videoCreatedChannel, nil
}