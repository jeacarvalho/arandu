package main

import (
	"arandu/internal/platform/version"
	"fmt"
)

func main() {
	fmt.Printf("Version: %s\n", version.Version)
	fmt.Printf("Commit: %s\n", version.Commit)
	fmt.Printf("BuildTime: %s\n", version.BuildTime)
}
