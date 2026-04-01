package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"unsafe"

	"aegis-inferno/internal/content"
	"aegis-inferno/internal/crypto"
)

//go:embed all:../../ui
var uiFS embed.FS

//go:embed all:../../content
var contentFS embed.FS

const appName = "AEGIS-Inferno"

func main() {
	// On Windows: handle firewall if needed
	if runtime.GOOS == "windows" {
		handleWindowsFirewall()
	}

	// Derive content key (obfuscated in binary)
	key := crypto.DeriveKey()

	// Content manager — decrypts and serves game data
	cm, err := content.NewManager(contentFS, key)
	if err != nil {
		log.Fatalf("Failed to init content: %v", err)
	}

	// Try to bind localhost — should work without firewall
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		// If even localhost fails, show error with firewall hint
		showError(
			"Cannot start game server.\n\n" +
				"Your firewall or antivirus may be blocking local connections.\n" +
				"Try running as Administrator, or add an exception for " + appName + ".",
		)
		os.Exit(1)
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
		log.Printf("%s running on http://%s", appName, addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Brief pause for server to start
	time.Sleep(200 * time.Millisecond)

	// Open browser
	url := fmt.Sprintf("http://%s", addr)
	openBrowser(url)

	fmt.Printf("\n  %s is running.\n", appName)
	fmt.Printf("  Open %s in your browser if it didn't open automatically.\n", url)
	fmt.Printf("  Press Ctrl+C to quit.\n\n")

	// Wait for interrupt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")
	server.Close()
}

// handleWindowsFirewall checks if we're running as admin,
// and if so, adds a firewall rule. If not, and it's the first run,
// offers to relaunch as admin.
func handleWindowsFirewall() {
	if isAdmin() {
		// We have admin rights — add firewall rule silently
		addFirewallRule()
		return
	}

	// Check if firewall rule already exists
	if firewallRuleExists() {
		return // All good
	}

	// First run without admin — show notice
	fmt.Println()
	fmt.Printf("  %s — First Launch\n", appName)
	fmt.Println("  ─────────────────────────────────────────")
	fmt.Println()
	fmt.Println("  Your firewall may block the game from running.")
	fmt.Println("  To fix this, we can add a firewall exception.")
	fmt.Println("  This requires Administrator privileges.")
	fmt.Println()
	fmt.Print("  Relaunch as Administrator? [Y/n]: ")

	var answer string
	fmt.Scanln(&answer)
	answer = strings.TrimSpace(strings.ToLower(answer))

	if answer == "" || answer == "y" || answer == "yes" || answer == "д" || answer == "да" {
		relaunchAsAdmin()
		os.Exit(0)
	}

	fmt.Println("  Continuing without admin rights. If the game doesn't work,")
	fmt.Println("  try running it as Administrator manually.")
	fmt.Println()
}

func isAdmin() bool {
	cmd := exec.Command("net", "session")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := cmd.Run()
	return err == nil
}

func firewallRuleExists() bool {
	cmd := exec.Command("netsh", "advfirewall", "firewall", "show", "rule", "name="+appName)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), appName)
}

func addFirewallRule() {
	exe, _ := os.Executable()
	cmd := exec.Command("netsh", "advfirewall", "firewall", "add", "rule",
		"name="+appName,
		"dir=in",
		"action=allow",
		"program="+exe,
		"enable=yes",
		"profile=private,public",
		"description=AEGIS:Inferno game server (localhost only)",
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Run(); err != nil {
		log.Printf("Warning: could not add firewall rule: %v", err)
	} else {
		log.Printf("Firewall rule added for %s", appName)
	}
}

func relaunchAsAdmin() {
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()

	// ShellExecute with "runas" verb = UAC prompt
	verbPtr, _ := syscall.UTF16PtrFromString("runas")
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(strings.Join(os.Args[1:], " "))

	err := shellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, syscall.SW_SHOWNORMAL)
	if err != nil {
		fmt.Printf("  Could not relaunch: %v\n", err)
	}
}

var (
	shell32          = syscall.NewLazyDLL("shell32.dll")
	procShellExecute = shell32.NewProc("ShellExecuteW")
)

func shellExecute(hwnd uintptr, verb, file, args, dir *uint16, show int32) error {
	ret, _, _ := procShellExecute.Call(
		hwnd,
		uintptr(unsafe.Pointer(verb)),
		uintptr(unsafe.Pointer(file)),
		uintptr(unsafe.Pointer(args)),
		uintptr(unsafe.Pointer(dir)),
		uintptr(show),
	)
	if ret <= 32 {
		return fmt.Errorf("ShellExecute failed with code %d", ret)
	}
	return nil
}

func showError(msg string) {
	if runtime.GOOS == "windows" {
		titlePtr, _ := syscall.UTF16PtrFromString(appName)
		msgPtr, _ := syscall.UTF16PtrFromString(msg)
		user32 := syscall.NewLazyDLL("user32.dll")
		msgBox := user32.NewProc("MessageBoxW")
		msgBox.Call(0, uintptr(unsafe.Pointer(msgPtr)), uintptr(unsafe.Pointer(titlePtr)), 0x10)
	} else {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", msg)
	}
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
