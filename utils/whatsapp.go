package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var token string
var phoneId string

type SendWithTemplateRequest struct {
	MessagingProduct string   `json:"messaging_product,omitempty"`
	To               string   `json:"to,omitempty"`
	Type             string   `json:"type,omitempty"`
	Template         Template `json:"template,omitempty"`
}

type TemplateLanguage struct {
	Code string `json:"code,omitempty"`
}
type TemplateParameters struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}

type Components struct {
	Type       string               `json:"type,omitempty"`
	Parameters []TemplateParameters `json:"parameters,omitempty"`
}

type Template struct {
	Name       string           `json:"name,omitempty"`
	Language   TemplateLanguage `json:"language,omitempty"`
	Components []Components     `json:"components,omitempty"`
}

type SendWithInteractiveRequest struct {
	MessagingProduct string      `json:"messaging_product,omitempty"`
	To               string      `json:"to,omitempty"`
	Type             string      `json:"type,omitempty"`
	Interactive      Interactive `json:"interactive,omitempty"`
}

type Interactive struct {
	Type   string            `json:"type,omitempty"`
	Header InteractiveHeader `json:"header,omitempty"`
	Body   InteractiveBody   `json:"body,omitempty"`
	Footer InteractiveFooter `json:"footer,omitempty"`
	Action InteractiveAction `json:"action,omitempty"`
}

type InteractiveHeader struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}

type InteractiveBody struct {
	Text string `json:"text,omitempty"`
}

type InteractiveFooter struct {
	Text string `json:"text,omitempty"`
}

type Button struct {
	Type  string `json:"type,omitempty"`
	Reply Reply  `json:"reply,omitempty"`
}

type Reply struct {
	ID    string `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
}

type CTAParameters struct {
	DisplayText string `json:"display_text,omitempty"`
	Url         string `json:"url,omitempty"`
}

type InteractiveAction struct {
	Button     string        `json:"button,omitempty"`
	Buttons    []Button      `json:"buttons,omitempty"`
	Name       string        `json:"name,omitempty"`
	Parameters CTAParameters `json:"parameters,omitempty"`
	Sections   []Section     `json:"sections,omitempty"`
}

type Section struct {
	Title string `json:"title,omitempty"`
	Rows  []Row  `json:"rows,omitempty"`
}

type Row struct {
	ID       string `json:"id,omitempty"`
	Title    string `json:"title,omitempty"`
	Metadata string `json:"description,omitempty"`
}

func InitWhatsApp() error {
	err := godotenv.Load(".env")

	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}
	phoneId = os.Getenv("WHATSAPP_PHONE_ID")
	token = os.Getenv("WHATSAPP_TOKEN")

	return nil
}

func SendWhatsappLoginURL(to string, token string) (res map[string]interface{}, err error) {
	return sendMessage(SendWithInteractiveRequest{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "interactive",
		Interactive: Interactive{
			Type:   "cta_url",
			Header: InteractiveHeader{Type: "text", Text: "Login"},
			Body:   InteractiveBody{Text: "Click the link below to access Beautifood"},
			Footer: InteractiveFooter{Text: "Beautifood App"},
			Action: InteractiveAction{
				Name: "cta_url",
				Parameters: CTAParameters{
					DisplayText: "Login",
					Url:         "https://beautifood.io",
				},
			},
		},
	})
}

func sendMessage(request any) (res map[string]interface{}, err error) {

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return res, err
	}

	body := bytes.NewReader(jsonRequest)

	endpoint := fmt.Sprintf("https://graph.facebook.com/%s/%s/messages", "v20.0", phoneId)
	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return res, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return res, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err := parseHTTPError(resp.Body)
		return res, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	err = json.Unmarshal(bodyBytes, &req)

	var b bytes.Buffer
	_, err = io.Copy(&b, resp.Body)

	if err != nil {
		return res, err
	}
	err = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&res)

	if err != nil {
		return res, err
	}

	return res, err
}

func parseHTTPError(body io.Reader) (err error) {
	var errRes map[string]map[string]interface{}
	err = json.NewDecoder(body).Decode(&errRes)
	if err != nil {
		return fmt.Errorf("unparsed error message")
	}
	msg := fmt.Sprintf("%s", errRes["error"]["message"])
	return errors.New(msg)
}

func IsWhatsAppUser(phoneNumber string) (bool, error) {
	endpoint := fmt.Sprintf("https://graph.facebook.com/%s/%s/messages", "v20.0", phoneId)

	// Create a test message request
	request := SendWithTemplateRequest{
		MessagingProduct: "whatsapp",
		To:               phoneNumber,
		Type:             "template",
		Template: Template{
			Name:     "hello_world", // Use a simple template that exists in your WhatsApp Business account
			Language: TemplateLanguage{Code: "en"},
		},
	}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(jsonRequest))
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// If status is 200, the number exists on WhatsApp
	// If we get an error about the number not being registered, the user is not on WhatsApp
	if resp.StatusCode == 200 {
		return true, nil
	}

	var errRes map[string]map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&errRes); err != nil {
		return false, err
	}

	// Check if error is specifically about number not being registered on WhatsApp
	if errMsg, ok := errRes["error"]["message"].(string); ok {
		if strings.Contains(strings.ToLower(errMsg), "not a whatsapp user") {
			return false, nil
		}
	}

	return false, fmt.Errorf("failed to check WhatsApp status: %v", errRes["error"]["message"])
}

func SendWhatsAppOTP(to string, code string) error {
	_, err := sendMessage(SendWithTemplateRequest{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "template",
		Template: Template{
			Name:     "otp_message", // Make sure this template is approved in your WhatsApp Business account
			Language: TemplateLanguage{Code: "en"},
			Components: []Components{
				{
					Type: "body",
					Parameters: []TemplateParameters{
						{
							Type: "text",
							Text: code,
						},
					},
				},
			},
		},
	})
	return err
}
