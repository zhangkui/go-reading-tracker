# Bug Reproduction

## 现象

A pre-canceled import request still reads the request body.

## 触发步骤

Call Import with an already canceled context and an observable reader.

## 完整错误信息

```text
TestImportChecksCanceledContextBeforeReading reports that input was read after cancellation.
```

## 复现命令

```bash
go test -buildvcs=false -count=1 -run "^TestImportChecksCanceledContextBeforeReading$" ./internal/reading
go test -buildvcs=false -count=1 ./...
go build -buildvcs=false ./...
```
