# go-reading-tracker

## 项目说明

go-reading-tracker 是一个使用 Go 标准库实现的个人阅读记录 Web API。服务使用内存存储，支持书架管理、阅读进度、标签搜索、月度阅读目标、批量导入和阅读统计，不依赖数据库、Redis 或消息队列。

## 标准命令

```bash
go build ./...
go test ./...
go vet ./...
go run ./cmd
```

服务默认监听 `:8080`，可通过环境变量 `ADDR` 修改。

## 使用方式

- `POST /books`：添加书籍
- `GET /books`：筛选并分页查询
- `GET|PUT|DELETE /books/{id}`：管理单本书籍
- `PUT /books/{id}/progress`：更新阅读进度
- `POST /goals`：设置月度阅读目标
- `GET /goals/progress`：查询目标进度
- `POST /imports`：批量导入阅读记录
- `GET /stats`：查看阅读统计

请求和响应均使用 JSON。
