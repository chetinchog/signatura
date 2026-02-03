# 🎨 Guía de Uso - Colección Postman de Signatura

## 📥 Instalación

### 1. Importar la Colección

#### Opción A: Importar desde archivo
1. Abre Postman
2. Haz clic en **Import** (botón superior izquierdo)
3. Selecciona el archivo `Signatura_API_Collection.postman_collection.json`
4. Haz clic en **Import**

#### Opción B: Importar desde URL
```
https://raw.githubusercontent.com/tu-repo/signatura-go/main/Signatura_API_Collection.postman_collection.json
```

### 2. Importar el Environment

1. En Postman, haz clic en **Import**
2. Selecciona el archivo `Signatura_Production.postman_environment.json`
3. Haz clic en **Import**

### 3. Configurar tu API Key

1. En Postman, selecciona el environment **🎨 Signatura - Production**
2. Haz clic en el ícono del ojo 👁️ en la esquina superior derecha
3. Haz clic en **Edit** junto al environment
4. Encuentra la variable `api_key`
5. Reemplaza `"your-api-key-here"` con tu API key real de Signatura
6. Haz clic en **Save**

## 🚀 Uso Rápido

### Crear tu Primer Documento

1. **Prepara tu PDF**
   - Asegúrate de tener un archivo PDF listo
   - El PDF debe ser válido y no estar corrupto

2. **Codificar a Base64**
   
   En la pestaña **Pre-request Script** del request "Create Document", el script automáticamente genera un PDF de ejemplo. Para usar tu propio PDF:

   **Opción 1: Online**
   - Usa [base64encode.org](https://www.base64encode.org/)
   - Sube tu PDF
   - Copia el resultado

   **Opción 2: Terminal (Linux/Mac)**
   ```bash
   base64 tu-documento.pdf > documento-base64.txt
   ```

   **Opción 3: PowerShell (Windows)**
   ```powershell
   [Convert]::ToBase64String([IO.File]::ReadAllBytes("tu-documento.pdf"))
   ```

3. **Ejecutar el Request**
   - Abre el request **✨ Create Document - Email + Phone Validation**
   - Modifica el body:
     - `title`: Cambia el título del documento
     - `file_content`: Pega tu PDF en base64
     - `signatures[0].validations.EM`: Email del firmante
     - `signatures[0].signer_name`: Nombre del firmante
   - Haz clic en **Send**

4. **Revisar la Respuesta**
   - Verás el `document_id` generado
   - Encontrarás la `signing_url` en `signatures[0].signing_url`
   - El `document_id` se guardará automáticamente en las variables de entorno

## 📚 Flujos Comunes

### Flujo 1: Documento Simple con Email

```
1. Create Document - Email + Phone Validation
   ↓
2. Get Document Status (para verificar)
   ↓
3. (Esperar a que el firmante firme)
   ↓
4. Get Document Status (verificar completado)
   ↓
5. Download Signed Document
```

### Flujo 2: Documento con Múltiples Firmantes

```
1. Create Document - Multiple Signers
   ↓
2. (Compartir URLs con cada firmante)
   ↓
3. Get Document Status (monitorear progreso)
   ↓
4. (Reenviar invitación si es necesario)
   → Resend Signature Invitation
   ↓
5. Get Document Status (verificar que todos firmaron)
   ↓
6. Download Signed Document
```

### Flujo 3: Listar y Gestionar Documentos

```
1. List Documents (ver todos los documentos)
   ↓
2. (Filtrar por status si es necesario)
   ↓
3. Get Document Status (para uno específico)
   ↓
4. Download / Cancel (según necesites)
```

## 🧪 Tests Automáticos

Cada request incluye tests automáticos que se ejecutan después de recibir la respuesta:

### Revisar Tests Ejecutados

1. Después de enviar un request, ve a la pestaña **Test Results**
2. Verás checkmarks ✓ verdes para tests exitosos
3. Verás X rojas para tests fallidos

### Tests Incluidos

- ✅ Validación de código de respuesta HTTP
- ✅ Verificación de estructura de respuesta
- ✅ Validación de datos requeridos
- ✅ Logging de información útil en la consola

### Ver Logs en Consola

1. Abre la **Postman Console** (View → Show Postman Console)
2. Envía un request
3. Verás logs detallados como:
   ```
   📝 Documento creado: doc_abc123
   🔗 URL de firma: https://connect.signatura.co/sign/xyz
   ```

## 🎯 Ejemplos por Caso de Uso

### Caso 1: Contrato Laboral

**Request:** Create Document - Email + Phone Validation

```json
{
  "title": "Contrato de Trabajo - María González",
  "file_content": "{{sample_pdf_base64}}",
  "signatures": [
    {
      "signer_name": "María González",
      "validations": {
        "EM": "maria.gonzalez@ejemplo.com",
        "PH": null
      },
      "invite_channel": ["EM"]
    }
  ],
  "metadata": {
    "department": "RRHH",
    "employee_id": "EMP-2024-001",
    "position": "Developer Senior"
  }
}
```

### Caso 2: NDA con Máxima Seguridad

**Request:** Create Document - Biometric Validation

```json
{
  "title": "NDA - Proyecto Confidencial Alpha",
  "file_content": "{{sample_pdf_base64}}",
  "signatures": [
    {
      "validations": {
        "BI": null
      }
    }
  ],
  "metadata": {
    "project": "alpha",
    "classification": "top_secret",
    "expires": "2025-12-31"
  }
}
```

### Caso 3: Acuerdo Comercial Multi-Parte

**Request:** Create Document - Multiple Signers

```json
{
  "title": "Alianza Estratégica 2024",
  "file_content": "{{sample_pdf_base64}}",
  "signatures": [
    {
      "signer_name": "CEO - Empresa A",
      "validations": {"EM": "ceo@empresa-a.com"},
      "invite_channel": ["EM"]
    },
    {
      "signer_name": "Director Legal - Empresa B",
      "validations": {
        "EM": "legal@empresa-b.com",
        "PH": "+5491134567890"
      },
      "invite_channel": ["EM"]
    },
    {
      "signer_name": "CFO - Empresa C",
      "validations": {"BI": null}
    }
  ],
  "metadata": {
    "contract_value": 1000000,
    "currency": "USD",
    "term_years": 3
  }
}
```

## 🔔 Configurar Webhooks

### 1. En tu Servidor

Crea un endpoint que acepte POST requests:

```javascript
app.post('/webhooks/signatura', (req, res) => {
  const event = req.body;
  console.log('Webhook recibido:', event);
  
  // Procesar evento
  if (event.notification_action === 'DS') {
    console.log('Documento firmado:', event.document_id);
  }
  
  res.status(200).send('OK');
});
```

### 2. En Signatura

1. Inicia sesión en [connect.signatura.co](https://connect.signatura.co)
2. Ve a **API → Webhooks**
3. Haz clic en **Crear Webhook**
4. Ingresa tu URL: `https://tu-servidor.com/webhooks/signatura`
5. Guarda

### 3. Probar con Postman

Los requests en la carpeta **🔔 Webhooks** muestran ejemplos de los payloads que recibirás.

## 🎨 Personalización

### Crear Variables Personalizadas

1. Ve al environment **🎨 Signatura - Production**
2. Agrega nuevas variables según necesites:
   - `default_signer_email`
   - `company_name`
   - `default_metadata`

### Usar Variables en Requests

```json
{
  "title": "Contrato - {{company_name}}",
  "signatures": [
    {
      "validations": {
        "EM": "{{default_signer_email}}"
      }
    }
  ]
}
```

## 🐛 Troubleshooting

### Error: "API key required"
- ✅ Verifica que hayas configurado `api_key` en el environment
- ✅ Asegúrate de haber seleccionado el environment correcto
- ✅ Revisa que el API key no tenga espacios al inicio o final

### Error: "Invalid file_content"
- ✅ Verifica que el PDF esté correctamente codificado en base64
- ✅ El PDF debe ser válido (sin errores)
- ✅ El tamaño recomendado es < 2MB

### Error: "Document not found"
- ✅ Verifica que el `document_id` sea correcto
- ✅ Asegúrate de estar usando la API key correcta
- ✅ El documento podría haber expirado

### Tests Fallan
- 🔍 Abre la Postman Console para ver detalles
- 🔍 Revisa el status code de la respuesta
- 🔍 Verifica el formato del JSON enviado

## 💡 Tips Avanzados

### 1. Ejecutar Colección Completa

Puedes ejecutar todos los requests en secuencia:

1. Haz clic derecho en la colección
2. Selecciona **Run collection**
3. Configura el orden y delays
4. Haz clic en **Run**

### 2. Exportar Resultados

Después de ejecutar requests:
1. Haz clic en **Runner**
2. Ve al historial de ejecuciones
3. Haz clic en **Export Results**

### 3. Ambiente de Testing

Crea un environment separado para testing:

1. Duplica **🎨 Signatura - Production**
2. Renómbralo a **🧪 Signatura - Testing**
3. Usa un API key diferente
4. Cambia el `base_url` si tienes un ambiente de pruebas

## 📊 Monitoreo

### Ver Uso de API

1. Abre **History** en Postman
2. Filtra por la colección de Signatura
3. Revisa tiempos de respuesta y errores

### Organizar Requests

Puedes crear subcarpetas:
- Por proyecto
- Por tipo de documento
- Por ambiente

## 🎓 Recursos Adicionales

- [Documentación Oficial de Signatura](https://docs.signatura.co)
- [API Reference](https://docs.signatura.co/reference)
- [Postman Learning Center](https://learning.postman.com)

---

¿Preguntas? Contacta a [help@signatura.co](mailto:help@signatura.co)
