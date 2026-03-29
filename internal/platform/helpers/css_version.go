package helpers

import (
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	cssVersion string
	cssOnce    sync.Once
)

func GetCSSVersion() string {
	cssOnce.Do(func() {
		info, err := os.Stat("web/static/css/style.css")
		if err != nil {
			cssVersion = fmt.Sprintf("dev_%d", time.Now().Unix())
			return
		}
		cssVersion = fmt.Sprintf("%d", info.ModTime().Unix())
	})
	return cssVersion
}
