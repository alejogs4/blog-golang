package application

import (
	"github.com/alejogs4/blog/src/post/domain/like"
	"github.com/alejogs4/blog/src/post/domain/post"
	"github.com/google/uuid"
)

// PostCommands list of commands/ use cases regarding blog posts
type PostCommands struct {
	postRepository post.PostRepository
}

func NewPostCommands(postRepository post.PostRepository) PostCommands {
	return PostCommands{postRepository}
}

// CreateNewPost add a new blog post
func (pc PostCommands) CreateNewPost(userID, title, content, picture string, tags []post.Tag) error {
	postUUID := uuid.New().String()
	newPost, error := post.CreateNewPost(postUUID, userID, title, content, picture, []post.Comment{}, tags, []like.Like{})
	if error != nil {
		return error
	}

	return pc.postRepository.CreatePost(newPost)
}

// CreateNewComment add new comment to a certain blog post validating comment content
func (pc PostCommands) CreateNewComment(userID, postID, content string) error {
	commentUUID := uuid.New().String()

	newComment, error := post.CreateNewComment(commentUUID, postID, userID, content)
	if error != nil {
		return error
	}

	return pc.postRepository.AddComment(newComment)
}

// RemovePostComment remove comment by its id validating if comment owner is the one removing it
func (pc PostCommands) RemovePostComment(id, userID string) error {
	postComment, error := pc.postRepository.GetPostCommentByID(id)
	if error != nil {
		return error
	}

	if postComment.UserID != userID {
		return post.ErrNoCommentOwner
	}

	return pc.postRepository.RemoveComment(postComment)
}

// AddLike add a new like to blog post, considering several cases and validations such as what happen if same like was already done
// or how handle a dislike when a like has been done already
func (pc PostCommands) AddLike(userID, postID, Type string) error {
	likeType, err := like.CreateNewLikeType(Type)
	if err != nil {
		return err
	}

	currentPost, error := pc.postRepository.GetPostByID(postID)
	if error != nil {
		return error
	}

	newLike, error := like.CreateNewLike(uuid.New().String(), postID, userID, likeType.GetTypeValue(), like.Active)
	if error != nil {
		return error
	}

	if currentPost.IsLikeAlreadyDone(newLike) {
		return pc.postRepository.RemoveLike(newLike)
	}

	reversedLike, err := newLike.SwitchType()
	if err != nil {
		return err
	}

	presentLike := currentPost.LookPresentUserLike(userID, reversedLike.Type)
	if presentLike.ID != "" {
		error := pc.postRepository.RemoveLike(presentLike)
		if error != nil {
			return error
		}
	}

	return pc.postRepository.AddLike(postID, newLike)
}
