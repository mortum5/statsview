# 🚀 Statsview

Statsview is a real-time Golang runtime stats visualization profiler. It is built top on another open-source project, [go-echarts](https://github.com/go-echarts/go-echarts), which helps statsview to show its graphs on the browser.

Fork of project [statsview](https://github.com/go-echarts/statsview)

[![Go Reference](https://pkg.go.dev/badge/github.com/mortum5/statsview.svg)](https://pkg.go.dev/github.com/mortum5/statsview)
[![Go Report Card](https://goreportcard.com/badge/github.com/mortum5/statsview)](https://goreportcard.com/report/github.com/mortum5/statsview)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)


## 🔰 Installation

```shell
$ go get -u github.com/mortum5/statsview/...
```

## 📝 Usage

Statsview is quite simple to use and all static assets have been packaged into the project which makes it possible to run offline. It's worth pointing out that statsview has integrated the standard `net/http/pprof` hence statsview will be the only profiler you need.

```golang
package main

import (
    "time"

    "github.com/mortum5/statsview"
	"github.com/mortum5/statsview/viewer"
)

func main() {
    // NewEmptyViewers() create empty viewers collection
	viewers := statsview.NewEmptyViewers()
    
    // NewDefaultViewers() create viewers collection with default viewers
	// viewers := statsview.NewDefaultViewers()

	viewers.Register(
		viewer.NewGoroutinesViewer(),
		viewer.NewHeapViewer(),
		viewer.NewStackViewer(),
	)

	mgr := statsview.New(viewers)

    // Start() runs a HTTP server at `localhost:18066` by default.
	go mgr.Start()

	// Stop() will shutdown the http server gracefully
	// mgr.Stop()

	// busy working....
	time.Sleep(time.Minute)
}

    // Visit your browser at http://localhost:18066/debug/statsview
    // Or debug as always via http://localhost:18066/debug/pprof, http://localhost:18066/debug/pprof/heap, ...
```

## ⚙️ Configuration

Statsview gets a variety of configurations for the users. Everyone could customize their favorite charts style.

```golang
// WithInterval sets the interval(in Millisecond) of collecting and pulling metrics
// default -> 2000
WithInterval(interval int)

// WithMaxPoints sets the maximum points of each chart series
// default -> 30
WithMaxPoints(n int)

// WithTemplate sets the rendered template which fetching stats from the server and
// handling the metrics data
WithTemplate(t string)

// WithAddr sets the listening address and link address
// default -> "localhost:18066"
WithAddr(addr string)

// WithLinkAddr sets the html link address
// default -> "localhost:18066"
WithLinkAddr(addr string)

// WithTimeFormat sets the time format for the line-chart Y-axis label
// default -> "15:04:05"
WithTimeFormat(s string)

// WithBrowserOpen start browser session and open url automatically
// default -> disabled
WithBrowserOpen()

// WithTheme sets the theme of the charts
// default -> Macarons
//
// Optional:
// * ThemeWesteros
// * ThemeMacarons
WithTheme(theme Theme)
```

#### Set the options

```golang
import (
    "github.com/mortum5/statsview"
    "github.com/mortum5/statsview/viewer"
)

viewers := statsview.NewDefaultViewers()

// set configurations before calling `statsview.New()` method
viewer.SetConfiguration(
    viewer.WithTheme(viewer.ThemeWesteros), 
    viewer.WithAddr("localhost:8087")
)

mgr := statsview.New(viewers)
go mgr.Start()
```

## 🗂 Viewers

Viewer is the abstraction of a Graph which in charge of collecting metrics from Runtime. Statsview provides some default viewers as below.

* `GCCPUFractionViewer`
* `GCNumViewer`
* `GCSizeViewer`
* `GoroutinesViewer`
* `HeapViewer`
* `StackViewer`

Viewer wraps a go-echarts [*charts.Line](https://github.com/go-echarts/go-echarts/blob/master/charts/line.go) instance that means all options/features on it could be used. To be honest, I think that is the most charming thing about this project.

## 🔖 Snapshot

#### ThemeMacarons(default)

![Macarons](https://user-images.githubusercontent.com/19553554/99491359-92d9f680-29a6-11eb-99c8-bc333cb90893.png)

#### ThemeWesteros

![Westeros](https://user-images.githubusercontent.com/19553554/99491179-42629900-29a6-11eb-852b-694662fcd3aa.png)

## 📄 License

MIT [©chenjiandongx](https://github.com/chenjiandongx)
MIT [©mortum5](https://github.com/mortum5)
