package bksqlite

import (
    "sync"
)

type SQLiteBackend struct {
    sync.Mutex
}

/*

// GetKeys returns a list of keys for a given path
func (sb *SQLiteBackend) GetKeys(bucket []string) ([]string, error) {
}

// CreateVal creates value with a given TTL and key in the bucket
// if the value already exists, returns AlreadyExistsError
func (sb *SQLiteBackend) CreateVal(bucket []string, key string, val []byte, ttl time.Duration) error {
}

// TouchVal updates the TTL of the key without changing the value
func (sb *SQLiteBackend) TouchVal(bucket []string, key string, ttl time.Duration) error {
}

// UpsertVal updates or inserts value with a given TTL into a bucket
// ForeverTTL for no TTL
func (sb *SQLiteBackend) UpsertVal(bucket []string, key string, val []byte, ttl time.Duration) error {
}

// GetVal return a value for a given key in the bucket
func (sb *SQLiteBackend) GetVal(path []string, key string) ([]byte, error) {
}

// GetValAndTTL returns value and TTL for a key in bucket
func (sb *SQLiteBackend) GetValAndTTL(bucket []string, key string) ([]byte, time.Duration, error) {
}

// DeleteKey deletes a key in a bucket
func (sb *SQLiteBackend) DeleteKey(bucket []string, key string) error {
}

// DeleteBucket deletes the bucket by a given path
func (sb *SQLiteBackend) DeleteBucket(path []string, bkt string) error {
}

// AcquireLock grabs a lock that will be released automatically in TTL
func (sb *SQLiteBackend) AcquireLock(token string, ttl time.Duration) error {
}

// ReleaseLock forces lock release before TTL
func (sb *SQLiteBackend) ReleaseLock(token string) error {
}

// CompareAndSwap implements compare ans swap operation for a key
func (sb *SQLiteBackend) CompareAndSwap(bucket []string, key string, val []byte, ttl time.Duration, prevVal []byte) ([]byte, error) {
}

// Close releases the resources taken up by this backend
func (sb *SQLiteBackend) Close() error {
}

*/
