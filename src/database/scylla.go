package database

import (
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/table"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
)

var cluster *gocql.ClusterConfig
var session gocqlx.Session

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

// Posts

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
	// SortKey: []string{"timestamp", "latitude", "longitude"},
})

func readPostScylla(id string) (*Post, error) {
	scyllaID, err := gocql.ParseUUID(id)
	if err != nil {
		return nil, err
	}

	scyllaPost := PostScylla{
		ID: scyllaID,
	}
	q := session.Query(postTable.Get()).BindStruct(scyllaPost)
	if err := q.GetRelease(&scyllaPost); err != nil {
		return nil, err
	}

	post := scyllaPost.toPost()
	return &post, nil
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

func updatePostScylla(id string, updateMap map[string]interface{}) error {
	scyllaID, err := gocql.ParseUUID(id)
	if err != nil {
		return err
	}

	keys := []string{}
	bindMap := qb.M{"id": scyllaID}

	setStatement := ""
	for key, val := range updateMap {
		if key != "id" {
			setStatement += fmt.Sprintf("%s=?,", key)
			keys = append(keys, key)
			bindMap[key] = val
		}
	}
	keys = append(keys, "id")
	setStatement = setStatement[:len(setStatement)-1]

	q := session.Query("UPDATE user_data.posts SET "+setStatement+" WHERE id=?", keys).BindMap(bindMap)
	if err := q.ExecRelease(); err != nil {
		return err
	}

	return nil
}

func deletePostScylla(id string) error {
	scyllaID, err := gocql.ParseUUID(id)
	if err != nil {
		return err
	}

	log.Print(postTable.Delete())

	q := session.Query(postTable.Delete()).BindStruct(PostScylla{ID: scyllaID})
	if err := q.ExecRelease(); err != nil {
		return err
	}

	return nil
}

// Replies

type ReplyScylla struct {
	PostID     gocql.UUID
	ID         gocql.UUID
	Body       string
	UserHandle string
	Timestamp  int64
}

func (r ReplyScylla) toReply() Reply {
	return Reply{
		PostID:     r.PostID.String(),
		ID:         r.ID.String(),
		Body:       r.Body,
		UserHandle: r.UserHandle,
		Timestamp:  time.UnixMilli(r.Timestamp),
	}
}

func convertReply(reply *Reply) (*ReplyScylla, error) {
	postId, err := gocql.ParseUUID(reply.PostID)
	if err != nil {
		return nil, err
	}

	id, err := gocql.ParseUUID(reply.ID)
	if err != nil {
		return nil, err
	}

	return &ReplyScylla{
		PostID:     postId,
		ID:         id,
		Body:       reply.Body,
		UserHandle: reply.UserHandle,
		Timestamp:  reply.Timestamp.UnixMilli(),
	}, nil
}

var replyTable = table.New(table.Metadata{
	Name: "user_data.replies",
	Columns: []string{
		"post_id", "id", "body", "user_handle", "timestamp",
	},
	PartKey: []string{"post_id"},
	SortKey: []string{"id"},
})

func readReplyScylla(postID string, id string) (*Reply, error) {
	scyllaPostID, err := gocql.ParseUUID(postID)
	if err != nil {
		return nil, err
	}

	scyllaID, err := gocql.ParseUUID(id)
	if err != nil {
		return nil, err
	}

	scyllaReply := ReplyScylla{
		PostID: scyllaPostID,
		ID:     scyllaID,
	}

	q := session.Query(replyTable.Get()).BindStruct(scyllaReply)
	if err := q.GetRelease(&scyllaReply); err != nil {
		return nil, err
	}

	reply := scyllaReply.toReply()
	return &reply, nil
}

func readRepliesScylla(postId string) ([]Reply, error) {
	postID, err := gocql.ParseUUID(postId)
	if err != nil {
		return nil, err
	}

	scyllaReplies := []ReplyScylla{}
	q := session.Query(replyTable.Select()).BindMap(qb.M{"post_id": postID})
	if err := q.SelectRelease(&scyllaReplies); err != nil {
		return nil, err
	}

	replies := []Reply{}
	for _, scyllaReply := range scyllaReplies {
		replies = append(replies, scyllaReply.toReply())
	}

	return replies, nil
}

func insertReplyScylla(reply *Reply) error {
	scyllaReply, err := convertReply(reply)
	if err != nil {
		return err
	}

	q := session.Query(replyTable.Insert()).BindStruct(scyllaReply)
	if err = q.ExecRelease(); err != nil {
		return err
	}
	return nil
}

func updateReplyScylla(postID string, id string, updateMap map[string]interface{}) error {
	scyllaPostID, err := gocql.ParseUUID(postID)
	if err != nil {
		return err
	}

	scyllaID, err := gocql.ParseUUID(id)
	if err != nil {
		return err
	}

	keys := []string{}
	bindMap := qb.M{"post_id": scyllaPostID, "id": scyllaID}

	setStatement := ""
	for key, val := range updateMap {
		if key != "id" && key != "post_id" {
			setStatement += fmt.Sprintf("%s=?,", key)
			keys = append(keys, key)
			bindMap[key] = val
		}
	}
	keys = append(keys, "post_id")
	keys = append(keys, "id")
	setStatement = setStatement[:len(setStatement)-1]

	q := session.Query("UPDATE user_data.replies SET "+setStatement+" WHERE post_id=? AND id=?", keys).BindMap(bindMap)
	log.Print(q)
	if err := q.ExecRelease(); err != nil {
		return err
	}

	return nil
}

func deleteReplyScylla(postID string, id string) error {
	scyllaPostID, err := gocql.ParseUUID(postID)
	if err != nil {
		return err
	}

	scyllaID, err := gocql.ParseUUID(id)
	if err != nil {
		return err
	}

	log.Print(replyTable.Delete())

	q := session.Query(replyTable.Delete()).BindStruct(ReplyScylla{PostID: scyllaPostID, ID: scyllaID})
	if err := q.ExecRelease(); err != nil {
		return err
	}

	return nil
}
