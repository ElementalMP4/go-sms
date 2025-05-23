package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

const baseURL = "http://192.168.9.1"

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

func main() {
	// Define CLI flags
	numberFlag := flag.String("number", "", "Destination phone number")
	textFlag := flag.String("text", "", "Text message content")
	flag.Parse()

	if *numberFlag == "" || *textFlag == "" {
		fmt.Println("Both --number and --text flags are required.")
		flag.Usage()
		os.Exit(1)
	}

	// Step 1: Fetch session and token
	resp, err := http.Get(baseURL + "/api/webserver/SesTokInfo")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var info SesTokInfo
	if err := xml.Unmarshal(body, &info); err != nil {
		panic(err)
	}

	// Step 2: Create the SMSRequest struct
	smsReq := SMSRequest{
		Index:    -1,
		Phones:   Phones{Phone: []string{*numberFlag}},
		Sca:      "",
		Content:  *textFlag,
		Length:   -1,
		Reserved: 1,
		Date:     -1,
	}

	xmlBody, err := xml.Marshal(smsReq)
	if err != nil {
		panic(err)
	}

	xmlBody = append([]byte(xml.Header), xmlBody...)

	// Step 3: Send SMS
	req, err := http.NewRequest("POST", baseURL+"/api/sms/send-sms", bytes.NewReader(xmlBody))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Cookie", fmt.Sprintf("SessionID=%s", info.SesInfo))
	req.Header.Set("__RequestVerificationToken", info.TokInfo)

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respMsg, _ := io.ReadAll(resp.Body)
	fmt.Println(string(respMsg))
}
