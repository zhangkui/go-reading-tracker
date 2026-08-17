# Bug Reproduction

## 现象

Monthly completed-book goals count the same month from other years.

## 触发步骤

Complete one book in January 2025 and one in January 2026, then query January 2026.

## 完整错误信息

```text
TestMonthlyBookGoalDoesNotCountSameMonthFromOtherYears counts 2 instead of 1.
```

## 复现命令

```bash
go test -buildvcs=false -count=1 -run "^TestMonthlyBookGoalDoesNotCountSameMonthFromOtherYears$" ./internal/reading
go test -buildvcs=false -count=1 ./...
go build -buildvcs=false ./...
```
