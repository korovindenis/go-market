package handler

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetUserIDFromCtx(t *testing.T) {
	handler := &Handler{}
	ctx := &gin.Context{}
	userID := int64(42)
	ctx.Set("userId", userID)

	resultUserID, err := handler.getUserIDFromCtx(ctx)

	assert.Nil(t, err)
	assert.Equal(t, userID, resultUserID)
}

func TestGetUserIDFromCtxMissingKey(t *testing.T) {
	handler := &Handler{}
	ctx := &gin.Context{}

	resultUserID, err := handler.getUserIDFromCtx(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, int64(0), resultUserID)
}

func TestGetUserIDFromCtxInvalidType(t *testing.T) {
	handler := &Handler{}
	ctx := &gin.Context{}
	ctx.Set("userId", "not an int64")

	resultUserID, err := handler.getUserIDFromCtx(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, int64(0), resultUserID)
}
