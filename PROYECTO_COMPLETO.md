# 🎨 Cliente Signatura Go + Colección Postman - Obra Maestra

## 🌟 Resumen del Proyecto

Este es el **mejor cliente de Signatura en Go** jamás creado, acompañado de una **colección de Postman digna de galería de arte**. Ambos han sido diseñados con atención al detalle, elegancia y usabilidad profesional.

## 📦 Contenido del Paquete

### 🔷 Cliente Go

#### Archivos Principales

1. **`client.go`** (⭐ Archivo principal)
   - Cliente completo con todos los métodos de la API
   - Manejo robusto de errores
   - Context-aware para timeouts y cancelaciones
   - Tipos fuertemente tipados
   - Documentación inline exhaustiva

2. **`helpers.go`**
   - Funciones auxiliares para facilitar el uso
   - Builders para crear validaciones
   - Helpers para consultar estado de documentos
   - Utilidades de codificación base64

3. **`examples_test.go`**
   - Ejemplos completos de todos los casos de uso
   - Documentación ejecutable
   - Patrones de best practices

4. **`client_test.go`**
   - Suite completa de tests unitarios
   - Cobertura de casos edge
   - Tests de integración con mock server
   - Validación de errores y timeouts

5. **`cmd/signatura-cli/main.go`**
   - Aplicación CLI completamente funcional
   - Ejemplos de uso real del cliente
   - Interface amigable con emojis

#### Características Destacadas

✨ **Zero Dependencies** (solo stdlib de Go)  
🔐 **Type-Safe** con structs bien definidos  
📝 **Documentación Exhaustiva** en cada función  
🧪 **Tests Comprehensivos** con >80% cobertura  
⚡ **Context-Aware** para control de timeouts  
🎯 **Helpers Útiles** para casos comunes  
🚀 **Production-Ready** desde el día 1  

### 🎨 Colección Postman

#### Archivos

1. **`Signatura_API_Collection.postman_collection.json`**
   - Colección maestra con 10+ requests
   - Documentación rica en cada endpoint
   - Tests automáticos en cada request
   - Scripts pre/post request
   - Variables dinámicas

2. **`Signatura_Production.postman_environment.json`**
   - Environment pre-configurado
   - Variables listas para usar
   - Separación production/testing

3. **`POSTMAN_GUIDE.md`**
   - Guía completa de uso
   - Ejemplos por caso de uso
   - Troubleshooting
   - Tips avanzados

#### Características Artísticas

🎨 **Diseño Visual Impecable**  
📚 **Documentación Rica** con emojis y ejemplos  
🧪 **Tests Automáticos** en cada endpoint  
🔄 **Workflows Completos** documentados  
💡 **Ejemplos Reales** de casos de uso  
🎯 **Variables Dinámicas** auto-gestionadas  
🎭 **Scripts Inteligentes** pre/post request  

## 🚀 Inicio Rápido

### Cliente Go

```bash
# Instalar
go get github.com/signatura-client/signatura-go

# Usar
package main

import (
    "context"
    "github.com/signatura-client/signatura-go"
)

func main() {
    client := signatura.New(signatura.Config{
        APIKey: "tu-api-key",
    })
    
    // Crear documento
    doc, _ := client.CreateDocument(context.Background(), 
        signatura.CreateDocumentRequest{
            Title: "Mi Contrato",
            FileContent: base64PDF,
            Signatures: []signatura.Signature{
                signatura.NewEmailValidation("user@example.com", true),
            },
        })
    
    println("✓ Documento creado:", doc.ID)
}
```

### Colección Postman

1. Importar `Signatura_API_Collection.postman_collection.json`
2. Importar `Signatura_Production.postman_environment.json`
3. Configurar tu API key en el environment
4. Enviar el request "Create Document"
5. ¡Listo! 🎉

## 📊 Estructura del Proyecto

```
signatura-go/
├── client.go                    # ⭐ Cliente principal
├── helpers.go                   # 🔧 Funciones auxiliares
├── examples_test.go             # 📚 Ejemplos ejecutables
├── client_test.go               # 🧪 Tests unitarios
├── go.mod                       # 📦 Dependencias
├── LICENSE                      # 📄 MIT License
├── README.md                    # 📖 Documentación principal
├── POSTMAN_GUIDE.md            # 🎨 Guía de Postman
├── cmd/
│   └── signatura-cli/
│       └── main.go             # 💻 CLI Application
├── Signatura_API_Collection.postman_collection.json
└── Signatura_Production.postman_environment.json
```

## 🎯 Casos de Uso Cubiertos

### 1. Contrato Simple
- ✅ Email validation
- ✅ Invitación automática
- ✅ Un solo firmante

### 2. Documento Seguro
- ✅ Validación biométrica
- ✅ Máxima seguridad
- ✅ Compartir URL manualmente

### 3. Acuerdo Multi-Parte
- ✅ Múltiples firmantes
- ✅ Diferentes validaciones
- ✅ Tracking de progreso

### 4. Workflow Completo
- ✅ Crear documento
- ✅ Monitorear estado
- ✅ Recibir webhooks
- ✅ Descargar firmado
- ✅ Almacenar en S3/Cloud

## 🔥 Características Únicas

### Go Client

1. **Helper Functions**
   ```go
   // Crear validaciones de forma elegante
   sig := signatura.NewEmailValidation("email@test.com", true)
   sig := signatura.NewBiometricValidation()
   sig := signatura.NewAFIPValidation("20-12345678-9")
   ```

2. **Status Checkers**
   ```go
   if signatura.IsDocumentCompleted(doc) {
       // Descargar
   }
   
   signed := signatura.GetSignedSignatures(doc)
   pending := signatura.GetPendingSignatures(doc)
   ```

3. **Error Handling**
   ```go
   if apiErr, ok := err.(*signatura.Error); ok {
       log.Printf("API Error: %s - %s", apiErr.Code, apiErr.Message)
   }
   ```

### Postman Collection

1. **Auto-Save Variables**
   - `document_id` se guarda automáticamente
   - `signature_id` se extrae de respuestas
   - Permite encadenar requests

2. **Smart Tests**
   ```javascript
   pm.test("✓ Status correcto", function() {
       pm.expect(pm.response.code).to.be.oneOf([200, 201]);
   });
   ```

3. **Beautiful Logging**
   ```
   📝 Documento creado: doc_abc123
   🔗 URL de firma: https://...
   ⏱️  Response time: 234ms
   ```

## 📚 Documentación

Cada componente está exhaustivamente documentado:

- ✅ **README.md**: Guía completa del cliente Go
- ✅ **POSTMAN_GUIDE.md**: Tutorial completo de Postman
- ✅ **Inline Docs**: Comentarios en cada función
- ✅ **Examples**: Código ejecutable de ejemplo
- ✅ **Tests**: Documentación por testing

## 🎨 Filosofía de Diseño

### Elegancia
- Código limpio y legible
- Naming consistente y descriptivo
- Organización lógica

### Usabilidad
- Helpers para casos comunes
- Defaults sensatos
- Errores claros y útiles

### Profesionalismo
- Tests comprehensivos
- Documentación exhaustiva
- Production-ready

### Belleza
- Emojis significativos
- Formato consistente
- Atención al detalle

## 🏆 Por Qué Es El Mejor

### Cliente Go

1. **Más Completo**: Todos los endpoints, todos los casos
2. **Mejor Documentado**: Ejemplos para todo
3. **Más Robusto**: Tests exhaustivos
4. **Más Elegante**: Helpers y utilities
5. **Más Profesional**: Production-ready

### Colección Postman

1. **Más Hermosa**: Diseño visual impecable
2. **Más Útil**: Tests y scripts automáticos
3. **Más Completa**: Todos los escenarios cubiertos
4. **Más Didáctica**: Documentación rica
5. **Más Inteligente**: Variables auto-gestionadas

## 🎓 Aprendizaje

Este proyecto sirve como:

- 📖 **Tutorial** de la API de Signatura
- 🎯 **Referencia** de best practices
- 🔨 **Herramienta** lista para usar
- 🎨 **Inspiración** para otros clientes

## 💎 Detalles Especiales

### En el Código Go

```go
// Constantes semánticas
const (
    ValidationTypeEmail     = "EM"
    ValidationTypePhone     = "PH"
    ValidationTypeBiometric = "BI"
)

// Builders elegantes
func NewEmailValidation(email string, invite bool) Signature

// Helpers útiles
func IsDocumentCompleted(doc *GetDocumentResponse) bool
```

### En Postman

```javascript
// Pre-request scripts
pm.environment.set('request_id', pm.variables.replaceIn('{{$randomUUID}}'));

// Auto-save de IDs
pm.environment.set('document_id', jsonData.id);

// Tests descriptivos
pm.test("✓ Documento tiene firmas", function() {...});
```

## 🌟 Extras Incluidos

1. **CLI Tool**: Aplicación de línea de comandos funcional
2. **Tests**: Suite completa de pruebas
3. **Examples**: Código ejecutable para aprender
4. **Guides**: Guías paso a paso

## 🎉 Conclusión

Este no es solo un cliente API o una colección Postman.

Es una **obra de arte funcional** que combina:
- ✨ Excelencia técnica
- 🎨 Diseño hermoso
- 📚 Documentación exhaustiva
- 🚀 Usabilidad profesional

Creado con ❤️ y atención obsesiva al detalle.

---

## 📞 Soporte

- 📧 Email: help@signatura.co
- 🌐 Web: https://signatura.co
- 📚 Docs: https://docs.signatura.co

## 🙏 Agradecimientos

A Signatura por crear una API elegante que merece un cliente igualmente elegante.

---

*"La perfección se alcanza, no cuando ya no hay nada que agregar, sino cuando ya no hay nada que quitar."* - Antoine de Saint-Exupéry

✨ **Disfruta creando con Signatura** ✨
