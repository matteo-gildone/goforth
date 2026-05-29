package goforth

// Object represents a single object in the game world.
type Object struct {
	ID          string
	Name        string
	Description string
	Takeable    bool
}

// NewObject creates a new object with the given id, name, description and takeable flag.
func NewObject(id, name, description string, takeable bool) *Object {
	return &Object{
		ID:          id,
		Name:        name,
		Description: description,
		Takeable:    takeable,
	}
}
