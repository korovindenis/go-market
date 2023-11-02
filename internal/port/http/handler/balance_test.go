package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/korovindenis/go-market/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestGetBalance(t *testing.T) {
	handler := &Handler{
		usecase: &mockUsecase{},
	}

	t.Run("ValidUserID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Set("userId", int64(1))

		handler.GetBalance(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var result entity.Balance
		err := c.ShouldBindJSON(&result)
		assert.Nil(t, err)
		assert.Equal(t, entity.Balance{Current: 1, Withdrawn: 100}, result)
	})

	t.Run("Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Set("userId", int64(2))

		handler.GetBalance(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestWithdrawBalance(t *testing.T) {
	handler := &Handler{
		usecase: &mockUsecase{},
	}

	t.Run("ValidOrder", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Set("userId", int64(1))

		balanceJSON := `{"order": 2, "sum": 50}`
		r := httptest.NewRequest(http.MethodPost, "/withdraw", strings.NewReader(balanceJSON))
		r.Header.Set("Content-Type", "application/json")
		c.Request = r

		handler.WithdrawBalance(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("InsufficientBalance", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Set("userId", int64(1))

		balanceJSON := `{"order": 1, "sum": 50}`
		r := httptest.NewRequest(http.MethodPost, "/withdraw", strings.NewReader(balanceJSON))
		r.Header.Set("Content-Type", "application/json")
		c.Request = r

		handler.WithdrawBalance(c)

		assert.Equal(t, http.StatusPaymentRequired, w.Code)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Set("userId", int64(1))

		balanceJSON := `{"order": "invalid", "sum": "invalid"}`
		r := httptest.NewRequest(http.MethodPost, "/withdraw", strings.NewReader(balanceJSON))
		r.Header.Set("Content-Type", "application/json")
		c.Request = r

		handler.WithdrawBalance(c)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})

	t.Run("Error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Set("userId", int64(2))

		balanceJSON := `{"order": 2, "sum": 50}`
		r := httptest.NewRequest(http.MethodPost, "/withdraw", strings.NewReader(balanceJSON))
		r.Header.Set("Content-Type", "application/json")
		c.Request = r

		handler.WithdrawBalance(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
