package main

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
)

type Release struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name string `json:"name"`
		URL  string `json:"browser_download_url"`
	} `json:"assets"`
}

func getLatestRosemary() (fetched error, rosemaryInfo SoftwareInfo) {
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Fetching latest Rosemary release!"+reset)
	const api = "https://api.github.com/repos/microlift/rosemary/releases/latest"

	var info SoftwareInfo
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
		return err, info
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

func downloadLatestRosemary(directory string) error {
	var worked, info = getLatestRosemary()
	if worked != nil {return worked}
	if !DownloadSoftware(directory, info) {return worked}
	return nil
}

func downloadLatestRoblox(directory string) error {
	var worked, info = getLatestRoblox()
	if worked != nil {return worked}
	if !DownloadSoftware(directory, info) {return worked}
	return nil
}

type EnvironmentInfo struct {
	StateFile	    string
	PrefixDir	    string
	RosemaryDir	    string
	RobloxDir	    string
	RosewaterFile	    string
	RosewaterIcons	    string
}

func UpdateEnvironment(envInfo EnvironmentInfo) {
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Updating environment info!"+reset)
	var stateFile string = envInfo.StateFile
	WriteLn(stateFile, 1, envInfo.StateFile)
	WriteLn(stateFile, 2, envInfo.PrefixDir)
	WriteLn(stateFile, 3, envInfo.RosemaryDir)
	WriteLn(stateFile, 4, envInfo.RobloxDir)
	WriteLn(stateFile, 5, envInfo.RosewaterFile)
	WriteLn(stateFile, 6, envInfo.RosewaterIcons)
}

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

	MakeDir(downloadsDir)
	MakeDir(prefixDir)
	MakeDir(stateDir)
	MakeDir(rosemaryDir)
	MakeDir(robloxDir)
	MakeDir(rosewaterIcons)

	MakeFile(stateFile)
	self, _ := os.Executable()
	CopyFile(self, rosewaterFile, 0755)
	envInfo.StateFile = stateFile
	envInfo.PrefixDir = prefixDir
	envInfo.RosemaryDir = rosemaryDir
	envInfo.RobloxDir = robloxDir
	envInfo.RosewaterFile = rosewaterFile
	envInfo.RosewaterIcons = rosewaterIcons

	UpdateEnvironment(envInfo)

	if downloadAfter {InitialDownloads()}
}

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

func getVersions() (bloxInfo SoftwareInfo, roseInfo SoftwareInfo) {
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Retrieving software versions"+reset)
	var worked, robloxInfo = getLatestRoblox()
	if worked != nil {panic(fmt.Errorf("There was an issue recovering info from roblox's website: %w", worked))}

	var worked2, rosemaryInfo = getLatestRosemary()
	if worked2 != nil {panic(fmt.Errorf("There was an issue recovering info from github's website: %w ", worked2))}

	return robloxInfo, rosemaryInfo
}

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

func InitialDownloads() {
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Downloading initial files"+reset)
	var success, envInfo = GetEnvironment()
	if !success {SetupEnvironment(true)}

	var robloxInfo, rosemaryInfo = getVersions()

	var robloxVersionExists, _ = FileExists(filepath.Join(envInfo.RobloxDir, robloxInfo.Name))
	if !robloxVersionExists {
		if err := downloadLatestRoblox(filepath.Join(envInfo.RobloxDir, robloxInfo.Name)); err != nil  {panic(fmt.Errorf("Could not download roblox: %w", err))}
	}

	var rosemaryVersionExists, _ = FileExists(filepath.Join(envInfo.RosemaryDir, rosemaryInfo.Name))
	if !rosemaryVersionExists {
		if err := downloadLatestRosemary(filepath.Join(envInfo.RosemaryDir, rosemaryInfo.Name)); err != nil  {panic(fmt.Errorf("Could not download rosemary: %w", err))}
	}
}

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
	cmd.Env = append(os.Environ(), "APPIMAGE_EXTRACT_AND_RUN=1")
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

func RunRoblox() {
	CheckForUpdates()
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Starting roblox"+reset)
	var success, envInfo = GetEnvironment()
	if !success {SetupEnvironment(true)}

	var _, rosemaryInfo = getVersions()

	var rosemaryExec = filepath.Join(envInfo.RosemaryDir, rosemaryInfo.Name)
	var robloxApp = filepath.Join(envInfo.PrefixDir,"root","Applications","RobloxPlayer.app")

	var rpc RPCQuery
	rpc.Action = "Playing a game on "
	rpc.Where = "Using Rosewater"
	RPCActivity(rpc)
	var envVariables = []string{
		"ROSEMARY_PREFIX="+filepath.Join(envInfo.PrefixDir,"root"),
		"ROSEMARY_DECORATIONS=server",
		"ROSEMARY_MENU=none",
		"ROSEMARY_STEAM=0",
	}
	cmd := exec.Command(rosemaryExec, "run", robloxApp)
	cmd.Env = append(os.Environ(), envVariables...)
	if err := cmd.Run(); err != nil {panic(fmt.Errorf("Failed to run the roblox player! %s", envInfo.PrefixDir))}
}

func Uninstall() {
	var userData = ReturnData()
	var rosewaterDir string = filepath.Join(userData.AppDir, "Rosewater")
	os.RemoveAll(rosewaterDir)

	var desktopFile = filepath.Join(userData.AppDir, "applications","rosewateruri.desktop")
	var desktopFile2 = filepath.Join(userData.AppDir, "applications","rosewater.desktop")
	os.Remove(desktopFile)
	os.Remove(desktopFile2)

	RunShell("update-desktop-database ~/.local/share/applications")
}

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

func RunHeadless(uri string) bool {
	fmt.Println(bold + green + "/ᐠ˵- ⩊ -˵マ: " + reset + blue + "Starting roblox" + reset)
	var success, envInfo = GetEnvironment()
	if !success {
		os.Exit(1)
	}

	var _, rosemaryInfo = getVersions()

	var rosemaryExec = filepath.Join(envInfo.RosemaryDir, rosemaryInfo.Name)
	var robloxExe = filepath.Join(envInfo.PrefixDir, "root", "Applications", "RobloxPlayer.app", "Contents", "MacOS", "RobloxPlayer")

	args := []string{"run", robloxExe, uri}
	cmd := exec.Command(rosemaryExec, args...)

	var rpc RPCQuery
	rpc.Action = "Playing a game on "
	rpc.Where = "Using Rosewater"

	RPCActivity(rpc)
	var envVariables = []string{
		"ROSEMARY_PREFIX="+filepath.Join(envInfo.PrefixDir,"root"),
		"ROSEMARY_DECORATIONS=server",
		"ROSEMARY_MENU=none",
		"ROSEMARY_STEAM=0",
	}
	cmd.Env = append(os.Environ(), envVariables...)

	var logfileLocation = filepath.Join(envInfo.PrefixDir, "client.log")
	logfile, _ := os.Create(logfileLocation)
	defer logfile.Close()

	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true} //yay process group yaayy
	cmd.WaitDelay = 5 * time.Second

	logPipeRead, logPipeWrite := io.Pipe()

	cmd.Stdout = logPipeWrite
	cmd.Stderr = logPipeWrite

	if err := cmd.Start(); err != nil {
		return false
	}

	var closed bool = false
	stop := func() {
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

	rules := []LogRule {
		{
			Match: "leaveUGCGameInternal",
			Handler: func(line string) {stop()},
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

	if closed {
		return true
	}

	return exited == nil
}


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
