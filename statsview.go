/*
Package statsview provide a real-time Golang runtime stats
visualization profiler. It is built top on another open-source project,
[go-echarts](https://github.com/go-echarts/go-echarts), which helps
statsview to show its graphs on the browser.
*/
package statsview

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/templates"
	"github.com/mortum5/statsview/statics"
	"github.com/mortum5/statsview/viewer"
	"github.com/pkg/browser"
	"github.com/rs/cors"
)

func init() {
	templates.PageTpl = `
		{{- define "page" }}
		<!DOCTYPE html>
		<html>
			{{- template "header" . }}
		<body>
		<div class="box"> {{- range .Charts }} {{ template "base" . }} {{- end }} </div>
		</body>
		</html>
		{{ end }}
		`
}

// Viewers represent collection of Viewer
type Viewers []viewer.Viewer

// NewDefaultViewers generate default collection that includes
// - GoroutinesViewer
// - HeapViewer
// - StackViewer
// - GCNumViewer
// - GCSizeViewer
// - GCCPUFractionViewer
func NewDefaultViewers() Viewers {
	return Viewers{
		viewer.NewGoroutinesViewer(),
		viewer.NewHeapViewer(),
		viewer.NewStackViewer(),
		viewer.NewGCNumViewer(),
		viewer.NewGCSizeViewer(),
		viewer.NewGCCPUFractionViewer(),
	}
}

// NewEmptyViewers returns empty collection without any Viewer
func NewEmptyViewers() Viewers {
	return Viewers{}
}

// Register adds Viewer to collection
func (v *Viewers) Register(views ...viewer.Viewer) {
	*v = append(*v, views...)
}

// ViewManager
type ViewManager struct {
	srv *http.Server

	Smgr   *viewer.StatsMgr
	Views  []viewer.Viewer
	Ctx    context.Context
	Cancel context.CancelFunc
}

// Start runs a http server and begin to collect metrics
func (vm *ViewManager) Start() error {
	if viewer.BrowserOpen() {
		browser.OpenURL(fmt.Sprintf("http://%s/debug/statsview", viewer.Addr()))
	}
	return vm.srv.ListenAndServe()
}

// Stop shutdown the http server gracefully
func (vm *ViewManager) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	vm.srv.Shutdown(ctx)
	vm.Cancel()
}

// New creates a new ViewManager instance
func New(viewers Viewers) *ViewManager {
	page := components.NewPage()
	page.PageTitle = "Statsview"
	page.AssetsHost = fmt.Sprintf("http://%s/debug/statsview/statics/", viewer.LinkAddr())
	page.Assets.JSAssets.Add("jquery.min.js")

	mgr := &ViewManager{
		srv: &http.Server{
			Addr:           viewer.Addr(),
			ReadTimeout:    time.Minute,
			WriteTimeout:   time.Minute,
			MaxHeaderBytes: 1 << 20,
		},
	}
	mgr.Ctx, mgr.Cancel = context.WithCancel(context.Background())
	mgr.Views = viewers

	smgr := viewer.NewStatsMgr(mgr.Ctx)
	for _, v := range mgr.Views {
		v.SetStatsMgr(smgr)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	for _, v := range mgr.Views {
		page.AddCharts(v.View())
		mux.HandleFunc("/debug/statsview/view/"+v.Name(), v.Serve)
	}

	mux.HandleFunc("/debug/statsview", func(w http.ResponseWriter, _ *http.Request) {
		page.Render(w)
	})

	staticsPrev := "/debug/statsview/statics/"
	mux.HandleFunc(staticsPrev+"echarts.min.js", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(statics.EchartJS))
	})

	mux.HandleFunc(staticsPrev+"jquery.min.js", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(statics.JqueryJS))
	})

	mux.HandleFunc(staticsPrev+"themes/westeros.js", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(statics.WesterosJS))
	})

	mux.HandleFunc(staticsPrev+"themes/macarons.js", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(statics.MacaronsJS))
	})

	mgr.srv.Handler = cors.AllowAll().Handler(mux)
	return mgr
}
