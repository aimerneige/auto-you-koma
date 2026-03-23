# Auto Yon Koma | オート四コマ (English)

**Auto Yon Koma** is a fully automated multi-agent 4-panel comic (Yon-Koma) generation system. Utilizing React for the frontend and Go (Gin) for the backend, this project leverages multimodal AI agents to collaboratively design characters, write scripts, direct storyboards, generate images, and composite the final comic strip from a single spark of inspiration!

[**English**] | [**简体中文**](./README.md)

---

## 🚀 Features

- **Multi-Agent Collaboration**: 8 specialized AI Agents (Director, Typesetter, Reviewer, etc.) work together systematically.
- **Interactive Previz**: Chat-style script and storyboard editing interface designed for seamless Human-in-The-Loop (HITL) intervention.
- **Highly Customizable Generators**: Built-in support for Gemini, Nano Banana 2, and flexibly layered integrations for other LLM & Diffusion engines.
- **Modern Tech Stack**: Fullstack architecture featuring Go (Gin) and React (Vite).
- **1-Click Deployment**: Polished Docker Compose setup for an out-of-the-box experience.

## 🛠 Tech Stack

- **Frontend**: React, TypeScript, Vite, Zustand
- **Backend**: Go, Gin, GORM
- **Database**: SQLite (default), optionally PostgreSQL / MySQL
- **Message Queue**: Go Channel (default), optionally Redis / RabbitMQ
- **Infrastructure**: Docker, Docker Compose

## 📦 Quick Start

### Prerequisites
- [Docker & Docker Compose](https://docs.docker.com/get-docker/) installed.
- Or modern Node.js (v20+) and Go (v1.23+) for local development.

### Start up with Docker
```bash
git clone https://github.com/aimerneige/auto-you-koma.git
cd auto-you-koma
docker-compose up -d --build
```
Then visit `http://localhost:3000` and start creating your first Yon-Koma!

## 🎨 Example

> 🚧 Work In Progress, AI-generated comic showcases will be displayed here soon. 🚧

## 🤝 Contributing

Feel free to submit a Pull Request or report an issue. Refer to the [Development Plan](./development_plan.md) to catch up on our roadmap. Let's build a flawless AI comic creation system together!

## 📄 License

This project is licensed under the **[MIT License](https://github.com/aimerneige/auto-you-koma/blob/master/LICENSE)**.
