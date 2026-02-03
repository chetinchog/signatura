package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	apiKey := os.Getenv("SIGNATURA_API_KEY")
	if apiKey == "" {
		log.Fatal("Error: SIGNATURA_API_KEY environment variable not set")
	}

	// Usar el primer documento ID de la lista
	documentID := "019c2019-98db-76a8-8f27-e2f3498be992"

	fmt.Println("🔍 DEBUG: GetDocument API Response")
	fmt.Printf("   Document ID: %s\n\n", documentID)

	url := fmt.Sprintf("https://connect.signatura.co/api/v2/documents/%s", documentID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("❌ Error en request: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	fmt.Printf("   Status Code: %d\n", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("❌ Error leyendo body: %v", err)
	}

	fmt.Printf("   Response Length: %d bytes\n\n", len(body))

	// Parsear y mostrar la estructura completa
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("❌ Error parseando JSON: %v\n", err)
		fmt.Printf("Raw response: %s\n", string(body))
	} else {
		fmt.Println("   ✅ JSON Response:")
		prettyJSON, _ := json.MarshalIndent(response, "   ", "  ")
		fmt.Printf("%s\n", string(prettyJSON))
	}
}
