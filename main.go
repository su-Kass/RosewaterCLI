package main

import (
	"os"
	"os/exec"
	"syscall"
	"unsafe" // imported for ONE type in this whole file btw
	"fmt"
	"github.com/hugolgst/rich-go/client" // smh client.logout!!
	"sort"
)

// code i stole, idek what this is other than it somehow works!!!
func isTerminal(f *os.File) bool {
	var t syscall.Termios
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(),syscall.TCGETS, uintptr(unsafe.Pointer(&t)))
	return errno == 0
}

// recognized terminals by the --not-headless flag
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

// code just to launch ANY terminal if someone double clicks it (because people will)
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

// menu that is shown when an option is not recognized!!!
// all of this stuff could just be consts or files buttt....
// that's boring and i'm lazy......
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

// alt menu for uri parser
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

// silly little uri parser menu, was used for when i was testing the URI parsing of macos's client, it's basically useless now though
// lol..
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

func fastFlagRemoveMenu() {
	var rpc RPCQuery
	rpc.Action = "Configuring "
	rpc.Where = "In RosewaterCLI - fastFlagRemoveMenu()"
	RPCActivity(rpc)

	var flagMap = GetFFlagMap()
	var fastFlagRemoveContext Context
	fastFlagRemoveContext.Header = "Add a fast flag"
	fastFlagRemoveContext.Content = []string{
		"You are currently removing a FFlag",
		"Input the FFlag you want to remove",
		"Example: FFlagDebugSkyGray",
		"Do not add quotations or anything!",
		"",
		"Current FFlags: ",
	}

	keys := make([]string, 0, len(flagMap))
	for k := range flagMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, flag := range keys {
		fastFlagRemoveContext.Content = append(fastFlagRemoveContext.Content, flag+": "+flagMap[flag])
	}

	var choice string = ContextPrompt(fastFlagRemoveContext)
	RemoveFFlagFromConfig(choice)
}

func fastFlagAddMenu() {
	var rpc RPCQuery
	rpc.Action = "Configuring "
	rpc.Where = "In RosewaterCLI - fastFlagAddMenu()"
	RPCActivity(rpc)

	var flagMap = GetFFlagMap()
	var fastFlagAddContext Context
	fastFlagAddContext.Header = "Add a fast flag"
	fastFlagAddContext.Content = []string{
		"You are currently adding a FFlag",
		"Input the FFlag you want to add",
		"Example: FFlagDebugSkyGray",
		"Do not add quotations or anything!",
		"",
		"Current FFlags: ",
	}

	keys := make([]string, 0, len(flagMap))
	for k := range flagMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, flag := range keys {
		fastFlagAddContext.Content = append(fastFlagAddContext.Content, flag+": "+flagMap[flag])
	}

	var choice1 string = ContextPrompt(fastFlagAddContext)

	var fastFlagDataContext Context
	fastFlagDataContext.Header = "Add a fast flag"
	fastFlagDataContext.Content = []string{
		"You are adding data to a FFlag!",
		"Input the data for the FFlag",
		"Example: True",
		"Do not add quotations or anything!",
		"",
		"Current FFlags: ",
	}

	keys1 := make([]string, 0, len(flagMap))
	for k := range flagMap {
		keys1 = append(keys1, k)
	}
	sort.Strings(keys1)

	for _, flag := range keys1 {
		fastFlagAddContext.Content = append(fastFlagAddContext.Content, flag+": "+flagMap[flag])
	}

	var choice2 string = ContextPrompt(fastFlagDataContext)

	AddFFlagtoConfig(choice1, choice2)
}

// menu for adding/removing fast flags
func fastFlagMenu() {
	var rpc RPCQuery
	rpc.Action = "Configuring "
	rpc.Where = "In RosewaterCLI - fastFlagMenu()"
	RPCActivity(rpc)

	var flagMap = GetFFlagMap()
	var fastFlagContext Context
	fastFlagContext.Header = "Rosewater Fast Flag Configuration"
	fastFlagContext.Content = []string{
		"This is the fflag config menu for Rosewater",
		"You can add an fflag with 'add'",
		"You can remove one with 'remove'",
		"both options take you to a new menu",
		"back goes back to the previous menu",
		"Current FFlags: ",
	}

	keys := make([]string, 0, len(flagMap))
	for k := range flagMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, flag := range keys {
		fastFlagContext.Content = append(fastFlagContext.Content, flag+": "+flagMap[flag])
	}

	fastFlagContext.Options = []string{
		"add",
		"remove",
		"back",
	}
	var choice string = ContextPrompt(fastFlagContext)

	switch choice{
		case "add", "0":
			fastFlagAddMenu()
		case "remove", "1":
			fastFlagRemoveMenu()
		case "back", "2":
			return
		case "exit":
			os.Exit(1)
		default: unknownOptionMenu(choice)
	}
	fastFlagMenu()
}

// user configuration menu
func usrConfigMenu() {
	var rpc RPCQuery
	rpc.Action = "Configuring "
	rpc.Where = "In RosewaterCLI - usrConfigMenu()"
	RPCActivity(rpc)

	var configContext Context
	configContext.Header = "Rosewater User Config"
	configContext.Content = []string{
		"This is the Rosewater User Configuration Menu",
		"It contains options for you to customize your user experience",
		"",
		"Some of the options should be familiar, you should search for the terms",
		"If you do not understand them!",
	}
	configContext.Options = []string{
		"rpc-menu",
		"rosestrap-menu",
		"fflag-menu",
		"back",
	}
	var answer string = ContextPrompt(configContext)

	switch answer {
		case "rpc-menu", "0":

		case "rosestrap-menu", "1":

		case "fflag-menu", "2":
			fastFlagMenu()
		case "back", "3":
			return
		default: unknownOptionMenu(answer)
	}
	usrConfigMenu()
}

// user advanced menu, we LOVE recursion!!!!
// just gives user more advanced options
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
		"do-setup-cfg",
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
		"do-setup-desktop",
		"do-setup-cfg",
		"debug-diagnose",
		"back",
	}

	var choice = ContextPrompt(advancedContext)
	switch choice {
		case "uri-setup", "0":
			SetupURI()
		case "parse-uri", "1":
			parseURIMenu()
		case "uninstall", "2":
			Uninstall()
		case "run-roblox", "3":
			RunHeadless("app")
		case "check-updates", "4":
			DiagnosticCheckVersions()
		case "do-downloads", "5":
			InitialSetup()
		case "do-setup", "6":
			InitialSetup()
		case "do-setup-env", "7":
			SetupEnvironment(false)
		case "do-setup-desktop", "8":
			SetupDesktop()
		case "do-setup-cfg", "9":
			SetupUserConfigFiles()
		case "debug-diagnose", "10":
			DiagnosticTool()
		case "back", "11":
			return
		case "exit":
			os.Exit(1)
		default:
			unknownOptionMenu(choice)
	}
	advancedMenu()
}

// simple user menu, if you dont know how to use this then like whatever yo!!!
// we love declaring variables multiple times, honestly i should probably pass the variables made here to the next instance to stop the silly memory leak-like behaviour
// oh well... this is just a shitty installer anyways :3
func startedUser() {
	defer client.Logout()
	var rpc RPCQuery
	rpc.Action = "Configuring "
	rpc.Where = "In RosewaterCLI - startedUser()"
	RPCActivity(rpc)

	var starterContext Context
	var starterHeader string = "Welcome to RosewaterCLI!"
	starterInfo := [7]string{
		"This is the interactive installer!",
		"You can install roblox with \"install\"",
		"You can also use the index next to the option...",
		"",
		"You can configure various things with \"config\"",
		"Furthermore, after install, you can run with \"run\"",
		"Exit with \"exit\", good luck!",
	}
	starterOptions := [7]string{
		"full-install",
		"update",
		"uninstall",
		"config",
		"run",
		"debug",
		"exit",
	}

	starterContext.Header = starterHeader
	starterContext.Content = starterInfo[:]
	starterContext.Options = starterOptions[:]

	var given = ContextPrompt(starterContext)

	switch given {
		case "install", "0":
			SetupEnvironment(false)
			InitialDownloads()
			InitialSetup()
		case "uninstall", "1":
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
		case "config", "2":
			usrConfigMenu()
		case "run", "3":
			RunHeadless("app")
		case "debug", "4":
			advancedMenu()
		case "exit", "5":
			os.Exit(1)
		default:
			unknownOptionMenu(given)
	}

	startedUser()
}

// lazy ass flags implementation because i dont care :sob:
func main() {
	flag, data := ReturnArguments()

	if flag == "--run-rosewater" {
		RunHeadless("app")
		os.Exit(1)
	}

	if flag == "--uri" {
		RPCLogin()
		RunHeadless(data)
		os.Exit(1)
	}

	// code just has a sort of loop break here for the launchTerminal() function lol
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
