package postrepository

import (
	"database/sql"
	"errors"

	"github.com/alejogs4/blog/src/post/domain/like"
	"github.com/alejogs4/blog/src/post/domain/post"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) PostgresRepository {
	return PostgresRepository{db}
}

func (postgres PostgresRepository) CreatePost(post post.Post) error {
	_, err := postgres.db.Exec(
		"INSERT INTO post(id, person_id, title, content, picture) VALUES($1, $2, $3, $4, $5)",
		post.ID, post.UserID, post.Title, post.Content, post.Picture,
	)

	if err != nil {
		return err
	}

	for _, tag := range post.Tags {
		_, err = postgres.db.Exec(
			"INSERT INTO post_tag(tag_id, post_id) VALUES($1, $2)",
			tag.GetID(), post.ID,
		)
	}

	return nil
}

func (postgres PostgresRepository) AddLike(postID string, like like.Like) error {
	_, err := postgres.db.Exec(
		"INSERT INTO post_like(id, person_id, post_id, state, type) VALUES($1, $2, $3, $4, $5)",
		like.ID, like.UserID, like.PostID, like.State.GetLikeState(), like.Type.GetTypeValue(),
	)

	return err
}

func (postgres PostgresRepository) RemoveLike(removedLike like.Like) error {
	_, err := postgres.db.Exec(
		"UPDATE post_like SET state=$1 WHERE person_id=$2 AND post_id=$3",
		like.Removed, removedLike.UserID, removedLike.PostID,
	)

	return err
}

func (postgres PostgresRepository) AddComment(comment post.Comment) error {
	_, err := postgres.db.Exec(
		"INSERT INTO comment(id, content, person_id, post_id, state) VALUES($1, $2, $3, $4, $5)",
		comment.ID, comment.Content, comment.UserID, comment.PostID, comment.State,
	)

	return err
}

func (postgres PostgresRepository) RemoveComment(comment post.Comment) error {
	_, err := postgres.db.Exec(
		"UPDATE comment SET state=$1 WHERE id=$2",
		post.RemovedComment, comment.ID,
	)

	return err
}

func (postgres PostgresRepository) GetPostCommentByID(id string) (post.Comment, error) {
	var comment post.Comment

	result := postgres.db.QueryRow("SELECT id, content, person_id, post_id, state FROM comment WHERE id=$1", id)
	err := result.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.PostID, &comment.State)

	if err != nil {
		return comment, post.ErrUnexistentComment
	}

	return comment, nil
}

func (postgres PostgresRepository) GetAllPosts() ([]post.PostsDTO, error) {
	// TODO: Optimize this query
	var posts []post.PostsDTO = []post.PostsDTO{}
	result, err := postgres.db.Query(`
		SELECT
		p.id, p.person_id, p.title, p.content, p.picture,
		(SELECT COUNT(l.id) FROM post_like AS l WHERE l.post_id=p.id AND state=$1 AND type=$2) as likes,
		(SELECT COUNT(l.id) FROM post_like AS l WHERE l.post_id=p.id AND state=$1 AND type=$3) as dislikes,
		(SELECT COUNT(c.id) FROM comment AS c WHERE c.post_id=p.id AND state=$4) as comments_count
		FROM post AS p
	`, like.Active, like.TLike, like.Dislike, post.ActiveComment)

	if err != nil {
		return nil, err
	}
	defer result.Close()

	for result.Next() {
		var newPost post.Post
		var likes int = 0
		var dislikes int = 0
		var comments int = 0

		err := result.Scan(
			&newPost.ID,
			&newPost.UserID,
			&newPost.Title,
			&newPost.Content,
			&newPost.Picture,
			&likes,
			&dislikes,
			&comments,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post.ToPostsDTO(newPost, likes, dislikes, comments))
	}

	return posts, nil
}

func (postgres PostgresRepository) GetPostLikes(postID string) ([]like.Like, error) {
	var likes []like.Like = []like.Like{}

	result, err := postgres.db.Query(`
	SELECT id, post_id, person_id, state, type
	FROM post_like
	WHERE post_id=$1 AND state=$2`,
		postID, like.Active,
	)
	if err != nil {
		return likes, err
	}

	for result.Next() {
		var newLike like.Like
		err := result.Scan(&newLike.ID, &newLike.PostID, &newLike.UserID, &newLike.State.Value, &newLike.Type.Value)
		if err != nil {
			return likes, err
		}

		likes = append(likes, newLike)
	}

	return likes, nil
}

func (postgres PostgresRepository) getPostTags(postID string) ([]post.Tag, error) {
	var tags []post.Tag = []post.Tag{}
	rows, err := postgres.db.Query(
		"SELECT t.id, t.content FROM tag t INNER JOIN post_tag pt ON pt.tag_id = t.id WHERE pt.post_id=$1", postID,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tag post.Tag

		err := rows.Scan(&tag.ID, &tag.Content)
		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (postgres PostgresRepository) getPostComments(postID string) ([]post.Comment, error) {
	var comments []post.Comment = []post.Comment{}

	rows, err := postgres.db.Query(
		"SELECT id, content, person_id, post_id, state FROM comment c WHERE post_id=$1", postID,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var comment post.Comment

		err := rows.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.PostID, &comment.State)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

func (postgres PostgresRepository) GetPostByID(postID string) (post.Post, error) {
	var returnedPost post.Post
	result := postgres.db.QueryRow(
		"SELECT id, person_id, title, content, picture FROM post WHERE id=$1", postID,
	)

	err := result.Scan(&returnedPost.ID, &returnedPost.UserID, &returnedPost.Title, &returnedPost.Content, &returnedPost.Picture)
	if errors.Is(err, sql.ErrNoRows) {
		return returnedPost, post.ErrNoFoundPost
	}

	if err != nil {
		return returnedPost, err
	}

	comments, err := postgres.getPostComments(returnedPost.ID)
	if err != nil {
		return returnedPost, err
	}

	likes, err := postgres.GetPostLikes(returnedPost.ID)
	if err != nil {
		return returnedPost, err
	}

	tags, err := postgres.getPostTags(returnedPost.ID)
	if err != nil {
		return returnedPost, err
	}

	returnedPost.Comments = comments
	returnedPost.Likes = likes
	returnedPost.Tags = tags

	return returnedPost, nil
}
