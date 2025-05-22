package system

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	jwtgo "github.com/xiaohongshu/PnSql/server/common/jwt"
	"github.com/xiaohongshu/PnSql/server/common/response"
	"github.com/xiaohongshu/PnSql/server/utils"
	"time"
)

type Login struct{}

type LoginParam struct {
	Username string `json:"account" form:"account"`
	Password string `json:"password" form:"password"`
	Id       string `json:"id" form:"id"`
	Code     string `json:"code" form:"code"`
}

var store = base64Captcha.DefaultMemStore

// 生成验证码
func (l *Login) Login(c *gin.Context) {
	digit := &base64Captcha.DriverDigit{Height: 80, Width: 240, Length: 6, MaxSkew: 0.8, DotCount: 120}
	captcha := base64Captcha.NewCaptcha(digit, store)
	id, baseURL, _, err := captcha.Generate()
	if err != nil {
		response.FailWithMessage("验证生成错误", c)
		return
	}
	response.OkWithData(map[string]any{"id": id, "baseURL": baseURL}, c)
}

// 登录接口
func (l *Login) ToLogin(c *gin.Context) {
	param := LoginParam{}
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithMessage("参数格式错误", c)
		return
	}

	//校验验证码
	//verify := store.Verify(param.Id, param.Code, true)
	//if !verify {
	//	response.FailWithMessage("验证码错误", c)
	//	return
	//}

	token := funcJwt(c, param)
	judgeUser, err := exampleLogin.SysUserService.GetSysUsersAccount(param.Username)
	if err != nil {
		response.FailWithMessage("校验用户错误", c)
		return
	}

	if judgeUser != nil {
		if judgeUser.Password == utils.Md5(param.Password) {
			if judgeUser.Enable == 1 {
				response.OkWithData(map[string]any{"user": param.Username, "token": token}, c)
				return
			} else {
				response.FailWithMessage("用户被冻结", c)
				return
			}
		}
	}
	response.FailWithMessage("用户或密码错误", c)
}

func funcJwt(c *gin.Context, param LoginParam) string {
	var bf, _ = utils.ParseDuration("2d") //token 缓存时间
	var ep, _ = utils.ParseDuration("7d") //token过期时间

	jwt := jwtgo.NewJWT()
	//生成 token
	var j = jwtgo.NewCustomClaims(
		param.Id, param.Username,
		int64(bf/time.Second),
		time.Now().Unix(),
		time.Now().Add(-1000*time.Second).Unix(),
		time.Now().Add(ep).Unix(),
	)
	token, _ := jwt.CreateToken(j)
	return token
}
