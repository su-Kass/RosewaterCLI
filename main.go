package main

import (
	"os"
	"os/exec"
	"syscall"
	"unsafe"
)

// code i stole because some people dont understand what a cli tool is and try to just run it
// if there is a better way of doing this, i would appreciate knowing about it
func isTerminal(f *os.File) bool {
	var t syscall.Termios
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(),
				       syscall.TCGETS, uintptr(unsafe.Pointer(&t)))
	return errno == 0
}

// code i stole because some people dont understand what a cli tool is and try to just run it
// if there is a better way of doing this, i would appreciate knowing about it
var terminals = []struct{ name, flag string }{
	{"x-terminal-emulator", "-e"},
	{"gnome-terminal", "--"},
	{"kgx", "--"},
	{"konsole", "-e"},
	{"xfce4-terminal", "-x"},
	{"alacritty", "-e"},
	{"kitty", ""},
	{"foot", ""},
	{"xterm", "-e"},
}

// code i stole because some people dont understand what a cli tool is and try to just run it
// if there is a better way of doing this, i would appreciate knowing about it
func ensureTerminal(ud UserData) bool {
	if isTerminal(os.Stdin) || isTerminal(os.Stdout) || os.Getenv("RELAUNCHED") != "" {
		return false
	}

	exe, err := os.Executable()
	if err != nil {
		return false
	}

	for _, t := range terminals {
		path, err := exec.LookPath(t.name)
		if err != nil {
			continue
		}
		var args []string
		if t.flag != "" {
			args = append(args, t.flag)
		}
		args = append(args, exe)
		args = append(args, os.Args[1:]...)

		cmd := exec.Command(path, args...)
		cmd.Dir = ud.AppDir
		cmd.Env = append(os.Environ(), "RELAUNCHED=1")
		if cmd.Start() == nil {
			return true
		}
	}
	return false
}


func startedUser() {
	var starterContext Context
	var starterHeader string = "Welcome to RosewaterCLI!"
	starterInfo := [7]string{
		"This is the interactive installer!",
		"You can install rosemary+roblox with \"full-install\"",
		"You can also use the index next to the option...",
		"",
		"You can also update roblox using \"update\"",
		"Furthermore, after install, you can run with \"run\"",
		"Exit with \"exit\", good luck!",
	}
	starterOptions := [5]string{
		"full-install",
		"update",
		"uninstall",
		"run",
		"exit",
	}

	starterContext.Header = starterHeader
	starterContext.Content = starterInfo[:]
	starterContext.Options = starterOptions[:]

	var given = ContextPrompt(starterContext)

	if given == "full-install" || given == "0" {
		SetupEnvironment(false)
		InitialDownloads()
		InitialSetup()
		startedUser()
	}
	if given == "update" || given == "1" {
		CheckForUpdates()
		startedUser()
	}
	if given == "uninstall" || given == "2" {
		var uninstallContext Context
		uninstallContext.Header = "Uninstall Warning"
		uninstallContext.Content = []string{
			"Are you sure you wish to uninstall Rosewater?",
			"Rosewater does not keep your data within a persistent location",
			"If you ever reinstall Rosewater, you will need to sign in again",
		}
		uninstallContext.Options = []string{
			"YES",
			"NO",
		}
		var choice = ContextPrompt(uninstallContext)
		if choice == "YES" || choice == "0" {
			Uninstall()
		} else {
			startedUser()
		}
	}
	if given == "run" || given == "3" {
		RunRoblox()
		startedUser()
	}
	if given == "exit" || given == "4" {
		os.Exit(1)
	}
	startedUser()
}

func main() {
	var flag, data = ReturnArguments()
	print(data)
	if flag == "" {
		var usrData = ReturnData()
		ensureTerminal(usrData)
		startedUser()
	}
}
