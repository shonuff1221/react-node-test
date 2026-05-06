# Solutions — PropVision Senior SWE Skill Test

## Challenge 1: Fix "Add New Notice" Page (Frontend)

**File:** `frontend/src/domains/notice/components/notice-form.tsx`
**File:** `frontend/src/domains/notice/pages/add-notice-page.tsx`

### Problem
When clicking 'Save' on the Add New Notice form, the description field was not submitted to the API.

### Root Cause
The form's description `TextField` was registered with React Hook Form as `content`:
```tsx
<TextField {...register('content')} />
```

However, the Zod validation schema (`NoticeFormSchema`), the TypeScript types (`NoticeFormProps`), the backend API, and the edit/view pages all expect the field to be named `description`. This caused:
1. Zod validation never validated the field (it validated `description`, but the form had `content`)
2. The error helper `errors.description` was always empty, masking the issue
3. The API received `{ content: "..." }` instead of `{ description: "..." }`, so the backend ignored the description value

### Fix
- Changed `register('content')` → `register('description')` in `notice-form.tsx`
- Updated `initialState` in `add-notice-page.tsx` from `content: ''` → `description: ''`

---

## Challenge 2: Complete Student CRUD Operations (Backend)

**File:** `backend/src/modules/students/students-controller.js`

### Problem
All five controller handler functions were empty stubs — `//write your code`.

### Approach
The service layer (`students-service.js`) and repository layer (`students-repository.js`) were fully implemented with proper error handling, validation, and database queries. The task was to wire the Express request/response cycle to the existing service functions, following the same patterns used in the staff controller.

### Implementation

| Endpoint | Method | Handler | Notes |
|---|---|---|---|
| `/api/v1/students` | GET | `handleGetAllStudents` | Returns `{ students: [...] }` to match frontend RTK Query types. Passes `req.query` for filtering (name, class, section, roll) |
| `/api/v1/students` | POST | `handleAddStudent` | Creates student + sends verification email via `req.body` |
| `/api/v1/students/:id` | GET | `handleGetStudentDetail` | Returns full student profile by ID |
| `/api/v1/students/:id` | PUT | `handleUpdateStudent` | Merges `req.params.id` into payload (frontend strips `id` from body, only sends via URL param) |
| `/api/v1/students/:id/status` | POST | `handleStudentStatus` | Toggles active/inactive, uses `req.user.id` from JWT middleware as reviewer |

### Key Details
- **Response shape**: The `GET /students` endpoint wraps results in `{ students }` to match the frontend `StudentData` type — the same pattern used by the staff controller
- **Update payload**: The frontend's `updateStudent` mutation destructures `{ id, ...payload }` and only sends `payload` in the body, so the controller merges `req.params.id` back in
- **Status reviewer**: `handleStudentStatus` extracts `req.user.id` from the JWT auth middleware as the reviewer, matching the staff status handler pattern

---

## Bonus: Go PDF Report Service

**Files:** `go-service/main.go`, `backend/src/modules/students/report-controller.js`, `backend/src/modules/students/sudents-router.js`

### What I Built
A standalone Go microservice that generates downloadable PDF reports for students, with a Node.js proxy route to integrate it into the existing API.

### Architecture
- **Go service** (`go-service/main.go`): Connects directly to PostgreSQL, fetches student details, generates a professional PDF using `gofpdf`, and streams it as a downloadable file
- **Node proxy** (`report-controller.js`): The existing Express app proxies `GET /api/v1/students/:id/report` to the Go service, with proper error handling for service unavailability (503) and missing students (404)
- **Route registration**: Added `/:id/report` route before `/:id` to prevent Express from matching the report path as a student ID

### Setup
```bash
# Terminal 1: Start Go service
cd go-service
export DATABASE_URL="postgresql://postgres:postgres@localhost:5432/school_mgmt?sslmode=disable"
go run main.go

# Terminal 2: Start Node backend (add GO_SERVICE_URL to .env)
cd backend
echo "GO_SERVICE_URL=http://localhost:8080" >> .env
npm start
```

### Why Go?
The README specified a Go microservice for PDF generation. Using Go for CPU-intensive tasks like PDF rendering is a good architectural pattern — it keeps the Node.js event loop free while the Go service handles the heavy lifting.
