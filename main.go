package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os"
	"strings"
)

type Content struct {
	Rendered string `json:"rendered"`
}

type Item struct {
	Content Content `json:"content"`
	Date string `json:"date"`
}


const (
	urlToFetch = "https://www.srbvoz.rs/wp-json/wp/v2/info_post?per_page=10"
	hashFile   = "lasthash.txt"
)

var (
	keywords = []string{"zemuna", "lazarevac", "mladenovac"}

	smtpHost = "smtp.gmail.com"
	smtpPort = "587"


	senderEmail    = os.Getenv("SMTP_ADDR")
	senderPassword = os.Getenv("SMTP_PASS")

	recipientEmail = senderEmail
)

func containsKeyword(s string, keywords []string) bool {
	s = strings.ToLower(s)
	for _, kw := range keywords {
		if strings.Contains(s, strings.ToLower(kw)) {
			return true
		}
	}
	return false
}

func loadLastHash() string {
	b, err := ioutil.ReadFile(hashFile)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

func saveLastHash(hash string) error {
	return ioutil.WriteFile(hashFile, []byte(hash), 0644)
}


func sendEmail(subject, message, date string) error {
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)
	message = strings.Replace(message, "<p>", "", -1)
	message = strings.Replace(message, "</p>", "", -1)

	msg := "From: " + senderEmail + "\r\n" +
		"To: " + recipientEmail + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		message + "\r\n"+
		date +"\r\n"

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{recipientEmail}, []byte(msg))
}

func main() {
	resp, err := http.Get(urlToFetch)
	if err != nil {
		fmt.Println("Fetch error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	var items []Item
	if err := json.Unmarshal(body, &items); err != nil {
		fmt.Println("JSON error:", err)
		return
	}

	lastHash := loadLastHash()
	sum := sha256.Sum256(body)
	h := fmt.Sprintf("%x", sum)
	if h != lastHash {

		for _, item := range items {
			content := item.Content.Rendered
			if containsKeyword(content, keywords) {
				date:=item.Date
				if err := sendEmail("Proveri Obavestenja Srbija Voz",content,date); err != nil {
					fmt.Println("Email error:", err)
					return
				}

				fmt.Println("Email sent")
				break  // Stop after sending one new email
			}
		}
		if err := saveLastHash(h); err != nil {
		fmt.Println("Save hash error:", err)}
	}
}
