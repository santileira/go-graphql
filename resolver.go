package go_graphql

import (
	"context"
	"database/sql"
	"github.com/santileira/go-graphql/database"
	"github.com/santileira/go-graphql/errors"
	"github.com/santileira/go-graphql/models"
	"time"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	db *sql.DB
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct {
	*Resolver
}

func (r *mutationResolver) CreateVideo(ctx context.Context, input NewVideo) (models.Video, error) {
	newVideo := models.Video{
		URL:       input.URL,
		Name:      input.Name,
		CreatedAt: time.Now().UTC(),
	}

	rows, err := database.LogAndQuery(r.db, "INSERT INTO videos (name, url, user_id, created_at) VALUES($1, $2, $3, $4) RETURNING id",
		input.Name, input.URL, input.UserID, newVideo.CreatedAt)
	defer rows.Close()

	if err != nil || !rows.Next() {
		return models.Video{}, err
	}
	if err := rows.Scan(&newVideo.ID); err != nil {
		errors.DebugPrintf(err)
		if errors.IsForeignKeyError(err) {
			return models.Video{}, errors.UserNotExist
		}
		return models.Video{}, errors.InternalServerError
	}

	return newVideo, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Videos(ctx context.Context, limit *int, offset *int) ([]models.Video, error) {
	var video models.Video
	var videos []models.Video

	rows, err := database.LogAndQuery(r.db, "SELECT id, name, url, created_at, user_id FROM videos ORDER BY created_at desc limit $1 offset $2", limit, offset)
	defer rows.Close()
	if err != nil {
		errors.DebugPrintf(err)
		return nil, errors.InternalServerError
	}
	for rows.Next() {
		if err := rows.Scan(&video.ID, &video.Name, &video.URL, &video.CreatedAt, &video.UserID); err != nil {
			errors.DebugPrintf(err)
			return nil, errors.InternalServerError
		}
		videos = append(videos, video)
	}

	return videos, nil
}
