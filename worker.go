package main

// import slop!!!
import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"debug/elf"
	"log"
	"github.com/hugolgst/rich-go/client"
	"strconv"
)

// god i really just... LOVE stealing god it makes me SO HAPPY
// ... i stole this from my counterpart..
type Release struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name string `json:"name"`
		URL  string `json:"browser_download_url"`
	} `json:"assets"`
}

// gets the latest rosemary build shit and returns as software info struct
func getLatestRosemary() (fetched error, rosemaryInfo SoftwareInfo) {
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Fetching latest Rosemary release!"+reset)// mf so giddy, who do you think you are?
	const api = "https://api.github.com/repos/microlift/rosemary/releases/latest" // have you ever had a dream that you.

	var info SoftwareInfo
	req, err := http.NewRequest("GET", api, nil) // no way???
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1) // watch me make the return better, you arent expecting this!
		return fmt.Errorf("Failed to get info: %w", err), info // wow im so smart, you saw that right?
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err, info
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "unexpected status:", resp.Status)
		return err, info
	}

	var rel Release
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err, info
	}

	var url string
	var name string
	for _, a := range rel.Assets {
		if strings.HasSuffix(strings.ToLower(a.Name), ".appimage") {
			url = a.URL
			name = a.Name
		}
	}

	nameParts := strings.Split(name, "-")
	if len(nameParts) < 2 {
		os.Exit(1)
	}
	var version string = nameParts[1]
	var filename string = version+".AppImage"

	info.Download = url
	info.Version = version
	info.Name = filename
	return nil, info
}

// does the same thing as the fucnkntion above it but with roblox :>
func getLatestRoblox() (worked error, robloxInfo SoftwareInfo) {
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Fetching latest Roblox release!"+reset)
	const api = "https://clientsettingscdn.roblox.com/v2/client-version/MacPlayer"
	data, err := ReadDomain(api)
	var info SoftwareInfo
	if err != nil {return err, info}
	var fileName string = data+"-RobloxPlayer.zip"
	var zipURL string = "https://setup.rbxcdn.com/mac/" + fileName
	info.Download = zipURL
	info.Name = fileName
	info.Version = data
	return nil, info
}

// download it boiii
func downloadLatestRosemary(directory string) error {
	var worked, info = getLatestRosemary()
	if worked != nil {return worked}
	var err = DownloadSoftware(directory, info)
	if err != nil {return worked}
	return nil
}

// download another thing boiiiiiii
func downloadLatestRoblox(directory string) error {
	var worked, info = getLatestRoblox()
	if worked != nil {return worked}
	var err = DownloadSoftware(directory, info)
	if err != nil {return worked}
	return nil
}

// struct for holding the info relating the rosewater environment
type EnvironmentInfo struct {
	// environment files
	StateFile	    string
	PrefixDir	    string
	RosemaryDir	    string
	RobloxDir	    string

	// rosewater files
	RosewaterFile	    string
	RosewaterIcons	    string

	// userfiles
	UserConfigFile	    string
	UserFFlagFile	    string
}

// just updates the file containing all the env info surrounding persistent env stuff
// if you're really smart you would just link all this stuff in every function
//
// i am not a smart person
func UpdateEnvironment(envInfo EnvironmentInfo) {
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Updating environment info!"+reset)
	var stateFile string = envInfo.StateFile
	WriteLn(stateFile, 1, envInfo.StateFile)
	WriteLn(stateFile, 2, envInfo.PrefixDir)
	WriteLn(stateFile, 3, envInfo.RosemaryDir)
	WriteLn(stateFile, 4, envInfo.RobloxDir)
	WriteLn(stateFile, 5, envInfo.RosewaterFile)
	WriteLn(stateFile, 6, envInfo.RosewaterIcons)
	WriteLn(stateFile, 7, envInfo.UserConfigFile)
	WriteLn(stateFile, 8, envInfo.UserFFlagFile)
}

// heh, set up my environa!
func SetupEnvironment(downloadAfter bool) {
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Setting up Rosewater environment"+reset)
	var envInfo EnvironmentInfo
	var userData = ReturnData()

	var downloadsDir string = filepath.Join(userData.AppDir, "Rosewater", "Downloads")
	var rosemaryDir string = filepath.Join(userData.AppDir, "Rosewater", "Downloads", "Rosemary")
	var robloxDir string = filepath.Join(userData.AppDir, "Rosewater", "Downloads", "Roblox")
	var rosewaterFile string = filepath.Join(userData.AppDir, "Rosewater", "RosewaterCLI.Headless")
	var rosewaterIcons string = filepath.Join(userData.AppDir, "Rosewater", "Icons")

	var prefixDir string = filepath.Join(userData.AppDir, "Rosewater", "Prefix")
	var stateDir string = filepath.Join(userData.AppDir, "Rosewater", "State")
	var stateFile string = filepath.Join(userData.AppDir, "Rosewater", "State", "scriptState")

	var fastFlagFile string = filepath.Join(userData.AppDir, "Rosewater", "State", "fastFlags")
	var configFile string = filepath.Join(userData.AppDir, "Rosewater", "State", "userConfig")

	MakeDir(downloadsDir)
	MakeDir(prefixDir)
	MakeDir(stateDir)
	MakeDir(rosemaryDir)
	MakeDir(robloxDir)
	MakeDir(rosewaterIcons)

	MakeFile(stateFile)
	MakeFile(fastFlagFile)
	MakeFile(configFile)

	self, _ := os.Executable()
	CopyFile(self, rosewaterFile, 0755)
	envInfo.StateFile = stateFile
	envInfo.PrefixDir = prefixDir
	envInfo.RosemaryDir = rosemaryDir
	envInfo.RobloxDir = robloxDir
	envInfo.RosewaterFile = rosewaterFile
	envInfo.RosewaterIcons = rosewaterIcons
	envInfo.UserConfigFile = configFile
	envInfo.UserFFlagFile = fastFlagFile

	UpdateEnvironment(envInfo)

	if downloadAfter {InitialDownloads()}
}

type UserPreference struct {
	RPC struct {
		Delay     int
		ShowGame  bool
		Enabled   bool
	}

	RoseStrap struct {
		AllowWinTrans bool
		AllowRPCQuery bool
		XOverWayland  bool
		Enabled	      bool
	}
}

// yes yes all these user preference things are very sloppy, brain hurty so i just wanted to do some
// mindless stuff
func GetUserPreference() UserPreference {
	var success, info = GetEnvironment()
	if !success {panic(fmt.Errorf("Could not get environment for user files!"))}

	var configFile = info.UserConfigFile
	var returnPreference UserPreference

	var RPCLine00, _ = ReadLn(configFile, 2)
	var RPCLine01, _ = ReadLn(configFile, 3)
	var RPCLine02, _ = ReadLn(configFile, 4)

	var RoseStrapLine00, _ = ReadLn(configFile, 7)
	var RoseStrapLine01, _ = ReadLn(configFile, 8)
	var RoseStrapLine02, _ = ReadLn(configFile, 9)
	var RoseStrapLine03, _ = ReadLn(configFile, 10)

	returnPreference.RPC.Delay, _ = strconv.Atoi(RPCLine00)
	returnPreference.RPC.ShowGame, _ = strconv.ParseBool(RPCLine01)
	returnPreference.RPC.Enabled, _ = strconv.ParseBool(RPCLine02)

	returnPreference.RoseStrap.AllowWinTrans, _ = strconv.ParseBool(RoseStrapLine00)
	returnPreference.RoseStrap.AllowRPCQuery, _ = strconv.ParseBool(RoseStrapLine01)
	returnPreference.RoseStrap.XOverWayland, _ = strconv.ParseBool(RoseStrapLine02)
	returnPreference.RoseStrap.Enabled, _ = strconv.ParseBool(RoseStrapLine03)

	return returnPreference
}

func UpdateUserPreference(newValues UserPreference) {
	var success, info = GetEnvironment()
	if !success {panic(fmt.Errorf("Could not get environment for user files!"))}

	var configFile = info.UserConfigFile
	WriteLn(configFile, 1, "-RPC SETTINGS- Delay, ShowGame, Enabled")
	WriteLn(configFile, 2, strconv.Itoa(newValues.RPC.Delay))
	WriteLn(configFile, 3, strconv.FormatBool(newValues.RPC.ShowGame))
	WriteLn(configFile, 4, strconv.FormatBool(newValues.RPC.Enabled))
	WriteLn(configFile, 5, "")

	WriteLn(configFile, 6, "-ROSESTRAP SETTINGS- AllowWindowTransform, AllowRPCQuery, XOverWayland, Enabled")
	WriteLn(configFile, 7, strconv.FormatBool(newValues.RoseStrap.AllowWinTrans))
	WriteLn(configFile, 8, strconv.FormatBool(newValues.RoseStrap.AllowRPCQuery))
	WriteLn(configFile, 9, strconv.FormatBool(newValues.RoseStrap.XOverWayland))
	WriteLn(configFile, 10, strconv.FormatBool(newValues.RoseStrap.Enabled))
}

func SetupUserConfigFiles() {
	var success, info = GetEnvironment()
	if !success {panic(fmt.Errorf("Could not get environment for user files!"))}

	var fastFlagFile = info.UserFFlagFile

	var preference UserPreference
	// defaults
	preference.RPC.Delay = 3
	preference.RPC.ShowGame = true
	preference.RPC.Enabled = true

	preference.RoseStrap.Enabled = false
	preference.RoseStrap.AllowRPCQuery = true
	preference.RoseStrap.AllowWinTrans = true
	preference.RoseStrap.XOverWayland = true

	WriteLn(fastFlagFile, 1, "{\"FFlagExample\": \"true\"}")

	UpdateUserPreference(preference)
}

// returns info globally about the environment, hopefully nothing ever changes ever!
// im sure ill never need new info!
func GetEnvironment() (sucess bool, data EnvironmentInfo) {
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Fetching environment details!"+reset)
	var userData = ReturnData()
	var envInfo EnvironmentInfo
	var checkFile string = filepath.Join(userData.AppDir, "Rosewater", "State", "scriptState")

	var testData, valid = ReadLn(checkFile, 1)
	if !valid {return false, envInfo}
	if testData != checkFile {return false, envInfo}

	var stateFile, _ = ReadLn(checkFile, 1)
	var prefixDir, _ = ReadLn(checkFile, 2)
	var rosemaryDir, _ = ReadLn(checkFile, 3)
	var robloxDir, _ = ReadLn(checkFile, 4)
	var rosewaterFile, _ = ReadLn(checkFile, 5)
	var rosewaterIcons, _ = ReadLn(checkFile, 6)

	envInfo.StateFile = stateFile
	envInfo.PrefixDir = prefixDir
	envInfo.RosemaryDir = rosemaryDir
	envInfo.RobloxDir = robloxDir
	envInfo.RosewaterFile = rosewaterFile
	envInfo.RosewaterIcons = rosewaterIcons
	return true, envInfo
}

// i wrote this because i was calling these in like every function
func getVersions() (bloxInfo SoftwareInfo, roseInfo SoftwareInfo) {
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Retrieving software versions"+reset)
	var worked, robloxInfo = getLatestRoblox()
	if worked != nil {panic(fmt.Errorf("There was an issue recovering info from roblox's website: %w", worked))}

	var worked2, rosemaryInfo = getLatestRosemary()
	if worked2 != nil {panic(fmt.Errorf("There was an issue recovering info from github's website: %w ", worked2))}

	return robloxInfo, rosemaryInfo
}

// checks for updates and installs them, honestly this thing suffers from scope creep
// why the hell you installing shit???
// you just information to me
func CheckForUpdates() {
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Checking for updates"+reset)
	var success, envInfo = GetEnvironment()
	if !success {SetupEnvironment(true)}

	var robloxInfo, rosemaryInfo = getVersions()

	var robloxVersionExists, _ = FileExists(filepath.Join(envInfo.RobloxDir, robloxInfo.Name))
	if robloxVersionExists {
		fmt.Println(bold+green+"Roblox is fully updated to "+robloxInfo.Version+reset)} else
	{
		fmt.Println(red+underline+"OUT OF DATE!"+reset+red+" : "+reset+yellow+"Updating Roblox to "+robloxInfo.Version+reset)
		ClearDir(envInfo.RobloxDir)
		if downloadLatestRoblox(filepath.Join(envInfo.RobloxDir, robloxInfo.Name)) != nil {panic(fmt.Errorf("Could not download roblox."))}
		fmt.Println(bold+green+"Roblox has been updated to "+robloxInfo.Version+reset)
	}

	var rosemaryVersionExists, _ = FileExists(filepath.Join(envInfo.RosemaryDir, rosemaryInfo.Name))
	if rosemaryVersionExists {
		fmt.Println(bold+green+"Rosemary is fully updated to "+rosemaryInfo.Version+reset)} else
	{
		fmt.Println(red+underline+"OUT OF DATE!"+reset+red+" : "+reset+yellow+"Updating Rosemary to "+rosemaryInfo.Version+reset)
		ClearDir(envInfo.RosemaryDir)
		if downloadLatestRosemary(filepath.Join(envInfo.RosemaryDir, rosemaryInfo.Name)) != nil  {panic(fmt.Errorf("Could not download rosemary"))}
		fmt.Println(bold+green+"Rosemary has been updated to "+rosemaryInfo.Version+reset)
	}
}

// waw diagnostic tool, it's great because i'll probably never use it myself, but others will
// manage to have the worst linux setup imaginable and break everything
//
// this is just the version part of it though
func DiagnosticCheckVersions() {
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Checking for updates"+reset)
	var success, envInfo = GetEnvironment()
	if !success {DiagnosticTool(); return}

	var robloxInfo, rosemaryInfo = getVersions()

	var robloxVersionExists, _ = FileExists(filepath.Join(envInfo.RobloxDir, robloxInfo.Name))
	if robloxVersionExists {
		fmt.Println(bold+green+"Roblox is fully updated to "+robloxInfo.Version+reset)} else
		{
			fmt.Println(red+underline+"OUT OF DATE!"+reset+red+" : "+reset+yellow+"Roblox should be at version "+robloxInfo.Version+reset)
		}

		var rosemaryVersionExists, _ = FileExists(filepath.Join(envInfo.RosemaryDir, rosemaryInfo.Name))
		if rosemaryVersionExists {
			fmt.Println(bold+green+"Rosemary is fully updated to "+rosemaryInfo.Version+reset)} else
			{
				fmt.Println(red+underline+"OUT OF DATE!"+reset+red+" : "+reset+yellow+"Rosemary should be at version "+rosemaryInfo.Version+reset)
			}
}

// fetch and download, woof woof :3
func InitialDownloads() {
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Downloading initial files"+reset)
	var success, envInfo = GetEnvironment()
	if !success {SetupEnvironment(true)}

	var robloxInfo, rosemaryInfo = getVersions()

	var robloxVersionExists, _ = FileExists(filepath.Join(envInfo.RobloxDir, robloxInfo.Name))
	if !robloxVersionExists {
		if err := downloadLatestRoblox(filepath.Join(envInfo.RobloxDir, robloxInfo.Name)); err != nil  {
			panic(fmt.Errorf("Could not download roblox: %w", err))
		}
	}

	var rosemaryVersionExists, _ = FileExists(filepath.Join(envInfo.RosemaryDir, rosemaryInfo.Name))
	if !rosemaryVersionExists {
		if err := downloadLatestRosemary(filepath.Join(envInfo.RosemaryDir, rosemaryInfo.Name)); err != nil  {
			panic(fmt.Errorf("Could not download rosemary: %w", err))
		}
	}
}

// this is the REAL full-install, step aside ladies, InitialSetup() is here.
// what do you think it does
func InitialSetup() {
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Setting up folders and unzipping"+reset)
	var success, envInfo = GetEnvironment()
	if !success {SetupEnvironment(true)}

	var robloxInfo, rosemaryInfo = getVersions()

	var robloxZip = filepath.Join(envInfo.RobloxDir, robloxInfo.Name)
	var rosemaryExec = filepath.Join(envInfo.RosemaryDir, rosemaryInfo.Name)

	var robloxVersionExists, _ = FileExists(robloxZip)
	if !robloxVersionExists {
		InitialDownloads()
	}

	var rosemaryVersionExists, _ = FileExists(rosemaryExec)
	if !rosemaryVersionExists {
		InitialDownloads()
	}

	os.Chmod(rosemaryExec, 0755)
	cmd := exec.Command(rosemaryExec, "prefix", "create", envInfo.PrefixDir)
	cmd.Env = append(os.Environ(), "APPIMAGE_EXTRACT_AND_RUN=1", "ROSEMARY_PREFIX_SETUP=download") //ex
	out, err := cmd.CombinedOutput()
	if err != nil {panic(fmt.Errorf("Failed to create prefix in %s\nCommand %s\n Reason: %w\n Out: %s", envInfo.PrefixDir, cmd, err, out))}

	Unzip(robloxZip, filepath.Join(envInfo.PrefixDir, "root", "Applications"))
	var iconAsset Asset
	iconAsset.Directory = "Icons"
	iconAsset.Asset = "AppIcon.icns"
	var icon = DownloadAsset(envInfo.RosewaterIcons, iconAsset)
	var robloxIcon = filepath.Join(envInfo.PrefixDir, "root", "Applications","RobloxPlayer.app","Contents","Resources","AppIcon.icns")
	os.Remove(robloxIcon)

	CopyFile(icon, robloxIcon, 0644)

	SetupURI()
	SetupDesktop()
}

// aw shit someone installed software and then didnt know how to use it :pensive:
// alternatively my shit broke because i'm an awful programmer
func Uninstall() {
	var userData = ReturnData()
	var rosewaterDir string = filepath.Join(userData.AppDir, "Rosewater")
	os.RemoveAll(rosewaterDir)

	var desktopFile = filepath.Join(userData.AppDir, "applications","rosewateruri.desktop")
	var desktopFile2 = filepath.Join(userData.AppDir, "applications","rosewater.desktop")
	os.Remove(desktopFile)
	os.Remove(desktopFile2)

	RunShell("update-desktop-database ~/.local/share/applications") // software is my passion
}

// yeah this just sets up a .desktop file
// but like, only the URI parsing part
//
// you know in retrospect i probably just should've added a check for a uri, and run headlessly if it exists
// oh well, i dont care lol..
func SetupURI() {
	var success, envInfo = GetEnvironment()
	if !success {SetupEnvironment(true)}

	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Setting up user .desktop handling for URI parsing"+reset)
	var userData = ReturnData()
	var desktopFile string = filepath.Join(userData.AppDir, "applications", "rosewateruri.desktop")
	MakeFile(desktopFile)

	WriteLn(desktopFile, 1, "[Desktop Entry]")
	WriteLn(desktopFile, 2, "Name=RosewaterURI")
	WriteLn(desktopFile, 3, "Type=Application")
	WriteLn(desktopFile, 4, "NoDisplay=true")
	WriteLn(desktopFile, 5, "Terminal=false")
	WriteLn(desktopFile, 6, "Exec="+envInfo.RosewaterFile+" --uri %u")
	WriteLn(desktopFile, 7, "MimeType=x-scheme-handler/roblox-player;")

	RunShell("update-desktop-database ~/.local/share/applications")
	RunShell("xdg-mime default "+desktopFile+" x-scheme-handler/roblox-player")
	RunShell("xdg-mime query default x-scheme-handler/roblox-player")
}

// sets up the actual user part of the desktop because it was easier to split it into 2 functions ig
// lmao
func SetupDesktop() {
	var success, envInfo = GetEnvironment()
	if !success {SetupEnvironment(true)}

	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Setting up desktop file for User"+reset)

	var asset Asset
	asset.Directory = "Icons"
	asset.Asset = "rosewatericon_256.png"

	var downloadedAsset = DownloadAsset(envInfo.RosewaterIcons, asset)

	var userData = ReturnData()
	var desktopFile string = filepath.Join(userData.AppDir, "applications", "rosewater.desktop")
	MakeFile(desktopFile)

	WriteLn(desktopFile, 1, "[Desktop Entry]")
	WriteLn(desktopFile, 2, "Name=Rosewater")
	WriteLn(desktopFile, 3, "Comment=Roblox on Linux through Rosemary!")
	WriteLn(desktopFile, 4, "Icon="+downloadedAsset)
	WriteLn(desktopFile, 5, "Type=Application")
	WriteLn(desktopFile, 6, "Terminal=false")
	WriteLn(desktopFile, 7, "Exec="+envInfo.RosewaterFile+" --run-rosewater")

	RunShell("desktop-file-validate "+desktopFile)
	RunShell("update-desktop-database ~/.local/share/applications")
}

// so this is actually pretty funny, basically rosemary's appimage requires libpulse to run, and that would be fine...
// however, it requires libpulse as an app image and genuinely it is unable to run webkitgtk without it
// libpulse is not bundled in the appimage
//
// meaning that webkit just would not work and roblox wouldn't display lots of menus, yeah, pretty bad
// this just finds a locally installed libpulse library and then returns it so that i can pass it in a variable to
// rosemary using LD_PRELOAD; probably this function wont need to exist soon? but we'll see
func findPulseCommon() (string, bool) {
	var dirs []string

	for _, location := range []string{"/usr/lib", "/usr/lib64", "/lib64", "/usr/lib/x86_64-linux-gnu", "/lib/x86_64-linux-gnu"} { // WHERE THE FUCK IS IT???

			real, err := filepath.EvalSymlinks(filepath.Join(location, "libpulse.so.0"))
			if err != nil {
				continue
			}

			origin := filepath.Dir(real)
			if f, err := elf.Open(real); err == nil {
				for _, tag := range []elf.DynTag{elf.DT_RUNPATH, elf.DT_RPATH} {
					vals, _ := f.DynString(tag)
					for _, v := range vals {
						for _, dir := range strings.Split(v, ":") {
							dir = strings.ReplaceAll(dir, "${ORIGIN}", origin)
							dir = strings.ReplaceAll(dir, "$ORIGIN", origin)
							dirs = append(dirs, dir)
						}
					}
				}
				f.Close()
			}

			dirs = append(dirs, filepath.Join(origin, "pulseaudio"))
			break
		}

		dirs = append(dirs,"/usr/lib/pulseaudio", "/usr/lib64/pulseaudio", "/usr/lib/x86_64-linux-gnu/pulseaudio")

		for _, dir := range dirs {
			matches, _ := filepath.Glob(filepath.Join(dir, "libpulsecommon-*.so")) // yoink, you're mine now'
			if len(matches) > 0 {
				return matches[0], true
			}
		}
		return "", false // fuck, well i guess no menus for you buddy!
}

type bloxstrapQuery struct {
	Command string `json:"command"`
	Data	struct {
		X	int `json:"x"`
		Y	int `json:"y"`
		Width	int `json:"width"`
		Height	int `json:"height"`
		ScaleW	int `json:"scaleWidth"`
		ScaleH	int `json:"scaleHeight"`
	} `json:"data"`
}

// runs roblox headlessly, i.e no terminal, also closes the game entirely on exit because i thought it was annoying
// it will also do weird rpc shit with logscan, but im lazy right now
//
// replaces RunRoblox() now because it was annoying to have two functions for starting
func RunHeadless(uri string) bool {
	defer client.Logout()
	fmt.Println(bold + green + "/ᐠ˵- ⩊ -˵マ: " + reset + blue + "Starting roblox" + reset)
	var success, envInfo = GetEnvironment()
	if !success {os.Exit(1)}

	var _, rosemaryInfo = getVersions()

	var rosemaryExec = filepath.Join(envInfo.RosemaryDir, rosemaryInfo.Name)
	var robloxExe = filepath.Join(envInfo.PrefixDir, "root", "Applications", "RobloxPlayer.app", "Contents", "MacOS", "RobloxPlayer")

	var args []string
	if uri == "app" {
		args = []string{"run", robloxExe}
	} else {
		args = []string{"run", robloxExe, uri}
	}

	cmd := exec.Command(rosemaryExec, args...)

	var lib, ok = findPulseCommon()
	if !ok {lib = ""}

	var envVariables = []string{
		"ROSEMARY_PREFIX="+filepath.Join(envInfo.PrefixDir,"root"),
		"ROSEMARY_DECORATIONS=server",
		"ROSEMARY_MENU=none",
		"ROSEMARY_STEAM=0",
		"LD_PRELOAD="+lib, // DRY is for people with not enough time on their paws!!!
	}
	var xCompose = []string {
		"WAYLAND_DISPLAY=",
		"GDK_BACKEND=x11",
		"SDL_VIDEODRIVER=x11",
	}

	const xSwap bool = false
	if xSwap {
		envVariables = append(envVariables, xCompose...)
	}

	cmd.Env = append(os.Environ(), envVariables...)

	var logfileLocation = filepath.Join(envInfo.PrefixDir, "client.log") // it's my loggerona!
	logfile, _ := os.Create(logfileLocation)
	defer logfile.Close()

	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} //yay process group yaayy
	cmd.WaitDelay = 5 * time.Second

	logPipeRead, logPipeWrite := io.Pipe() // pipe so we can read it

	cmd.Stdout = logPipeWrite
	cmd.Stderr = logPipeWrite // STDs

	f, _ := os.Create(logfileLocation+"2")
	defer f.Close()
	log.SetOutput(f)

	if err := cmd.Start(); err != nil {return false}

	var closed bool = false
	stop := func() {
		if uri != "app" {
			defer client.Logout()
			if closed {
				return
			}
			closed = true

			processID := cmd.Process.Pid
			syscall.Kill(-processID, syscall.SIGTERM)

			time.AfterFunc(time.Second*5, func() {
				syscall.Kill(-processID, syscall.SIGKILL)
			})
		}
	}

	setJoining := func() {
		if closed {return}

		var joiningRPC RPCQuery
		joiningRPC.Action = "Joining a game in "
		joiningRPC.Where = "Establishing a connection"

		RPCActivity(joiningRPC)
	}
	setJoined := func() {
		if closed {return}

		var playingRPC RPCQuery
		playingRPC.Action = "Playing a game on "
		playingRPC.Where = "Running Roblox on Linux through Rosemary"

		RPCActivity(playingRPC)
	}
	setTeleporting := func() {
		if closed {return}

		var teleportingRPC RPCQuery
		teleportingRPC.Action = "Swapping places in "
		teleportingRPC.Where = "Teleporting..."

		RPCActivity(teleportingRPC)
	}
	bloxstrapSetWindow := func(line string) {
		if !xSwap {return}
		start := strings.Index(line, "{")
		if start == -1 {
			return
		}

		var bloxOut bloxstrapQuery
		if err := json.Unmarshal([]byte(line[start:]), &bloxOut); err != nil {
			log.Printf("Could not unpack json object! Reason: %v\n", err)
			return
		}

		var windowQuery WindowQueryData
		windowQuery.Data_PosX = bloxOut.Data.X
		windowQuery.Data_PosY = bloxOut.Data.Y
		windowQuery.Data_Height = bloxOut.Data.Height
		windowQuery.Data_Width = bloxOut.Data.Width

		if err := SetWindow(cmd.Process.Pid, windowQuery); err != nil {
			log.Println(err)
		}
	}
	// {"command":"SetWindow","data":{"y":339,"x":454}}
	rules := []LogRule {
		{
			Match: "leaveUGCGameInternal",
			Handler: func(line string) {stop()},
		},
		{
			Match: " ! Joining game '",
			Handler: func(line string) {setJoining()},
		},
		{
			Match: "initiateTeleport",
			Handler: func(line string) {setTeleporting()},
		},
		{
			Match: " serverId: ",
			Handler: func(line string) {setJoined()},
		},
		{
			Match: "[BloxstrapRPC]",
			Handler: func(line string) {bloxstrapSetWindow(line)},
		},
	}

	logRun := make(chan struct{})
	go func() {
		defer close(logRun)
		LogScan(logPipeRead, logfile, rules)
	}()
	exited := cmd.Wait()
	logPipeWrite.Close()
	<-logRun
	if uri != "app" {
		if closed {return true}
		return exited == nil
	}
	return true
}

// actual DiagnosticTool function, you know, just incase the environment fucking implodes again
func DiagnosticTool() {
	var success, envInfo = GetEnvironment()
	fmt.Println(bold+underline+red+"Debug Info Begin"+reset)
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Diagnostic Information:"+reset)
	var debugSuffix string = bold+yellow+"Debug: "+reset+red
	fmt.Println(bold+yellow+"Debug: "+reset+red+"Checking Info"+reset)
	if !success {fmt.Println(debugSuffix+"There is no environment for the Diagnostic Tool to examine"+reset); return}

	// check for roblox zip
	var file, err = SearchForFiletype(envInfo.RobloxDir, "zip")
	ErrHandle(err,false)
	if file != "" {fmt.Println(debugSuffix+"Roblox is installed at "+envInfo.RobloxDir+" under the name: "+file+reset)} else
	{fmt.Println(debugSuffix+"Roblox is not downloaded."+reset)}

	// check for rosemary image
	file, err = SearchForFiletype(envInfo.RosemaryDir, "AppImage")
	ErrHandle(err,false)
	if file != "" {fmt.Println(debugSuffix+"Rosemary is installed at "+envInfo.RosemaryDir+" under the name: "+file+reset)} else
	{fmt.Println(debugSuffix+"Rosemary is not downloaded."+reset)}

	var _, read = ReadLn(envInfo.StateFile, 1)
	if !read {
		fmt.Println(debugSuffix+"Statefile is not created!"+reset)
	} else {
		fmt.Println(debugSuffix+"Statefile info:"+reset)
		fmt.Println(debugSuffix+"StateFile="+envInfo.StateFile+reset)
		fmt.Println(debugSuffix+"RosewaterFile="+envInfo.RosewaterFile+reset)
		fmt.Println(debugSuffix+"PrefixDir="+envInfo.PrefixDir+reset)
		fmt.Println(debugSuffix+"RosemaryDir="+envInfo.RosemaryDir+reset)
		fmt.Println(debugSuffix+"RobloxDir="+envInfo.RobloxDir+reset)
	}
	var searchDir string = filepath.Join(envInfo.PrefixDir,"root","Applications","RobloxPlayer.app","Contents")
	file, err = SearchForFiletype(searchDir, "plist")

	if file != "" && err == nil {
		fmt.Println(debugSuffix+"Roblox is properly installed to the prefix"+reset)
	} else {
		fmt.Println(debugSuffix+"There is an issue with the roblox client installation at "+searchDir+reset)
	}

	fmt.Println(bold+underline+red+"Debug Info End"+reset)
}
