package winterminal

import (
	"claude-squad-windows/cmd"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// WindowsTerminalSession represents a managed Windows Terminal session
type WindowsTerminalSession struct {
	// The name of the Windows Terminal session
	sessionName string
	// The sanitized name used for Windows Terminal commands
	sanitizedName string
	// The program to run in the terminal
	program string
	
	// State management
	isRunning bool
	process   *os.Process
	mutex     sync.RWMutex
	
	// Content monitoring
	lastContent string
	contentMutex sync.RWMutex
	
	// Command executor interface
	cmdExec cmd.Executor
}

// NewWindowsTerminalSession creates a new Windows Terminal session
func NewWindowsTerminalSession(sessionName, program string) *WindowsTerminalSession {
	sanitizedName := sanitizeName(sessionName)
	return &WindowsTerminalSession{
		sessionName:   sessionName,
		sanitizedName: sanitizedName,
		program:       program,
		cmdExec:       cmd.MakeExecutor(),
	}
}

// NewWindowsTerminalSessionWithDeps creates a new session with custom executor
func NewWindowsTerminalSessionWithDeps(sessionName, program string, cmdExec cmd.Executor) *WindowsTerminalSession {
	sanitizedName := sanitizeName(sessionName)
	return &WindowsTerminalSession{
		sessionName:   sessionName,
		sanitizedName: sanitizedName,
		program:       program,
		cmdExec:       cmdExec,
	}
}

// sanitizeName cleans up the session name for Windows Terminal use
func sanitizeName(name string) string {
	// Replace spaces and special characters with underscores
	str := strings.ReplaceAll(name, " ", "_")
	str = strings.ReplaceAll(str, "-", "_")
	str = strings.ReplaceAll(str, ".", "_")
	return "claudesquad_" + str
}

// Start creates and starts a new Windows Terminal session
func (w *WindowsTerminalSession) Start(workDir string) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.isRunning {
		return fmt.Errorf("Windows Terminal session already running: %s", w.sanitizedName)
	}

	// Create new Windows Terminal tab with the specified program
	// Using wt.exe (Windows Terminal command line)
	var cmd *exec.Cmd
	if workDir != "" {
		// Start Windows Terminal in specific directory with the program
		cmd = exec.Command("wt", "-w", "0", "new-tab", "--title", w.sessionName, "-d", workDir, w.program)
	} else {
		// Start Windows Terminal with the program in current directory
		cmd = exec.Command("wt", "-w", "0", "new-tab", "--title", w.sessionName, w.program)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting Windows Terminal session: %w", err)
	}

	w.process = cmd.Process
	w.isRunning = true

	// Give it a moment to start
	time.Sleep(500 * time.Millisecond)

	return nil
}

// Restore attempts to reattach to an existing session
func (w *WindowsTerminalSession) Restore() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// For Windows Terminal, we'll attempt to focus the existing tab
	// This is a simplified implementation
	w.isRunning = w.DoesSessionExist()
	return nil
}

// DoesSessionExist checks if the Windows Terminal session exists
func (w *WindowsTerminalSession) DoesSessionExist() bool {
	// Check if there's a Windows Terminal process with our title
	// This is a simplified check - in a full implementation, you'd query
	// the actual Windows Terminal process for tabs with the specific title
	
	// For now, we'll check if our stored process is still running
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	
	if w.process == nil {
		return false
	}
	
	// Check if process is still running
	err := w.process.Signal(os.Signal(nil))
	return err == nil
}

// SendKeys sends keystrokes to the Windows Terminal session
func (w *WindowsTerminalSession) SendKeys(keys string) error {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	if !w.isRunning {
		return fmt.Errorf("Windows Terminal session is not running")
	}

	// In a full implementation, you would use Windows APIs to send keystrokes
	// For now, this is a placeholder that would need platform-specific code
	fmt.Printf("[DEBUG] Would send keys to Windows Terminal: %s\n", keys)
	return nil
}

// TapEnter sends an enter keystroke to the session
func (w *WindowsTerminalSession) TapEnter() error {
	return w.SendKeys("\r")
}

// HasUpdated checks if the terminal content has changed
func (w *WindowsTerminalSession) HasUpdated() (updated bool, hasPrompt bool) {
	// This would require Windows APIs to capture terminal content
	// For now, return placeholder values
	return false, false
}

// CapturePaneContent captures the current content of the terminal
func (w *WindowsTerminalSession) CapturePaneContent() (string, error) {
	// This would require Windows APIs to capture terminal content
	// For now, return placeholder content
	return "[Windows Terminal Content Placeholder]", nil
}

// CapturePaneContentWithOptions captures content with specific options
func (w *WindowsTerminalSession) CapturePaneContentWithOptions(start, end string) (string, error) {
	// This would require Windows APIs to capture terminal content with options
	// For now, return placeholder content
	return "[Windows Terminal Full History Placeholder]", nil
}

// SetDetachedSize sets the terminal size when detached
func (w *WindowsTerminalSession) SetDetachedSize(width, height int) error {
	// Windows Terminal size management would go here
	fmt.Printf("[DEBUG] Would set Windows Terminal size to %dx%d\n", width, height)
	return nil
}

// Attach creates an interactive connection to the session
func (w *WindowsTerminalSession) Attach() (chan struct{}, error) {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	if !w.isRunning {
		return nil, fmt.Errorf("Windows Terminal session is not running")
	}

	// Create a channel to signal detachment
	detachChan := make(chan struct{})

	// In a full implementation, this would create an interactive PTY connection
	// For now, we'll simulate with a goroutine
	go func() {
		// Simulate attachment - in real implementation, this would handle I/O
		fmt.Printf("[DEBUG] Attached to Windows Terminal session: %s\n", w.sessionName)
		// Wait for detachment signal
		<-detachChan
		fmt.Printf("[DEBUG] Detached from Windows Terminal session: %s\n", w.sessionName)
	}()

	return detachChan, nil
}

// DetachSafely safely detaches from the session without terminating it
func (w *WindowsTerminalSession) DetachSafely() error {
	// In Windows Terminal, we don't need to explicitly detach
	// The tab continues running independently
	fmt.Printf("[DEBUG] Safely detached from Windows Terminal session: %s\n", w.sessionName)
	return nil
}

// Close terminates the Windows Terminal session
func (w *WindowsTerminalSession) Close() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if !w.isRunning {
		return nil
	}

	var errs []error

	// Terminate the process if it exists
	if w.process != nil {
		if err := w.process.Kill(); err != nil {
			errs = append(errs, fmt.Errorf("error killing Windows Terminal process: %w", err))
		}
		w.process = nil
	}

	w.isRunning = false

	if len(errs) > 0 {
		return errs[0] // Return first error
	}
	return nil
}

// CleanupSessions cleans up all Windows Terminal sessions created by claude-squad
func CleanupSessions(cmdExec cmd.Executor) error {
	// In a full implementation, this would enumerate and close Windows Terminal tabs
	// that match our session naming pattern (claudesquad_*)
	fmt.Println("[DEBUG] Would cleanup Windows Terminal sessions matching 'claudesquad_*'")
	return nil
}