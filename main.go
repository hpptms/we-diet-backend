package main

import (
	"fmt"
	"os"
	"time"

	"my-gin-app/controller"
	"my-gin-app/database"
	migrate "my-gin-app/database/migrate"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var identityKey = "id"

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func jwtMiddleware() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "example zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*database.User); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &database.User{
				UserName: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBindWith(&loginVals, binding.JSON); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			// ここでユーザー認証ロジックを実装（例: DB照会）
			if userID == "admin" && password == "password" {
				return &database.User{
					UserName: userID,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// ここでパーミッション管理ロジックを実装
			if v, ok := data.(*database.User); ok && v.UserName == "admin" {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
}

func main() {
	// .env読み込み
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(".envファイルの読み込みに失敗しました")
	}

	// DB接続
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tokyo",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("DB接続失敗: " + err.Error())
	}
	database.SetDB(db)

	// マイグレーション
	if err := migrate.Migrate(db); err != nil {
		panic("DBマイグレーション失敗: " + err.Error())
	}

	r := gin.Default()

	authMiddleware, err := jwtMiddleware()
	if err != nil {
		panic("JWT Error:" + err.Error())
	}

	// ログインエンドポイント
	r.POST("/login", authMiddleware.LoginHandler)

	// Googleログインエンドポイント
	r.GET("/auth/google/login", controller.GoogleLogin)
	r.GET("/auth/google/callback", controller.GoogleCallback)

	// Facebookログインエンドポイント
	r.GET("/auth/facebook/login", controller.FacebookLogin)
	r.GET("/auth/facebook/callback", controller.FacebookCallback)

	// TikTokログインエンドポイント
	r.GET("/auth/tiktok/login", controller.TikTokLogin)
	r.GET("/auth/tiktok/callback", controller.TikTokCallback)

	// メール登録API
	r.POST("/register/mail", controller.MailRegister)

	// 認証が必要なAPIグループ
	auth := r.Group("/api")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", func(c *gin.Context) {
			claims := jwt.ExtractClaims(c)
			user, _ := c.Get(identityKey)
			c.JSON(200, gin.H{
				"userID":   claims[identityKey],
				"userName": user.(*database.User).UserName,
				"text":     "Hello from Gin API!",
			})
		})
	}

	// Reactのビルドファイルを /static で配信
	r.Static("/static", "../frontend/build/static")

	// ルートやその他のパスは index.html を返す
	r.NoRoute(func(c *gin.Context) {
		c.File("../frontend/build/index.html")
	})

	r.Run("0.0.0.0:8080")
}
