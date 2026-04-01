package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"aegis-inferno/internal/content"
	"aegis-inferno/internal/crypto"
)

//go:embed all:../../ui
var uiFS embed.FS

//go:embed all:../../content
var contentFS embed.FS

func main() {
	// Derive content key (obfuscated in binary)
	key := crypto.DeriveKey()

	// Content manager — decrypts and serves game data
	cm, err := content.NewManager(contentFS, key)
	if err != nil {
		log.Fatalf("Failed to init content: %v", err)
	}

	// Find free port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatalf("Failed to find free port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	// HTTP mux
	mux := http.NewServeMux()

	// Serve UI (static files)
	uiSub, err := fs.Sub(uiFS, "ui")
	if err != nil {
		log.Fatalf("Failed to get ui subfs: %v", err)
	}
	mux.Handle("/", http.FileServer(http.FS(uiSub)))

	// API: get content (decrypted)
	mux.HandleFunc("/api/content/", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Path[len("/api/content/"):]
		data, err := cm.Get(name)
		if err != nil {
			http.Error(w, "not found", 404)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		w.Write(data)
	})

	// API: save/load game state
	mux.HandleFunc("/api/save", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "POST only", 405)
			return
		}
		// Read body, encrypt, save to user directory
		// TODO: implement save system
		w.Write([]byte(`{"ok":true}`))
	})

	mux.HandleFunc("/api/load", func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement load system
		http.Error(w, "no save found", 404)
	})

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	server := &http.Server{Addr: addr, Handler: mux}

	// Start server
	go func() {
		log.Printf("AEGIS:Inferno running on http://%s", addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Open browser
	url := fmt.Sprintf("http://%s", addr)
	openBrowser(url)

	// Wait for interrupt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")
	server.Close()
}

func openBrowser(url string) {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default:
		cmd = "xdg-open"
		args = []string{url}
	}
	p, _ := os.StartProcess(cmd, append([]string{cmd}, args...), &os.ProcAttr{
		Dir: ".",
		Env: os.Environ(),
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	})
	if p != nil {
		p.Release()
	}
}
