# HITL Yon Koma - 编译运行测试指南

## 前置要求

- Go 1.26+
- Node.js 18+
- npm 或 yarn

## 后端编译运行

### 1. 进入后端目录

```bash
cd backend
```

### 2. 下载依赖

```bash
go mod tidy
```

### 3. 编译

```bash
go build -o server ./cmd/server/
```

### 4. 运行

```bash
./server
```

或直接运行：

```bash
go run cmd/server/main.go
```

### 5. 测试 API

健康检查：

```bash
curl http://localhost:8080/health
```

创建角色：

```bash
curl -X POST http://localhost:8080/api/v1/characters \
  -H "Content-Type: application/json" \
  -d '{"name": "测试角色", "visual_prompt": "蓝色头发"}'
```

## 前端编译运行

### 1. 进入前端目录

```bash
cd frontend
```

### 2. 安装依赖

```bash
npm install
```

### 3. 开发模式运行

```bash
npm run dev
```

访问 http://localhost:3000

### 4. 生产构建

```bash
npm run build
```

构建产物在 `dist/` 目录

## 完整工作流测试

### 1. 启动后端

```bash
cd backend
go run cmd/server/main.go
```

### 2. 启动前端（新终端）

```bash
cd frontend
npm run dev
```

### 3. 访问界面

打开浏览器访问 http://localhost:3000

### 4. 测试流程

1. **创建角色**: 点击 Characters → New Character → 填写信息 → Save
2. **生成参考图**: 编辑角色页面 → Generate Reference Sheet → 确认保存
3. **创建项目**: Projects → New Project → 填写标题和概要
4. **生成剧本**: 进入项目 → Generate Script → 编辑台词 → Confirm
5. **生成分镜**: Generate Storyboard → 调整参数 → Confirm
6. **渲染图片**: Start Rendering → 质量检查 → Regenerate 或 Confirm
7. **导出**: Export 按钮下载成品

## API 端点汇总

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /api/v1/characters | 创建角色 |
| GET | /api/v1/characters | 角色列表 |
| POST | /api/v1/characters/:id/generate-reference | 生成参考图 |
| POST | /api/v1/characters/:id/confirm-reference | 确认参考图 |
| POST | /api/v1/projects | 创建项目 |
| GET | /api/v1/projects | 项目列表 |
| POST | /api/v1/projects/:id/generate-script | 生成剧本 |
| POST | /api/v1/projects/:id/generate-storyboard | 生成分镜 |
| POST | /api/v1/projects/:id/render | 开始渲染 |
| POST | /api/v1/projects/:id/render/regenerate | 重绘单格 |
| POST | /api/v1/projects/:id/render/confirm | 确认渲染 |
| GET | /api/v1/projects/:id/export | 导出项目 |

## 注意事项

- 当前使用 Mock LLM 生成器，返回模拟数据
- 实际部署时需要配置真实的 Gemini 和 Nano Banana API
- 数据库文件会自动创建在 `backend/data/auto-you-koma.db`