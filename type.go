package main

import "encoding/xml"

type Config struct {
	ApiBase      string   `json:"apiBase"`
	Ignored      []string `json:"ignore"`
	OllamaBase   string   `json:"ollamaBase"`
	Model        string   `json:"model"`
	SystemPrompt string   `json:"systemPrompt"`
}

type SesTokInfo struct {
	XMLName xml.Name `xml:"response"`
	SesInfo string   `xml:"SesInfo"`
	TokInfo string   `xml:"TokInfo"`
}

type SMSRequest struct {
	XMLName  xml.Name `xml:"request"`
	Index    int      `xml:"Index"`
	Phones   Phones   `xml:"Phones"`
	Sca      string   `xml:"Sca"`
	Content  string   `xml:"Content"`
	Length   int      `xml:"Length"`
	Reserved int      `xml:"Reserved"`
	Date     int      `xml:"Date"`
}

type Phones struct {
	Phone []string `xml:"Phone"`
}

type OKResponse struct {
	XMLName xml.Name `xml:"response"`
	Value   string   `xml:",chardata"`
}

type ErrorResponse struct {
	XMLName xml.Name `xml:"error"`
	Code    string   `xml:"code"`
	Message string   `xml:"message"`
}

type SMSCountResponse struct {
	XMLName      xml.Name `xml:"response"`
	LocalUnread  int      `xml:"LocalUnread"`
	LocalInbox   int      `xml:"LocalInbox"`
	LocalOutbox  int      `xml:"LocalOutbox"`
	LocalDraft   int      `xml:"LocalDraft"`
	LocalDeleted int      `xml:"LocalDeleted"`
	SimUnread    int      `xml:"SimUnread"`
	SimInbox     int      `xml:"SimInbox"`
	SimOutbox    int      `xml:"SimOutbox"`
	SimDraft     int      `xml:"SimDraft"`
	LocalMax     int      `xml:"LocalMax"`
	SimMax       int      `xml:"SimMax"`
	SimUsed      int      `xml:"SimUsed"`
	NewMsg       int      `xml:"NewMsg"`
}

type SMSListRequest struct {
	XMLName         xml.Name `xml:"request"`
	PageIndex       int      `xml:"PageIndex"`
	ReadCount       int      `xml:"ReadCount"`
	BoxType         int      `xml:"BoxType"` // 1 = Inbox
	SortType        int      `xml:"SortType"`
	Ascending       int      `xml:"Ascending"`
	UnreadPreferred int      `xml:"UnreadPreferred"`
}

type SMSMessage struct {
	Index   int    `xml:"Index"`
	Phone   string `xml:"Phone"`
	Content string `xml:"Content"`
	Date    string `xml:"Date"`
	Read    int    `xml:"Smstat"`
}

type SMSListResponse struct {
	XMLName  xml.Name     `xml:"response"`
	Messages []SMSMessage `xml:"Messages>Message"`
}

type StreamChunk struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Done bool `json:"done"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ChatResponse struct {
	Message Message `json:"message"`
	Done    bool    `json:"done"`
}
