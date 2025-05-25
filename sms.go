package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

func readUnreadMessages() {
	info, err := getSessionToken()
	if err != nil {
		fmt.Printf("Error fetching token: %v\n", err)
		return
	}

	reqBody := SMSListRequest{
		PageIndex:       1,
		ReadCount:       20,
		BoxType:         1,
		SortType:        0,
		Ascending:       0,
		UnreadPreferred: 1,
	}

	xmlPayload, err := xml.Marshal(reqBody)
	if err != nil {
		fmt.Printf("Error marshalling SMS list request: %v\n", err)
		return
	}
	xmlPayload = append([]byte(xml.Header), xmlPayload...)

	req, err := http.NewRequest("POST", config.ApiBase+"/api/sms/sms-list", bytes.NewReader(xmlPayload))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Cookie", fmt.Sprintf("SessionID=%s", info.SesInfo))
	req.Header.Set("__RequestVerificationToken", info.TokInfo)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending SMS list request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	var smsResp SMSListResponse
	if err := xml.Unmarshal(body, &smsResp); err != nil {
		fmt.Printf("Error unmarshalling SMS list response: %v\n", err)
		fmt.Println("Raw response:", string(body))
		return
	}

	if len(smsResp.Messages) == 0 {
		fmt.Println("📭 No new unread messages.")
		return
	}

	fmt.Println("📩 Unread SMS Messages:")
	for _, msg := range smsResp.Messages {
		if msg.Read == 1 {
			continue
		}
		if contains(config.Ignored, msg.Phone) {
			fmt.Printf("Ignoring message from %s\n", msg.Phone)
			continue
		}
		fmt.Printf("Seen? %d From: %s\nDate: %s\nMessage: %s\n", msg.Read, msg.Phone, msg.Date, msg.Content)
		markMessageRead(msg.Index)
		go replyToMessage(msg.Phone, msg.Content)
	}
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func markMessageRead(index int) {
	info, err := getSessionToken()
	if err != nil {
		fmt.Printf("Error fetching token: %v\n", err)
		return
	}

	type SetReadRequest struct {
		XMLName xml.Name `xml:"request"`
		Index   int      `xml:"Index"`
	}

	payload := SetReadRequest{
		Index: index,
	}

	xmlBody, err := xml.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling request: %v\n", err)
		return
	}
	xmlBody = append([]byte(xml.Header), xmlBody...)

	req, err := http.NewRequest("POST", config.ApiBase+"/api/sms/set-read", bytes.NewReader(xmlBody))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Cookie", fmt.Sprintf("SessionID=%s", info.SesInfo))
	req.Header.Set("__RequestVerificationToken", info.TokInfo)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error marking message unread: %v\n", err)
		return
	}
	defer resp.Body.Close()

	respData, _ := io.ReadAll(resp.Body)
	logAPIResponse(fmt.Sprintf("Mark message %d as unread", index), respData)
}

func sendSms(number string, text string) {
	info, err := getSessionToken()
	if err != nil {
		fmt.Printf("Error fetching token: %v\n", err)
		return
	}

	smsReq := SMSRequest{
		Index:    -1,
		Phones:   Phones{Phone: []string{number}},
		Sca:      "",
		Content:  text,
		Length:   -1,
		Reserved: 1,
		Date:     -1,
	}

	xmlBody, err := xml.Marshal(smsReq)
	if err != nil {
		fmt.Printf("Error marshalling SMS request: %v\n", err)
		return
	}

	xmlBody = append([]byte(xml.Header), xmlBody...)

	req, err := http.NewRequest("POST", config.ApiBase+"/api/sms/send-sms", bytes.NewReader(xmlBody))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Cookie", fmt.Sprintf("SessionID=%s", info.SesInfo))
	req.Header.Set("__RequestVerificationToken", info.TokInfo)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending SMS: %v\n", err)
		return
	}
	defer resp.Body.Close()

	respMsg, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading SMS response: %v\n", err)
		return
	}

	logAPIResponse(fmt.Sprintf("Send SMS to %s", number), respMsg)
}

func pollSMSCount(callback func()) {
	for {
		start := time.Now()

		resp, err := http.Get(config.ApiBase + "/api/sms/sms-count")
		if err != nil {
			fmt.Printf("Error polling SMS count: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Printf("Error reading response: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}

		var count SMSCountResponse
		if err := xml.Unmarshal(body, &count); err != nil {
			fmt.Printf("Error parsing XML: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}

		if count.LocalUnread > 0 {
			callback()
		}

		elapsed := time.Since(start)
		if elapsed < time.Second {
			time.Sleep(time.Second - elapsed)
		}
	}
}
