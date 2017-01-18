package main

import (
	"testing"
)

func TestRobotsURLSuccessful(t *testing.T) {
	expectedURL := "http://my-cool-domain.com/robots.txt"
	result := "http://my-cool-domain.com/robots.txt"

	if result != expectedURL {
		t.Fatal("Expected " + expectedURL + " but got " + result)
	}
}
