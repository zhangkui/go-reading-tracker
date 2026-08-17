# Bug Reproduction

## 现象

Updating progress to the final page leaves the book in reading state.

## 触发步骤

Create a 240-page book and update current page from 239 to 240.

## 完整错误信息

```text
TestProgressAtLastPageCompletesBook reports status reading and missing completion time.
```

## 复现命令

```bash
go test -buildvcs=false -count=1 -run "^TestProgressAtLastPageCompletesBook$" ./internal/reading
go test -buildvcs=false -count=1 ./...
go build -buildvcs=false ./...
```
