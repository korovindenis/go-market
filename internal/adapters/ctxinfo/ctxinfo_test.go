package ctxinfo

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetUserIDFromCtx(t *testing.T) {
	handler := &Ctxinfo{}
	ctx := &gin.Context{}
	userID := int64(42)
	ctx.Set("userId", userID)

	resultUserID, err := handler.GetUserIDFromCtx(ctx)

	assert.Nil(t, err)
	assert.Equal(t, userID, resultUserID)
}

func TestGetUserIDFromCtxMissingKey(t *testing.T) {
	handler := &Ctxinfo{}
	ctx := &gin.Context{}

	resultUserID, err := handler.GetUserIDFromCtx(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, int64(0), resultUserID)
}

func TestGetUserIDFromCtxInvalidType(t *testing.T) {
	handler := &Ctxinfo{}
	ctx := &gin.Context{}
	ctx.Set("userId", "not an int64")

	resultUserID, err := handler.GetUserIDFromCtx(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, int64(0), resultUserID)
}
