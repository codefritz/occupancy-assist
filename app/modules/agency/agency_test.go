package agency

import "testing"

// test to test the occupied function
func TestOccupied(t *testing.T) {
	if asString(true) != "belegt" {
		t.Error("Expected belegt, got ", asString(true))
	}
	if asString(false) != "frei" {
		t.Error("Expected frei, got ", asString(false))
	}
}
