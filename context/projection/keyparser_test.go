package projection

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestKeyParser(t *testing.T) {

	key, err := ParseKey("123"); assert.Nil(t, err); assert.Equal(t, []string{"123"}, key)
	key, err = ParseKey("a.b"); assert.Nil(t, err); assert.Equal(t, []string{"a", "b"}, key)
	key, err = ParseKey("hello.world.123"); assert.Nil(t, err); assert.Equal(t, []string{"hello", "world", "123"}, key)

	key, err = ParseKey(""); assert.NotNil(t, err)
	key, err = ParseKey("."); assert.NotNil(t, err)
	key, err = ParseKey(".world"); assert.NotNil(t, err)
	key, err = ParseKey(".world.123"); assert.NotNil(t, err)
	key, err = ParseKey("world.12/3"); assert.NotNil(t, err)
	key, err = ParseKey("world.12."); assert.NotNil(t, err)


	key, err = ParseKey("hello[0].world.123"); assert.Nil(t, err); assert.Equal(t, []string{"hello", "0", "world", "123"}, key)
	key, err = ParseKey("hello[a][b].w[orld].123"); assert.Nil(t, err); assert.Equal(t, []string{"hello", "a", "b", "w", "orld", "123"}, key)

	key, err = ParseKey("world[0]12."); assert.NotNil(t, err)
	key, err = ParseKey("world[0]..12."); assert.NotNil(t, err)

}