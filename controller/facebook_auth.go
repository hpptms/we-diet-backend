package controller

import (
	"context"
	"encoding/json"
	"my-gin-app/database"
	"my-gin-app/database/model"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	facebookOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("FACEBOOK_OAUTH_REDIRECT_URL"),
		ClientID:     os.Getenv("FACEBOOK_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("FACEBOOK_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"email", "public_profile"},
		Endpoint:     facebook.Endpoint,
	}
	facebookOauthStateString = "random" // 本番ではランダム生成・検証推奨
)

type facebookUserInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
}

func FacebookLogin(c *gin.Context) {
	url := facebookOauthConfig.AuthCodeURL(facebookOauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func FacebookCallback(c *gin.Context) {
	state := c.Query("state")
	if state != facebookOauthStateString {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OAuth state"})
		return
	}

	code := c.Query("code")
	token, err := facebookOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	client := facebookOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://graph.facebook.com/me?fields=id,name,email,picture")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var userInfo facebookUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user info"})
		return
	}

	db := database.GetDB()
	var user model.User
	result := db.Where("facebook_id = ?", userInfo.ID).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		// 新規ユーザー作成
		user = model.User{
			FacebookID:      userInfo.ID,
			FacebookEmail:   userInfo.Email,
			FacebookName:    userInfo.Name,
			FacebookPicture: userInfo.Picture.Data.URL,
		}
		db.Create(&user)
	} else if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}

	// セッションやJWT発行処理（ここでは簡易的にユーザー情報返却）
	c.JSON(http.StatusOK, gin.H{
		"message": "Facebook login success",
		"user": gin.H{
			"id":      user.ID,
			"name":    user.FacebookName,
			"email":   user.FacebookEmail,
			"picture": user.FacebookPicture,
		},
	})
}
