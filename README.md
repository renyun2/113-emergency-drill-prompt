# 应急预案与演练记录管理系统

基于 **Go (Gin + GORM) + Vue 3 (Vite + Element Plus) + PostgreSQL 16** 的应急预案、版本审批、演练计划与记录、问题整改及合规统计演示系统，使用 **Docker Compose** 一键启动。

## 功能概览

1. **预案管理**：预案名称、类型（火灾 / 危化品泄漏 / 地震 / 停电 / 自然灾害）、适用场景；正文为 **Markdown**，按版本存储，历史版本可查阅。
2. **预案审批**：修订后提交审批（编制人 / 编制日期、审批人 / 批准日期）；**批准**后自动设为当前生效版本，其余版本取消「当前生效」标记。
3. **演练计划**：演练类型（桌面推演 / 实战演练 / 联合演练）、计划日期、地点、组织部门、参与范围、演练目标、**提前通知相关部门**（文本字段）。
4. **演练记录**：完成后登记实际参演人数、演练时长、过程描述、问题清单、评估结论（优秀 / 良好 / 一般）、演练照片路径（逗号分隔）。
5. **问题整改**：从演练问题生成整改项（责任人、措施、期限），状态 `pending` / `done`，形成闭环。
6. **合规统计**：指定年度内已完成演练次数、各预案是否满足「每年至少 1 次」、部门覆盖比例、整改完成率。

## 数据表

与 `init.sql` 一致：`emergency_plans`、`plan_versions`、`drills`、`drill_issues`、`rectifications`。

演示数据：**10 份预案**、多版本与审批状态、**2024–2025 年演练**记录及部分 2025 年计划中演练、问题与整改台账。

## 快速启动

```bash
cd 113_emergency-drill
docker compose up -d --build
```

浏览器访问：**http://localhost:8080**

> 若本机 8080 已被占用，可设置环境变量 `EMERGENCY_DRILL_WEB_PORT`，例如 `8090:80` 映射：

```bash
set EMERGENCY_DRILL_WEB_PORT=8090
docker compose up -d --build
```

## 服务说明

| 服务 | 说明 |
|------|------|
| `db` | PostgreSQL 16，首次启动执行 `init.sql` |
| `api` | Gin，监听容器内 `8080`，路由前缀 `/api` |
| `web` | Nginx 托管前端静态文件，并把 `/api` 反代到 `api` 服务 |

默认数据库：`emergency_drill`，用户 `edrill` / 密码 `edrill`（仅用于本地演示）。

## API 摘录

- `GET /api/plans` — 预案列表  
- `GET /api/plans/:id/versions` — 某预案全部版本  
- `POST /api/plans/:id/versions` — 新建草稿版本  
- `POST /api/plan-versions/:vid/submit` — 提交审批  
- `POST /api/plan-versions/:vid/approve` — 批准为当前生效（body: `approver`，可选 `approved_date`）  
- `GET /api/drills` — 查询参数：`plan_id`、`year`、`status`（`planned` / `completed`）  
- `PATCH /api/drills/:id` — 更新计划或录入完成记录  
- `POST /api/drills/:id/issues` — 演练问题  
- `POST /api/issues/:issueId/rectifications` — 从问题建整改项  
- `PATCH /api/rectifications/:id` — 更新整改（`status` 为 `done` 时可带 `completed_at`）  
- `GET /api/stats/compliance?year=2026` — 合规统计  

健康检查：`GET /healthz`

## 本地开发（可选）

**后端**（需本机 PostgreSQL 或 port-forward 到容器库）：

```bash
cd backend
set DATABASE_DSN=host=localhost user=edrill password=edrill dbname=emergency_drill port=5432 sslmode=disable TimeZone=Asia/Shanghai
go run ./cmd/server
```

**前端**（`vite.config.js` 已将 `/api` 代理到 `127.0.0.1:8080`）：

```bash
cd frontend
npm install
npm run dev
```

## 目录结构

```
113_emergency-drill/
├── docker-compose.yml
├── init.sql
├── README.md
├── backend/          # Go API
└── frontend/         # Vue3 + Vite
```

## 说明

- 演练「照片路径」为演示用字符串，未实现真实文件上传；生产环境可对接对象存储并保存 URL。  
- Markdown 预览使用 `marked`，仅用于内部管理界面，勿直接渲染不可信外部内容。  
- 合规统计中「部门覆盖」以已完成演练的 `org_dept` 去重数为分子，分母为历史去重部门数与常数 8 的较大值，便于演示比例计算。
