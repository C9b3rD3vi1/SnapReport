# SnapReport

Convert screenshots into professional PDF reports.

## Tech Stack

- **Frontend:** React, TypeScript, Vite, TailwindCSS, shadcn/ui
- **Backend:** Go 1.24+, Chi Router, SQLite
- **PDF:** HTML + CSS → Headless Chromium (via chromedp)

## Development

### Prerequisites

- Go 1.24+
- Node.js 22+
- Chromium (for PDF generation)

### Backend

```bash
make run
```

Server starts on `http://localhost:8080`.

### Frontend

```bash
cd web
npm run dev
```

Dev server starts on `http://localhost:5173` with API proxy to backend.

### Docker

```bash
make docker-up
```

### Testing

```bash
# Backend tests
make test

# Frontend tests
cd web && npm test
```

## API

All responses follow the format:

```json
{ "success": true, "data": {}, "message": "..." }
```

```json
{ "success": false, "error": "..." }
```

### Health

```
GET /health
```

### Uploads

#### Upload files

```
POST /api/v1/uploads
Content-Type: multipart/form-data

files: (binary, PNG/JPG/WEBP, max 10 MB each)
```

Response `201`:

```json
{
  "success": true,
  "data": [
    {
      "id": "abc123",
      "filename": "abc123.png",
      "original_name": "screenshot.png",
      "mime_type": "image/png",
      "size": 123456,
      "title": "",
      "description": "",
      "order_index": 0,
      "created_at": "2026-07-01T12:00:00Z"
    }
  ]
}
```

#### List uploads

```
GET /api/v1/uploads
```

#### Delete upload

```
DELETE /api/v1/uploads/{id}
```

### Reports

#### Create report

```
POST /api/v1/reports
Content-Type: application/json

{
  "title": "Test Report",
  "project": "SnapReport",
  "company": "ACME",
  "author": "John",
  "version": "1.0",
  "uploads": [
    {
      "id": "abc123",
      "title": "Login Screen",
      "description": "The login page with credentials form",
      "notes": "",
      "order_index": 0
    }
  ]
}
```

Response `201`:

```json
{
  "success": true,
  "data": {
    "id": "report-xyz",
    "title": "Test Report",
    "project": "SnapReport",
    "company": "ACME",
    "author": "John",
    "version": "1.0",
    "status": "completed",
    "pdf_path": "./generated/report-xyz.pdf",
    "created_at": "2026-07-01T12:00:00Z",
    "updated_at": "2026-07-01T12:00:00Z"
  }
}
```

#### Get report

```
GET /api/v1/reports/{id}
```

Response includes report details and associated uploads.

#### Download PDF

```
GET /api/v1/reports/{id}/download
```

Returns the generated PDF as a binary download.

## Project Structure

```
cmd/server/           # Go entry point
internal/
  api/                # Router
  config/             # Environment configuration
  handlers/           # HTTP handlers
  middleware/         # Logging, CORS
  models/             # Data models
  repository/         # SQLite persistence layer
  services/           # Business logic (upload, report)
    image/            # Image validation
    report/           # Report creation + PDF generation
  pdf/                # HTML template + Chromium PDF generator
  utils/              # Response helpers, file utilities
web/                  # React frontend
  src/
    components/       # Reusable UI components
    pages/            # UploadPage, ReportEditorPage
    services/         # API client
    types/            # TypeScript interfaces
    utils/            # cn() helper
uploads/              # Temporary uploaded files
generated/            # Generated PDFs
```
