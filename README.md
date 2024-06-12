### EventWatcher
#### Overview
EventWatcher is an open-source library designed for monitoring Windows Event Logs in real-time. It provides a robust and efficient solution for tracking and reacting to system events, application logs, and other important event sources. This library is particularly useful for developers and system administrators who need to monitor event logs for debugging, auditing, and system management purposes.

#### Usage
To use the EventWatcher library, you need to:
1. Create an `EventNotifier` instance.
2. Add event watchers for the logs you are interested in.
3. Listen for event data on the `EventLogChannel`.
4. Ensure a graceful shutdown by properly closing the `EventNotifier`.

#### Installation
To install the EventWatcher library, run:

```golang
go get github.com/auuunya/eventwatcher
```

#### Example

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
			fmt.Printf("val: %v\n",val)
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

#### Contribution
Contributions are welcome! Feel free to open issues or submit pull requests on the GitHub repository.

#### License
This project is licensed under the MIT License. See the LICENSE file for details.