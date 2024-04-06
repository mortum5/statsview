package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mortum5/statsview"
	customView "github.com/mortum5/statsview/example/viewer"
	"github.com/mortum5/statsview/viewer"
)

func main() {

	viewer.SetConfiguration(
		viewer.WithAddr("localhost:18066"),
		viewer.WithBrowserOpen(),
	)

	viewers := statsview.NewEmptyViewers()
	static := customView.NewStaticViewer()

	viewers.Register(
		viewer.NewGoroutinesViewer(),
		viewer.NewHeapViewer(),
		viewer.NewStackViewer(),
		static,
	)

	mgr := statsview.New(viewers)

	go func() {
		if err := mgr.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	log.Println("Gracefull shutdown")
	mgr.Stop()

}
