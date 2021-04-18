package model

type Email struct {
	From string `json:"from,omitempty"`
	// To: REQUIRED Recipient email address. Multiple addresses are comma separated. Max 50.
	To string `json:"to,omitempty"`
	// Cc recipient email address. Multiple addresses are comma separated. Max 50.
	Cc string `json:"cc,omitempty"`
	// Bcc recipient email address. Multiple addresses are comma separated. Max 50.
	Bcc string `json:"bcc,omitempty"`
	// Subject: Email subject
	Subject string `json:"subject,omitempty"`
	// Tag: Email tag that allows you to categorize outgoing emails and get detailed statistics.
	Tag string `json:"tag,omitempty"`
	// HtmlBody: HTML email message. REQUIRED, If no TextBody specified
	HtmlBody string `json:"html_body,omitempty"`
	// TextBody: Plain text email message. REQUIRED, If no HtmlBody specified
	TextBody string `json:"text_body,omitempty"`
	// ReplyTo: Reply To override email address. Defaults to the Reply To set in the sender signature.
	ReplyTo string `json:"reply_to,omitempty"`
	// Headers: List of custom headers to include.
	Headers []Header `json:"headers,omitempty"`
	// TrackOpens: Activate open tracking for this email.
	TrackOpens bool `json:"track_opens,omitempty"`
	// Attachments: List of attachments
	Attachments []Attachment `json:"attachments,omitempty"`
	// Metadata: metadata
	Metadata map[string]string `json:"metdata,omitempty"`
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
