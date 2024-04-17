package agency

import "testing"

// test to test the occupied function
func TestOccupied(t *testing.T) {
	if occupied(true) != "belegt" {
		t.Error("Expected belegt, got ", occupied(true))
	}
	if occupied(false) != "frei" {
		t.Error("Expected frei, got ", occupied(false))
	}
}
