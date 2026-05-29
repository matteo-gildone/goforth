package goforth

import (
	"testing"
)

func TestNewObject(t *testing.T) {
	objectId := "sword"
	objectName := "Sword"
	objectDescription := "An Elven sword"

	o := NewObject(objectId, objectName, objectDescription, true)

	if o.ID != objectId {
		t.Errorf("want: %q, got: %q", objectId, o.ID)
	}

	if o.Name != objectName {
		t.Errorf("want: %q, got: %q", objectName, o.Name)
	}

	if o.Description != objectDescription {
		t.Errorf("want: %q, got: %q", objectDescription, o.Description)
	}

	if !o.Takeable {
		t.Errorf("want: %v, got: %v", true, o.Takeable)
	}
}
