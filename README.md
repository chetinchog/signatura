# Signatura Go Client

![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue)
![License](https://img.shields.io/badge/license-MIT-green)
![Coverage](https://img.shields.io/badge/coverage-90.5%25-brightgreen)
![Tests](https://img.shields.io/badge/tests-50%2B%20passing-success)

Un cliente Go elegante, robusto y completo para la API de Signatura. Integra firmas electrónicas con validación de identidad de forma sencilla y profesional.

## ✨ Características

- 🔐 **Autenticación segura** con API keys
- 📝 **Creación de documentos** con múltiples firmantes
- ✅ **Validaciones múltiples**: Email, Teléfono, Biométrica, AFIP
- 📊 **Gestión completa** de documentos y firmas
- 🔔 **Soporte para webhooks** con tipos fuertemente tipados
- ⚡ **Context-aware** para timeouts y cancelaciones
- 🎯 **Type-safe** con structs bien definidos
- 🧪 **90.5% cobertura de tests** con más de 50 casos de prueba
- 🚀 **Zero dependencies** (solo stdlib)

## 📦 Instalación

```bash
go get github.com/chetinchog/signatura
```

## 🚀 Inicio Rápido

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/chetinchog/signatura"
)

func main() {
    // Crear cliente
    client := signatura.New(signatura.Config{
        APIKey: "tu-api-key-aqui",
    })
    
    // Codificar PDF
    pdfContent, err := signatura.EncodeFileToBase64("contrato.pdf")
    if err != nil {
        log.Fatal(err)
    }
    
    // Crear documento con firma electrónica
    doc, err := client.CreateDocument(context.Background(), signatura.CreateDocumentRequest{
        Title:       "Contrato de Trabajo - Juan Pérez",
        FileContent: pdfContent,
        Signatures: []signatura.Signature{
            {
                SignerName: "Juan Pérez",
                Validations: signatura.Validation{
                    Email: signatura.String("juan@ejemplo.com"),
                    Phone: nil, // El firmante ingresará su teléfono
                },
                InviteChannel: []string{signatura.InviteChannelEmail},
            },
        },
    })
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("✓ Documento creado: %s\n", doc.ID)
    fmt.Printf("  URL de firma: %s\n", doc.Signatures[0].URL)
}
```

## 📚 Casos de Uso

### 1️⃣ Documento con Validación Email + Teléfono

```go
client := signatura.New(signatura.Config{APIKey: "tu-api-key"})

doc, err := client.CreateDocument(context.Background(), signatura.CreateDocumentRequest{
    Title:       "Contrato con Doble Validación",
    FileContent: base64PDF,
    Signatures: []signatura.Signature{
        {
            Validations: signatura.Validation{
                Email: signatura.String("cliente@empresa.com"),
                Phone: nil, // Se solicitará al firmante
            },
            InviteChannel: []string{signatura.InviteChannelEmail},
        },
    },
})
```

### 2️⃣ Documento con Validación Biométrica

```go
doc, err := client.CreateDocument(context.Background(), signatura.CreateDocumentRequest{
    Title:       "Acuerdo de Confidencialidad",
    FileContent: base64PDF,
    Signatures: []signatura.Signature{
        signatura.NewBiometricValidation(),
    },
})
if err != nil {
    log.Fatal(err)
}

// La URL de firma estará disponible inmediatamente
fmt.Println("URL para firmar:", doc.Signatures[0].URL)
```

### 3️⃣ Múltiples Firmantes

```go
doc, err := client.CreateDocument(context.Background(), signatura.CreateDocumentRequest{
    Title:       "Acuerdo Tripartito",
    FileContent: base64PDF,
    Signatures: []signatura.Signature{
        {
            SignerName: "Parte A",
            Validations: signatura.Validation{
                Email: signatura.String("parte-a@empresa.com"),
            },
            InviteChannel: []string{signatura.InviteChannelEmail},
        },
        {
            SignerName: "Parte B", 
            Validations: signatura.Validation{
                Email: signatura.String("parte-b@empresa.com"),
                Phone: signatura.String("+5491134567890"),
            },
            InviteChannel: []string{signatura.InviteChannelEmail},
        },
        signatura.NewBiometricValidation(), // Parte C - sin invitación
    },
    Metadata: map[string]interface{}{
        "tipo_contrato": "tripartito",
        "referencia":    "CONT-2024-123",
    },
})
```

### 4️⃣ Consultar Estado de Documento

```go
doc, err := client.GetDocument(context.Background(), documentID)
if err != nil {
    log.Fatal(err)
}

// Helpers útiles
if signatura.IsDocumentCompleted(doc) {
    fmt.Println("✓ Documento completamente firmado!")
}

signed := signatura.GetSignedSignatures(doc)
pending := signatura.GetPendingSignatures(doc)
declined := signatura.GetDeclinedSignatures(doc)

fmt.Printf("Firmadas: %d, Pendientes: %d, Rechazadas: %d\n", 
    len(signed), len(pending), len(declined))
```

### 5️⃣ Listar Documentos

```go
// Obtener documentos completados de los últimos 30 días
thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

result, err := client.ListDocuments(context.Background(), signatura.ListDocumentsParams{
    Status:       signatura.DocumentStatusCompleted,
    Limit:        100,
    CreatedAfter: &thirtyDaysAgo,
})
if err != nil {
    log.Fatal(err)
}

for _, doc := range result.Documents {
    fmt.Printf("- %s (estado: %s)\n", doc.Title, doc.Status)
    if doc.CreationDate != nil {
        fmt.Printf("  Creado: %s\n", doc.CreationDate.Format("02/01/2006"))
    }
}
```

### 6️⃣ Descargar Documento Firmado

```go
// Verificar que esté completado
doc, err := client.GetDocument(context.Background(), documentID)
if err != nil {
    log.Fatal(err)
}

if !signatura.IsDocumentCompleted(doc) {
    log.Fatal("El documento aún no está completo")
}

// Descargar PDF firmado
pdfBytes, err := client.DownloadDocument(context.Background(), documentID)
if err != nil {
    log.Fatal(err)
}

// Guardar a archivo
if err := os.WriteFile("contrato-firmado.pdf", pdfBytes, 0644); err != nil {
    log.Fatalf("Error guardando PDF: %v", err)
}
```

### 7️⃣ Webhooks

```go
// En tu handler HTTP
func webhookHandler(w http.ResponseWriter, r *http.Request) {
    var event signatura.WebhookEvent
    
    if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
        http.Error(w, "Invalid payload", http.StatusBadRequest)
        return
    }
    
    switch event.NotificationAction {
    case signatura.WebhookActionDocumentSigned:
        // Firmante completó su firma
        fmt.Printf("Documento %s firmado por %s\n",
            event.DocumentID, event.SignatureID)

        // Opcional: consultar estado actualizado
        doc, err := client.GetDocument(context.Background(), event.DocumentID)
        if err != nil {
            log.Printf("Error obteniendo documento: %v", err)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        if signatura.IsDocumentCompleted(doc) {
            // Todos firmaron - procesar documento
            log.Printf("Documento %s completado, descargando...", doc.ID)

            pdfData, err := client.DownloadDocument(context.Background(), doc.ID)
            if err != nil {
                log.Printf("Error descargando: %v", err)
                return
            }

            // Guardar o procesar el PDF según tu lógica de negocio
            _ = pdfData // Tu lógica aquí
        }

    case signatura.WebhookActionSignatureDeclined:
        // Firmante rechazó la firma
        log.Printf("Firma rechazada - Documento: %s, Firma: %s",
            event.DocumentID, event.SignatureID)

        // Implementa tu lógica de manejo de rechazo aquí

    case signatura.WebhookActionDocumentChange:
        // Cambió el estado del documento
        log.Printf("Documento %s cambió a: %s",
            event.DocumentID, event.NewStatus)

        if event.NewStatus == signatura.DocumentStatusCompleted {
            log.Printf("Documento %s completado", event.DocumentID)
            // Implementa tu lógica de notificación aquí
        }
    }
    
    w.WriteHeader(http.StatusOK)
}
```

### 8️⃣ Configuración Avanzada

```go
import "net/http"

// Cliente con configuración personalizada
httpClient := &http.Client{
    Timeout: 60 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
    },
}

client := signatura.New(signatura.Config{
    APIKey:     os.Getenv("SIGNATURA_API_KEY"),
    BaseURL:    signatura.DefaultBaseURL,
    HTTPClient: httpClient,
})

// Usar context para timeout individual
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

doc, err := client.CreateDocument(ctx, request)
```

## 🔑 Tipos de Validación

| Tipo | Código | Descripción |
|------|--------|-------------|
| **Email** | `EM` | Validación por correo electrónico |
| **Teléfono** | `PH` | Validación por SMS |
| **Biométrica** | `BI` | Validación por reconocimiento facial |
| **AFIP** | `AF` | Validación con clave fiscal AFIP |

### Estructura de Validaciones en Respuestas

Cuando la API devuelve información sobre validaciones, usa una estructura anidada:

```go
type ValidationValue struct {
    Validated bool        `json:"validated"`  // Si la validación fue completada
    Value     interface{} `json:"value"`      // El valor validado (email, phone, etc.)
}

type ValidationResponse struct {
    Email      *ValidationValue `json:"EM,omitempty"`
    Phone      *ValidationValue `json:"PH,omitempty"`
    Biometric  *ValidationValue `json:"BI,omitempty"`
    AFIP       *ValidationValue `json:"AF,omitempty"`
}
```

Ejemplo de respuesta:
```json
{
  "validations": {
    "EM": {
      "validated": true,
      "value": "usuario@ejemplo.com"
    },
    "PH": {
      "validated": false,
      "value": null
    }
  }
}
```

## 📡 Canales de Invitación

- `EM` - Email (se envía enlace de firma por correo)
- `SM` - SMS (se envía enlace de firma por SMS)

Si no especificas `invite_channel`, debes compartir manualmente la URL de firma.

## 📊 Estados de Documento

| Estado | Código | Descripción |
|--------|--------|-------------|
| **Pendiente** | `PE` | Esperando firmas |
| **Completado** | `CO` | Todas las firmas completadas |
| **Cancelado** | `CA` | Documento cancelado |

## ✍️ Estados de Firma

| Estado | Código | Descripción |
|--------|--------|-------------|
| **Invitado** | `IN` | Firmante invitado, pendiente de acción |
| **Firmado** | `SI` | Firma completada exitosamente |
| **Rechazado** | `DE` | Firmante rechazó la firma |
| **Pendiente** | `PE` | En proceso de firma |

## 🔔 Eventos de Webhook

| Evento | Código | Descripción |
|--------|--------|-------------|
| **Documento Firmado** | `DS` | Un firmante completó su firma |
| **Firma Rechazada** | `SD` | Un firmante rechazó la firma |
| **Cambio de Estado** | `DC` | El documento cambió de estado |

## 🛡️ Mejores Prácticas

### ✅ Recomendado

```go
// 1. Usar webhooks en lugar de polling
// ❌ NO hagas esto
for {
    doc, _ := client.GetDocument(ctx, docID)
    if doc.Status == signatura.DocumentStatusCompleted {
        break
    }
    time.Sleep(5 * time.Minute) // Desperdicia llamadas API
}

// ✅ SÍ haz esto - Configura un webhook handler
func webhookHandler(w http.ResponseWriter, r *http.Request) {
    var event signatura.WebhookEvent
    if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
        http.Error(w, "Invalid payload", http.StatusBadRequest)
        return
    }

    switch event.NotificationAction {
    case signatura.WebhookActionDocumentSigned:
        // Verificar si el documento está completo
        doc, _ := client.GetDocument(context.Background(), event.DocumentID)
        if signatura.IsDocumentCompleted(doc) {
            processCompletedDocument(doc.ID)
        }
    }

    w.WriteHeader(http.StatusOK)
}

// Registrar el webhook en tu servidor
http.HandleFunc("/webhooks/signatura", webhookHandler)
```

```go
// 2. Usar context para timeouts
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

doc, err := client.CreateDocument(ctx, request)
```

```go
// 3. Manejar errores apropiadamente
doc, err := client.GetDocument(ctx, docID)
if err != nil {
    if apiErr, ok := err.(*signatura.Error); ok {
        fmt.Printf("API Error: %s - %s\n", apiErr.Code, apiErr.Message)
    } else {
        fmt.Printf("Error: %v\n", err)
    }
    return
}
```

```go
// 4. Usar helpers para código más limpio
if signatura.IsDocumentCompleted(doc) {
    pdfData, err := client.DownloadDocument(ctx, doc.ID)
    if err != nil {
        log.Printf("Error downloading: %v", err)
        return
    }

    // Guardar el PDF
    if err := os.WriteFile("contrato-firmado.pdf", pdfData, 0644); err != nil {
        log.Printf("Error saving PDF: %v", err)
    }
}
```

## 🧪 Testing

```go
// Ejemplo de test usando solo stdlib (sin dependencias externas)
func TestCreateDocument(t *testing.T) {
    client := signatura.New(signatura.Config{
        APIKey: os.Getenv("SIGNATURA_TEST_API_KEY"),
    })

    base64Content, err := signatura.EncodeFileToBase64("testdata/sample.pdf")
    if err != nil {
        t.Fatalf("Failed to encode file: %v", err)
    }

    doc, err := client.CreateDocument(context.Background(), signatura.CreateDocumentRequest{
        Title:       "Test Document",
        FileContent: base64Content,
        Signatures: []signatura.Signature{
            signatura.NewBiometricValidation(),
        },
    })

    if err != nil {
        t.Fatalf("CreateDocument() error = %v", err)
    }

    if doc.ID == "" {
        t.Error("Expected non-empty document ID")
    }

    if doc.Status != signatura.DocumentStatusPending {
        t.Errorf("Status = %v, want %v", doc.Status, signatura.DocumentStatusPending)
    }
}
```

## 📖 Documentación Adicional

### En este Repositorio

- [CHANGELOG.md](CHANGELOG.md) - Historial de cambios y versiones
- [SECURITY.md](SECURITY.md) - Política de seguridad y mejores prácticas
- [CLAUDE.md](CLAUDE.md) - Guía para asistentes AI

### Documentación de Signatura

- [Guía de Inicio Rápido de Signatura](https://docs.signatura.co/docs/intro)
- [Límites de API](https://docs.signatura.co/docs/rate-limiting)
- [Webhooks](https://docs.signatura.co/docs/webhooks)
- [Seguridad](https://docs.signatura.co/docs/security)

## 🤝 Contribuir

Las contribuciones son bienvenidas. Por favor:

1. Fork el repositorio
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📄 Licencia

MIT License - ver [LICENSE](LICENSE) para detalles.

## 💬 Soporte

- 📧 Email: help@signatura.co
- 🌐 Website: https://signatura.co
- 📚 Docs: https://docs.signatura.co

---

Hecho con ❤️ para desarrolladores que valoran código elegante y bien documentado.
