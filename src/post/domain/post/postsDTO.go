package post

type PostsDTO struct {
	ID            string `json:"id"`
	UserID        string `json:"user_id"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	Picture       string `json:"picture"`
	Likes         int    `json:"likes"`
	Dislikes      int    `json:"dislikes"`
	CommentsCount int    `json:"comments_count"`
}

func ToPostsDTO(post Post, likes, dislikes, commentsCount int) PostsDTO {
	return PostsDTO{
		ID:            post.ID,
		UserID:        post.UserID,
		Title:         post.Title,
		Content:       post.Content,
		Picture:       post.Picture,
		Likes:         likes,
		Dislikes:      dislikes,
		CommentsCount: commentsCount,
	}
}
