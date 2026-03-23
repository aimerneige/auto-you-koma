# Auto Yon Koma 开发计划

## 一、项目总览

Auto Yon Koma 是一个 AI 多智能体驱动的四格漫画全自动生成系统。前端使用 **React**，后端使用 **Go (Gin)**，通过 8 个专门的 AI Agent 协作完成从角色设定、剧本分镜到图像生成与排版的全流程漫画创作。

### 技术栈

| 层级     | 技术选型                                                                           |
| -------- | ---------------------------------------------------------------------------------- |
| 前端     | React + TypeScript + Vite                                                          |
| 后端     | Go + Gin                                                                           |
| 数据库   | 抽象接口，默认 SQLite，可选 PostgreSQL / MySQL                                     |
| 消息队列 | 抽象接口，默认 Go Channel，可选 Redis / RabbitMQ                                   |
| AI 模型  | 抽象适配层，默认 Gemini / Nano Banana 2，可选 Deepseek / GPT / Stable Diffusion 等 |
| 部署     | Docker 容器化                                                                      |

---

## 二、系统架构设计

### 2.1 整体架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                        Frontend (React)                         │
│  ┌──────────┬──────────┬──────────┬──────────┬────────────────┐ │
│  │  Auth    │ Character│ Script   │ Previz   │ Render/Export  │ │
│  │  Module  │ Manager  │ Editor   │ Chat UI  │ Settings       │ │
│  └──────────┴──────────┴──────────┴──────────┴────────────────┘ │
└────────────────────────────┬────────────────────────────────────┘
                             │ REST API + WebSocket (SSE)
┌────────────────────────────▼────────────────────────────────────┐
│                        Backend (Go + Gin)                       │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                     API Layer (Gin Router)                │   │
│  │  /api/v1/auth  /api/v1/characters  /api/v1/projects ...  │   │
│  └──────────────────────┬───────────────────────────────────┘   │
│  ┌──────────────────────▼───────────────────────────────────┐   │
│  │                   Service Layer                           │   │
│  │  AuthService │ CharacterService │ ProjectService │ ...    │   │
│  └──────────────────────┬───────────────────────────────────┘   │
│  ┌──────────────────────▼───────────────────────────────────┐   │
│  │              Agent Orchestrator (Pipeline)                │   │
│  │  ┌─────────┬─────────┬─────────┬─────────┬────────────┐  │   │
│  │  │Character│Script   │Story    │Asset    │Generation  │  │   │
│  │  │Agent    │Agent    │Board    │Manager  │Engine      │  │   │
│  │  │         │         │Agent    │Agent    │Agent       │  │   │
│  │  ├─────────┼─────────┼─────────┼─────────┼────────────┤  │   │
│  │  │QC       │Type     │Compos   │         │            │  │   │
│  │  │Reviewer │setter   │itor    │         │            │  │   │
│  │  │Agent    │Agent    │Agent    │         │            │  │   │
│  │  └─────────┴─────────┴─────────┴─────────┴────────────┘  │   │
│  └──────────────────────┬───────────────────────────────────┘   │
│  ┌──────────────────────▼───────────────────────────────────┐   │
│  │              Abstraction Layer                            │   │
│  │  ┌──────────────┐  ┌─────────────┐  ┌─────────────────┐ │   │
│  │  │ LLM Adapter  │  │ Queue       │  │ DB Repository   │ │   │
│  │  │ Interface    │  │ Interface   │  │ Interface       │ │   │
│  │  │              │  │             │  │                 │ │   │
│  │  │ ◆ Gemini     │  │ ◆ Channel   │  │ ◆ SQLite        │ │   │
│  │  │ ◆ NanoBanana │  │ ◆ Redis     │  │ ◆ PostgreSQL    │ │   │
│  │  │ ◆ OpenAI     │  │ ◆ RabbitMQ  │  │ ◆ MySQL         │ │   │
│  │  │ ◆ Deepseek   │  │             │  │                 │ │   │
│  │  │ ◆ StableDiff │  │             │  │                 │ │   │
│  │  └──────────────┘  └─────────────┘  └─────────────────┘ │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
                             │
              ┌──────────────▼──────────────┐
              │    File Storage (Images)    │
              │    /data/assets/characters/ │
              │    /data/assets/outputs/    │
              └─────────────────────────────┘
```

### 2.2 核心抽象接口设计

#### 2.2.1 LLM 适配层

Agent 调度器与 LLM 完全解耦。LLM 对 Agent 来说只是一个**将需求（Prompt）转化为结果（Response）的抽象算法**。

```go
// LLM 文本生成接口
type TextGenerator interface {
    Generate(ctx context.Context, req TextRequest) (*TextResponse, error)
    GenerateStream(ctx context.Context, req TextRequest) (<-chan TextChunk, error)
}

// LLM 图像生成接口
type ImageGenerator interface {
    GenerateImage(ctx context.Context, req ImageRequest) (*ImageResponse, error)
}

// 多模态分析接口（用于 QC Agent）
type VisionAnalyzer interface {
    Analyze(ctx context.Context, req VisionRequest) (*VisionResponse, error)
}
```

默认实现：
- `GeminiTextGenerator` — 调用 Google Gemini API
- `NanoBananaImageGenerator` — 调用 Nano Banana 2 API
- `OpenAITextGenerator` — 调用 OpenAI GPT API
- `DeepseekTextGenerator` — 调用 Deepseek API
- `StableDiffusionImageGenerator` — 调用 Stable Diffusion API

#### 2.2.2 消息队列抽象

```go
type TaskQueue interface {
    Publish(ctx context.Context, topic string, message *TaskMessage) error
    Subscribe(ctx context.Context, topic string) (<-chan *TaskMessage, error)
    Ack(ctx context.Context, messageID string) error
}
```

默认实现：
- `ChannelQueue` — 基于 Go Channel 的进程内队列
- `RedisQueue` — 基于 Redis Stream
- `RabbitMQQueue` — 基于 RabbitMQ

#### 2.2.3 数据库 Repository 抽象

```go
type CharacterRepository interface {
    Create(ctx context.Context, character *Character) error
    GetByID(ctx context.Context, id string) (*Character, error)
    List(ctx context.Context, filter CharacterFilter) ([]*Character, error)
    Update(ctx context.Context, character *Character) error
    Delete(ctx context.Context, id string) error
    Search(ctx context.Context, query string) ([]*Character, error)
}

type ProjectRepository interface { /* 类似结构 */ }
type UserRepository interface { /* 类似结构 */ }
```

默认实现：
- `SQLiteCharacterRepo` / `PostgresCharacterRepo` / `MySQLCharacterRepo`

---

## 三、后端详细设计

### 3.1 目录结构

```
backend/
├── cmd/
│   └── server/
│       └── main.go                 # 入口，初始化依赖注入
├── internal/
│   ├── config/
│   │   └── config.go               # 配置加载（YAML/ENV）
│   ├── middleware/
│   │   ├── auth.go                  # JWT 认证中间件
│   │   ├── ratelimit.go             # 用量限制中间件
│   │   └── cors.go                  # 跨域中间件
│   ├── handler/                     # Gin Handler（API 层）
│   │   ├── auth_handler.go
│   │   ├── character_handler.go
│   │   ├── project_handler.go
│   │   └── pipeline_handler.go
│   ├── service/                     # 业务逻辑层
│   │   ├── auth_service.go
│   │   ├── character_service.go
│   │   ├── project_service.go
│   │   └── pipeline_service.go
│   ├── repository/                  # 数据访问层（抽象 + 实现）
│   │   ├── interfaces.go            # 所有 Repository 接口定义
│   │   ├── sqlite/
│   │   │   ├── character_repo.go
│   │   │   ├── project_repo.go
│   │   │   └── user_repo.go
│   │   ├── postgres/
│   │   │   └── ...
│   │   └── mysql/
│   │       └── ...
│   ├── model/                       # 数据模型
│   │   ├── user.go
│   │   ├── character.go
│   │   ├── project.go
│   │   ├── script.go
│   │   └── storyboard.go
│   ├── agent/                       # AI Agent 系统
│   │   ├── interfaces.go            # Agent 抽象接口
│   │   ├── orchestrator.go          # Agent 编排调度器
│   │   ├── character_agent.go       # Agent 1 - 设定主笔
│   │   ├── scriptwriter_agent.go    # Agent 2 - 剧本总监
│   │   ├── storyboard_agent.go      # Agent 3 - 分镜导演
│   │   ├── asset_manager_agent.go   # Agent 4 - 视觉资产管家
│   │   ├── generation_agent.go      # Agent 5 - 主笔画师
│   │   ├── qc_agent.go              # Agent 6 - 质量审校
│   │   ├── typesetter_agent.go      # Agent 7 - 嵌字排版师
│   │   └── compositor_agent.go      # Agent 8 - 拼版总集
│   ├── llm/                         # LLM 适配层
│   │   ├── interfaces.go            # TextGenerator / ImageGenerator / VisionAnalyzer
│   │   ├── gemini.go
│   │   ├── nanobanana.go
│   │   ├── openai.go
│   │   ├── deepseek.go
│   │   └── stablediffusion.go
│   ├── queue/                       # 消息队列
│   │   ├── interfaces.go
│   │   ├── channel_queue.go
│   │   ├── redis_queue.go
│   │   └── rabbitmq_queue.go
│   └── auth/                        # 认证模块
│       ├── jwt.go
│       ├── totp.go                  # 2FA TOTP
│       └── email.go                 # 邮件验证
├── migrations/                      # 数据库迁移文件
│   ├── 001_create_users.sql
│   ├── 002_create_characters.sql
│   └── 003_create_projects.sql
├── go.mod
├── go.sum
├── Makefile
└── Dockerfile
```

### 3.2 数据模型设计

#### 3.2.1 用户模型 (User)

```
users
├── id              UUID   PK
├── email           TEXT   UNIQUE NOT NULL
├── password_hash   TEXT   NOT NULL
├── display_name    TEXT
├── totp_secret     TEXT            -- 2FA 密钥（加密存储）
├── totp_enabled    BOOLEAN DEFAULT FALSE
├── quota_limit     INT    DEFAULT 100   -- 每日 AI 调用额度
├── quota_used      INT    DEFAULT 0
├── created_at      TIMESTAMP
└── updated_at      TIMESTAMP

user_usage_logs
├── id              UUID   PK
├── user_id         UUID   FK -> users.id
├── action          TEXT            -- 操作类型 (generate_script / generate_image / ...)
├── model_used      TEXT            -- 使用的模型
├── tokens_consumed INT             -- token 消耗量
├── cost_estimate   DECIMAL         -- 预估费用
├── request_payload JSONB           -- 请求内容摘要
├── created_at      TIMESTAMP
└── project_id      UUID   FK -> projects.id (nullable)
```

#### 3.2.2 角色模型 (Character)

角色资产分为**文字维度**和**图片维度**两个层面。

```
characters
├── id              UUID   PK
├── user_id         UUID   FK -> users.id
├── name            TEXT   NOT NULL
├── name_jp         TEXT            -- 日文名
├── gender          TEXT
├── age             TEXT
├── personality     TEXT            -- 性格关键词 (JSON Array)
├── backstory       TEXT            -- 背景故事
├── visual_prompt   TEXT            -- 基础外貌描述 Prompt
├── tags            TEXT            -- 分类标签 (JSON Array, 用于检索)
├── category        TEXT            -- 分类 (如：原创 / LoveLive / ...)
├── created_at      TIMESTAMP
└── updated_at      TIMESTAMP

character_variants                  -- 角色变体（不同服装/性格偏向）
├── id              UUID   PK
├── character_id    UUID   FK -> characters.id
├── variant_name    TEXT            -- 变体名 (如："夏日制服"、"战斗装")
├── personality_mod TEXT            -- 性格偏向修正
├── visual_prompt   TEXT            -- 该变体的外貌 Prompt 覆盖
├── created_at      TIMESTAMP
└── updated_at      TIMESTAMP

character_images                    -- 角色图片资产
├── id              UUID   PK
├── character_id    UUID   FK -> characters.id
├── variant_id      UUID   FK -> character_variants.id (nullable)
├── image_type      TEXT            -- 类型: avatar / full_body / chibi / expression / reference
├── file_path       TEXT   NOT NULL -- 图片存储路径 (相对于 data 目录)
├── description     TEXT            -- 图片描述
├── is_primary      BOOLEAN DEFAULT FALSE -- 是否为主展示图
├── created_at      TIMESTAMP
└── updated_at      TIMESTAMP

character_groups                    -- 角色组合/分组
├── id              UUID   PK
├── user_id         UUID   FK -> users.id
├── group_name      TEXT   NOT NULL
├── description     TEXT
├── created_at      TIMESTAMP
└── updated_at      TIMESTAMP

character_group_members
├── group_id        UUID   FK -> character_groups.id
├── character_id    UUID   FK -> characters.id
└── sort_order      INT   DEFAULT 0
```

#### 3.2.3 项目模型 (Project)

```
projects
├── id              UUID   PK
├── user_id         UUID   FK -> users.id
├── title           TEXT   NOT NULL
├── mode            TEXT            -- standalone / serialized
├── status          TEXT            -- draft / scripted / previewed / rendering / done
├── synopsis        TEXT            -- 用户输入的情景梗概
├── created_at      TIMESTAMP
└── updated_at      TIMESTAMP

project_characters                  -- 项目使用的角色
├── project_id      UUID   FK -> projects.id
├── character_id    UUID   FK -> characters.id
└── variant_id      UUID   FK -> character_variants.id (nullable)

scripts                             -- AI 生成的剧本
├── id              UUID   PK
├── project_id      UUID   FK -> projects.id
├── episode_number  INT    DEFAULT 1
├── content         JSONB           -- Script_Outline.json
├── version         INT    DEFAULT 1
├── created_at      TIMESTAMP
└── updated_at      TIMESTAMP

storyboards                         -- 分镜脚本
├── id              UUID   PK
├── script_id       UUID   FK -> scripts.id
├── content         JSONB           -- Storyboard_Render_Payload.json
├── version         INT    DEFAULT 1
├── created_at      TIMESTAMP
└── updated_at      TIMESTAMP

render_tasks                        -- 渲染任务
├── id              UUID   PK
├── project_id      UUID   FK -> projects.id
├── storyboard_id   UUID   FK -> storyboards.id
├── export_type     TEXT            -- native_text / clean_plate
├── layout          TEXT            -- 2x2 / 1x4
├── image_width     INT
├── image_height    INT
├── status          TEXT            -- queued / rendering / qc_check / done / failed
├── output_paths    JSONB           -- 输出的图片路径列表
├── error_message   TEXT
├── created_at      TIMESTAMP
└── updated_at      TIMESTAMP
```

### 3.3 API 设计

#### 认证模块

| 方法 | 路径                        | 描述                       |
| ---- | --------------------------- | -------------------------- |
| POST | `/api/v1/auth/register`     | 邮箱注册                   |
| POST | `/api/v1/auth/login`        | 邮箱密码登录               |
| POST | `/api/v1/auth/verify-email` | 邮箱验证                   |
| POST | `/api/v1/auth/2fa/setup`    | 启用 2FA，返回 TOTP 二维码 |
| POST | `/api/v1/auth/2fa/verify`   | 验证 2FA 码                |
| GET  | `/api/v1/auth/me`           | 获取当前用户信息及用量     |

#### 角色管理

| 方法   | 路径                                   | 描述                                 |
| ------ | -------------------------------------- | ------------------------------------ |
| GET    | `/api/v1/characters`                   | 角色列表（支持分类/标签/关键词检索） |
| POST   | `/api/v1/characters`                   | 创建角色                             |
| GET    | `/api/v1/characters/:id`               | 角色详情                             |
| PUT    | `/api/v1/characters/:id`               | 更新角色文字信息                     |
| DELETE | `/api/v1/characters/:id`               | 删除角色                             |
| POST   | `/api/v1/characters/:id/variants`      | 添加角色变体                         |
| PUT    | `/api/v1/characters/:id/variants/:vid` | 更新变体                             |
| DELETE | `/api/v1/characters/:id/variants/:vid` | 删除变体                             |
| POST   | `/api/v1/characters/:id/images`        | 上传角色图片 (multipart/form-data)   |
| DELETE | `/api/v1/characters/:id/images/:imgId` | 删除角色图片                         |
| GET    | `/api/v1/character-groups`             | 角色分组列表                         |
| POST   | `/api/v1/character-groups`             | 创建角色分组                         |
| PUT    | `/api/v1/character-groups/:id/members` | 管理分组成员                         |

#### 项目/创作管线

| 方法   | 路径                                       | 描述                     |
| ------ | ------------------------------------------ | ------------------------ |
| GET    | `/api/v1/projects`                         | 项目列表                 |
| POST   | `/api/v1/projects`                         | 创建项目                 |
| GET    | `/api/v1/projects/:id`                     | 项目详情                 |
| PUT    | `/api/v1/projects/:id`                     | 更新项目                 |
| DELETE | `/api/v1/projects/:id`                     | 删除项目                 |
| POST   | `/api/v1/projects/:id/generate-script`     | 触发 Agent 1-2：生成剧本 |
| PUT    | `/api/v1/projects/:id/scripts/:sid`        | 用户修改剧本             |
| POST   | `/api/v1/projects/:id/generate-storyboard` | 触发 Agent 3：生成分镜   |
| PUT    | `/api/v1/projects/:id/storyboards/:sbid`   | 用户修改分镜             |
| GET    | `/api/v1/projects/:id/previz`              | 获取聊天式预览数据       |
| POST   | `/api/v1/projects/:id/render`              | 触发 Agent 4-8：渲染请求 |
| GET    | `/api/v1/projects/:id/render-status`       | 渲染状态查询             |
| GET    | `/api/v1/projects/:id/outputs`             | 获取最终产出图片         |

#### 系统配置

| 方法 | 路径             | 描述                   |
| ---- | ---------------- | ---------------------- |
| GET  | `/api/v1/models` | 获取可用模型列表       |
| GET  | `/api/v1/usage`  | 获取当前用户的用量统计 |

#### WebSocket / SSE

| 路径                             | 描述                                         |
| -------------------------------- | -------------------------------------------- |
| `WS /api/v1/projects/:id/stream` | 管线进度实时推送（剧本生成进度、渲染进度等） |

### 3.4 Agent 编排管线

```
用户输入 (synopsis + characters)
         │
  ┌──────▼──────┐
  │ Character   │  Agent 1：检索/实例化角色，拼接完整 Prompt
  │ Agent       │
  └──────┬──────┘
         │ Character_Profile.json
  ┌──────▼──────┐
  │ Scriptwriter│  Agent 2：生成剧本 (standalone 或 serialized)
  │ Agent       │
  └──────┬──────┘
         │ Script_Outline.json
  ┌──────▼──────┐
  │ Storyboard  │  Agent 3：生成机器可读的分镜 JSON
  │ Agent       │
  └──────┬──────┘
         │ Storyboard_Render_Payload.json
   ══════╧══════════════════════════════
         │ ← HITL Checkpoint (聊天式预览)
         │   用户在此修改台词/镜头/节奏
   ══════╤══════════════════════════════
         │ (用户确认后继续)
  ┌──────▼──────┐
  │ Asset       │  Agent 4：分发参考图、Seed、风格指令
  │ Manager     │
  └──────┬──────┘
         │
  ┌──────▼──────┐
  │ Generation  │  Agent 5：调用图像生成模型
  │ Engine      │
  └──────┬──────┘
         │ 生成的图片
  ┌──────▼──────┐
  │ QC Reviewer │  Agent 6：多模态审校，不合格则回退重绘
  │ Agent       │
  └──────┬──────┘
         │ 审校通过的图片
  ┌──────▼──────┐
  │ Typesetter  │  Agent 7：嵌字/净版路由
  │ Agent       │
  └──────┬──────┘
         │
  ┌──────▼──────┐
  │ Compositor  │  Agent 8：按模板排版合成最终图片
  │ Agent       │
  └──────┴──────┘
         │
    最终四格漫画输出
```

每个 Agent 实现统一接口：

```go
type Agent interface {
    Name() string
    Execute(ctx context.Context, input *PipelinePayload) (*PipelinePayload, error)
}
```

Orchestrator 通过 TaskQueue 调度 Agent 执行顺序，支持：
- 顺序执行（阶段间）
- 并行执行（同阶段内，如 4 格图并行渲染）
- 错误回退（QC 失败后重新触发 Generation Engine）

---

## 四、前端详细设计

### 4.1 目录结构

```
frontend/
├── public/
│   └── index.html
├── src/
│   ├── main.tsx                    # 入口
│   ├── App.tsx                     # 路由配置
│   ├── api/                        # API 调用封装
│   │   ├── client.ts               # Axios 实例 + 拦截器
│   │   ├── auth.ts
│   │   ├── characters.ts
│   │   └── projects.ts
│   ├── store/                      # 状态管理 (Zustand)
│   │   ├── authStore.ts
│   │   ├── characterStore.ts
│   │   └── projectStore.ts
│   ├── hooks/                      # 自定义 Hooks
│   │   ├── useWebSocket.ts
│   │   └── useAuth.ts
│   ├── components/                 # 通用组件
│   │   ├── Layout/
│   │   │   ├── Sidebar.tsx
│   │   │   ├── Header.tsx
│   │   │   └── MainLayout.tsx
│   │   ├── Auth/
│   │   │   ├── LoginForm.tsx
│   │   │   ├── RegisterForm.tsx
│   │   │   └── TwoFactorSetup.tsx
│   │   ├── Character/
│   │   │   ├── CharacterCard.tsx
│   │   │   ├── CharacterForm.tsx
│   │   │   ├── CharacterGallery.tsx
│   │   │   ├── VariantEditor.tsx
│   │   │   └── ImageUploader.tsx
│   │   ├── Project/
│   │   │   ├── ProjectCard.tsx
│   │   │   ├── ProjectWizard.tsx   # 创作向导
│   │   │   └── ScriptEditor.tsx
│   │   ├── Previz/
│   │   │   ├── ChatBubble.tsx      # 聊天气泡
│   │   │   ├── ChatTimeline.tsx    # 聊天式剧情预览
│   │   │   └── PrevizToolbar.tsx   # 修改工具栏
│   │   ├── Render/
│   │   │   ├── LayoutSelector.tsx  # 2x2 / 1x4 布局选择
│   │   │   ├── SizeConfigurator.tsx # 图片尺寸设置
│   │   │   ├── RenderProgress.tsx  # 渲染进度
│   │   │   └── OutputViewer.tsx    # 最终成品查看
│   │   └── Common/
│   │       ├── SearchBar.tsx
│   │       ├── TagFilter.tsx
│   │       └── UsageIndicator.tsx  # 用量指示器
│   ├── pages/
│   │   ├── LoginPage.tsx
│   │   ├── DashboardPage.tsx       # 仪表板/项目总览
│   │   ├── CharacterLibraryPage.tsx
│   │   ├── CharacterDetailPage.tsx
│   │   ├── ProjectCreatePage.tsx
│   │   ├── PrevizPage.tsx          # 聊天式预览页
│   │   ├── RenderPage.tsx          # 渲染设置与结果页
│   │   └── SettingsPage.tsx        # 用户设置/模型选择/用量
│   ├── styles/
│   │   ├── globals.css
│   │   ├── variables.css
│   │   └── components/
│   ├── types/                      # TypeScript 类型定义
│   │   ├── character.ts
│   │   ├── project.ts
│   │   └── api.ts
│   └── utils/
│       ├── format.ts
│       └── validation.ts
├── package.json
├── tsconfig.json
├── vite.config.ts
└── Dockerfile
```

### 4.2 页面流程

```
登录/注册 → 仪表板（项目列表）
                │
     ┌──────────┼──────────┐
     ▼          ▼          ▼
  角色资产库  新建项目    设置
     │          │
     │    ┌─────▼─────┐
     │    │ 选择模式   │  (独立四格 / 连载四格)
     │    │ 输入灵感   │
     │    │ 选择角色   │
     │    └─────┬─────┘
     │          │
     │    ┌─────▼─────┐
     │    │ 聊天式预览 │  ← 核心交互页面
     │    │ (Previz)   │  角色头像 + 台词气泡
     │    │ 修改/确认  │
     │    └─────┬─────┘
     │          │ 确认分镜
     │    ┌─────▼─────┐
     │    │ 渲染设置   │  选择布局/尺寸/有字无字
     │    │ 一键渲染   │
     │    └─────┬─────┘
     │          │
     │    ┌─────▼─────┐
     │    │ 成品查看   │  预览/下载最终图片
     │    │ 下载导出   │
     │    └───────────┘
     │
     ▼
  角色管理
  ├── 创建/编辑角色（文字信息）
  ├── 上传角色图片（立绘/Q版/表情）
  ├── 管理角色变体（服装/性格偏向）
  └── 角色分组与检索
```

### 4.3 聊天式预览 (Previz) 界面说明

Previz 界面是系统的**核心人机协同页面**，设计如下：

- **外观**：类似即时通讯软件（微信/LINE），左右交替的聊天气泡
- **角色头像**：使用数据库中角色 `character_images` 表里 `image_type = 'avatar'` 的图片
- **消息气泡**：
  - 每个气泡对应分镜脚本中的一句台词
  - 旁白/场景描述以系统通知样式居中展示
  - 格间分隔用视觉分割线标识（第一格 / 第二格 / 第三格 / 第四格）
- **编辑功能**：
  - 点击气泡可直接**内联编辑**台词文本
  - 侧边栏显示当前格的镜头信息（景别、视角），支持下拉修改
  - "重新生成本格"按钮 — 仅对 AI 请求重新生成单格
  - "重新生成全部"按钮 — 重新生成整个分镜
- **连载模式**：底部有"上一话 / 下一话"导航

---

## 五、核心数据契约 (JSON Schema)

### 5.1 Character_Profile.json

```json
{
  "character_id": "uuid",
  "name": "高坂穂乃果",
  "name_jp": "高坂穂乃果",
  "base": {
    "gender": "female",
    "age": "16",
    "personality": ["元気", "前向き", "天然"],
    "visual_prompt": "orange hair, side ponytail, blue eyes, medium height",
    "backstory": "音ノ木坂学院の2年生..."
  },
  "active_variant": {
    "variant_name": "冬季制服",
    "personality_mod": null,
    "visual_prompt_override": "wearing winter school uniform, gray blazer, red ribbon"
  },
  "reference_images": [
    { "type": "avatar", "path": "/data/assets/characters/honoka/avatar.png" },
    { "type": "full_body", "path": "/data/assets/characters/honoka/winter_uniform.png" }
  ]
}
```

### 5.2 Script_Outline.json

```json
{
  "project_id": "uuid",
  "mode": "standalone",
  "episode": 1,
  "total_episodes": 1,
  "title": "穂乃果的便当大作战",
  "synopsis": "穂乃果尝试给海未做便当但一切都出了问题...",
  "panels": [
    {
      "panel_number": 1,
      "structure": "起",
      "scene_description": "校园屋顶，午休时间",
      "characters_involved": ["honoka", "umi"],
      "dialogue": [
        { "character": "honoka", "line": "海未ちゃん！今日はお弁当作ったよ！" },
        { "character": "umi", "line": "穂乃果が料理...？少し不安だけど..." }
      ],
      "narration": null
    },
    { "panel_number": 2, "structure": "承", "..." : "..." },
    { "panel_number": 3, "structure": "转", "..." : "..." },
    { "panel_number": 4, "structure": "合", "..." : "..." }
  ]
}
```

### 5.3 Storyboard_Render_Payload.json

```json
{
  "storyboard_id": "uuid",
  "export_type": "native_text",
  "layout": "2x2",
  "image_width": 1024,
  "image_height": 1024,
  "panels": [
    {
      "panel_number": 1,
      "composition": {
        "shot_type": "medium_shot",
        "angle": "eye_level",
        "atmosphere": "明るい日差し、桜の花びらが舞う屋上"
      },
      "characters": [
        {
          "character_id": "honoka_uuid",
          "position": "left",
          "action": "proudly holding a wrapped bento box with both hands",
          "expression": "beaming smile, sparkling eyes"
        },
        {
          "character_id": "umi_uuid",
          "position": "right",
          "action": "leaning back slightly with arms crossed",
          "expression": "nervous smile, sweat drop"
        }
      ],
      "positive_prompt": "2girls, school rooftop, cherry blossoms, lunchtime...",
      "negative_prompt": "bad anatomy, extra limbs...",
      "dialogue_overlay": [
        {
          "order": 1,
          "character": "honoka",
          "text": "海未ちゃん！今日はお弁当作ったよ！",
          "bubble_type": "excited_burst",
          "position_hint": "top_left"
        },
        {
          "order": 2,
          "character": "umi",
          "text": "穂乃果が料理...？少し不安だけど...",
          "bubble_type": "thought_cloud",
          "position_hint": "top_right"
        }
      ],
      "sfx": [
        { "text": "ジャーン", "position_hint": "center" }
      ],
      "seed": 42
    }
  ]
}
```

---

## 六、MVP 分期开发计划

### Phase 0：项目基础搭建（预计 3 天）

- [x] 后端 Go 项目脚手架搭建（Gin + 目录结构 + 配置加载模块）
- [x] 前端 React + Vite + TypeScript 项目初始化
- [x] 定义所有抽象接口（LLM / Queue / Repository）
- [x] 实现 SQLite Repository 基础版
- [x] 实现 Channel Queue 基础版
- [x] Docker Compose 配置（前端 + 后端 + 可选 DB）
- [x] CI 基础配置

### Phase 1 (P0)：核心创作管线（预计 10 天）

#### 1.1 用户认证系统（3 天）
- [ ] 邮箱注册/登录 API
- [ ] JWT Token 生成与校验中间件
- [ ] 2FA (TOTP) 设置与验证
- [ ] 用量记录与额度限制中间件
- [ ] 前端登录/注册/2FA 页面

#### 1.2 角色管理系统（3 天）
- [ ] 角色 CRUD API + 图片上传/管理
- [ ] 角色变体管理
- [ ] 角色分组与分类
- [ ] 标签/关键词搜索
- [ ] 前端角色资产库页面（卡片列表 + 搜索/筛选 + 详情页 + 图片画廊）

#### 1.3 剧本生成管线（2 天）
- [ ] Gemini TextGenerator 适配器
- [ ] Agent 1 (Character Agent) 实现
- [ ] Agent 2 (Scriptwriter Agent) 实现 — 独立四格模式
- [ ] Agent 3 (Storyboard Agent) 实现
- [ ] 前端创作向导（选模式 + 输入灵感 + 选角色 + 触发生成）

#### 1.4 聊天式预览 (Previz) 与人工干预（2 天）
- [ ] Previz API（将 Storyboard JSON 转为聊天展示格式）
- [ ] 前端 Previz 聊天界面（角色头像 + 台词气泡 + 格间分隔）
- [ ] 台词内联编辑 + 镜头参数修改
- [ ] 单格/全部重新生成功能
- [ ] WebSocket 实时推送生成进度

#### 1.5 图像生成与基础排版（2 天 - 不含 1.6）
- [ ] Nano Banana 2 ImageGenerator 适配器
- [ ] Agent 5 (Generation Engine) 实现
- [ ] Agent 8 (Compositor) 基础实现 — 2x2 / 1x4 模板合成
- [ ] 前端渲染设置页（布局选择 + 尺寸配置）
- [ ] 前端成品查看与下载页

### Phase 2 (P1)：增强功能（预计 5 天）

- [ ] Agent 2 连载模式支持（serialized）
- [ ] Agent 6 (QC Reviewer) 实现 — 多模态审校 + 自动重绘
- [ ] Agent 7 (Typesetter) 实现 — 有字版 / 无字净版双模式
- [ ] Agent 4 (Asset Manager) 实现 — 统筹参考图、固定 Seed、风格控制
- [ ] 前端连载管理（话数导航、连载时间线）
- [ ] 输出模式切换 UI（AI 嵌字 / 无字净版）

### Phase 3 (P2)：高级功能（预计 5 天）

- [ ] 角色服装扩展槽位系统完善
- [ ] 视觉一致性管理（固定 Seed + 参考图自动分发）
- [ ] 工程文件导出（分层 PSD / 项目 JSON）
- [ ] 更多 LLM 适配器（OpenAI / Deepseek / Stable Diffusion）
- [ ] 更多队列实现（Redis Queue / RabbitMQ Queue）
- [ ] 更多数据库实现（PostgreSQL / MySQL）
- [ ] 生产环境 Docker 部署优化

---

## 七、配置文件设计

后端使用 YAML 配置文件 + 环境变量覆盖：

```yaml
# config.yaml
server:
  port: 8080
  mode: debug  # debug / release

database:
  driver: sqlite   # sqlite / postgres / mysql
  sqlite:
    path: ./data/auto-you-koma.db
  postgres:
    host: localhost
    port: 5432
    dbname: auto_you_koma
    user: postgres
    password: ""

queue:
  driver: channel  # channel / redis / rabbitmq
  redis:
    addr: localhost:6379
  rabbitmq:
    url: amqp://guest:guest@localhost:5672/

llm:
  text:
    provider: gemini  # gemini / openai / deepseek
    gemini:
      api_key: ${GEMINI_API_KEY}
      model: gemini-3-flash
    openai:
      api_key: ${OPENAI_API_KEY}
      model: gpt-4
    deepseek:
      api_key: ${DEEPSEEK_API_KEY}
      model: deepseek-chat
  image:
    provider: nanobanana  # nanobanana / stablediffusion
    nanobanana:
      api_key: ${NANOBANANA_API_KEY}
    stablediffusion:
      api_url: http://localhost:7860

auth:
  jwt_secret: ${JWT_SECRET}
  token_expiry: 24h
  email:
    smtp_host: smtp.gmail.com
    smtp_port: 587
    username: ${SMTP_USER}
    password: ${SMTP_PASS}

storage:
  base_path: ./data/assets  # 图片资产的基础存储路径

quota:
  default_daily_limit: 100   # 每日默认 AI 调用额度
```

---

## 八、Docker 部署方案

```yaml
# docker-compose.yml
version: '3.8'

services:
  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    volumes:
      - ./data:/app/data          # 持久化数据库与图片资产
      - ./config.yaml:/app/config.yaml
    env_file:
      - .env
    depends_on:
      - redis  # 可选

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    ports:
      - "3000:80"
    depends_on:
      - backend

  redis:  # 可选：仅在 queue.driver=redis 时需要
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

volumes:
  redis_data:
```

---

## 九、开发规范

### 9.1 代码规范

- **后端**：遵循 Go 官方代码规范，使用 `golangci-lint` 进行静态检查
- **前端**：ESLint + Prettier，使用 TypeScript 严格模式
- **Git**：使用 Conventional Commits 规范（`feat:` / `fix:` / `chore:` 等）
- **分支**：`main`（稳定）→ `develop`（开发）→ `feature/*`（功能分支）

### 9.2 接口联调

- 前后端分离开发
- 前端初期依赖后端提供的 **Mock JSON 数据**
- 优先跑通"聊天式预览"的交互逻辑和无字/有字版本切换逻辑
- 后端开发联调时使用 `vim` 编辑 JSON 测试数据

### 9.3 测试策略

- **后端**：Go 标准测试框架，重点测试 Service 层和 Agent 逻辑
- **前端**：Vitest + React Testing Library
- **集成测试**：Docker Compose 启动完整环境后运行 E2E 测试
