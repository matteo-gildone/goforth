package goforth

import (
	"testing"
)

func TestWorld_NewObject(t *testing.T) {
	objectId := "sword"
	objectName := "Sword"
	o := NewObject(objectId, objectName)

	if o.ID != objectId {
		t.Errorf("want: %q, got: %q", "entrance", o.ID)
	}

	if o.Name != objectName {
		t.Errorf("want: %q, got: %q", objectName, o.Name)
	}
}
