package middlewares

import (
	"net/http"
	"strconv"
	"strings"
	"testproject/task004/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			utils.SendError(c, http.StatusUnauthorized, "Authorization header required")
			return
		}

		token, err := utils.ValidateToken(strings.Split(tokenString, "Bearer ")[1])
		if err != nil || !token.Valid {
			utils.SendError(c, http.StatusUnauthorized, "Invalid token")
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userid, resolveUserErr := strconv.ParseUint(claims["id"].(string), 10, 32)
		if resolveUserErr != nil {
			utils.SendError(c, http.StatusUnauthorized, "reslove id err")
			return
		}
		c.Set("userID", userid)
		c.Next()
	}
}
