# Bug Reproduction

## 现象

Multi-tag search returns books that contain only one requested tag.

## 触发步骤

Create books tagged go+api, go only, and api only; search with tag=go and tag=api.

## 完整错误信息

```text
TestSearchRequiresEveryRequestedTag returns Total=3 instead of Total=1.
```

## 复现命令

```bash
go test -buildvcs=false -count=1 -run "^TestSearchRequiresEveryRequestedTag$" ./internal/reading
go test -buildvcs=false -count=1 ./...
go build -buildvcs=false ./...
```
