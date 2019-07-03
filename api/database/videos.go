package database

import "github.com/santileira/go-graphql/api/models"

var videos []*models.Video

func init() {
	videos = make([]*models.Video, 0)
}

func Videos() []*models.Video {
	return videos
}

func AddVideo(video *models.Video) {
	videos = append(videos, video)
}
