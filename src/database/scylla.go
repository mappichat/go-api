package database

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/table"
	"github.com/scylladb/gocqlx/v2"
)

var cluster *gocql.ClusterConfig
var session gocqlx.Session

type PostScylla struct {
	ID         gocql.UUID
	Title      string
	Body       string
	UserHandle string
	Timestamp  int64
	Latitude   float32
	Longitude  float32
	Level      int8
	ReplyCount int32
	UpVotes    int32
	DownVotes  int32
}

func (p PostScylla) toPost() Post {
	return Post{
		ID:         p.ID.String(),
		Title:      p.Title,
		Body:       p.Body,
		UserHandle: p.UserHandle,
		Timestamp:  time.UnixMilli(p.Timestamp),
		Latitude:   p.Latitude,
		Longitude:  p.Longitude,
		Level:      p.Level,
		ReplyCount: p.ReplyCount,
		UpVotes:    p.UpVotes,
		DownVotes:  p.DownVotes,
	}
}

func convertPost(post *Post) (*PostScylla, error) {
	id, err := gocql.ParseUUID(post.ID)
	if err != nil {
		return nil, err
	}

	return &PostScylla{
		ID:         id,
		Title:      post.Title,
		Body:       post.Body,
		UserHandle: post.UserHandle,
		Timestamp:  post.Timestamp.UnixMilli(),
		Latitude:   post.Latitude,
		Longitude:  post.Longitude,
		Level:      post.Level,
		ReplyCount: post.ReplyCount,
		UpVotes:    post.UpVotes,
		DownVotes:  post.DownVotes,
	}, nil
}

var postTable = table.New(table.Metadata{
	Name: "user_data.posts",
	Columns: []string{
		"id", "title", "body", "user_handle",
		"timestamp", "latitude", "longitude",
		"level", "reply_count",
		"up_votes", "down_votes",
	},
	PartKey: []string{"id"},
	SortKey: []string{"timestamp", "latitude", "longitude"},
})

func initializeScylla(hosts []string) error {
	var err error
	// Create gocql cluster.
	cluster = gocql.NewCluster(hosts...)
	// Wrap session on creation, gocqlx session embeds gocql.Session pointer.
	session, err = gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		return err
	}

	return nil
}

func insertPostScylla(post *Post) error {
	scyllaPost, err := convertPost(post)
	if err != nil {
		return err
	}

	q := session.Query(postTable.Insert()).BindStruct(scyllaPost)
	if err = q.ExecRelease(); err != nil {
		return err
	}
	return nil
}

func readPostsScylla() ([]Post, error) {
	scyllaPosts := []PostScylla{}
	// stmt, names := postTable.Select()
	q := session.Query(`SELECT * FROM user_data.posts`, []string{})
	if err := q.SelectRelease(&scyllaPosts); err != nil {
		return nil, err
	}

	posts := []Post{}
	for _, scyllaPost := range scyllaPosts {
		posts = append(posts, scyllaPost.toPost())
	}

	return posts, nil
}
