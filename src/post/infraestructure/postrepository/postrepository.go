package postrepository

import (
	"github.com/alejogs4/blog/src/post/domain/like"
	"github.com/alejogs4/blog/src/post/domain/post"
	"github.com/alejogs4/blog/src/shared/infraestructure/database"
)

type PostgresRepository struct{}

func (postgres PostgresRepository) CreatePost(post post.Post) error {
	_, err := database.PostgresDB.Exec(
		"INSERT INTO post(id, person_id, title, content, picture) VALUES($1, $2, $3, $4, $5)",
		post.ID, post.UserID, post.Title, post.Content, post.Picture,
	)

	if err != nil {
		return err
	}

	for _, tag := range post.Tags {
		_, err = database.PostgresDB.Exec(
			"INSERT INTO post_tag(tag_id, post_id) VALUES($1, $2)",
			tag.GetID(), post.ID,
		)
	}

	return nil
}

func (postgres PostgresRepository) AddLike(postID string, like like.Like) error {
	_, err := database.PostgresDB.Exec(
		"INSERT INTO post_like(id, person_id, post_id, state, type) VALUES($1, $2, $3, $4, $5)",
		like.ID, like.UserID, like.PostID, like.State.GetLikeState(), like.Type.GetTypeValue(),
	)

	return err
}

func (postgres PostgresRepository) RemoveLike(removedLike like.Like) error {
	_, err := database.PostgresDB.Exec(
		"UPDATE post_like SET state=$1 WHERE person_id=$2 AND post_id=$3",
		like.Removed, removedLike.UserID, removedLike.PostID,
	)

	return err
}

func (postgres PostgresRepository) AddComment(comment post.Comment) error {
	_, err := database.PostgresDB.Exec(
		"INSERT INTO comment(id, content, person_id, post_id, state) VALUES($1, $2, $3, $4, $5)",
		comment.ID, comment.Content, comment.UserID, comment.PostID, comment.State,
	)

	return err
}

func (postgres PostgresRepository) RemoveComment(comment post.Comment) error {
	_, err := database.PostgresDB.Exec(
		"UPDATE comment SET state=$1 WHERE id=$2",
		post.RemovedComment, comment.ID,
	)

	return err
}

func (postgres PostgresRepository) GetPostCommentByID(id string) (post.Comment, error) {
	var comment post.Comment

	result := database.PostgresDB.QueryRow("SELECT id, content, person_id, post_id, state FROM comment WHERE post_id=$1", id)
	err := result.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.PostID, &comment.State)

	if err != nil {
		return comment, post.ErrUnexistentComment
	}

	return comment, nil
}

func (postgres PostgresRepository) GetAllPosts() ([]post.Post, error) {
	// TODO: Optimize this query
	var posts []post.Post = []post.Post{}
	result, err := database.PostgresDB.Query(`
		SELECT p.id, p.person_id, p.title, p.content, p.picture FROM post AS p
	`)
	defer result.Close()

	if err != nil {
		return posts, err
	}

	for result.Next() {
		var newPost post.Post

		err := result.Scan(
			&newPost.ID,
			&newPost.UserID,
			&newPost.Title,
			&newPost.Content,
			&newPost.Picture,
		)
		if err != nil {
			return posts, err
		}
		newPost.Comments = []post.Comment{}
		newPost.Tags = []post.Tag{}
		newPost.Likes = []like.Like{}

		posts = append(posts, newPost)
	}

	return posts, nil
}

func (postgres PostgresRepository) GetPostLikes(postID string) ([]like.Like, error) {
	var likes []like.Like

	result, err := database.PostgresDB.Query(`
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

func (postgres PostgresRepository) GetPostByID(postID string) (post.Post, error) {
	return post.Post{}, nil
}
