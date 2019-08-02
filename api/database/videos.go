package database

import "github.com/santileira/go-graphql/api/models"

// videos collection in memory
var videos []*models.Video

// init creates storage in memory
func init() {
	videos = make([]*models.Video, 0)
}

func Videos() []*models.Video {
	return videos
}

func AddVideo(video *models.Video) {
	videos = append(videos, video)
}
