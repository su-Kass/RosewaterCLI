package main

// rosewater is library slop gluecode confirmed
import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
	"net/url"
	"strings"
	"bufio"
)

// god i love structs...
type UserData struct {
	HmeDir	     string	// user home dir
	CurDir	     string	// current dir
	TmpDir	     string	// temp dir
	AppDir	     string	// application dir
	SysType	     string	// reports system type
	CpuNum	     int	// reports amount of cpu cores
	Username     string	// user username
}

// returns the file args cause i love arguably useless functions in my code
func ReturnArguments() (flag string, data string) {
	if len(os.Args) > 1 {flag = os.Args[1]}
	if len(os.Args) > 2 {data = os.Args[2]}
	return flag, data
}

// i dont even remember if i use this wait, let me check...
// yeah okay i use it twice in the diagnostic tool, for literally no reason
// god i am the smartest programmer ever
func ErrHandle(err error, exit bool) error {
	if err != nil {
		if exit {
			panic(err)
		}
	}
	return err
}

// i dont like searching for directories so i just use this for the xdg-mime and desktop stuff
// hopefully nobody explodes from it
func RunShell(command string) error {
	cmd := exec.Command("sh", "-c", command)

	if err := cmd.Run(); err != nil {
		var errout = fmt.Errorf("You had an error: %s", err)
		fmt.Println(bold+red+"↶(•̀ᴗ•́)Ψ :"+reset+yellow+" %w"+reset, errout)
		return err
	}
	return nil
}

// yeah ummm, chief i stole so much of this code so i barely know how it works
// but downloads a file from a url
// honestly net/http is a pretty sick library
func DownloadFileURL(url string, fileName string) error {
	tr := http.DefaultTransport.(*http.Transport).Clone() // trans
	tr.ResponseHeaderTimeout = 30 * time.Second // timeout
	client := &http.Client{Transport: tr}// transport.... transgender...

	resp, err := client.Get(url) // i guess bro
	if err != nil {
		fmt.Println(bold + red + "↶(•̀ᴗ•́)Ψ :" + reset + yellow + " %w Failed to download file from " + url + reset, err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(bold + red + "↶(•̀ᴗ•́)Ψ :" + reset + yellow + " Failed to download file from " + url + reset)
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	tmp := fileName + ".part" // mmmm temp files, i didnt get the memo on where this even exists but it works soooooooo...
	file, err := os.Create(tmp)
	if err != nil {
		fmt.Println(bold + red + "↶(•̀ᴗ•́)Ψ :" + reset + yellow + " %w Failed to create " + tmp + reset, err)
		return err
	}

	n, err := io.Copy(file, resp.Body) // wow!!! also while im here im gonna make this a single line
	if cerr := file.Close(); err == nil {err = cerr} // much better
	if err == nil && resp.ContentLength > 0 && n != resp.ContentLength {err = fmt.Errorf("truncated download: got %d of %d bytes", n, resp.ContentLength)} // fuck you!
	if err != nil { // nevermind i love you, i could never forget you \n
		os.Remove(tmp)
		fmt.Printf(bold + red + "↶(•̀ᴗ•́)Ψ :" + reset + yellow + " An error occured during the download: %v\n" + reset, err)
		return err
	}

	return os.Rename(tmp, fileName)
}

const assetRepo string = "https://raw.githubusercontent.com/su-Kass/RosewaterCLI/main/" // heh, it's me!

// now you may ask yourself, why dont i just pass these directly to the functions
// and to that i say.... ummm... this looks nicer and is object oriented
// hah! take that!
type Asset struct {
	Directory string
	Asset 	  string
}

func DownloadAsset(directory string, asset Asset) string {
	assetURL, err := url.JoinPath(assetRepo, asset.Directory, asset.Asset)
	if err != nil {
		fmt.Printf(bold+red+"↶(•̀ᴗ•́)Ψ :"+reset+yellow+" Bad asset URL: %v"+reset+"\n", err)
		return ""
	}

	fileName := filepath.Join(directory, asset.Asset)
	if err := DownloadFileURL(assetURL, fileName); err != nil {return ""} // honestly new lines are overrated
	return fileName
}

// yeah i dont even remember why i do it this way, couldnt i just decode it without
// nevermind...
type clientVersion struct {
	Version             string `json:"version"`
	ClientVersionUpload string `json:"clientVersionUpload"`
}

func ReadDomain(url string) (version string, err error) {
	fmt.Println(bold+red+"↶(•̀ᴗ•́)Ψ :"+reset+yellow+" Reading data from "+url+reset)

	client := &http.Client{Timeout: 15 * time.Second} // you know im sure someone will eventually take 15 seconds to connect to the roblox domains
	resp, err := client.Get(url)
	if err != nil {return "", err}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {return "", fmt.Errorf("unexpected status: %s", resp.Status)} // why did i even add new lines for these? am i stupid?

	var cv clientVersion
	if err := json.NewDecoder(resp.Body).Decode(&cv); err != nil {return "", err}
	if cv.ClientVersionUpload == "" {return "", fmt.Errorf("Missing ClientVersionUpload")}

	return cv.ClientVersionUpload, nil
}

// heh, it's my loggerona!
// this is actually a keylogger, fear rosewater
type LogRule struct {
	Match   string
	Handler func(line string)
}

// ok wait i actually dont remember how i use this, give me a second to write actual docs for it
// LogScan() takes io.read and io.Writer inputs, along with an array of rules to check for
// It will scan the input and check for the rules (LogRule struct), before outputting it to
// the output file/writer
//
// thank god i can go back to shitposting
func LogScan(r io.Reader, out io.Writer, rules []LogRule) {
	reader := bufio.NewReader(r)
	for {
		line, err := reader.ReadString('\n')
		if len(line) > 0 {

			if out != nil {io.WriteString(out, line)}
			for _, rule := range rules {if strings.Contains(line, rule.Match) {rule.Handler(line)}
			}
		}
		if err != nil {return}
	}
}

// software stuffs, why dont i just construct the name from the version in the downloads script?
// oh wait i do but in only one of them?
// smartest programmer ALIVE
type SoftwareInfo struct {
	Download    string
	Version     string
	Name        string
}

// look man if you're reading this you have the source code, do you really need a comment for this
func DownloadSoftware(directory string, info SoftwareInfo) error {
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Downloading "+info.Name+reset)
	var err = DownloadFileURL(info.Download, directory)
	if err != nil {return err}
	return err
}

// returns user data with a struct because ime stealing your data, to vibecode using claude fear me!!! :3
func ReturnData() UserData {
	fmt.Println(bold+red+"↶(•̀ᴗ•́)Ψ :"+reset+yellow+" Fetching user Data!"+reset)
	var userData UserData
	var currentDir, _ = os.Getwd()

	var systemType string = runtime.GOOS
	var tempFolder string = os.TempDir()
	var cpuNum int = runtime.NumCPU()

	var userName string
	var applicationDir string
	var env1 string = os.Getenv("XDG_DATA_HOME") // << feeling cute, might package flatpak later :3
	if env1 == "" {
		var dir, err = os.UserHomeDir()
		if err == nil {applicationDir = filepath.Join(dir, ".local", "share")} else {panic(fmt.Errorf("Could not find user's home directory: %w", err))}
	} else {applicationDir = env1}
	userData.CurDir = currentDir
	userData.SysType = systemType
	userData.TmpDir = tempFolder
	userData.CpuNum = cpuNum
	userData.Username = userName
	userData.AppDir = applicationDir
	return userData
}
