# SnapReport

**Technical Report Builder** — Turn screenshots into professional PDF reports.

## Tech Stack

- **Frontend:** React, TypeScript, Vite, TailwindCSS, shadcn/ui
- **Backend:** Go 1.26+, Chi Router, SQLite
- **PDF:** HTML + CSS → Headless Chromium (chromedp)

## Quick Start

```bash
# Backend
make run                          # Starts on :8080

# Frontend (separate terminal)
cd web && npm run dev             # Starts on :5173 with API proxy

# Docker
make docker-up                    # Production build
```

## Features

- Drag-and-drop upload + paste from clipboard
- Auto-title from filename
- Thumbnail generation (300px previews)
- Quick PDF generation (2 clicks)
- Drag-and-drop reordering with keyboard support
- Inline metadata editing (title, description)
- Professional PDF with cover page, figure numbers, headers, footers
- Report history with search and delete
- Configurable security: rate limiting, CORS, security headers

## API

### Health

```
GET /health
```

Response: `{ "success": true, "data": { "status": "healthy", "version": "1.0.0", "uptime": "5m2s" } }`

### Uploads

```
POST   /api/v1/uploads           # Upload files (multipart, field: "files")
GET    /api/v1/uploads           # List uploads
DELETE /api/v1/uploads/{id}      # Delete upload
```

Upload response includes `thumbnail_url` for preview.

### Reports

```
GET    /api/v1/reports            # List reports
POST   /api/v1/reports            # Create report + generate PDF
GET    /api/v1/reports/{id}       # Get report details + uploads
GET    /api/v1/reports/{id}/download  # Download PDF
DELETE /api/v1/reports/{id}       # Delete report
```

### Response Format

```json
{ "success": true, "data": { ... } }
{ "success": false, "error": "message" }
```

### Error Codes

| Code | HTTP | When |
|---|---|---|
| `VALIDATION_ERROR` | 400 | Missing/invalid fields |
| `NOT_FOUND` | 404 | Resource not found |
| `PAYLOAD_TOO_LARGE` | 413 | Upload exceeds size limit |
| `RATE_LIMITED` | 429 | Too many requests |
| `INTERNAL_ERROR` | 500 | Server error |

## Configuration

| Env | Default | Description |
|---|---|---|
| `PORT` | 8080 | Server port |
| `DATABASE_PATH` | ./data/snapreport.db | SQLite database path |
| `UPLOAD_DIR` | ./uploads | Uploaded file storage |
| `PDF_DIR` | ./generated | Generated PDF storage |
| `MAX_UPLOAD_MB` | 10 | Max file size per upload |
| `MAX_REQUEST_BODY_MB` | 50 | Max total request body |
| `ALLOWED_ORIGINS` | * | CORS allowed origins (comma-separated) |
| `ENVIRONMENT` | development | Runtime environment |

## Project Structure

```
cmd/server/            # Entry point
internal/
  api/                 # Router
  config/              # Environment configuration
  fileutil/            # File operations + ID generation
  handler/             # HTTP handlers
  imgvalidator/        # Image MIME validation via magic bytes
  middleware/          # Logging, CORS, rate limiting, security headers
  model/               # Data models
  pdf/                 # PDF generation (chromedp + HTML template)
  repository/          # SQLite persistence
  response/            # JSON response helpers
  service/             # Interfaces + sentinel errors
  services/            # Business logic (upload, thumbnail, report)
web/
  src/
    components/        # Reusable UI components
    pages/             # UploadPage, ReportEditorPage, ReportHistoryPage, ReportDetailPage
    services/          # Axios API client
    types/             # TypeScript interfaces
```

## Testing

```bash
make test              # Backend tests
cd web && npm test     # Frontend tests
```
