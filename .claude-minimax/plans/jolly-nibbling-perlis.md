# Plan: HITL Yon Koma - Continue Development

## Context
I have completed initial setup:
- Task 1.1: Backend Go + Gin project with basic server
- Task 1.2: Frontend React + Vite + TypeScript project

Now need to continue with remaining tasks from task.md.

## Implementation Plan

### Phase 1: Database & Core Interfaces
1. **Task 1.3: Database Schema (SQLite)**
   - Create tables: users, projects, characters, scripts, storyboards, panels
   - Key: `characters` table needs `reference_sheet_url` field

2. **Task 1.4: Core Interface Abstraction**
   - Define `LLMTextClient` interface (Gemini)
   - Define `LLMImageClient` interface (Nano Banana 2) with `ReferenceImage` and `ImagePromptWeight` fields

### Phase 2: Character System
3. **Task 2.1: Character CRUD API & UI**
   - Backend: character CRUD endpoints
   - Frontend: character list and detail pages

4. **Task 2.2: Reference Sheet Generation (HITL Node 1)**
   - Backend: API to generate character reference sheet
   - Frontend: "Generate Reference Sheet" button, confirm/retry flow

5. **Task 2.3: Asset Storage Middleware**
   - Local file system or OSS upload for reference images

### Phase 3: Script & Storyboard Pipeline
6. **Task 3.1: Script Generation (HITL Node 2)**
   - Backend: Generate 4-panel script from prompt + characters
   - Frontend: Script editor UI with editable JSON

7. **Task 3.2: Storyboard Generation (HITL Node 3)**
   - Backend: Convert script to detailed storyboard JSON
   - Frontend: Previz chat-style preview UI

### Phase 4: Image Generation Engine
8. **Task 4.1: Render Task Assembler**
   - Backend: Assemble reference images with storyboard prompts

9. **Task 4.2: Single Panel Generation**
   - Backend: Call image generation API for each panel

10. **Task 4.3: Quality Check UI (HITL Node 4)**
    - Frontend: Render approval UI with "Regenerate Panel" option

### Phase 5: Layout & Export
11. **Task 5.1: Text Coordinate Generation**
    - Backend: Calculate speech bubble positions

12. **Task 5.2: Image Compositing**
    - Backend: Overlay text on images, stitch 2x2 or 1x4 layout

13. **Task 5.3: Export**
    - Backend: Download API for final output

## Critical Files
- Backend: `backend/cmd/server/main.go`
- Frontend: `frontend/src/App.tsx`
- Demo design: `demo/demo-anime.html`

## Verification
- Run backend: `cd backend && go run cmd/server/main.go`
- Run frontend: `cd frontend && npm run dev`
- Test API endpoints with curl or frontend UI