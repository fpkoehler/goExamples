package main

import (
        "fmt"
        "io"
        "log"
        "net/http"
)

func main() {
        url := "https://www.footballdb.com/scores/index.html?lg=NFL&yr=2024&type=reg&wk=1"
        // Create a custom HTTP client
        client := &http.Client{}

        // Create a new request
        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
                log.Fatalf("Failed to create request: %v", err)
        }

        // Add headers (e.g., User-Agent)
        req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; MyGolangBot/1.0)")

        // Make the request
        resp, err := client.Do(req)
        if err != nil {
                log.Fatalf("Failed to make request: %v", err)
        }
        defer resp.Body.Close()

        // Check the status code
        if resp.StatusCode == http.StatusForbidden {
                fmt.Printf("Received 403 Forbidden. Headers sent: %v\n", req.Header)
                fmt.Printf("Response headers: %v\n", resp.Header)
                body, _ := io.ReadAll(resp.Body) // Read body for more details
                fmt.Printf("Response body: %s\n", string(body))
                return
        } else if resp.StatusCode != http.StatusOK {
                log.Fatalf("Unexpected status code: %d", resp.StatusCode)
        }

        // Read and print the response body
        body, err := io.ReadAll(resp.Body)
        if err != nil {
                log.Fatalf("Failed to read response body: %v", err)
        }
        fmt.Println("Response:", string(body))
}
