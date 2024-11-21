package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var BOLD_USER string
var BOLD_PASS string
var BOLD_SERVID string

func InitMacroKiosk() error {
	err := godotenv.Load(".env")

	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	BOLD_USER = os.Getenv("BOLD_USER")
	BOLD_PASS = os.Getenv("BOLD_PASS")
	BOLD_SERVID = os.Getenv("BOLD_SERVID")

	return nil
}

func SendSMS(phoneNumber string, code string) error {
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
