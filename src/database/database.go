package database

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/mappichat/go-api.git/src/utils"
	"github.com/uber/h3-go/v3"
)

type Post struct {
	ID         string    `json:"id"`
	Tile       string    `json:"tile"`
	Title      string    `json:"title"`
	Body       string    `json:"body"`
	AccountId  string    `json:"account_id"`
	UserHandle string    `json:"user_handle"`
	TimeStamp  time.Time `json:"time_stamp"`
	Latitude   float32   `json:"latitude"`
	Longitude  float32   `json:"longitude"`
	Level      int8      `json:"level"`
	ReplyCount int32     `json:"reply_count"`
	UpVotes    int32     `json:"up_votes"`
	DownVotes  int32     `json:"down_votes"`
}

type Reply struct {
	PostID     string    `json:"post_id"`
	ID         string    `json:"id"`
	Body       string    `json:"body"`
	AccountId  string    `json:"account_id"`
	UserHandle string    `json:"user_handle"`
	TimeStamp  time.Time `json:"time_stamp"`
}

type Vote struct {
	PostID    string    `json:"post_id"`
	AccountId string    `json:"account_id"`
	Up        bool      `json:"up"`
	Level     int8      `json:"level"`
	TimeStamp time.Time `json:"time_stamp"`
}

func Initialize(hosts []string) error {
	return initializeScylla(hosts)
}

// Posts

func ReadPost(id string, tile string) (*Post, error) {
	partKey := h3.ToString(h3.ToParent(h3.FromString(tile), partitionResolution))
	dest := PostScylla{}
	if err := readOneScylla(
		tableNames.Posts,
		[]string{"part_key=?", "tile=?", "id=?"},
		[]string{"part_key", "tile", "id"},
		map[string]interface{}{"part_key": partKey, "tile": tile, "id": id},
		&dest,
		false,
	); err != nil {
		return nil, err
	}

	post := dest.toPost()
	return &post, nil
}

func ReadPosts(level int8, tiles []string) ([]Post, error) {
	partKeys := map[string]bool{}
	for _, tile := range tiles {
		key := h3.ToString(h3.ToParent(h3.FromString(tile), partitionResolution))
		partKeys[key] = true
	}

	keyMarks := make([]string, len(partKeys))
	names := make([]string, len(partKeys)+len(tiles)+1)
	bindMap := map[string]interface{}{}
	i := 0
	for key := range partKeys {
		keyMarks[i] = "?"
		names[i] = key
		bindMap[key] = key
		i++
	}

	names[i] = "level"
	bindMap["level"] = level
	i++

	tileMarks := make([]string, len(tiles))
	for j, tile := range tiles {
		tileMarks[j] = "?"
		bindMap[tile] = tile
		names[i] = tile
		i++
	}

	dest := []PostScylla{}
	if err := readManyScylla(
		tableNames.Posts,
		[]string{fmt.Sprintf("part_key IN (%s)", strings.Join(keyMarks, ",")), "level=?", fmt.Sprintf("tile IN (%s)", strings.Join(tileMarks, "?"))},
		names,
		bindMap,
		&dest,
		false,
	); err != nil {
		return nil, err
	}

	posts := make([]Post, len(dest))
	for j, scyllaPost := range dest {
		posts[j] = scyllaPost.toPost()
	}

	return posts, nil
}

func InsertPost(post *Post) error {
	scyllaPost, err := convertPost(post)
	if err != nil {
		return err
	}
	bindMap, err := utils.DecodeSnakeCase(*scyllaPost)
	if err != nil {
		return err
	}
	if err = insertOneScylla(
		tableNames.Posts,
		bindMap,
		false,
	); err != nil {
		return err
	}
	return nil
}

func UpdatePost(id string, tile string, level int8, accountId string, updateMap map[string]interface{}) error {
	scyllaID, err := gocql.ParseUUID(id)
	if err != nil {
		return err
	}

	names := make([]string, len(updateMap)+5)
	setOps := make([]string, len(updateMap))
	i := 0
	for k := range updateMap {
		names[i] = k
		setOps[i] = fmt.Sprintf("%s=?", k)
		i++
	}

	updateMap["part_key"] = h3.ToString(h3.ToParent(h3.FromString(tile), partitionResolution))
	names[i] = "part_key"
	updateMap["level"] = level
	names[i+1] = "level"
	updateMap["tile"] = tile
	names[i+2] = "tile"
	updateMap["id"] = scyllaID
	names[i+3] = "id"
	updateMap["account_id"] = accountId
	names[i+4] = "account_id"

	if err = updateScylla(
		tableNames.Posts,
		setOps,
		[]string{"part_key=?", "level=?", "tile=?", "id=?", "account_id=?"},
		names,
		updateMap,
		false,
	); err != nil {
		return err
	}
	return nil
}

func DeletePost(id string, tile string, level int8, accountId string) error {
	scyllaID, err := gocql.ParseUUID(id)
	if err != nil {
		return err
	}
	partKey := h3.ToString(h3.ToParent(h3.FromString(tile), partitionResolution))
	if err = deleteScylla(
		tableNames.Posts,
		[]string{"part_key=?", "id=?", "level=?", "tile=?", "account_id=?"},
		[]string{"part_key", "id", "level", "tile", "account_id"},
		map[string]interface{}{"part_key": partKey, "id": scyllaID, "level": level, "tile": tile, "account_id": accountId},
		false,
	); err != nil {
		return err
	}
	return nil
}

// Replies

func ReadReply(postID string, id string) (*Reply, error) {
	return nil, nil
}

func ReadReplies(postId string) ([]Reply, error) {
	return nil, nil
}

func InsertReply(reply *Reply) error {
	return nil
}

func UpdateReply(postID string, id string, accountId string, updateMap map[string]interface{}) error {
	return nil
}

func DeleteReply(postID string, id string, accountId string) error {
	return nil
}

// Votes

func ReadVote(postID string, accountId string, level int8) (*Vote, error) {
	return nil, nil
}

func ReadVotes(postId string) ([]Vote, error) {
	return nil, nil
}

func InsertVote(vote *Vote) error {
	return nil
}

func UpdateVote(postID string, accountId string, level int8, updateMap map[string]interface{}) error {
	return nil
}

func DeleteVote(postID string, accountId string, level int8) error {
	return nil
}
