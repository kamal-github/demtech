package model

type EmailRequest struct {
	Source               string      `json:"Source" binding:"required,email"`
	Destination          Destination `json:"Destination"`
	Message              Message     `json:"Message"`
	ConfigurationSetName string      `json:"ConfigurationSetName,omitempty"`
	ReplyToAddresses     []string    `json:"ReplyToAddresses,omitempty" binding:"omitempty,dive,email"`
	ReturnPath           string      `json:"ReturnPath,omitempty" binding:"email"`
	ReturnPathArn        string      `json:"ReturnPathArn,omitempty"`
	SourceArn            string      `json:"SourceArn,omitempty"`
	Tags                 []Tag       `json:"Tags,omitempty"`
}

type Destination struct {
	ToAddresses  []string `json:"ToAddresses" binding:"required,dive,email,max=50"`
	CcAddresses  []string `json:"CcAddresses" binding:"dive,email,max=50"`
	BccAddresses []string `json:"BccAddresses" binding:"dive,email,max=50"`
}

func (d Destination) All() []string {
	return append(append(d.ToAddresses, d.CcAddresses...), d.BccAddresses...)
}

type Subject struct {
	Data string `json:"Data" binding:"required"`
}

type TextBody struct {
	Data string `json:"Data" binding:"required"`
}

type Body struct {
	Text TextBody `json:"Text"`
}

type Message struct {
	Subject Subject `json:"Subject"`
	Body    Body    `json:"Body"`
}

type Tag struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
}

// SESResponse represents a successful response type
type SESResponse struct {
	MessageID string `json:"MessageId"`
}

// SESError represents an AWS-style error response
type SESError struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
}

func (err *SESError) Error() string {
	return err.Code + ": " + err.Message
}
