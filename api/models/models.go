package models

import (
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/santileira/go-graphql/api/errors"

	"io"
	"strconv"
	"time"
)

type Screenshot struct {
	ID      int    `json:"id"`
	VideoID int    `json:"videoId"`
	URL     string `json:"url"`
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Video struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UserID      int       `json:"-"`
	URL         string    `json:"url"`
	CreatedAt   time.Time `json:"createdAt"`
	Related     []*Video  `json:"related"`
}

func MarshalID(id int) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.Quote(fmt.Sprintf("%d", id)))
	})
}

func UnmarshalID(v interface{}) (int, error) {
	id, ok := v.(string)
	if !ok {
		return 0, fmt.Errorf("ids must be strings")
	}
	i, e := strconv.Atoi(id)
	return int(i), e
}

const layoutFormat = "2006-01-02 15:04:05"

func MarshalDateTime(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		io.WriteString(w, strconv.Quote(t.Format(layoutFormat)))
	})
}

func UnmarshalDateTime(v interface{}) (time.Time, error) {
	if tmpStr, ok := v.(string); ok {
		tmpTime, _ := time.Parse(layoutFormat, tmpStr)
		return tmpTime, nil
	}
	return time.Time{}, errors.TimeStampError
}
