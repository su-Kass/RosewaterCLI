package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type release struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name string `json:"name"`
		URL  string `json:"browser_download_url"`
	} `json:"assets"`
}

func getLatestRosemary() (fetched bool, rosemaryInfo SoftwareInfo) {
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Fetching latest Rosemary release!"+reset)
	const api = "https://api.github.com/repos/microlift/rosemary/releases/latest"

	var info SoftwareInfo
	req, err := http.NewRequest("GET", api, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
		return false, info
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false, info
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintln(os.Stderr, "unexpected status:", resp.Status)
		return false, info
	}

	var rel release
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false, info
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
	return true, info
}

func getLatestRoblox() (worked bool, robloxInfo SoftwareInfo) {
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Fetching latest Roblox release!"+reset)
	const api = "https://clientsettingscdn.roblox.com/v2/client-version/MacPlayer"
	data, err := ReadDomain(api)
	var info SoftwareInfo
	if err != nil {return false, info}
	var fileName string = data+"-RobloxPlayer.zip"
	var zipURL string = "https://setup.rbxcdn.com/mac/" + fileName
	info.Download = zipURL
	info.Name = fileName
	info.Version = data
	return true, info
}

func downloadLatestRosemary(directory string) bool {
	var worked, info = getLatestRosemary()
	if !worked {return false}
	if !DownloadSoftware(directory, info) {return false}
	return true
}

func downloadLatestRoblox(directory string) bool {
	var worked, info = getLatestRoblox()
	if !worked {return false}
	if !DownloadSoftware(directory, info) {return false}
	return true
}

type EnvironmentInfo struct {
	StateFile	    string
	PrefixDir	    string
	RosemaryDir	    string
	RobloxDir	    string
	RosewaterFile	    string
}

func UpdateEnvironment(envInfo EnvironmentInfo) {
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Updating environment info!"+reset)
	var stateFile string = envInfo.StateFile
	WriteLn(stateFile, 1, envInfo.StateFile)
	WriteLn(stateFile, 2, envInfo.PrefixDir)
	WriteLn(stateFile, 3, envInfo.RosemaryDir)
	WriteLn(stateFile, 4, envInfo.RobloxDir)
	WriteLn(stateFile, 5, envInfo.RosewaterFile)
}

func SetupEnvironment(downloadAfter bool) {
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Setting up Rosewater environment"+reset)
	var envInfo EnvironmentInfo
	var userData = ReturnData()

	var downloadsDir string = filepath.Join(userData.AppDir, "Rosewater", "Downloads")
	var rosemaryDir string = filepath.Join(userData.AppDir, "Rosewater", "Downloads", "Rosemary")
	var robloxDir string = filepath.Join(userData.AppDir, "Rosewater", "Downloads", "Roblox")
	var rosewaterFile string = filepath.Join(userData.AppDir, "Rosewater", "RosewaterCLI.Headless")

	var prefixDir string = filepath.Join(userData.AppDir, "Rosewater", "Prefix")
	var stateDir string = filepath.Join(userData.AppDir, "Rosewater", "State")
	var stateFile string = filepath.Join(userData.AppDir, "Rosewater", "State", "scriptState")

	MakeDir(downloadsDir)
	MakeDir(prefixDir)
	MakeDir(stateDir)
	MakeDir(rosemaryDir)
	MakeDir(robloxDir)

	MakeFile(stateFile)
	self, _ := os.Executable()
	CopyFile(self, rosewaterFile)
	envInfo.StateFile = stateFile
	envInfo.PrefixDir = prefixDir
	envInfo.RosemaryDir = rosemaryDir
	envInfo.RobloxDir = robloxDir
	envInfo.RosewaterFile = rosewaterFile
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

	envInfo.StateFile = stateFile
	envInfo.PrefixDir = prefixDir
	envInfo.RosemaryDir = rosemaryDir
	envInfo.RobloxDir = robloxDir
	envInfo.RosewaterFile = rosewaterFile

	return true, envInfo
}

func clearDir(dir string) {
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Cleaning "+dir+reset)
	os.RemoveAll(dir)
	os.MkdirAll(dir, 0755)
}

func getVersions() (bloxInfo SoftwareInfo, roseInfo SoftwareInfo) {
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Retrieving software versions"+reset)
	var worked, robloxInfo = getLatestRoblox()
	if !worked {panic(fmt.Errorf("There was an issue recovering info from roblox's website!"))}

	var worked2, rosemaryInfo = getLatestRosemary()
	if !worked2 {panic(fmt.Errorf("There was an issue recovering info from github's website!"))}

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
		clearDir(envInfo.RobloxDir)
		if !downloadLatestRoblox(filepath.Join(envInfo.RobloxDir, robloxInfo.Name)) {panic(fmt.Errorf("Could not download roblox."))}
		fmt.Println(bold+green+"Roblox has been updated to "+robloxInfo.Version+reset)
	}

	var rosemaryVersionExists, _ = FileExists(filepath.Join(envInfo.RosemaryDir, rosemaryInfo.Name))
	if rosemaryVersionExists {
		fmt.Println(bold+green+"Rosemary is fully updated to "+rosemaryInfo.Version+reset)} else
	{
		fmt.Println(red+underline+"OUT OF DATE!"+reset+red+" : "+reset+yellow+"Updating Rosemary to "+rosemaryInfo.Version+reset)
		clearDir(envInfo.RosemaryDir)
		if !downloadLatestRosemary(filepath.Join(envInfo.RosemaryDir, rosemaryInfo.Name)) {panic(fmt.Errorf("Could not download rosemary"))}
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
		if !downloadLatestRoblox(filepath.Join(envInfo.RobloxDir, robloxInfo.Name)) {panic(fmt.Errorf("Could not download roblox."))}
	}

	var rosemaryVersionExists, _ = FileExists(filepath.Join(envInfo.RosemaryDir, rosemaryInfo.Name))
	if !rosemaryVersionExists {
		if !downloadLatestRosemary(filepath.Join(envInfo.RosemaryDir, rosemaryInfo.Name)) {panic(fmt.Errorf("Could not download rosemary"))}
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
	if err := cmd.Run(); err != nil {panic(fmt.Errorf("Failed to create prefix in %s", envInfo.PrefixDir))}

	Unzip(robloxZip, filepath.Join(envInfo.PrefixDir, "root", "Applications"))
}

func RunRoblox() {
	CheckForUpdates()
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Starting roblox"+reset)
	var success, envInfo = GetEnvironment()
	if !success {SetupEnvironment(true)}

	var _, rosemaryInfo = getVersions()

	var rosemaryExec = filepath.Join(envInfo.RosemaryDir, rosemaryInfo.Name)
	var robloxApp = filepath.Join(envInfo.PrefixDir,"root","Applications","RobloxPlayer.app")

	cmd := exec.Command(rosemaryExec, "run", robloxApp)
	cmd.Env = append(os.Environ(), "ROSEMARY_PREFIX="+filepath.Join(envInfo.PrefixDir,"root"))
	if err := cmd.Run(); err != nil {panic(fmt.Errorf("Failed to run the roblox player! %s", envInfo.PrefixDir))}
}

func Uninstall() {
	var userData = ReturnData()
	var rosewaterDir string = filepath.Join(userData.AppDir, "Rosewater")
	os.RemoveAll(rosewaterDir)

	var desktopFile = filepath.Join(userData.AppDir, "applications","rosewater.desktop")
	os.Remove(desktopFile)

	RunShell("update-desktop-database ~/.local/share/applications")
}

func SetupDesktop() {
	var success, envInfo = GetEnvironment()
	if !success {SetupEnvironment(true)}

	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Setting up user .desktop handling for URI parsing"+reset)
	var userData = ReturnData()
	var desktopFile string = filepath.Join(userData.AppDir, "applications", "rosewater.desktop")
	MakeFile(desktopFile)

	WriteLn(desktopFile, 1, "[Desktop Entry]")
	WriteLn(desktopFile, 2, "Name=Rosewater")
	WriteLn(desktopFile, 3, "Type=Application")
	WriteLn(desktopFile, 4, "NoDisplay=true")
	WriteLn(desktopFile, 5, "Terminal=false")
	WriteLn(desktopFile, 6, "Exec="+envInfo.RosewaterFile+" --uri %u")
	WriteLn(desktopFile, 7, "MimeType=x-scheme-handler/roblox-player;")

	RunShell("update-desktop-database ~/.local/share/applications")
	RunShell("xdg-mime default "+desktopFile+" x-scheme-handler/roblox-player")
	RunShell("xdg-mime query default x-scheme-handler/roblox-player")
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

	cmd.Env = append(os.Environ(), "ROSEMARY_PREFIX="+filepath.Join(envInfo.PrefixDir, "root"))

	var logfileLocation = filepath.Join(envInfo.PrefixDir, "client.log")
	logfile, _ := os.Create(logfileLocation)
	defer logfile.Close()

	cmd.Stdout = logfile
	cmd.Stderr = logfile

	if err := cmd.Run(); err != nil {
		return false
	}
	return true
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
