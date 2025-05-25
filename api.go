package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

func getSessionToken() (SesTokInfo, error) {
	resp, err := http.Get(config.ApiBase + "/api/webserver/SesTokInfo")
	if err != nil {
		return SesTokInfo{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return SesTokInfo{}, err
	}

	var info SesTokInfo
	if err := xml.Unmarshal(body, &info); err != nil {
		return SesTokInfo{}, err
	}
	return info, nil
}

func logAPIResponse(action string, respData []byte) {
	var okResp OKResponse
	var errResp ErrorResponse

	if err := xml.Unmarshal(respData, &okResp); err == nil && okResp.Value != "" {
		fmt.Printf("✅ [%s] Success: %s\n", action, okResp.Value)
		return
	}

	if err := xml.Unmarshal(respData, &errResp); err == nil && errResp.Code != "" {
		fmt.Printf("❌ [%s] Error Code: %s - %s\n", action, errResp.Code, errResp.Message)
		return
	}

	fmt.Printf("⚠️ [%s] Unknown response:\n%s\n", action, string(respData))
}
