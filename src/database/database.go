package database

import (
	"time"
)

type Post struct {
	ID         string
	Title      string
	Body       string
	UserHandle string
	Timestamp  time.Time
	Latitude   float32
	Longitude  float32
	Level      int8
	ReplyCount int32
	UpVotes    int32
	DownVotes  int32
}

type Reply struct {
	PostID     string
	ID         string
	Body       string
	UserHandle string
	Timestamp  time.Time
}

type Vote struct {
	PostId     string
	UserHandle string
	Up         bool
}

func Initialize(hosts []string) error {
	return initializeScylla(hosts)
}

// Posts

func ReadPost(id string) (*Post, error) {
	return readPostScylla(id)
}

func ReadPosts(level int8, latitude float32, longitude float32, latitudeDelta float32, longitudeDelta float32) ([]Post, error) {
	return readPostsScylla()
}

func InsertPost(post *Post) error {
	return insertPostScylla(post)
}

func UpdatePost(id string, updateMap map[string]interface{}) error {
	return updatePostScylla(id, updateMap)
}

func DeletePost(id string) error {
	return deletePostScylla(id)
}

// Replies

func ReadReply(postID string, id string) (*Reply, error) {
	return readReplyScylla(postID, id)
}

func ReadReplies(postId string) ([]Reply, error) {
	return readRepliesScylla(postId)
}

func InsertReply(reply *Reply) error {
	return insertReplyScylla(reply)
}

func UpdateReply(postID string, id string, updateMap map[string]interface{}) error {
	return updateReplyScylla(postID, id, updateMap)
}

func DeleteReply(postID string, id string) error {
	return deleteReplyScylla(postID, id)
}
