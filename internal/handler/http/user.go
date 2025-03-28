package http

import (
	"chat/internal/model"
	"chat/internal/service"
	myerrors "chat/pkg/errors"
	"chat/pkg/response"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

// UserRegister 用户注册
func (h *UserHandler) UserRegister(c *gin.Context) {
	var user model.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, response.Fail(myerrors.ErrCodeParamInvalid, "参数错误"))
		return
	}
	if err := h.UserService.Register(&user); err != nil {
		if errors.Is(myerrors.ErrUserExisted, err) {
			c.JSON(http.StatusBadRequest, response.Fail(myerrors.ErrCodeUserExisted, "用户已存在"))
			return
		} else {
			c.JSON(http.StatusBadRequest, response.Fail(myerrors.ErrCodeRegisterFail, "注册失败"))
			return
		}
	}
	c.JSON(http.StatusOK, response.Success("注册成功"))

}
