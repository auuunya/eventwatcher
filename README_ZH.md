### EventWatcher
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
    ...
)

func main() {
	ctx := context.TODO()
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
			val := eventwatcher.ParseEventLogData(ch)
			fmt.Printf("val: %v\n", val)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

```

#### Windows powershell add event
```Powershell
Write-EventLog -LogName "Application" -Source "TestSource" -EventID 1000 -EntryType Information -Message "Application Test Info"
```
#### Windows cmd add event
```cmd
eventcreate /ID 10001 /L APPLICATION /T INFORMATION  /SO MYEVENTSOURCE /D "Test Application Infomation"
```

#### 贡献
欢迎贡献代码！请随时在 GitHub 仓库上提出问题或提交拉取请求。

#### 许可
此项目根据 MIT 许可发布。详情请参阅 LICENSE 文件。