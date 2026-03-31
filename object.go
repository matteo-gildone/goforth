package goforth

// Object represents a single object in the game world.
type Object struct {
	ID   string
	Name string
}

// NewObject creates a new object
func NewObject(id, name string) *Object {
	return &Object{
		ID:   id,
		Name: name,
	}
}
