# BENZHI Evaluation Guide

## 项目说明

- 项目：`zhangkui/go-reading-tracker`
- 用途：内存型 Go 阅读追踪 Web API。
- Go 工具链：`golang:1.22`
- 前端工具链：无

## 标准构建、运行和测试命令

```bash
cd '/app' && GOTOOLCHAIN=local go build ./...
cd '/app' && GOTOOLCHAIN=local go run ./cmd
cd '/app' && GOTOOLCHAIN=local go test ./...
```

## Docker 构建和进入容器

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh go-reading-tracker-bug4-candidate linux/amd64
./build_benzhi_docker.sh go-reading-tracker-bug4-candidate linux/arm64
docker run --rm -it --platform linux/amd64 go-reading-tracker-bug4-candidate
```

## 题目验证命令

```bash
go test -buildvcs=false -count=1 -run "^TestMonthlyBookGoalDoesNotCountSameMonthFromOtherYears$" ./internal/reading
go test -buildvcs=false -count=1 ./...
go build -buildvcs=false ./...
```

## Bug 复现

参见 `BUG_REPRO.md`。
