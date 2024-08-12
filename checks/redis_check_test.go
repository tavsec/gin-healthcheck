package checks

import (
	"errors"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisCheck_Pass(t *testing.T) {
	// create a mock redis client
	mockClient, mock := redismock.NewClientMock()
	mock.ExpectPing().SetVal("PONG")

	// create a RedisCheck instance using the mock client
	redisCheck := NewRedisCheck(mockClient)

	// call Pass() method and assert that it returns true
	assert.True(t, redisCheck.Pass())
}
func TestRedisCheck_Fail(t *testing.T) {
	// create a mock redis client
	mockClient, mock := redismock.NewClientMock()
	mock.ExpectPing().SetErr(errors.New("ping failed"))

	// create a RedisCheck instance using the mock client
	redisCheck := NewRedisCheck(mockClient)

	// call Pass() method and assert that it returns false
	assert.False(t, redisCheck.Pass())
}

func TestRedisCheck_Name(t *testing.T) {
	// create a mock redis client
	mockClient, _ := redismock.NewClientMock()

	// create a RedisCheck instance using the mock client
	redisCheck := NewRedisCheck(mockClient)

	// call Name() method and assert that it returns "redis"
	assert.Equal(t, "redis", redisCheck.Name())

}
