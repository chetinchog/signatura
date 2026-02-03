package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chetinchog/signatura"
)

func main() {
	// Obtener API key de variable de entorno
	apiKey := os.Getenv("SIGNATURA_API_KEY")
	if apiKey == "" {
		log.Fatal("Error: SIGNATURA_API_KEY environment variable not set")
	}

	// Crear cliente
	client := signatura.New(signatura.Config{
		APIKey: apiKey,
	})

	ctx := context.Background()

	fmt.Println("🔐 Cliente Signatura inicializado correctamente")
	fmt.Println("📋 Ejecutando pruebas con API real...\n")

	// Test 1: Listar documentos
	fmt.Println("1️⃣  Test: Listar documentos")
	fmt.Println("   Obteniendo últimos 10 documentos...")

	result, err := client.ListDocuments(ctx, signatura.ListDocumentsParams{
		Limit: 10,
	})

	if err != nil {
		log.Fatalf("❌ Error listando documentos: %v", err)
	}

	fmt.Printf("   ✅ Éxito! Total de documentos: %d\n", result.TotalCount)
	fmt.Printf("   📄 Documentos en esta página: %d\n\n", len(result.Documents))

	// Mostrar algunos documentos
	if len(result.Documents) > 0 {
		fmt.Println("   📝 Primeros documentos:")
		for i, doc := range result.Documents {
			if i >= 3 {
				break
			}
			fmt.Printf("      - ID: %s\n", doc.ID)
			fmt.Printf("        Título: %s\n", doc.Title)
			fmt.Printf("        Estado: %s\n", doc.Status)
			if doc.CreationDate != nil {
				fmt.Printf("        Creado: %s\n", doc.CreationDate.Format("02/01/2006 15:04"))
			}
			fmt.Println()
		}
	}

	// Test 2: Filtrar documentos completados
	fmt.Println("2️⃣  Test: Listar documentos completados")

	completedResult, err := client.ListDocuments(ctx, signatura.ListDocumentsParams{
		Status: signatura.DocumentStatusCompleted,
		Limit:  5,
	})

	if err != nil {
		log.Fatalf("❌ Error listando documentos completados: %v", err)
	}

	fmt.Printf("   ✅ Documentos completados: %d\n\n", len(completedResult.Documents))

	// Test 3: Listar documentos pendientes
	fmt.Println("3️⃣  Test: Listar documentos pendientes")

	pendingResult, err := client.ListDocuments(ctx, signatura.ListDocumentsParams{
		Status: signatura.DocumentStatusPending,
		Limit:  5,
	})

	if err != nil {
		log.Fatalf("❌ Error listando documentos pendientes: %v", err)
	}

	fmt.Printf("   ✅ Documentos pendientes: %d\n\n", len(pendingResult.Documents))

	// Test 4: Obtener un documento específico (si existe alguno)
	if len(result.Documents) > 0 {
		fmt.Println("4️⃣  Test: Obtener documento específico")
		testDocID := result.Documents[0].ID
		fmt.Printf("   Obteniendo documento: %s\n", testDocID)

		doc, err := client.GetDocument(ctx, testDocID)
		if err != nil {
			log.Fatalf("❌ Error obteniendo documento: %v", err)
		}

		fmt.Printf("   ✅ Documento obtenido correctamente\n")
		fmt.Printf("      Título: %s\n", doc.Title)
		fmt.Printf("      Estado: %s\n", doc.Status)
		fmt.Printf("      Firmas totales: %d\n", len(doc.Signatures))

		// Usar helpers para analizar firmas
		if len(doc.Signatures) > 0 {
			signed := signatura.GetSignedSignatures(doc)
			pending := signatura.GetPendingSignatures(doc)
			declined := signatura.GetDeclinedSignatures(doc)

			fmt.Printf("      Firmadas: %d | Pendientes: %d | Rechazadas: %d\n",
				len(signed), len(pending), len(declined))

			// Mostrar helpers de estado
			if signatura.IsDocumentCompleted(doc) {
				fmt.Println("      ✅ Estado: Completado")
			} else if signatura.IsDocumentPending(doc) {
				fmt.Println("      ⏳ Estado: Pendiente")
			} else if signatura.IsDocumentCanceled(doc) {
				fmt.Println("      ❌ Estado: Cancelado")
			}
		}
		fmt.Println()
	}

	// Test 5: Filtrar por fecha
	fmt.Println("5️⃣  Test: Listar documentos de los últimos 30 días")
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	recentResult, err := client.ListDocuments(ctx, signatura.ListDocumentsParams{
		CreatedAfter: &thirtyDaysAgo,
		Limit:        20,
	})

	if err != nil {
		log.Fatalf("❌ Error listando documentos recientes: %v", err)
	}

	fmt.Printf("   ✅ Documentos de los últimos 30 días: %d\n\n", len(recentResult.Documents))

	// Resumen final
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("📊 RESUMEN DE PRUEBAS")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("✅ Todas las pruebas pasaron exitosamente")
	fmt.Printf("📈 Total de documentos en la cuenta: %d\n", result.TotalCount)
	fmt.Printf("✅ Completados: %d\n", len(completedResult.Documents))
	fmt.Printf("⏳ Pendientes: %d\n", len(pendingResult.Documents))
	fmt.Println()
	fmt.Println("🎉 El cliente funciona correctamente con la API real!")
}
