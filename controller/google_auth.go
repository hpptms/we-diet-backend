package controller

import (
	"context"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2api "google.golang.org/api/oauth2/v2"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"my-gin-app/database"
	"my-gin-app/database/model"
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_OAUTH_REDIRECT_URL"),
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
	oauthStateString = "random" // 本番ではランダム生成・検証推奨
)

func GoogleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	state := c.Query("state")
	if state != oauthStateString {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OAuth state"})
		return
	}

	code := c.Query("code")
	_, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	oauth2Service, err := oauth2api.NewService(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create oauth2 service"})
		return
	}

	userinfo, err := oauth2api.NewUserinfoService(oauth2Service).Get().Do()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	db := database.GetDB()
	var user model.User
	result := db.Where("google_id = ?", userinfo.Id).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		// 新規ユーザー作成
		user = model.User{
			GoogleID: userinfo.Id,
			Email:    userinfo.Email,
			Name:     userinfo.Name,
			Picture:  userinfo.Picture,
		}
		db.Create(&user)
	} else if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}

	// セッションやJWT発行処理（ここでは簡易的にユーザー情報返却）
	c.JSON(http.StatusOK, gin.H{
		"message": "Google login success",
		"user": gin.H{
			"id":      user.ID,
			"name":    user.Name,
			"email":   user.Email,
			"picture": user.Picture,
		},
	})
}
