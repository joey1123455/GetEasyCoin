package utils

import (
	"testing"
	"time"
)

func TestCacheSetGet(t *testing.T) {
	cache := NewCache()

	// Test setting a value
	cache.Set("key1", "value1", 1*time.Second)

	// Test getting the value
	value, found := cache.Get("key1")
	if !found {
		t.Error("Expected to find key1")
	}
	if value != "value1" {
		t.Errorf("Expected value1, got %v", value)
	}

	// Wait for the item to expire
	time.Sleep(2 * time.Second)

	// Test getting the value after expiration
	value, found = cache.Get("key1")
	if found {
		t.Error("Expected not to find key1 after expiration")
	}
}

func TestCacheDelete(t *testing.T) {
	cache := NewCache()

	// Test setting a value
	cache.Set("key2", "value2", 1*time.Second)

	// Test deleting the value
	cache.Delete("key2")

	// Test getting the value after deletion
	_, found := cache.Get("key2")
	if found {
		t.Error("Expected not to find key2 after deletion")
	}
}

func TestCacheClear(t *testing.T) {
	cache := NewCache()

	// Test setting values
	cache.Set("key3", "value3", 1*time.Second)
	cache.Set("key4", "value4", 1*time.Second)

	// Test clearing the cache
	cache.Clear()

	// Test getting the values after clearing
	_, found := cache.Get("key3")
	if found {
		t.Error("Expected not to find key3 after clearing")
	}
	_, found = cache.Get("key4")
	if found {
		t.Error("Expected not to find key4 after clearing")
	}

}
