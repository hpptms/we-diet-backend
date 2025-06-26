package controller

import (
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"

	pb "my-gin-app/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

func MailRegister(c *gin.Context) {
	// Protobufバイナリでリクエストを受け取る
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Data(http.StatusBadRequest, "application/x-protobuf", []byte{})
		return
	}
	var req pb.MailRegisterRequest
	if err = proto.Unmarshal(body, &req); err != nil || req.Email == "" {
		// エラー時はprotobufで返す
		res := &pb.MailRegisterResponse{Message: "メールアドレスが必要です"}
		data, _ := proto.Marshal(res)
		c.Data(http.StatusBadRequest, "application/x-protobuf", data)
		return
	}

	// リダイレクト用URLを生成（仮: 本番はトークン等を含める）
	redirectURL := "http://localhost:5173/register/complete?email=" + req.Email

	// メール送信処理（SMTP）
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	mailFrom := os.Getenv("MAIL_FROM")

	// メール本文
	subject := "【サービス登録】メールアドレス確認"
	bodyStr := "メールアドレス確認用リンク: " + redirectURL
	msg := "From: " + mailFrom + "\r\n" +
		"To: " + req.Email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n\r\n" +
		bodyStr

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, mailFrom, []string{req.Email}, []byte(msg))
	if err != nil {
		// エラー内容を標準出力に出す
		log.Printf("smtp.SendMail error: %v", err)
		panic(err)
		res := &pb.MailRegisterResponse{Message: "メール送信に失敗しました"}
		data, _ := proto.Marshal(res)
		c.Data(http.StatusInternalServerError, "application/x-protobuf", data)
		return
	}

	res := &pb.MailRegisterResponse{Message: "確認メールを送信しました。メールをご確認ください。"}
	data, _ := proto.Marshal(res)
	c.Data(http.StatusOK, "application/x-protobuf", data)
}
