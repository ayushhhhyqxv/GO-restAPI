package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func main() {
	err := godotenv.Load()
	if err!=nil {
		log.Fatalf("Error while reading env file: %v ",err)
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortstr := os.Getenv("SMTP_PORT")
	senderEmail := os.Getenv("SENDER_EMAIL")
	senderName := os.Getenv("SENDER_NAME")
	senderPassword := os.Getenv("SENDER_PASSWORD")
	receiverEmail := os.Getenv("RECEIVER_EMAIL")
	attachment := `localfilepath`

	smptPort,err := strconv.Atoi(smtpPortstr)
	if err!=nil{
		log.Fatalf("Invalid SMTP port")
	}

	m := gomail.NewMessage()
	m.SetHeader("From",m.FormatAddress(senderEmail,senderName))
	m.SetHeader("To",receiverEmail)
	m.SetHeader("Subject","Email with attachment ! ")

	m.SetBody("text/html",`
		<h3> Prototype Creation </h3>
	`)

	if attachment!=""{
		m.Attach(attachment)
	}

	d:= gomail.NewDialer(smtpHost,smptPort,senderEmail,senderPassword)

	if err:= d.DialAndSend(m);err!=nil{
		log.Fatalf("Failed to send mail: %v",err)
	}

	fmt.Println("Email was sent successfully")
}