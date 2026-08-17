# Bug Reproduction

## 现象

Equivalent ISBN forms are accepted as separate books.

## 触发步骤

Add 978-1-4919-5059-2, then add 9781491950592 or a whitespace-padded equivalent.

## 完整错误信息

```text
TestNormalizedISBNIsUnique reports expected duplicate ISBN error, got nil.
```

## 复现命令

```bash
go test -buildvcs=false -count=1 -run "^TestNormalizedISBNIsUnique$" ./internal/reading
go test -buildvcs=false -count=1 ./...
go build -buildvcs=false ./...
```
