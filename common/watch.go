package common

import (
	"os"
	"regexp"
	"time"

	"github.com/radovskyb/watcher"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Watch() {
	w := watcher.New()
	w.SetMaxEvents(1)

	r := regexp.MustCompile("(?i)(json$)|(toml$)")
	w.AddFilterHook(watcher.RegexFilterHook(r, false))

	go func() {
		for {
			select {
			case event := <-w.Event:
				log.Warning("Recv file changed event: ", event)
				os.Exit(0)
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch test_folder recursively for changes.
	if err := w.AddRecursive(viper.GetString("work")); err != nil {
		log.Fatalln(err)
	}

	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}
