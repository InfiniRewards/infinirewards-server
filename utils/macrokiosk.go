package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var BOLD_USER string
var BOLD_PASS string
var BOLD_SERVID string
var MOCK_OTP_DISCORD_URL string

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	BOLD_USER = os.Getenv("BOLD_USER")
	BOLD_PASS = os.Getenv("BOLD_PASS")
	BOLD_SERVID = os.Getenv("BOLD_SERVID")
	MOCK_OTP_DISCORD_URL = os.Getenv("MOCK_OTP_DISCORD_URL")
}

func SendSMS(phoneNumber string, code string) error {
	if phoneNumber == "+60143382537" {
		jsonData, err := json.Marshal(struct {
			UserName string `json:"username"`
			Content  string `json:"content"`
		}{UserName: phoneNumber, Content: fmt.Sprintf("JOMFI: Thanks for registering JomFi. Your OTP is %s. Valid for 15 minutes.", code)})
		if err != nil {
			return err
		}
		request, err := http.NewRequest("POST", MOCK_OTP_DISCORD_URL, bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}
		request.Header.Set("Content-Type", "application/json; charset=UTF-8")
		netClient := &http.Client{
			Timeout: time.Second * 10,
		}
		resp, err := netClient.Do(request)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		return nil
	}
	httpposturl := "https://www.etracker.cc/bulksms/send"

	var jsonData = []byte(fmt.Sprintf(`{
		"user":%s,
    "pass":%s,
    "type": "0",
    "to":"%s",
    "from":"JomFi",
    "text": "JOMFI: Thanks for registering JomFi. Your OTP is %s. Valid for 15 minutes.",
    "servid": %s,
    "title":"JomFi",
    "details":"1"
	}`, BOLD_USER, BOLD_PASS, strings.Split(phoneNumber, "+")[1], code, BOLD_SERVID))
	request, err := http.NewRequest("POST", httpposturl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := io.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
	return nil
}
