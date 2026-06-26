package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)


//  Structs for JSON Marshaling & Unmarshaling ---

//GeminiRequest matches the JSON structure required by the Gemini API for requests 

type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}


// GeminiResponse matches the JSON structure returned by the Gemini API.

type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content ResponseContent `json:"content"`
}

type ResponseContent struct {
	Parts []Part `json:"parts"`
}


// Audit Function

// auditFile sends a single source code file's contents to the Gemini API for auditing.

func auditFile(apiKey string, file CodeFile)(string, error) {
	//Construct the endpoint url,  including the API key as a query parameter 
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s", apiKey)
    
	//Build the audit prompt 
	prompt := fmt.Sprintf(
		"You are an expert software security auditor. Review the following code from the file '%s' and identify any bugs, security vulnerabilities, or performance issues. Provide clear explanations and suggest fixed versions of the code where appropriate.\n\nCode:\n```\n%s\n```",
		file.Path,
		file.Content,
	)

	//Create a request Payload structure 
	reqPayload := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
                 {Text: prompt},
				} ,
			},
		},
	}

	//json.Marshal converts our Go struct into a slice of JSON bytes ([]byte)
	jsonData, err := json.Marshal(reqPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	//Create a new HTTP POST request 
	//bytes.NewBuffer turns our byte slice into an io.Reader stream that the HTTP client can read from 
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
			return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	//Set the content-type header to tell the Gemini server we are sending JSON 
	req.Header.Set("Content-Type", "application/json")


	// Initialize an HTTP client with a 30-second timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	//Send The Request 
	resp, err := client.Do(req)
	if err != nil {
			return "", fmt.Errorf("API request failed: %w", err)
	}

	//Defer closing the response Body. 
	//This ensures the connection is closed when auditFile returns, preventing resource leaks 
	defer resp.Body.Close()

	// io.ReadAll reads the entire stream of bytes from the response body into memory
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// If the server returns a non-200 OK status, read and report the error response
	if resp.StatusCode != http.StatusOK {
				return "", fmt.Errorf("API returned non-200 status %d: %s", resp.StatusCode, string(respBody)) 
	}

	// Create a variable to hold the decoded response
	var geminiResp GeminiResponse

		// json.Unmarshal decodes the JSON bytes into our Go struct.
	// We pass a pointer (&geminiResp) so the function can write directly into our variable's memory.

	err = json.Unmarshal(respBody, &geminiResp)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %w", err)
	}

	//Validate the response structure and extract the generated text 
	if len(geminiResp.Candidates) == 0 ||
	 len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("received empty response candidates from Gemini API")
	 }

	 //Return the generated audit text from the first candidate 
	 return geminiResp.Candidates[0].Content.Parts[0].Text, nil 
}