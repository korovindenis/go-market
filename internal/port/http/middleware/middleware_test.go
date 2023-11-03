package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-market/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

type mockConfig struct{}

func (c *mockConfig) GetTokenName() string {
	return "token"
}

type mockAuth struct{}

func (a *mockAuth) CheckToken(user entity.User, tokenString string) error {
	return nil
}

func (a *mockAuth) GetUserFromToken(tokenString string) (entity.User, error) {
	return entity.User{}, nil
}

func setupGinTest() *gin.Context {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = &http.Request{}
	ctx.Request.Header = make(http.Header)
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Set("config", &mockConfig{})
	ctx.Set("auth", &mockAuth{})
	return ctx
}

func TestCheckMethod(t *testing.T) {
	ctx := setupGinTest()
	ctx.Request.Method = http.MethodPost

	middleware, _ := New(&mockConfig{}, &mockAuth{})
	handler := middleware.CheckMethod()

	handler(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())
}

func TestCheckContentTypeJSON(t *testing.T) {
	ctx := setupGinTest()

	middleware, _ := New(&mockConfig{}, &mockAuth{})
	handler := middleware.CheckContentTypeJSON()

	handler(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())
}

func TestCheckContentTypeText(t *testing.T) {
	ctx := setupGinTest()
	ctx.Request.Header.Set("Content-Type", "text/plain")

	middleware, _ := New(&mockConfig{}, &mockAuth{})
	handler := middleware.CheckContentTypeText()

	handler(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())
}

func TestCheckAuth(t *testing.T) {
	ctx := setupGinTest()
	ctx.Request.AddCookie(&http.Cookie{Name: "token", Value: "testToken"})

	middleware, _ := New(&mockConfig{}, &mockAuth{})
	handler := middleware.CheckAuth()

	handler(ctx)

	assert.Equal(t, http.StatusOK, ctx.Writer.Status())
}

func TestAddUserInfoToCtx(t *testing.T) {
	ctx := setupGinTest()
	ctx.Request.AddCookie(&http.Cookie{Name: "token", Value: "testToken"})

	middleware, _ := New(&mockConfig{}, &mockAuth{})
	handler := middleware.AddUserInfoToCtx()

	handler(ctx)

	userID, _ := ctx.Get("userId")
	assert.NotNil(t, userID)
}
