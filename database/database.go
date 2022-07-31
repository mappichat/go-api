package database

import "time"

type Reply struct {
	Body       string
	UserHandle string
	Timestamp  time.Time
}

type Post struct {
	Id         string
	Title      string
	Body       string
	UserHandle string
	Timestamp  time.Time
	Latitude   float32
	Longitude  float32
	Level      int8
	Replies    []Reply
	UpVotes    []string // maps to user handles
	DownVotes  []string
}

func InsertPost(post *Post) error {
	return InsertPostScylla(post)
}

func ReadPosts(latitude float32, longitude float32, latitudeDelta float32, longitudeDelta float32) ([]Post, error) {
	return ReadPostsScylla()
}
