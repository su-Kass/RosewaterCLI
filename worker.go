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

type SoftwareInfo struct {
	Download    string
	Version     string
	Name        string
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

func downloadSoftware(directory string, info SoftwareInfo) bool {
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Downloading "+info.Name+reset)
	if !DownloadFileURL(info.Download, directory) {return false}
	return true
}

func downloadLatestRosemary(directory string) bool {
	var worked, info = getLatestRosemary()
	if !worked {return false}
	if !downloadSoftware(directory, info) {return false}
	return true
}

func downloadLatestRoblox(directory string) bool {
	var worked, info = getLatestRoblox()
	if !worked {return false}
	if !downloadSoftware(directory, info) {return false}
	return true
}

type EnvironmentInfo struct {
	StateFile	    string
	PrefixDir	    string
	RosemaryDir	    string
	RobloxDir	    string
}

func UpdateEnvironment(envInfo EnvironmentInfo) {
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Updating environment info!"+reset)
	var stateFile string = envInfo.StateFile
	WriteLn(stateFile, 1, envInfo.StateFile)
	WriteLn(stateFile, 2, envInfo.PrefixDir)
	WriteLn(stateFile, 3, envInfo.RosemaryDir)
	WriteLn(stateFile, 4, envInfo.RobloxDir)
}

func SetupEnvironment(downloadAfter bool) {
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Setting up Rosewater environment"+reset)
	var envInfo EnvironmentInfo
	var userData = ReturnData()

	var downloadsDir string = filepath.Join(userData.AppDir, "Rosewater", "Downloads")
	var rosemaryDir string = filepath.Join(userData.AppDir, "Rosewater", "Downloads", "Rosemary")
	var robloxDir string = filepath.Join(userData.AppDir, "Rosewater", "Downloads", "Roblox")

	var prefixDir string = filepath.Join(userData.AppDir, "Rosewater", "Prefix")
	var stateDir string = filepath.Join(userData.AppDir, "Rosewater", "State")
	var stateFile string = filepath.Join(userData.AppDir, "Rosewater", "State", "scriptState")

	MakeDir(downloadsDir)
	MakeDir(prefixDir)
	MakeDir(stateDir)
	MakeDir(rosemaryDir)
	MakeDir(robloxDir)

	MakeFile(stateFile)
	envInfo.StateFile = stateFile
	envInfo.PrefixDir = prefixDir
	envInfo.RosemaryDir = rosemaryDir
	envInfo.RobloxDir = robloxDir
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

	envInfo.StateFile = stateFile
	envInfo.PrefixDir = prefixDir
	envInfo.RosemaryDir = rosemaryDir
	envInfo.RobloxDir = robloxDir

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
	if err := cmd.Run(); err != nil {panic(fmt.Errorf("Failed to run the roblox player! %s", envInfo.PrefixDir))}
}

func Uninstall() {
	var userData = ReturnData()
	var rosewaterDir string = filepath.Join(userData.AppDir, "Rosewater")
	os.RemoveAll(rosewaterDir)
}
