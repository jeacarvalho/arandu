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
	info, err := os.Stat("web/static/css/tailwind.css")
	if err != nil {
		return fmt.Sprintf("dev_%d", time.Now().Unix())
	}
	return fmt.Sprintf("%d", info.ModTime().Unix())
}
