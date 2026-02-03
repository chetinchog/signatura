package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/chetinchog/signatura"
)

func main() {
	apiKey := os.Getenv("SIGNATURA_API_KEY")
	if apiKey == "" {
		log.Fatal("Error: SIGNATURA_API_KEY environment variable not set")
	}

	fmt.Println("🔍 DEBUG: Comparando request directa vs cliente\n")

	// Test 1: Request directa con http.Client
	fmt.Println("1️⃣  Request HTTP directa")
	fmt.Println("   URL: https://connect.signatura.co/api/v2/documents")

	req, err := http.NewRequest("GET", "https://connect.signatura.co/api/v2/documents", nil)
	if err != nil {
		log.Fatal(err)
	}

	// Probar diferentes formatos de autenticación
	fmt.Println("   Header Authorization: Bearer " + apiKey[:20] + "...")

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("   ❌ Error en request: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	fmt.Printf("   Status Code: %d\n", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("   ❌ Error leyendo body: %v", err)
	}

	fmt.Printf("   Response Length: %d bytes\n", len(body))

	// Parsear respuesta
	var directResponse map[string]interface{}
	if err := json.Unmarshal(body, &directResponse); err != nil {
		fmt.Printf("   ❌ Error parseando JSON: %v\n", err)
		fmt.Printf("   Raw response: %s\n", string(body[:min(500, len(body))]))
	} else {
		fmt.Println("   ✅ JSON válido")
		prettyJSON, _ := json.MarshalIndent(directResponse, "   ", "  ")
		fmt.Printf("   Response:\n   %s\n\n", string(prettyJSON))
	}

	// Test 2: Usando el cliente
	fmt.Println("2️⃣  Usando el cliente Signatura")

	signaturaClient := signatura.New(signatura.Config{
		APIKey: apiKey,
	})

	result, err := signaturaClient.ListDocuments(context.Background(), signatura.ListDocumentsParams{
		Limit: 10,
	})

	if err != nil {
		fmt.Printf("   ❌ Error del cliente: %v\n", err)
		if apiErr, ok := err.(*signatura.Error); ok {
			fmt.Printf("   API Error Code: %s\n", apiErr.Code)
			fmt.Printf("   API Error Message: %s\n", apiErr.Message)
		}
	} else {
		fmt.Printf("   ✅ Éxito\n")
		fmt.Printf("   Total Count: %d\n", result.TotalCount)
		fmt.Printf("   Documents: %d\n", len(result.Documents))

		if len(result.Documents) > 0 {
			fmt.Println("   Primer documento:")
			doc, _ := json.MarshalIndent(result.Documents[0], "   ", "  ")
			fmt.Printf("   %s\n", string(doc))
		}
	}

	// Test 3: Verificar headers que está enviando el cliente
	fmt.Println("\n3️⃣  Verificando headers del cliente")
	fmt.Println("   BaseURL configurado:", signatura.DefaultBaseURL)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
