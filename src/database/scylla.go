package database

// import (
// 	"fmt"
// 	"strings"
// 	"time"

// 	"github.com/gocql/gocql"
// 	"github.com/mappichat/go-api.git/src/utils"
// 	"github.com/scylladb/gocqlx/v2"
// 	"github.com/uber/h3-go/v3"
// )

// var cluster *gocql.ClusterConfig
// var session gocqlx.Session
// var partitionResolution int = 1

// func initializeScylla(hosts []string) error {
// 	var err error
// 	// Create gocql cluster.
// 	cluster = gocql.NewCluster(hosts...)
// 	// Wrap session on creation, gocqlx session embeds gocql.Session pointer.
// 	session, err = gocqlx.WrapSession(cluster.CreateSession())
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// // Posts

// type PostScylla struct {
// 	ID         gocql.UUID
// 	PartKey    string
// 	Tile       string
// 	Title      string
// 	Body       string
// 	AccountId  string
// 	UserHandle string
// 	TimeStamp  int64
// 	Latitude   float32
// 	Longitude  float32
// 	Level      int8
// 	ReplyCount int32
// 	UpVotes    int32
// 	DownVotes  int32
// }

// func (p PostScylla) toPost() Post {
// 	return Post{
// 		ID:         p.ID.String(),
// 		Tile:       p.Tile,
// 		Title:      p.Title,
// 		Body:       p.Body,
// 		AccountId:  p.AccountId,
// 		UserHandle: p.UserHandle,
// 		TimeStamp:  time.UnixMilli(p.TimeStamp),
// 		Latitude:   p.Latitude,
// 		Longitude:  p.Longitude,
// 		Level:      p.Level,
// 		ReplyCount: p.ReplyCount,
// 		UpVotes:    p.UpVotes,
// 		DownVotes:  p.DownVotes,
// 	}
// }

// func convertPost(post *Post) (*PostScylla, error) {
// 	id, err := gocql.ParseUUID(post.ID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	partKey := h3.ToString(h3.ToParent(h3.FromString(post.Tile), partitionResolution))

// 	return &PostScylla{
// 		ID:         id,
// 		PartKey:    partKey,
// 		Tile:       post.Tile,
// 		Title:      post.Title,
// 		Body:       post.Body,
// 		AccountId:  post.AccountId,
// 		UserHandle: post.UserHandle,
// 		TimeStamp:  post.TimeStamp.UnixMilli(),
// 		Latitude:   post.Latitude,
// 		Longitude:  post.Longitude,
// 		Level:      post.Level,
// 		ReplyCount: post.ReplyCount,
// 		UpVotes:    post.UpVotes,
// 		DownVotes:  post.DownVotes,
// 	}, nil
// }

// type ReplyScylla struct {
// 	PostID     gocql.UUID
// 	ID         gocql.UUID
// 	Body       string
// 	AccountId  string
// 	UserHandle string
// 	Tile       string
// 	Latitude   float32
// 	Longitude  float32
// 	TimeStamp  int64
// }

// func (r ReplyScylla) toReply() Reply {
// 	return Reply{
// 		PostID:     r.PostID.String(),
// 		ID:         r.ID.String(),
// 		Body:       r.Body,
// 		AccountId:  r.AccountId,
// 		UserHandle: r.UserHandle,
// 		Tile:       r.Tile,
// 		Latitude:   r.Latitude,
// 		Longitude:  r.Longitude,
// 		TimeStamp:  time.UnixMilli(r.TimeStamp),
// 	}
// }

// func convertReply(reply *Reply) (*ReplyScylla, error) {
// 	postId, err := gocql.ParseUUID(reply.PostID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	id, err := gocql.ParseUUID(reply.ID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &ReplyScylla{
// 		PostID:     postId,
// 		ID:         id,
// 		Body:       reply.Body,
// 		AccountId:  reply.AccountId,
// 		UserHandle: reply.UserHandle,
// 		Tile:       reply.Tile,
// 		Latitude:   reply.Latitude,
// 		Longitude:  reply.Longitude,
// 		TimeStamp:  reply.TimeStamp.UnixMilli(),
// 	}, nil
// }

// type VoteScylla struct {
// 	PostID    gocql.UUID
// 	AccountId string
// 	Up        bool
// 	Level     int8
// 	Tile      string
// 	Latitude  float32
// 	Longitude float32
// 	TimeStamp int64
// }

// func (v VoteScylla) toVote() Vote {
// 	return Vote{
// 		PostID:    v.PostID.String(),
// 		AccountId: v.AccountId,
// 		Up:        v.Up,
// 		Level:     v.Level,
// 		Tile:      v.Tile,
// 		Latitude:  v.Latitude,
// 		Longitude: v.Longitude,
// 		TimeStamp: time.UnixMilli(v.TimeStamp),
// 	}
// }

// func convertVote(vote *Vote) (*VoteScylla, error) {
// 	postId, err := gocql.ParseUUID(vote.PostID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &VoteScylla{
// 		PostID:    postId,
// 		AccountId: vote.AccountId,
// 		Up:        vote.Up,
// 		Level:     vote.Level,
// 		Tile:      vote.Tile,
// 		Latitude:  vote.Latitude,
// 		Longitude: vote.Longitude,
// 		TimeStamp: vote.TimeStamp.UnixMilli(),
// 	}, nil
// }

// // Helpers

// var tableNames = struct {
// 	Posts   string
// 	Replies string
// 	Votes   string
// }{
// 	Posts:   "user_data.posts",
// 	Replies: "user_data.replies",
// 	Votes:   "user_data.votes",
// }

// const readOneFormat = "GET * FROM %s WHERE %s %s"

// func readOneScylla(table string, ops []string, names []string, bindMap map[string]interface{}, dest interface{}, filtering bool) error {
// 	filterString := ""
// 	if filtering {
// 		filterString = "ALLOW FILTERING"
// 	}

// 	q := session.Query(
// 		fmt.Sprintf(readOneFormat, table, strings.Join(ops, " AND "), filterString),
// 		names,
// 	).BindMap(bindMap)

// 	if err := q.GetRelease(dest); err != nil {
// 		return err
// 	}

// 	return nil
// }

// const readManyFormat = "SELECT * FROM %s WHERE %s %s"

// func readManyScylla(table string, ops []string, names []string, bindMap map[string]interface{}, dest interface{}, filtering bool) error {
// 	filterString := ""
// 	if filtering {
// 		filterString = "ALLOW FILTERING"
// 	}

// 	q := session.Query(
// 		fmt.Sprintf(readManyFormat, table, strings.Join(ops, " AND "), filterString),
// 		names,
// 	).BindMap(bindMap)

// 	if err := q.SelectRelease(dest); err != nil {
// 		return err
// 	}

// 	return nil
// }

// const insertOneFormat = "INSERT INTO %s (%s) VALUES (%s) %s"

// func insertOneScylla(table string, bindMap map[string]interface{}, filtering bool) error {
// 	filterString := ""
// 	if filtering {
// 		filterString = "ALLOW FILTERING"
// 	}

// 	names := make([]string, len(bindMap))
// 	vals := make([]string, len(bindMap))
// 	i := 0
// 	for k := range bindMap {
// 		names[i] = k
// 		vals[i] = "?"
// 		i++
// 	}
// 	q := session.Query(
// 		fmt.Sprintf(insertOneFormat, table, strings.Join(names, ","), strings.Join(vals, ","), filterString),
// 		names,
// 	).BindMap(bindMap)

// 	if err := q.ExecRelease(); err != nil {
// 		return err
// 	}

// 	return nil
// }

// const updateFormat = "UPDATE %s SET %s WHERE %s %s"

// func updateScylla(table string, setOps []string, ops []string, names []string, bindMap map[string]interface{}, filtering bool) error {
// 	filterString := ""
// 	if filtering {
// 		filterString = "ALLOW FILTERING"
// 	}

// 	q := session.Query(
// 		fmt.Sprintf(updateFormat, table, strings.Join(setOps, ","), strings.Join(ops, " AND "), filterString),
// 		names,
// 	).BindMap(bindMap)

// 	if err := q.ExecRelease(); err != nil {
// 		return err
// 	}

// 	return nil
// }

// const deleteFormat = "DELETE FROM %s WHERE %s IF EXISTS %s"

// func deleteScylla(table string, ops []string, names []string, bindMap map[string]interface{}, filtering bool) error {
// 	filterString := ""
// 	if filtering {
// 		filterString = "ALLOW FILTERING"
// 	}

// 	q := session.Query(
// 		fmt.Sprintf(deleteFormat, table, strings.Join(ops, " AND "), filterString),
// 		names,
// 	).BindMap(bindMap)

// 	if err := q.ExecRelease(); err != nil {
// 		return err
// 	}

// 	return nil
// }

// // db funcs

// func readPostScylla(id string, tile string) (*Post, error) {
// 	partKey := h3.ToString(h3.ToParent(h3.FromString(tile), partitionResolution))
// 	dest := PostScylla{}
// 	if err := readOneScylla(
// 		tableNames.Posts,
// 		[]string{"part_key=?", "tile=?", "id=?"},
// 		[]string{"part_key", "tile", "id"},
// 		map[string]interface{}{"part_key": partKey, "tile": tile, "id": id},
// 		&dest,
// 		false,
// 	); err != nil {
// 		return nil, err
// 	}

// 	post := dest.toPost()
// 	return &post, nil
// }

// func readPostsScylla(level int8, tiles []string) ([]Post, error) {
// 	partKeys := map[string]bool{}
// 	for _, tile := range tiles {
// 		key := h3.ToString(h3.ToParent(h3.FromString(tile), partitionResolution))
// 		partKeys[key] = true
// 	}

// 	keyMarks := make([]string, len(partKeys))
// 	names := make([]string, len(partKeys)+len(tiles)+1)
// 	bindMap := map[string]interface{}{}
// 	i := 0
// 	for key := range partKeys {
// 		keyMarks[i] = "?"
// 		names[i] = key
// 		bindMap[key] = key
// 		i++
// 	}

// 	names[i] = "level"
// 	bindMap["level"] = level
// 	i++

// 	tileMarks := make([]string, len(tiles))
// 	for j, tile := range tiles {
// 		tileMarks[j] = "?"
// 		bindMap[tile] = tile
// 		names[i] = tile
// 		i++
// 	}

// 	dest := []PostScylla{}
// 	if err := readManyScylla(
// 		tableNames.Posts,
// 		[]string{fmt.Sprintf("part_key IN (%s)", strings.Join(keyMarks, ",")), "level=?", fmt.Sprintf("tile IN (%s)", strings.Join(tileMarks, "?"))},
// 		names,
// 		bindMap,
// 		&dest,
// 		false,
// 	); err != nil {
// 		return nil, err
// 	}

// 	posts := make([]Post, len(dest))
// 	for j, scyllaPost := range dest {
// 		posts[j] = scyllaPost.toPost()
// 	}

// 	return posts, nil
// }

// func insertPostScylla(post *Post) error {
// 	scyllaPost, err := convertPost(post)
// 	if err != nil {
// 		return err
// 	}
// 	bindMap, err := utils.DecodeSnakeCase(*scyllaPost)
// 	if err != nil {
// 		return err
// 	}
// 	if err = insertOneScylla(
// 		tableNames.Posts,
// 		bindMap,
// 		false,
// 	); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func updatePostScylla(id string, tile string, level int8, accountId string, updateMap map[string]interface{}) error {
// 	scyllaID, err := gocql.ParseUUID(id)
// 	if err != nil {
// 		return err
// 	}

// 	names := make([]string, len(updateMap)+5)
// 	setOps := make([]string, len(updateMap))
// 	i := 0
// 	for k := range updateMap {
// 		names[i] = k
// 		setOps[i] = fmt.Sprintf("%s=?", k)
// 		i++
// 	}

// 	updateMap["part_key"] = h3.ToString(h3.ToParent(h3.FromString(tile), partitionResolution))
// 	names[i] = "part_key"
// 	updateMap["level"] = level
// 	names[i+1] = "level"
// 	updateMap["tile"] = tile
// 	names[i+2] = "tile"
// 	updateMap["id"] = scyllaID
// 	names[i+3] = "id"
// 	updateMap["account_id"] = accountId
// 	names[i+4] = "account_id"

// 	if err = updateScylla(
// 		tableNames.Posts,
// 		setOps,
// 		[]string{"part_key=?", "level=?", "tile=?", "id=?", "account_id=?"},
// 		names,
// 		updateMap,
// 		false,
// 	); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func deletePostScylla(id string, tile string, level int8, accountId string) error {
// 	scyllaID, err := gocql.ParseUUID(id)
// 	if err != nil {
// 		return err
// 	}
// 	partKey := h3.ToString(h3.ToParent(h3.FromString(tile), partitionResolution))
// 	if err = deleteScylla(
// 		tableNames.Posts,
// 		[]string{"part_key=?", "id=?", "level=?", "tile=?", "account_id=?"},
// 		[]string{"part_key", "id", "level", "tile", "account_id"},
// 		map[string]interface{}{"part_key": partKey, "id": scyllaID, "level": level, "tile": tile, "account_id": accountId},
// 		false,
// 	); err != nil {
// 		return err
// 	}
// 	return nil
// }

// // Replies

// func readReplyScylla(postID string, id string) (*Reply, error) {
// 	pid, err := gocql.ParseUUID(postID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	rid, err := gocql.ParseUUID(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	dest := ReplyScylla{}
// 	if err := readOneScylla(
// 		tableNames.Replies,
// 		[]string{"post_id=?", "id=?"},
// 		[]string{"post_id", "id"},
// 		map[string]interface{}{"post_id": pid, "id": rid},
// 		&dest,
// 		false,
// 	); err != nil {
// 		return nil, err
// 	}

// 	reply := dest.toReply()
// 	return &reply, nil
// }

// func readRepliesScylla(postId string) ([]Reply, error) {
// 	pid, err := gocql.ParseUUID(postId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	dest := []ReplyScylla{}
// 	if err := readManyScylla(
// 		tableNames.Replies,
// 		[]string{"post_id=?"},
// 		[]string{"post_id"},
// 		map[string]interface{}{"post_id": pid},
// 		&dest,
// 		false,
// 	); err != nil {
// 		return nil, err
// 	}

// 	replies := make([]Reply, len(dest))
// 	for i, reply := range dest {
// 		replies[i] = reply.toReply()
// 	}
// 	return replies, nil
// }

// func insertReplyScylla(reply *Reply) error {
// 	scyllaReply, err := convertReply(reply)
// 	if err != nil {
// 		return err
// 	}
// 	bindMap, err := utils.DecodeSnakeCase(*scyllaReply)
// 	if err != nil {
// 		return err
// 	}
// 	if err = insertOneScylla(
// 		tableNames.Replies,
// 		bindMap,
// 		false,
// 	); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func updateReplyScylla(postID string, id string, accountId string, updateMap map[string]interface{}) error {
// 	pid, err := gocql.ParseUUID(postID)
// 	if err != nil {
// 		return err
// 	}
// 	rid, err := gocql.ParseUUID(id)
// 	if err != nil {
// 		return err
// 	}

// 	names := make([]string, len(updateMap)+3)
// 	setOps := make([]string, len(updateMap))
// 	i := 0
// 	for k := range updateMap {
// 		names[i] = k
// 		setOps[i] = fmt.Sprintf("%s=?", k)
// 		i++
// 	}

// 	updateMap["post_id"] = pid
// 	names[i] = "post_id"
// 	updateMap["id"] = rid
// 	names[i+1] = "id"
// 	updateMap["account_id"] = accountId
// 	names[i+2] = "account_id"

// 	if err = updateScylla(
// 		tableNames.Replies,
// 		setOps,
// 		[]string{"post_id=?", "id=?", "account_id=?"},
// 		names,
// 		updateMap,
// 		false,
// 	); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func deleteReplyScylla(postID string, id string, accountId string) error {
// 	pid, err := gocql.ParseUUID(postID)
// 	if err != nil {
// 		return err
// 	}
// 	rid, err := gocql.ParseUUID(id)
// 	if err != nil {
// 		return err
// 	}

// 	if err = deleteScylla(
// 		tableNames.Replies,
// 		[]string{"post_id=?", "id=?", "account_id=?"},
// 		[]string{"post_id", "id", "account_id"},
// 		map[string]interface{}{"post_id": pid, "id": rid, "account_id": accountId},
// 		false,
// 	); err != nil {
// 		return err
// 	}
// 	return nil
// }

// // Votes

// func readVoteScylla(postID string, accountId string, level int8) (*Vote, error) {
// 	pid, err := gocql.ParseUUID(postID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	dest := VoteScylla{}
// 	if err = readOneScylla(
// 		tableNames.Votes,
// 		[]string{"post_id=?", "account_id=?", "level=?"},
// 		[]string{"post_id", "account_id", "level"},
// 		map[string]interface{}{"post_id": pid, "account_id": accountId, "level": level},
// 		&dest,
// 		false,
// 	); err != nil {
// 		return nil, err
// 	}

// 	vote := dest.toVote()
// 	return &vote, nil
// }

// func readVotesScylla(postID string) ([]Vote, error) {
// 	pid, err := gocql.ParseUUID(postID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	dest := []VoteScylla{}
// 	if err = readManyScylla(
// 		tableNames.Votes,
// 		[]string{"post_id=?"},
// 		[]string{"post_id"},
// 		map[string]interface{}{"post_id": pid},
// 		&dest,
// 		false,
// 	); err != nil {
// 		return nil, err
// 	}

// 	votes := make([]Vote, len(dest))
// 	for i, vote := range dest {
// 		votes[i] = vote.toVote()
// 	}
// 	return votes, nil
// }

// func insertVoteScylla(vote *Vote) error {
// 	scyllaVote, err := convertVote(vote)
// 	if err != nil {
// 		return err
// 	}
// 	bindMap, err := utils.DecodeSnakeCase(*scyllaVote)
// 	if err != nil {
// 		return err
// 	}
// 	if err = insertOneScylla(
// 		tableNames.Votes,
// 		bindMap,
// 		false,
// 	); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func updateVoteScylla(postID string, level int8, accountId string, updateMap map[string]interface{}) error {
// 	pid, err := gocql.ParseUUID(postID)
// 	if err != nil {
// 		return err
// 	}

// 	names := make([]string, len(updateMap)+3)
// 	setOps := make([]string, len(updateMap))
// 	i := 0
// 	for k := range updateMap {
// 		names[i] = k
// 		setOps[i] = fmt.Sprintf("%s=?", k)
// 		i++
// 	}

// 	updateMap["post_id"] = pid
// 	names[i] = "post_id"
// 	updateMap["level"] = level
// 	names[i+1] = "level"
// 	updateMap["account_id"] = accountId
// 	names[i+2] = "account_id"

// 	if err = updateScylla(
// 		tableNames.Votes,
// 		setOps,
// 		[]string{"post_id=?", "level=?", "account_id=?"},
// 		names,
// 		updateMap,
// 		false,
// 	); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func deleteVoteScylla(postID string, level int8, accountId string) error {
// 	pid, err := gocql.ParseUUID(postID)
// 	if err != nil {
// 		return err
// 	}

// 	if err = deleteScylla(
// 		tableNames.Votes,
// 		[]string{"post_id=?", "level=?", "account_id=?"},
// 		[]string{"post_id", "level", "account_id"},
// 		map[string]interface{}{"post_id": pid, "level": level, "account_id": accountId},
// 		false,
// 	); err != nil {
// 		return err
// 	}
// 	return nil
// }
