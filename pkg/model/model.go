package model

type Email struct {
	From string `json:",omitempty"`
	// To: REQUIRED Recipient email address. Multiple addresses are comma separated. Max 50.
	To string `json:",omitempty"`
	// Cc recipient email address. Multiple addresses are comma separated. Max 50.
	Cc string `json:",omitempty"`
	// Bcc recipient email address. Multiple addresses are comma separated. Max 50.
	Bcc string `json:",omitempty"`
	// Subject: Email subject
	Subject string `json:",omitempty"`
	// Tag: Email tag that allows you to categorize outgoing emails and get detailed statistics.
	Tag string `json:",omitempty"`
	// HtmlBody: HTML email message. REQUIRED, If no TextBody specified
	HtmlBody string `json:",omitempty"`
	// TextBody: Plain text email message. REQUIRED, If no HtmlBody specified
	TextBody string `json:",omitempty"`
	// ReplyTo: Reply To override email address. Defaults to the Reply To set in the sender signature.
	ReplyTo string `json:",omitempty"`
	// Headers: List of custom headers to include.
	Headers []Header `json:",omitempty"`
	// TrackOpens: Activate open tracking for this email.
	TrackOpens bool `json:",omitempty"`
	// Attachments: List of attachments
	Attachments []Attachment `json:",omitempty"`
	// Metadata: metadata
	Metadata map[string]string `json:",omitempty"`
}

// Header - an email header
type Header struct {
	// Name: header name
	Name string
	// Value: header value
	Value string
}

// Attachment is an optional encoded file to send along with an email
type Attachment struct {
	// Name: attachment name
	Name string
	// Content: Base64 encoded attachment data
	Content string
	// ContentType: attachment MIME type
	ContentType string
	// ContentId: populate for inlining images with the images cid
	ContentID string `json:",omitempty"`
}
