package http

type SlackRequest struct {
	Message string `validate:"required" json:"message"`
}

type SMSRequest struct {
	Message  string `validate:"required" json:"message"`
	Receiver string `validate:"required,e164" json:"receiver"`
}

type EmailRequest struct {
	Message  string `validate:"required" json:"message"`
	Receiver string `validate:"required,email" json:"receiver"`
}
