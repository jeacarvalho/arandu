//go:build ignore

// Test runner parallelo - executa todos os testes usando todos os núcleos
package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

func main() {
	fmt.Println("🧪 Arandu Test Runner (Paralelo)")
	fmt.Println("==============================")
	fmt.Printf("📊 Started at: %s\n", time.Now().Format("15:04:05"))

	// Run all test stages in parallel where possible
	var wg sync.WaitGroup

	// Stage 1: Unit tests (parallel)
	wg.Add(1)
	go func() {
		defer wg.Done()
		runUnitTests()
	}()

	// Stage 2: E2E tests (needs server, so sequential after unit starts)
	wg.Add(1)
	go func() {
		defer wg.Done()
		runE2ETests()
	}()

	wg.Wait()

	fmt.Println("\n✅ Todos os testes concluídos!")
}

func runUnitTests() {
	fmt.Println("\n🧪 Executando Unit Tests (Go)...")
	fmt.Println("==================")

	cmd := exec.Command("go", "list", "./...", "-f", "{{.ImportPath}}")
	cmd.Dir = getProjectDir()
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("❌ Erro ao listar pacotes: %v\n", err)
		return
	}

	// Filter packages
	var packages []string
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.Contains(line, "tests/e2e") && 
		   !strings.Contains(line, "/scripts/") && !strings.Contains(line, "/cmd/") {
			packages = append(packages, line)
		}
	}

	// Run tests in parallel using all cores
	args := []string{"test", "-p", fmt.Sprintf("%d", getCores())}
	args = append(args, packages...)

	cmd = exec.Command("go", args...)
	cmd.Dir = getProjectDir()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("❌ Unit tests failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Unit tests concluídos")
}

func runE2ETests() {
	fmt.Println("\n🌐 Executando E2E Tests...")

	// Deploy server first
	fmt.Println("🚀 Iniciando servidor para E2E...")
	cmd := exec.Command("./scripts/safe_deploy.sh")
	cmd.Dir = getProjectDir()
	cmd.Stdout = os.Stderr // suppress output
	err := cmd.Run()
	if err != nil {
		fmt.Printf("❌ Erro ao iniciar servidor: %v\n", err)
		return
	}
	time.Sleep(2 * time.Second)

	// Run E2E
	cmd = exec.Command("bash", "tests/run_e2e.sh")
	cmd.Dir = getProjectDir()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("❌ E2E tests failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ E2E tests concluídos")
}

func getProjectDir() string {
	// Get directory of this file
	dir, _ := os.Getwd()
	return dir
}

func getCores() int {
	// Try to detect number of CPU cores
	cmd := exec.Command("nproc")
	cmd.Dir = getProjectDir()
	output, err := cmd.Output()
	if err == nil {
		var n int
		fmt.Sscanf(string(output), "%d", &n)
		if n > 0 {
			return n
		}
	}
	// Fallback
	return 4
}