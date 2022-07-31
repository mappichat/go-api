package database

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/qb"
	"github.com/scylladb/gocqlx/table"
	"github.com/scylladb/gocqlx/v2"
)

var cluster *gocql.ClusterConfig
var session gocqlx.Session

type PostScylla struct {
	Id              gocql.UUID
	Title           string
	Body            string
	UserHandle      string
	Timestamp       int64
	Latitude        float32
	Longitude       float32
	Level           int8
	ReplyData       []map[string]string // {"body": "", "user_handle": ""}
	ReplyTimestamps []int64
	UpVotes         []string // maps to user handles
	DownVotes       []string
}

func (p PostScylla) toPost() Post {
	replies := []Reply{}
	for i, reply := range p.ReplyData {
		replies = append(replies, Reply{
			Body:       reply["body"],
			UserHandle: reply["user_handle"],
			Timestamp:  time.UnixMilli(p.ReplyTimestamps[i]),
		})
	}

	return Post{
		Id:         p.Id.String(),
		Title:      p.Title,
		Body:       p.Body,
		UserHandle: p.UserHandle,
		Timestamp:  time.UnixMilli(p.Timestamp),
		Latitude:   p.Latitude,
		Longitude:  p.Longitude,
		Level:      p.Level,
		Replies:    replies,
		UpVotes:    p.UpVotes,
		DownVotes:  p.DownVotes,
	}
}

func convertPost(post *Post) (*PostScylla, error) {
	id, err := gocql.ParseUUID(post.Id)
	if err != nil {
		return nil, err
	}
	replyData := []map[string]string{}
	replyTimestamps := []int64{}

	for _, reply := range post.Replies {
		replyData = append(replyData, map[string]string{"body": reply.Body, "user_handle": reply.UserHandle})
		replyTimestamps = append(replyTimestamps, reply.Timestamp.UnixMilli())
	}

	return &PostScylla{
		Id:              id,
		Title:           post.Title,
		Body:            post.Body,
		UserHandle:      post.UserHandle,
		Timestamp:       post.Timestamp.UnixMilli(),
		Latitude:        post.Latitude,
		Longitude:       post.Longitude,
		Level:           post.Level,
		ReplyData:       replyData,
		ReplyTimestamps: replyTimestamps,
		UpVotes:         post.UpVotes,
		DownVotes:       post.DownVotes,
	}, nil
}

var postTable = table.New(table.Metadata{
	Name: "user_data.posts",
	Columns: []string{
		"id", "title", "body", "account",
		"timestamp", "latitude", "longitude",
		"level", "reply_data", "reply_timestamps",
		"up_votes", "down_votes",
	},
	PartKey: []string{"id"},
	SortKey: []string{"timestamp", "latitude", "longitude"},
})

var postTableCreateCQL = `CREATE TABLE IF NOT EXISTS user_data.posts (
	id uuid PRIMARY KEY,
	title text,
	body text,
	user_handle string,
	timestamp timestamp,
	latitude float,
	longitude float,
	level int
	reply_data list<map<string, string>>
	reply_timestamps list<timestamp>
	up_votes set<string>
	down_votes set<string>
	)`

func Initialize(hosts []string) error {
	var err error
	// Create gocql cluster.
	cluster = gocql.NewCluster(hosts...)
	// Wrap session on creation, gocqlx session embeds gocql.Session pointer.
	session, err = gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		return err
	}

	err = session.ExecStmt(postTableCreateCQL)

	if err != nil {
		return err
	}

	return nil
}

func InsertPostScylla(post *Post) error {
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

func ReadPostsScylla() ([]Post, error) {
	var scyllaPosts []PostScylla
	q := session.Query(postTable.Select()).BindMap(qb.M{})
	if err := q.SelectRelease(&scyllaPosts); err != nil {
		return nil, err
	}

	posts := []Post{}
	for _, scyllaPost := range scyllaPosts {
		posts = append(posts, scyllaPost.toPost())
	}

	return posts, nil
}
