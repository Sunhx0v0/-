package middleware

//  认证的中间件
import (
	"fmt"
	"net/http"
	"strings"
	"webServer/middleware/webjwt"
	"webServer/models"

	"github.com/gin-gonic/gin"
)

// 用于登录认证的中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 authorization header
		tokenString := c.GetHeader("Authorization")

		fmt.Print("请求token", tokenString)

		//validate token formate   验证格式
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			c.Abort()
			return
		}

		tokenString = tokenString[7:] //截取字符    提取有效的字符

		token, claims, err := webjwt.ParseToken(tokenString)

		// 解析失败或无效token
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			c.Abort()
			return
		}

		//token通过验证, 获取claims中的UserID
		userPhone := claims.UserId
		userInfo := models.SelectAll(userPhone)

		//判断用户是否存在
		if models.IsTelephoneExists(userPhone) { // 在数据库查找手机号码是否存在
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 401,
				"data": nil,
				"msg":  "权限不足",
			})
			c.Abort() // 将请求抛弃
			return
		}

		//用户存在 将user信息写入上下文
		c.Set("user", userInfo)

		c.Next()
	}
}
