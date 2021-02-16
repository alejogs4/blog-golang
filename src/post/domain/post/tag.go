package post

type Tag struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

func CreateNewTag(id, content string) (Tag, error) {
	if id == "" || content == "" {
		return Tag{}, ErrInvalidTagInformation
	}

	return Tag{id, content}, nil
}

func (t Tag) GetID() string {
	return t.ID
}

func (t Tag) GetContent() string {
	return t.Content
}
