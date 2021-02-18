package post

import "strings"

type Tag struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

// CreateNewTag .
func CreateNewTag(id, content string) (Tag, error) {
	normalizedID := strings.TrimSpace(id)
	normalizedContent := strings.TrimSpace(content)

	if normalizedID == "" || normalizedContent == "" {
		return Tag{}, ErrInvalidTagInformation
	}

	return Tag{normalizedID, normalizedContent}, nil
}

// GetID .
func (t Tag) GetID() string {
	return t.ID
}

// GetContent .
func (t Tag) GetContent() string {
	return t.Content
}
