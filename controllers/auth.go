package controllers

import (
	"fmt"
	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
	myjwt "traffic-go/middleware/jwt"
	"traffic-go/models"
)

type UserInfo struct{
	Username string `json:"username"`
	Password string `json:"password"`
}
type LoginResult struct{
	Token string `json:"token"`
	models.User
}
func RegisterUser(c *gin.Context){
	var userInfo UserInfo
	if err := c.ShouldBindJSON(&userInfo);err !=nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 400, "error": err.Error(),"success":false})
		return
	}
	user := models.User{Username: userInfo.Username,Password: userInfo.Password}
	if err := models.DB.Where("username = ?", userInfo.Username).First(&user).Error; err != nil { //检验用户名是否重复
		models.DB.Create(&user)
		c.JSON(http.StatusOK,gin.H{"status":200,"success":true,"msg":"注册成功"})
		return
	}else{
		c.JSON(http.StatusAccepted,gin.H{"status":202,"success":false,"msg":"账号已占用"})
	}

}

func LoginUser(c *gin.Context){ //jwt验证
	var userInfo UserInfo
	if err :=c.ShouldBindJSON(&userInfo);err !=nil {
		c.JSON(http.StatusBadRequest,gin.H{"status": 400, "msg": err.Error(),"success":false})
		return
	}
	var user models.User
	models.DB.Where("username = ?",userInfo.Username).First(&user)
	if user.Password == userInfo.Password{
		generateToken(c, user)
	}else {
		c.JSON(http.StatusOK,gin.H{"status":200,"msg":"密码不正确，请重新输入","success":false})
		return
	}
}
func generateToken(c *gin.Context,user models.User){
	j := &myjwt.JWT{
		[]byte("newtrekWang"), // 有疑问
	}
	claims := myjwt.CustomClaims{
		strconv.FormatInt(user.ID,10),
		user.Username,
		user.Password,
		jwtgo.StandardClaims{
			NotBefore: int64(time.Now().Unix() - 1000), // 签名生效时间
			ExpiresAt: int64(time.Now().Unix() + 3600), // 过期时间 一小时
			Issuer:    "newtrekWang",                   //签名的发行者
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    err.Error(),
		})
		return
	}
	fmt.Println(token)
	data := LoginResult{
		User:  user,
		Token: token,
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
		"msg":    "登录成功！",
		"data":   data,
		"success":true,
	})
	return
}

// 测试接口
func GetDataByTime(c *gin.Context) {
	claims := c.MustGet("claims").(*myjwt.CustomClaims)
	if claims != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 0,
			"msg":    "token有效",
			"data":   claims,
			"success":true,
		})
	}
}
