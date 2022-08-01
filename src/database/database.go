package database

import "time"

type Reply struct {
	Body       string
	UserHandle string
	Timestamp  time.Time
}

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

func Initialize(hosts []string) error {
	return initializeScylla(hosts)
}

func InsertPost(post *Post) error {
	return insertPostScylla(post)
}

func ReadPosts(level int8, latitude float32, longitude float32, latitudeDelta float32, longitudeDelta float32) ([]Post, error) {
	return readPostsScylla()
}
