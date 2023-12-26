// Work with context
package ctxinfo

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Ctxinfo struct {
}

func New() (*Ctxinfo, error) {
	return &Ctxinfo{}, nil
}

func (u *Ctxinfo) GetUserIDFromCtx(ctx *gin.Context) (int64, error) {
	userIDRaw, ok := ctx.Get("userId")
	if !ok {
		return 0, fmt.Errorf("%s", "GetUserIDFromCtx Get UserId")
	}
	userID, ok := userIDRaw.(int64)
	if !ok {
		return 0, fmt.Errorf("%s", "GetUserIDFromCtx userIDRaw")
	}

	return userID, nil
}
