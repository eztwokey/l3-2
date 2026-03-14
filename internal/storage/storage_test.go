package storage

import (
	"errors"
	"testing"
)

func TestCacheKey(t *testing.T) {
	key := cacheKey("abc123")
	expected := "short:abc123"

	if key != expected {
		t.Errorf("expected %q, got %q", expected, key)
	}
}

func TestIsUniqueViolation_DuplicateKey(t *testing.T) {
	err := errors.New("pq: duplicate key value violates unique constraint")
	if !isUniqueViolation(err) {
		t.Error("expected true for duplicate key error")
	}
}

func TestIsUniqueViolation_Code23505(t *testing.T) {
	err := errors.New("ERROR: 23505 unique_violation")
	if !isUniqueViolation(err) {
		t.Error("expected true for 23505 error")
	}
}

func TestIsUniqueViolation_OtherError(t *testing.T) {
	err := errors.New("connection refused")
	if isUniqueViolation(err) {
		t.Error("expected false for non-unique error")
	}
}

func TestIsUniqueViolation_Nil(t *testing.T) {
	if isUniqueViolation(nil) {
		t.Error("expected false for nil error")
	}
}
