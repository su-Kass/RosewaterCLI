package main

import (
	"os"
	"os/exec"
	"syscall"
	"unsafe"
	"fmt"
	"github.com/hugolgst/rich-go/client"
)

// Code i stole, idek what this is!
func isTerminal(f *os.File) bool {
	var t syscall.Termios
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(),syscall.TCGETS, uintptr(unsafe.Pointer(&t)))
	return errno == 0
}

// Recognized terminals by the --not-headless flag
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

// Code just to launch ANY terminal if someone double clicks it (because people will)
func launchTerminal() bool {
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
		args = append(args, exe, "--not-headless")
		args = append(args, os.Args[1:]...)

		if exec.Command(path, args...).Start() == nil {
			return true
		}
	}
	return false
}

func unknownOptionMenu(command string) {
	var unknownMenuContext Context
	unknownMenuContext.Header = "Unknown command!"
	unknownMenuContext.Content = []string {
		"Unknown command: "+command,
		"Could not find a match for "+command,
		"",
		"Commands are listed under '..Options..'",
		"Example is below, respond with 'OK' to continue",
	}
	unknownMenuContext.Options = []string{
		"OK",
	}

	var choice = ContextPrompt(unknownMenuContext)
	if choice != "OK" || choice == "0" {
		fmt.Println(red+underline+"INFO"+reset+red+" : "+reset+yellow+"User error, read the interactive CLI info."+reset)
		os.Exit(1)
	}
}

func failedURIMenu() {
	var parserURIContext Context
	parserURIContext.Header = "Rosewater URI Parser"
	parserURIContext.Content = []string{
		"Your URI was invalid!",
		"URIs are time-sensitive, that could be why",
		"Would you like to try again?",
	}
	parserURIContext.Options = []string{
		"YES",
		"NO",
	}
	var choice string = ContextPrompt(parserURIContext)
	if choice == "YES" || choice == "0" {
		parseURIMenu()
	}
	advancedMenu()
}

func parseURIMenu() {
	var rpc RPCQuery
	rpc.Action = "Parsing URI "
	rpc.Where = "In RosewaterCLI - parseURIMenu()"
	RPCActivity(rpc)

	var parserURIContext Context
	parserURIContext.Header = "Rosewater URI Parser"
	parserURIContext.Content = []string{
		"This is the URI Parsing tool!",
		"URIs are time-sensitive, so make sure you paste your URI",
		"Before the URI expires and is invalid",
		"",
		"Paste your URI and then press enter",
	}
	parserURIContext.Options = []string{
		"PASTE YOUR URI BELOW!",
	}
	var uri = ContextPrompt(parserURIContext)
	var passed = RunHeadless(uri)
	if !passed {
		failedURIMenu()
	}
	advancedMenu()
}

func advancedMenu() {
	var rpc RPCQuery
	rpc.Action = "Configuring "
	rpc.Where = "In RosewaterCLI - advancedMenu()"
	RPCActivity(rpc)

	var advancedContext Context
	advancedContext.Header = "Rosewater Developer Debugging"
	advancedContext.Content = []string{
		"This is the advanced menu for RosewaterCLI, if you do not know what you are doing",
		"Messing with the functions in here could result in you needing to reinstall!",
		"If you are having an issue with Rosewater, it is suggested you write an issue in github",
		"",
		"Command info:",
		"uri-setup -- sets up desktop files for URI parsing",
		"run-uri -- opens a dialog for you to parse and run a URI from a paste",
		"uninstall -- duplicate of uninstall",
		"run-roblox -- opens the roblox client",
		"check-updates -- checks for updates",
		"do-downloads -- does downloads for Roblox and Rosemary",
		"do-setup -- sets up Rosewater",
		"do-setup-env -- sets up specifically the Rosewater environment",
		"do-setup-upd-env -- updates the Rosewater environment",
		"debug-diagnose -- Gives all information from the Rosewater environment",
		"back -- goes back to the previous menu",
	}

	advancedContext.Options = []string{
		"uri-setup",
		"run-uri",
		"uninstall",
		"run-roblox",
		"check-updates",
		"do-downloads",
		"do-setup",
		"do-setup-env",
		"debug-diagnose",
		"back",
	}

	var choice = ContextPrompt(advancedContext)
	var exit bool = false
	switch choice {
		case "uri-setup", "0":
			SetupDesktop()
		case "parse-uri", "1":
			parseURIMenu()
		case "uninstall", "2":
			Uninstall()
		case "run-roblox", "3":
			RunRoblox()
		case "check-updates", "4":
			DiagnosticCheckVersions()
		case "do-downloads", "5":
			InitialSetup()
		case "do-setup", "6":
			InitialSetup()
		case "do-setup-env", "7":
			SetupEnvironment(false)
		case "debug-diagnose", "8":
			DiagnosticTool()
		case "back", "9":
			exit = true
		default:
			unknownOptionMenu(choice)
	}
	if !exit {
		advancedMenu()
	}
}

func startedUser() {
	var rpc RPCQuery
	rpc.Action = "Configuring "
	rpc.Where = "In RosewaterCLI - startedUser()"
	RPCActivity(rpc)

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
	starterOptions := [6]string{
		"full-install",
		"update",
		"uninstall",
		"run",
		"debug",
		"exit",
	}

	starterContext.Header = starterHeader
	starterContext.Content = starterInfo[:]
	starterContext.Options = starterOptions[:]

	var given = ContextPrompt(starterContext)

	switch given {
		case "full-install", "0":
			SetupEnvironment(false)
			InitialDownloads()
			InitialSetup()
		case "update", "1":
			CheckForUpdates()
		case "uninstall", "2":
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
			}
		case "run", "3":
			RunRoblox()
		case "debug", "4":
			advancedMenu()
		case "exit", "5":
			os.Exit(1)
		default:
			unknownOptionMenu(given)
	}

	startedUser()
}

func main() {
	flag, data := ReturnArguments()
	if flag == "--uri" {
		RPCLogin()
		RunHeadless(data)
		defer client.Logout()
		os.Exit(1)
	}

	if flag == "--not-headless" {
		RPCLogin()
		startedUser()
		return
	}

	if !isTerminal(os.Stdin) && !isTerminal(os.Stdout) {
		if launchTerminal() {
			return
		}
	}
	if flag == "" {
		RPCLogin()
		startedUser()
	}
	defer client.Logout()
}
