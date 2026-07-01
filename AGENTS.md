# AGENTS.md

> AI Development Guide for **SnapReport**
>
> This document defines the engineering standards, architecture, coding conventions, and development rules that every AI coding assistant (ChatGPT, Claude, Gemini, Cursor, Cline, Roo Code, GitHub Copilot, etc.) **MUST** follow when contributing to this repository.

---

# Project Overview

SnapReport is a lightweight web application that converts screenshots into professional PDF reports.

The primary goal is to provide a fast, intuitive workflow for developers, QA engineers, IT teams, consultants, and cybersecurity professionals to generate documentation without manually editing Word documents.

The application prioritizes:

* Simplicity
* Performance
* Clean Architecture
* Maintainability
* Modular design
* Professional UI
* High-quality generated PDFs

---

# Core Engineering Principles

Every generated solution must follow these principles.

## Keep It Simple

Always choose the simplest implementation that satisfies the requirement.

Avoid unnecessary abstraction.

Avoid premature optimization.

---

## Readability First

Code is read more than it is written.

Prioritize

* readable code
* meaningful names
* small functions
* explicit logic

Never sacrifice readability for cleverness.

---

## Single Responsibility Principle

Every

* package
* struct
* component
* function

should have one responsibility.

---

## Composition over Inheritance

Prefer composition.

Avoid unnecessary interfaces.

---

## Build Small Components

Large files should be broken into reusable components.

Target sizes:

Frontend

* Components <250 lines

Backend

* Files <300 lines

Functions

* <40 lines whenever possible

---

# Technology Stack

## Frontend

* React
* TypeScript
* Vite
* TailwindCSS
* shadcn/ui
* React Hook Form
* Zod
* dnd-kit
* Axios

---

## Backend

* Go 1.24+
* Chi Router
* SQLite
* HTML Templates
* Headless Chromium (PDF)

---

## Future

* PostgreSQL
* MinIO
* Redis
* OCR
* AI Services

Do not implement future technologies unless explicitly requested.

---

# Architecture

The project follows **Clean Architecture**.

```text
HTTP

↓

Handler

↓

Service

↓

Repository

↓

Storage
```

Never bypass layers.

Handlers must never directly access repositories.

---

# Backend Folder Structure

```text
internal/

api/

config/

handlers/

middleware/

models/

repository/

services/

image/

pdf/

report/

templates/

utils/
```

Never introduce random folders.

---

# Frontend Folder Structure

```text
src/

components/

pages/

hooks/

services/

types/

utils/

assets/
```

Organize components by feature.

Avoid deeply nested folders.

---

# API Design

Follow REST conventions.

Examples

GET

POST

PUT

DELETE

Never use verbs in endpoints.

Correct

```text
POST /reports

GET /reports/{id}
```

Wrong

```text
/createReport

/getReports
```

---

# JSON Response Format

Success

```json
{
  "success": true,
  "data": {},
  "message": "Operation completed."
}
```

Failure

```json
{
  "success": false,
  "error": "Validation failed."
}
```

Maintain this format consistently.

---

# Error Handling

Never ignore errors.

Never panic.

Always return descriptive errors.

Wrap errors where appropriate.

Example

```go
return fmt.Errorf("generate pdf: %w", err)
```

---

# Logging

Use structured logging.

Never use fmt.Println().

Log levels

* INFO
* WARN
* ERROR

Never log

* passwords
* secrets
* tokens

---

# Dependency Injection

Inject dependencies.

Never instantiate services inside handlers.

Correct

```go
handler := NewReportHandler(service)
```

Wrong

```go
service := NewReportService()
```

inside every request.

---

# Database Rules

SQLite is used only for persistence.

Repositories own database logic.

Services must never execute SQL.

---

# Validation

Validate

* request body
* file size
* MIME type
* required fields

Validation occurs before business logic.

---

# File Upload Rules

Allowed

* PNG
* JPG
* JPEG
* WEBP

Maximum upload size should be configurable.

Reject invalid MIME types.

Never trust file extensions.

---

# Temporary Files

Uploaded files are temporary.

Delete them after report generation unless persistence is requested.

Never leave orphan files.

---

# PDF Generation

Always generate PDF using

HTML

↓

CSS

↓

Headless Chromium

Avoid manually drawing PDFs.

Templates should remain editable.

---

# Frontend Rules

Prefer

Functional Components

React Hooks

Never use class components.

---

# Styling Rules

Use TailwindCSS only.

Avoid inline styles.

Avoid custom CSS unless absolutely necessary.

Reuse shadcn components.

---

# Component Design

Each component should

* have one purpose
* receive typed props
* avoid unnecessary state

Extract reusable UI.

---

# State Management

Use React state.

Do not introduce Redux.

If global state becomes necessary

Use Context.

---

# Forms

Use

React Hook Form

*

Zod

Every form must validate both

client-side

and

server-side.

---

# Naming Conventions

## Go

Packages

lowercase

Structs

PascalCase

Interfaces

Reader

Writer

Generator

Avoid

IReportService

---

## TypeScript

Components

PascalCase

Variables

camelCase

Constants

UPPER_SNAKE_CASE

---

# Function Design

Functions should

* perform one task
* return early
* avoid nesting
* avoid side effects

Maximum preferred length

40 lines

---

# Comments

Only explain

WHY

Never explain

WHAT

Bad

```go
// Increment i

i++
```

Good

```go
// Retry because Chromium occasionally fails during startup.
```

---

# Security

Always

Validate input

Escape HTML

Limit uploads

Sanitize filenames

Protect against path traversal

Protect against XSS

Protect against CSRF (future)

Never trust user input.

---

# Performance

Optimize only when necessary.

Avoid

Premature caching

Premature concurrency

Premature optimization

Measure first.

---

# Accessibility

Every interactive element should

* have labels
* support keyboard navigation
* provide focus states

Images should contain alt text.

---

# Testing

Backend

Table-driven tests.

Frontend

React Testing Library.

Write tests for

* services
* repositories
* validators

Avoid testing implementation details.

---

# Git Commit Convention

Use Conventional Commits.

Examples

```text
feat: add screenshot upload

fix: correct pdf page numbering

refactor: simplify report service

docs: update README

test: add upload validation tests
```

---

# Pull Request Checklist

Every generated code change must

* compile
* pass linting
* pass tests
* follow architecture
* include documentation if needed

---

# Things AI Must Never Do

Do NOT

* rewrite unrelated files
* introduce unnecessary frameworks
* add dependencies without reason
* change project architecture
* over-engineer solutions
* create giant files
* duplicate code
* ignore errors
* hardcode secrets
* commit generated binaries
* invent APIs that do not exist

---

# Preferred Development Workflow

For every feature:

1. Understand the requirement.
2. Identify impacted layers.
3. Implement backend.
4. Implement frontend.
5. Add validation.
6. Handle errors.
7. Add tests.
8. Update documentation.
9. Ensure the project builds successfully.

Never skip validation.

---

# Definition of Done

A feature is considered complete only if:

* It compiles successfully.
* It follows Clean Architecture.
* It passes all tests.
* It includes proper validation.
* It handles errors gracefully.
* It has no obvious code duplication.
* It follows project naming conventions.
* It is documented where appropriate.
* It does not introduce unnecessary complexity.

---

# Guiding Philosophy

> **Build software that is simple to understand, easy to maintain, and pleasant to extend.**

Every contribution should improve the codebase—not merely add functionality.

When multiple solutions exist, prefer the one that is:

1. Simpler
2. More readable
3. Easier to maintain
4. Easier to test
5. Easier to extend
6. Consistent with the existing architecture
