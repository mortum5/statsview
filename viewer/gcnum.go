package viewer

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

const (
	// VGCNum is the name of GCNumViewer
	VGCNum = "gcnum"
)

// GCNumViewer collects the GC number metric via `runtime.ReadMemStats()`
type GCNumViewer struct {
	smgr  *StatsMgr
	graph *charts.Line
}

// NewGCNumViewer returns the GCNumViewer instance
// Series: GcNum
func NewGCNumViewer() Viewer {
	graph := NewBasicView(VGCNum)
	graph.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "GC Number"}),
		charts.WithYAxisOpts(opts.YAxis{Name: "Num"}),
	)
	graph.AddSeries("GcNum", []opts.LineData{})

	return &GCNumViewer{graph: graph}
}

func (vr *GCNumViewer) SetStatsMgr(smgr *StatsMgr) {
	vr.smgr = smgr
}

func (vr *GCNumViewer) Name() string {
	return VGCNum
}

func (vr *GCNumViewer) View() *charts.Line {
	return vr.graph
}

func (vr *GCNumViewer) Serve(w http.ResponseWriter, _ *http.Request) {
	vr.smgr.Tick()

	memstats.mu.RLock()
	metrics := Metrics{
		Values: []float64{float64(memstats.Stats.NumGC)},
		Time:   time.Unix(vr.smgr.GetTime(), 0).Format(TimeFormat()),
	}
	memstats.mu.RUnlock()

	bs, _ := json.Marshal(metrics)
	w.Write(bs)
}
