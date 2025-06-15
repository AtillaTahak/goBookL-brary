package test

import (
	"testing"
	"time"

	"github.com/AtillaTahaK/gobooklibrary/pkg/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RedisCacheTestSuite struct {
	suite.Suite
	cache *cache.RedisCache
}

func (suite *RedisCacheTestSuite) SetupTest() {
	// Use a test Redis instance or mock
	suite.cache = cache.NewRedisCache("localhost:6379", "", 1) // Use DB 1 for testing
}

func (suite *RedisCacheTestSuite) TearDownTest() {
	if suite.cache != nil {
		suite.cache.FlushAll() // Clean up test data
		suite.cache.Close()
	}
}

func (suite *RedisCacheTestSuite) TestSetAndGet() {
	testData := map[string]interface{}{
		"name": "Test Book",
		"id":   123,
	}

	// Test Set
	err := suite.cache.Set("test:book:1", testData, 5*time.Minute)
	if err != nil {
		// Skip test if Redis is not available
		suite.T().Skip("Redis not available, skipping test")
		return
	}

	// Test Get
	var retrieved map[string]interface{}
	err = suite.cache.Get("test:book:1", &retrieved)
	suite.NoError(err)
	suite.Equal("Test Book", retrieved["name"])
	suite.Equal(float64(123), retrieved["id"]) // JSON unmarshals numbers as float64
}

func (suite *RedisCacheTestSuite) TestGetNonExistentKey() {
	var data interface{}
	err := suite.cache.Get("non:existent:key", &data)
	suite.Error(err)
	suite.Contains(err.Error(), "key not found")
}

func (suite *RedisCacheTestSuite) TestDelete() {
	testData := "test value"

	// Set a key
	err := suite.cache.Set("test:delete", testData, 5*time.Minute)
	if err != nil {
		suite.T().Skip("Redis not available, skipping test")
		return
	}

	// Verify it exists
	exists, err := suite.cache.Exists("test:delete")
	suite.NoError(err)
	suite.True(exists)

	// Delete it
	err = suite.cache.Delete("test:delete")
	suite.NoError(err)

	// Verify it's gone
	exists, err = suite.cache.Exists("test:delete")
	suite.NoError(err)
	suite.False(exists)
}

func (suite *RedisCacheTestSuite) TestExpiration() {
	testData := "expiring value"

	// Set with short expiration
	err := suite.cache.Set("test:expire", testData, 1*time.Second)
	if err != nil {
		suite.T().Skip("Redis not available, skipping test")
		return
	}

	// Should exist immediately
	exists, err := suite.cache.Exists("test:expire")
	suite.NoError(err)
	suite.True(exists)

	// Wait for expiration
	time.Sleep(2 * time.Second)

	// Should be gone
	exists, err = suite.cache.Exists("test:expire")
	suite.NoError(err)
	suite.False(exists)
}

func (suite *RedisCacheTestSuite) TestIncrement() {
	// Test increment
	val, err := suite.cache.Incr("test:counter")
	if err != nil {
		suite.T().Skip("Redis not available, skipping test")
		return
	}
	suite.Equal(int64(1), val)

	// Increment again
	val, err = suite.cache.Incr("test:counter")
	suite.NoError(err)
	suite.Equal(int64(2), val)

	// Test IncrBy
	val, err = suite.cache.IncrBy("test:counter", 5)
	suite.NoError(err)
	suite.Equal(int64(7), val)
}

func (suite *RedisCacheTestSuite) TestSetNX() {
	// First SetNX should succeed
	success, err := suite.cache.SetNX("test:setnx", "value1", 5*time.Minute)
	if err != nil {
		suite.T().Skip("Redis not available, skipping test")
		return
	}
	suite.True(success)

	// Second SetNX should fail (key exists)
	success, err = suite.cache.SetNX("test:setnx", "value2", 5*time.Minute)
	suite.NoError(err)
	suite.False(success)

	// Verify original value is still there
	var value string
	err = suite.cache.Get("test:setnx", &value)
	suite.NoError(err)
	suite.Equal("value1", value)
}

func (suite *RedisCacheTestSuite) TestKeys() {
	// Set some test keys
	testKeys := []string{"test:pattern:1", "test:pattern:2", "test:other:1"}
	for _, key := range testKeys {
		err := suite.cache.Set(key, "value", 5*time.Minute)
		if err != nil {
			suite.T().Skip("Redis not available, skipping test")
			return
		}
	}

	// Search for pattern
	keys, err := suite.cache.Keys("test:pattern:*")
	suite.NoError(err)
	suite.Len(keys, 2)
	suite.Contains(keys, "test:pattern:1")
	suite.Contains(keys, "test:pattern:2")
}

func (suite *RedisCacheTestSuite) TestTTL() {
	// Set key with expiration
	err := suite.cache.Set("test:ttl", "value", 10*time.Second)
	if err != nil {
		suite.T().Skip("Redis not available, skipping test")
		return
	}

	// Check TTL
	ttl, err := suite.cache.TTL("test:ttl")
	suite.NoError(err)
	suite.True(ttl > 5*time.Second) // Should be around 10 seconds
	suite.True(ttl <= 10*time.Second)
}

func (suite *RedisCacheTestSuite) TestPing() {
	err := suite.cache.Ping()
	if err != nil {
		suite.T().Skip("Redis not available, skipping test")
		return
	}
	suite.NoError(err)
}

func (suite *RedisCacheTestSuite) TestGetStats() {
	stats, err := suite.cache.GetStats()
	if err != nil {
		suite.T().Skip("Redis not available, skipping test")
		return
	}
	suite.NoError(err)
	suite.True(stats.Connected)
}

// Benchmark tests
func BenchmarkRedisSet(b *testing.B) {
	redisCache := cache.NewRedisCache("localhost:6379", "", 1)
	defer redisCache.Close()

	testData := map[string]interface{}{
		"name": "Benchmark Book",
		"id":   999,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := redisCache.Set("bench:key", testData, 5*time.Minute)
		if err != nil {
			b.Skip("Redis not available")
		}
	}
}

func BenchmarkRedisGet(b *testing.B) {
	redisCache := cache.NewRedisCache("localhost:6379", "", 1)
	defer redisCache.Close()

	testData := map[string]interface{}{
		"name": "Benchmark Book",
		"id":   999,
	}

	// Setup
	err := redisCache.Set("bench:key", testData, 5*time.Minute)
	if err != nil {
		b.Skip("Redis not available")
	}

	var retrieved map[string]interface{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := redisCache.Get("bench:key", &retrieved)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Test with mock Redis when real Redis is not available
func TestRedisCache_MockScenarios(t *testing.T) {
	// Test graceful handling when Redis is not available
	redisCache := cache.NewRedisCache("localhost:9999", "", 0) // Invalid port
	defer redisCache.Close()

	// These should not panic, but return errors
	err := redisCache.Set("test", "value", time.Minute)
	assert.Error(t, err)

	var value string
	err = redisCache.Get("test", &value)
	assert.Error(t, err)

	err = redisCache.Ping()
	assert.Error(t, err)
}

func TestRedisCacheTestSuite(t *testing.T) {
	suite.Run(t, new(RedisCacheTestSuite))
}
