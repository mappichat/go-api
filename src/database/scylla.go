package database

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
	"github.com/uber/h3-go/v3"
)

var cluster *gocql.ClusterConfig
var session gocqlx.Session
var partitionResolution int = 1

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
	PartKey    string
	Tile       string
	Title      string
	Body       string
	AccountId  string
	UserHandle string
	TimeStamp  int64
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
		Tile:       p.Tile,
		Title:      p.Title,
		Body:       p.Body,
		AccountId:  p.AccountId,
		UserHandle: p.UserHandle,
		TimeStamp:  time.UnixMilli(p.TimeStamp),
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

	partKey := h3.ToString(h3.ToParent(h3.FromString(post.Tile), partitionResolution))

	return &PostScylla{
		ID:         id,
		PartKey:    partKey,
		Tile:       post.Tile,
		Title:      post.Title,
		Body:       post.Body,
		AccountId:  post.AccountId,
		UserHandle: post.UserHandle,
		TimeStamp:  post.TimeStamp.UnixMilli(),
		Latitude:   post.Latitude,
		Longitude:  post.Longitude,
		Level:      post.Level,
		ReplyCount: post.ReplyCount,
		UpVotes:    post.UpVotes,
		DownVotes:  post.DownVotes,
	}, nil
}

type ReplyScylla struct {
	PostID     gocql.UUID
	ID         gocql.UUID
	Body       string
	AccountId  string
	UserHandle string
	TimeStamp  int64
}

func (r ReplyScylla) toReply() Reply {
	return Reply{
		PostID:     r.PostID.String(),
		ID:         r.ID.String(),
		Body:       r.Body,
		AccountId:  r.AccountId,
		UserHandle: r.UserHandle,
		TimeStamp:  time.UnixMilli(r.TimeStamp),
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
		AccountId:  reply.AccountId,
		UserHandle: reply.UserHandle,
		TimeStamp:  reply.TimeStamp.UnixMilli(),
	}, nil
}

type VoteScylla struct {
	PostID    gocql.UUID
	AccountId string
	Up        bool
	Level     int8
	TimeStamp int64
}

func (v VoteScylla) toVote() Vote {
	return Vote{
		PostID:    v.PostID.String(),
		AccountId: v.AccountId,
		Up:        v.Up,
		Level:     v.Level,
		TimeStamp: time.UnixMilli(v.TimeStamp),
	}
}

func convertVote(vote *Vote) (*VoteScylla, error) {
	postId, err := gocql.ParseUUID(vote.PostID)
	if err != nil {
		return nil, err
	}

	return &VoteScylla{
		PostID:    postId,
		AccountId: vote.AccountId,
		Up:        vote.Up,
		Level:     vote.Level,
		TimeStamp: vote.TimeStamp.UnixMilli(),
	}, nil
}

// Helpers

var tableNames = struct {
	Posts   string
	Replies string
	Votes   string
}{
	Posts:   "user_data.posts",
	Replies: "user_data.replies",
	Votes:   "user_data.votes",
}

const readOneFormat = "GET * FROM %s WHERE %s"

func readOneScylla(table string, ops []string, names []string, bindMap map[string]interface{}, dest interface{}) error {
	q := session.Query(
		fmt.Sprintf(readOneFormat, table, strings.Join(ops, " AND ")),
		names,
	).BindMap(bindMap)

	if err := q.GetRelease(dest); err != nil {
		return err
	}

	return nil
}

const readManyFormat = "SELECT * FROM %s WHERE %s"

func readManyScylla(table string, ops []string, names []string, bindMap map[string]interface{}, dest interface{}) error {
	log.Print(fmt.Sprintf(readManyFormat, table, strings.Join(ops, " AND ")))
	q := session.Query(
		fmt.Sprintf(readManyFormat, table, strings.Join(ops, " AND ")),
		names,
	).BindMap(bindMap)

	if err := q.SelectRelease(dest); err != nil {
		return err
	}

	return nil
}

const insertOneFormat = "INSERT INTO %s (%s) VALUES (%s)"

func insertOneScylla(table string, bindMap map[string]interface{}) error {
	names := make([]string, len(bindMap))
	vals := make([]string, len(bindMap))
	i := 0
	for k := range bindMap {
		names[i] = k
		vals[i] = "?"
		i++
	}

	log.Print(fmt.Sprintf(insertOneFormat, table, strings.Join(names, ","), strings.Join(vals, ",")))
	q := session.Query(
		fmt.Sprintf(insertOneFormat, table, strings.Join(names, ","), strings.Join(vals, ",")),
		names,
	).BindMap(bindMap)

	if err := q.ExecRelease(); err != nil {
		return err
	}

	return nil
}

const updateFormat = "UPDATE %s SET %s WHERE %s"

func updateScylla(table string, setOps []string, ops []string, names []string, bindMap map[string]interface{}) error {
	log.Print(fmt.Sprintf(updateFormat, table, strings.Join(setOps, ","), strings.Join(ops, " AND ")))
	q := session.Query(
		fmt.Sprintf(updateFormat, table, strings.Join(setOps, ","), strings.Join(ops, " AND ")),
		names,
	).BindMap(bindMap)

	if err := q.ExecRelease(); err != nil {
		return err
	}

	return nil
}

const deleteFormat = "DELETE FROM %s WHERE %s IF EXISTS"

func deleteScylla(table string, ops []string, names []string, bindMap map[string]interface{}) error {
	log.Print(fmt.Sprintf(deleteFormat, table, strings.Join(ops, " AND ")))
	q := session.Query(
		fmt.Sprintf(deleteFormat, table, strings.Join(ops, " AND ")),
		names,
	).BindMap(bindMap)

	if err := q.ExecRelease(); err != nil {
		return err
	}

	return nil
}
