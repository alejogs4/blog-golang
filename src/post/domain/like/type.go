package like

import "strings"

const (
	// Dislike literal value to represent dislike
	Dislike = "DISLIKE"
	// TLike literal value to represent like
	TLike = "LIKE"
)

// Type value object to represent if it's like or dislike
type Type struct {
	Value string
}

// CreateNewLikeType factory function to create a valid like type
func CreateNewLikeType(value string) (Type, error) {
	normalizedValue := strings.ToUpper(strings.TrimSpace(value))
	if Dislike != normalizedValue && TLike != normalizedValue {
		return Type{}, ErrInvalidLikeType
	}

	return Type{value}, nil
}

// Switch like type, like if it was previously dislike or dislike otherwise
func (t *Type) Switch() Type {
	var lookType string
	if t.GetTypeValue() == TLike {
		lookType = Dislike
	} else {
		lookType = TLike
	}

	return Type{lookType}
}

// Equals value object function to see if another type it's equal as the current one
func (t *Type) Equals(anotherType Type) bool {
	return t.Value == anotherType.GetTypeValue()
}

// GetTypeValue getter function for type value
func (t *Type) GetTypeValue() string {
	return t.Value
}
