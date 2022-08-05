package database

import "time"

type Account struct {
	ID         string `json:"id" db:"id"`
	UserHandle string `json:"user_handle" db:"user_handle"`
	Email      string `json:"email" db:"email"`
}

type Post struct {
	ID        string    `json:"id" db:"id"`
	AccountId string    `json:"account_id" db:"account_id"`
	Title     string    `json:"title" db:"title"`
	Body      string    `json:"body" db:"body"`
	Latitude  float64   `json:"latitude" db:"latitude"`
	Longitude float64   `json:"longitude" db:"longitude"`
	Level     int8      `json:"level" db:"post_level"`
	TimeStamp time.Time `json:"time_stamp"  db:"time_stamp"`
}

type Reply struct {
	ID        string    `json:"id" db:"id"`
	PostID    string    `json:"post_id" db:"post_id"`
	AccountId string    `json:"account_id" db:"account_id"`
	Body      string    `json:"body" db:"body"`
	Latitude  float64   `json:"latitude" db:"latitude"`
	Longitude float64   `json:"longitude" db:"longitude"`
	TimeStamp time.Time `json:"time_stamp" db:"time_stamp"`
}

type Vote struct {
	PostID     string    `json:"post_id" db:"post_id"`
	AccountId  string    `json:"account_id" db:"account_id"`
	VoteWeight float64   `json:"vote_weight" db:"vote_weight"`
	Level      int8      `json:"level" db:"vote_level"`
	Latitude   float64   `json:"latitude" db:"latitude"`
	Longitude  float64   `json:"longitude" db:"longitude"`
	TimeStamp  time.Time `json:"time_stamp" db:"time_stamp"`
}

// func Initialize(connectString string) error {
// 	// return initializeScylla(hosts)
// 	return SqlInitialize(connectString)
// }

// // Posts

// func ReadPost(id string) (*Post, error) {
// 	return readPostPostgres(id)
// 	// return readPostScylla(id, tile)
// }

// func ReadPosts(tiles []string) ([]Post, error) {
// 	return readPostsPostgres(tiles)
// 	// return readPostsScylla(level, tiles)
// }

// func InsertPost(post *Post) error {
// 	return insertPostPostgres(post)
// 	// return insertPostScylla(post)
// }

// func UpdatePost(id string, newPost *Post) error {
// 	updatePostPostgres(id, accountId)
// 	// return updatePostScylla(id, tile, level, accountId, updateMap)
// }

// func DeletePost(id string, tile string, level int8, accountId string) error {
// 	return deletePostScylla(id, tile, level, accountId)
// }

// // Replies

// func ReadReply(postID string, id string) (*Reply, error) {
// 	return readReplyScylla(postID, id)
// }

// func ReadReplies(postId string) ([]Reply, error) {
// 	return readRepliesScylla(postId)
// }

// func InsertReply(reply *Reply) error {
// 	return insertReplyScylla(reply)
// }

// func UpdateReply(postID string, id string, accountId string, updateMap map[string]interface{}) error {
// 	return updateReplyScylla(postID, id, accountId, updateMap)
// }

// func DeleteReply(postID string, id string, accountId string) error {
// 	return deleteReplyScylla(postID, id, accountId)
// }

// // Votes

// func ReadVote(postID string, accountId string, level int8) (*Vote, error) {
// 	return readVoteScylla(postID, accountId, level)
// }

// func ReadVotes(postID string) ([]Vote, error) {
// 	return readVotesScylla(postID)
// }

// func InsertVote(vote *Vote) error {
// 	return insertVoteScylla(vote)
// }

// func UpdateVote(postID string, level int8, accountId string, updateMap map[string]interface{}) error {
// 	return updateVoteScylla(postID, level, accountId, updateMap)
// }

// func DeleteVote(postID string, level int8, accountId string) error {
// 	return deleteVoteScylla(postID, level, accountId)
// }
