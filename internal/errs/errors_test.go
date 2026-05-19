package errs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorWrappers(t *testing.T) {
	assert.ErrorIs(t, InvalidInput("bad"), ErrInvalidInput)
	assert.ErrorIs(t, NotFound("missing"), ErrNotFound)
	assert.ErrorIs(t, Upstream("failed: %s", "x"), ErrUpstream)
	assert.ErrorIs(t, Unauthorized("nope"), ErrUnauthorized)
	assert.ErrorIs(t, Forbidden("stop"), ErrForbidden)
	assert.ErrorIs(t, Conflict("dup"), ErrConflict)
	assert.False(t, errors.Is(InvalidInput("bad"), ErrNotFound))
}
