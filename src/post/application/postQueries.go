package application

import "github.com/alejogs4/blog/src/post/domain/post"

type PostQueries struct {
	postRepository post.PostRepository
}

func NewPostQueries(postRepository post.PostRepository) PostQueries {
	return PostQueries{postRepository}
}

func (query PostQueries) GetAllPosts() ([]post.PostsDTO, error) {
	return query.postRepository.GetAllPosts()
}
