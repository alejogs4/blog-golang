package valueobject

// ValueObject interface must be implemented by any aggregate value object to verify structural equality between two value objects
type ValueObject interface {
	Equals(other ValueObject) bool
}
