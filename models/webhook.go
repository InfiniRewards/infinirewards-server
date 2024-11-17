package models

type WebhookResponse struct {
	Object string  `json:"object,omitempty"`
	Entry  []Entry `json:"entry,omitempty"`
}

type Entry struct {
	ID      string   `json:"id,omitempty"`
	Changes []Change `json:"changes,omitempty"`
}

type Change struct {
	Value ChangeValue `json:"value,omitempty"`
	Field string      `json:"field,omitempty"`
}

type ChangeValue struct {
	Contacts         []Contact `json:"contacts,omitempty"`
	Errors           []Error   `json:"errors,omitempty"`
	MessagingProduct string    `json:"messaging_product,omitempty"`
	Metadata         Metadata  `json:"metadata,omitempty"`
	Messages         []Message `json:"messages,omitempty"`
	Statuses         []Status  `json:"statuses,omitempty"`
}

type Contact struct {
	Profile Profile `json:"profile,omitempty"`
	WaID    string  `json:"wa_id,omitempty"`
}

type Profile struct {
	Name string `json:"name,omitempty"`
}

type Error struct {
	Code      int       `json:"code,omitempty"`
	Title     string    `json:"title,omitempty"`
	Message   string    `json:"message,omitempty"`
	ErrorData ErrorData `json:"error_data,omitempty"`
}

type ErrorData struct {
	Details string `json:"details,omitempty"`
}

type Metadata struct {
	DisplayPhoneNumber string `json:"display_phone_number,omitempty"`
	PhoneNumberID      string `json:"phone_number_id,omitempty"`
}

type Message struct {
	From      string `json:"from,omitempty"`
	ID        string `json:"id,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	Text      Text   `json:"text,omitempty"`
	Type      string `json:"type,omitempty"`
}

type Text struct {
	Body string `json:"body,omitempty"`
}

type Status struct {
	ID                    string       `json:"id,omitempty"`
	BizOpaqueCallbackData string       `json:"biz_opaque_callback_data,omitempty"`
	Conversation          Conversation `json:"conversation,omitempty"`
	Errors                []Error      `json:"errors,omitempty"`
	Pricing               Pricing      `json:"pricing,omitempty"`
	RecipientID           string       `json:"recipient_id,omitempty"`
	Status                string       `json:"status,omitempty"`
	Timestamp             uint64       `json:"timestamp,omitempty"`
}

type Conversation struct {
	ID                  string `json:"id,omitempty"`
	Origin              Origin `json:"origin,omitempty"`
	ExpirationTimestamp uint64 `json:"expiration_timestamp,omitempty"`
}

type Origin struct {
	Type string `json:"type,omitempty"`
}

type Pricing struct {
	Category     string `json:"category,omitempty"`
	PricingModel string `json:"pricing_model,omitempty"`
}
