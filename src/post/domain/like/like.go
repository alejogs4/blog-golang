package like

// Like is all the information alongside a like
type Like struct {
	ID     string
	PostID string
	UserID string
	Type   Type
	State  State
}

// CreateNewLike factory function to create a right like
func CreateNewLike(id, postID, userID, Type, state string) (Like, error) {
	if id == "" || postID == "" || userID == "" {
		return Like{}, ErrBadLikeContent
	}

	likeType, err := CreateNewLikeType(Type)
	if err != nil {
		return Like{}, err
	}

	likeState, err := CreateNewLikeState(state)
	if err != nil {
		return Like{}, err
	}

	return Like{
		ID:     id,
		PostID: postID,
		UserID: userID,
		Type:   likeType,
		State:  likeState,
	}, nil
}

// SwitchType .
func (like *Like) SwitchType() (Like, error) {
	switchedType := like.Type.Switch()
	return CreateNewLike(like.ID, like.PostID, like.UserID, switchedType.GetTypeValue(), like.State.GetLikeState())
}

// Equals .
func (like *Like) Equals(another Like) bool {
	if like.PostID != another.PostID || like.UserID != another.UserID || !like.Type.Equals(another.Type) || !like.State.Equals(another.State) {
		return false
	}

	return true
}
