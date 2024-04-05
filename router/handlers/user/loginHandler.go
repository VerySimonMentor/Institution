package user

import "github.com/gin-gonic/gin"

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(ctx *gin.Context) {

}
