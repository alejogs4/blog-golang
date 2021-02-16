package like

import "strings"

const (
	// Active like state
	Active = "ACTIVE"
	// Removed like state
	Removed = "REMOVED"
)

// State value object that represent current like state
type State struct {
	Value string `json:"value"`
}

// CreateNewLikeState factory function to create a right like state
func CreateNewLikeState(value string) (State, error) {
	normalizedValue := strings.ToUpper(strings.TrimSpace(value))
	if normalizedValue != Active && normalizedValue != Removed {
		return State{}, ErrInvalidLikeState
	}

	return State{normalizedValue}, nil
}

// Equals value object equals function to see if it's equal to another like state
func (l State) Equals(another State) bool {
	return l.Value == another.GetLikeState()
}

// GetLikeState .
func (l State) GetLikeState() string {
	return l.Value
}
