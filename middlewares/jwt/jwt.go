package jwt

import (
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)
var (
	TokenExpired     error  = errors.New("Token is expired")
	TokenNotValidYet error  = errors.New("Token not active yet")
	TokenMalformed   error  = errors.New("That's not even a token")
	TokenInvalid     error  = errors.New("Couldn't handle this token:")
	SignKey          string = "newtrekWang"
)
func JWTAuth() gin.HandlerFunc{
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == ""{
			c.JSON(http.StatusOK,gin.H{
				"status": -1,
				"msg": "请求未携带token，无权限访问",
				"success":false,
			})
			c.Abort()
			return
		}
		log.Print("get token: ", token)
		j := NewJWT()
		claims,err := j.ParseToken(token)
		if err != nil{
			if err == TokenExpired{
				c.JSON(http.StatusOK,gin.H{
					"status":-1,
					"msg":	"授权已过期",
					"success":false,
				})
				c.Abort()
				return
			}
			// 继续交由下一个路由处理,并将解析出的信息传递下去
			c.Set("claims", claims)
		}

	}
}
// 自定义载荷 claims  在payload中 包含实体信息 额外信息等
type CustomClaims struct {
	ID    string `json:"userId"`
	Username  string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

type JWT struct {
	SigningKey []byte
}

// 新建jwt实例
func NewJWT() *JWT {
	return  &JWT{
		[]byte(GetSignKey()), // 一个比特流
	}
}

// 获取signkey
func GetSignKey() string {
	return SignKey
}

func SetSignKey(key string) string {
	SignKey = key
	return SignKey
}

//生成一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

//解析token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

//更新token
func(j *JWT) RefreshToken(tokenString string) (string,error){
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0,0)
	}
	token,err := jwt.ParseWithClaims(tokenString,&CustomClaims{},func(token *jwt.Token) (interface{},error){
		return j.SigningKey,nil
	})
	if err != nil {
		return "",err
	}
	if claims,ok := token.Claims.(*CustomClaims); ok &&  token.Valid{
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "",TokenInvalid
}
