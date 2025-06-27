package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/oauth2"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"my-gin-app/database"
	"my-gin-app/database/model"
)

var (
	tiktokOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("TIKTOK_OAUTH_REDIRECT_URL"),
		ClientID:     os.Getenv("TIKTOK_OAUTH_CLIENT_KEY"),
		ClientSecret: os.Getenv("TIKTOK_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"user.info.basic"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.tiktok.com/v2/auth/authorize/",
			TokenURL: "https://open-api.tiktok.com/oauth/access_token/",
		},
	}
	tiktokOauthStateString = "random" // 本番ではランダム生成・検証推奨
)

type tiktokUserInfo struct {
	Data struct {
		OpenID   string `json:"open_id"`
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	} `json:"data"`
}

func TikTokLogin(c *gin.Context) {
	url := tiktokOauthConfig.AuthCodeURL(tiktokOauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func TikTokCallback(c *gin.Context) {
	state := c.Query("state")
	if state != tiktokOauthStateString {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OAuth state"})
		return
	}

	code := c.Query("code")
	token, err := tiktokOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	client := tiktokOauthConfig.Client(context.Background(), token)
	// TikTokのユーザー情報API
	resp, err := client.Get("https://open-api.tiktok.com/oauth/userinfo/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var userInfo tiktokUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user info"})
		return
	}

	db := database.GetDB()
	// OpenIDをintに変換（失敗時は0）
	tiktokIDInt := 0
	if idInt, err := strconv.Atoi(userInfo.Data.OpenID); err == nil {
		tiktokIDInt = idInt
	}

	// Userテーブルに登録
	var user model.User
	result := db.Where("service_name = ? AND service_id = ?", "tiktok", tiktokIDInt).First(&user)
	if result.Error == gorm.ErrRecordNotFound {
		user = model.User{
			UserName:    userInfo.Data.Nickname,
			Password:    "",
			Subscribe:   false,
			Permission:  0,
			Picture:     userInfo.Data.Avatar,
			ServiceName: "tiktok",
			ServiceID:   tiktokIDInt,
		}
		db.Create(&user)
	} else if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}

	// TiktokUserテーブルに登録
	var tUser model.TikTokUser
	tResult := db.Where("tik_tok_id = ?", userInfo.Data.OpenID).First(&tUser)
	if tResult.Error == gorm.ErrRecordNotFound {
		tUser = model.TikTokUser{
			TikTokID:     userInfo.Data.OpenID,
			TikTokName:   userInfo.Data.Nickname,
			TikTokAvatar: userInfo.Data.Avatar,
		}
		db.Create(&tUser)
	} else if tResult.Error != nil && tResult.Error != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}

	// セッションやJWT発行処理（ここでは簡易的にユーザー情報返却）
	c.JSON(http.StatusOK, gin.H{
		"message": "TikTok login success",
		"user": gin.H{
			"id":     user.ID,
			"name":   user.UserName,
			"avatar": user.Picture,
		},
	})
}
