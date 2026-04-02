package goforth

import (
	"testing"
)

func TestNewObject(t *testing.T) {
	objectId := "sword"
	objectName := "Sword"
	o := NewObject(objectId, objectName)

	if o.ID != objectId {
		t.Errorf("want: %q, got: %q", objectId, o.ID)
	}

	if o.Name != objectName {
		t.Errorf("want: %q, got: %q", objectName, o.Name)
	}
}
