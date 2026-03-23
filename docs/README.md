# Auto Yon Koma | オート四コマ

**Auto Yon Koma** 是一款全自动的多智能体四格漫画（Yon-Koma）生成系统。它前端使用 React，后端使用 Go (Gin)，借助多模态 AI 模型进行角色设计、剧本编写、分镜编排、图像生成与排版，将简单的灵感一键转化为高质量的有分镜四格视觉故事！

[**English**](./README_en.md) | [**简体中文**]

---

## 🚀 特性 (Features)

- **多智能体协同**：拥有 8 个专职的 AI Agent（主笔、导演、排版等）各司其职。
- **互动式预览**：聊天式的剧本与分镜编辑页面，方便人工干预调整。
- **高度定制化的渲染策略**：支持灵活的多引擎接入，内置 Gemini 与 Nano Banana 2 等生成链路。
- **现代化技术栈**：极速的前后端分离架构 (Go Gin + React Vite)。
- **一键式部署**：完善的多阶段 Docker 支持，开箱即用。

## 🛠 开发栈 (Tech Stack)

- **前端**: React, TypeScript, Vite, Zustand
- **后端**: Go, Gin, GORM
- **数据库**: SQLite (默认), 可选 PostgreSQL / MySQL
- **消息队列**: Go Channel (默认), 可选 Redis / RabbitMQ
- **部署**: Docker, Docker Compose

## 📦 快速开始 (Quick Start)

### 前置要求
- [Docker & Docker Compose](https://docs.docker.com/get-docker/) 安装完毕。
- 或安装了最新的 Node.js (v20+) 与 Go (v1.23+)。

### 使用 Docker 启动
```bash
git clone https://github.com/aimerneige/auto-you-koma.git
cd auto-you-koma
docker-compose up -d --build
```
然后访问 `http://localhost:3000` 即可开始使用！

## 🎨 示例展示 (Example)

> 🚧 功能开发中，即将在此处展示 AI 自动生成的短片漫画。 🚧

## 🤝 贡献 (Contributing)

欢迎提交 Pull Request 或报告 Issue。请参考 [开发计划](./development_plan.md) 了解我们的进度与路线图。让我们一起构建更好用的 AI 漫画创作系统！

## 📄 许可证 (License)

本项目采用 **[MIT License](https://github.com/aimerneige/auto-you-koma/blob/master/LICENSE)**。
