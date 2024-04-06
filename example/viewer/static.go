package viewer

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/mortum5/statsview/viewer"
)

const (
	// VStatic is the name of StaticViewer
	VStatic = "static"
)

// StaticViewer show value from 0 to 10
type StaticViewer struct {
	smgr  *viewer.StatsMgr
	graph *charts.Line
}

// NewStaticViewer returns the StaticViewer instance
// Series: Count
func NewStaticViewer() viewer.Viewer {
	graph := viewer.NewBasicView(VStatic)
	graph.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Static count"}),
		charts.WithYAxisOpts(opts.YAxis{Name: "Num"}),
	)
	graph.AddSeries("Count", []opts.LineData{})

	return &StaticViewer{graph: graph}
}

func (vs *StaticViewer) SetStatsMgr(smgr *viewer.StatsMgr) {
	vs.smgr = smgr
}

func (vs *StaticViewer) Name() string {
	return VStatic
}

func (vs *StaticViewer) View() *charts.Line {
	return vs.graph
}

var i = 0

func (vs *StaticViewer) Serve(w http.ResponseWriter, _ *http.Request) {
	vs.smgr.Tick()

	metrics := viewer.Metrics{
		Values: []float64{float64(i % 10)},
		Time:   time.Unix(vs.smgr.GetTime(), 0).Format(viewer.TimeFormat()),
	}

	i++

	bs, _ := json.Marshal(metrics)
	w.Write(bs)
}
