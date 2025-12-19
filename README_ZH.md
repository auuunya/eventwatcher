### EventWatcher
[![Go Reference](https://pkg.go.dev/badge/github.com/auuunya/eventwatcher.svg)](https://pkg.go.dev/github.com/auuunya/eventwatcher) [![CI](https://github.com/auuunya/eventwatcher/actions/workflows/ci.yml/badge.svg)](https://github.com/auuunya/eventwatcher/actions/workflows/ci.yml)
#### 概述
EventWatcher 是一个开源库，设计用于实时监控 Windows 事件日志。它为系统事件、应用程序日志和其他重要事件源的跟踪和响应提供了一种健壮且高效的解决方案。这个库对于需要监控事件日志以进行调试、审计和系统管理的开发人员和系统管理员来说特别有用。

#### 使用方法
要使用 EventWatcher 库，您需要：

1. 创建一个 EventNotifier 实例。
2. 添加您感兴趣的日志的事件监视器。
3. 在 EventLogChannel 上监听事件数据。
4. 通过适当关闭 EventNotifier 确保优雅地关闭程序。

#### 安装
要安装 EventWatcher 库，请运行：

```go
go get github.com/auuunya/eventwatcher
```

#### 示例

```golang
package main

import (
	"github.com/auuunya/eventwatcher"
)

func main() {
	ctx := context.Background()
	notify := eventwatcher.NewEventNotifier(ctx)
	defer notify.Close()

	channels := []string{"Application", "System", "Microsoft-Windows-Kernel-Dump/Operational"}
	for _, channel := range channels {
		err := notify.AddWatcher(channel)
		if err != nil {
			continue
		}
	}

	go func() {
		for ch := range notify.EventLogChannel {
			fmt.Printf("event entry: %v\n", ch)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

```

#### Windows powershell add event
```Powershell
Write-EventLog -LogName "Application" -Source "TestSource" -EventID 1 -EntryType Information -Message "Application Test Info"
```
#### Windows cmd add event
```cmd
eventcreate /ID 1 /L APPLICATION /T INFORMATION  /SO MYEVENTSOURCE /D "Test Application Infomation"
```

#### 跨平台支持
- **Windows：** 使用 Windows 原生事件日志 API（保持原有行为）。与 Windows 相关的测试和实现均使用 `//go:build windows` build tag。
- **macOS / Linux：** 为类 Unix 平台提供基于 `fsnotify` 的文件监控实现。调用 `AddWatcher(path)`，其中 `path` 为文件路径，写入该文件即可触发事件。
- **说明：** 在非 Windows 平台上，Windows 特有 API 会返回未实现错误（not-implemented），建议使用 Unix watcher 进行跨平台监控。

#### 运行测试与性能分析
- 运行全部测试：`go test ./...`
- 仅运行 Unix watcher 测试（macOS/Linux）：`go test -run TestEventWatcherUnixFile -v`
- 检查启动内存：`go test -run TestMemSpike -v`（会记录 runtime.MemStats 的启动前/后数据）。

#### 贡献
欢迎贡献代码！请随时在 GitHub 仓库上提出问题或提交拉取请求。

#### 许可
此项目根据 MIT 许可发布。详情请参阅 LICENSE 文件。
