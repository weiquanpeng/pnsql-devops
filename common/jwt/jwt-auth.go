package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaohongshu/PnSql/server/common/response"
	"github.com/xiaohongshu/PnSql/server/dao/service"
	"github.com/xiaohongshu/PnSql/server/global"
)

var jwtService = service.GroupApp.SystemServer.Jwt

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.FailWithData(801, "", "请求未携带 token", c)
			c.Abort()
			return
		}
		if token == global.P_cfg.System.Token {
			// 你可以自定义 claims，这里示例直接放字符串
			c.Set("claims", "universal-token")
			c.Next()
			return
		}

		// 黑名单校验
		judgmentJwt, err := jwtService.GetJwtCount(token)
		if judgmentJwt > 0 {
			response.FailWithData(801, "", "token 已被加入黑名单", c)
			c.Abort()
			return
		}

		myJwt := NewJWT()
		claims, err := myJwt.ParseToken(token)
		if err != nil {
			response.FailWithData(801, "", "token 失效", c)
			c.Abort()
			return
		}

		////token 续期逻辑（已注释，看你需求要不要启用）
		//if claims.ExpiresAt-time.Now().Unix() < claims.BufferTime {
		//	newTime, _ := utils.ParseDuration("7d")
		//	claims.ExpiresAt = time.Now().Add(newTime).Unix()
		//	newToken, _ := myJwt.CreateTokenByOldToken(token, *claims)
		//	c.Header("PanGu-token", newToken)
		//	err := jwtService.CreateJwt(&system.Jwt{Jwt: token})
		//	if err != nil {
		//		return
		//	}
		//}

		c.Set("claims", claims)
		c.Next()
	}
}
