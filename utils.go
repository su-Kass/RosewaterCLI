package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"io"
	"io/fs"
	"time"
	"errors"
	"encoding/json"
	"archive/zip"
)

type UserData struct {
	HmeDir	     string	// user home dir
	CurDir	     string	// current dir
	TmpDir	     string	// temp dir
	AppDir	     string	// application dir
	SysType	     string	// reports system type
	CpuNum	     string	// reports amount of cpu cores
	Username     string	// user username
}

func ReturnArguments() (flag string, data string) {
	if len(os.Args) > 1 {
		flag = os.Args[1]
	}
	if len(os.Args) > 2 {
		data = os.Args[2]
	}
	return flag, data
}


func ErrHandle(err error, exit bool) bool {
	var errBool bool = false
	if err != nil {
		if exit {
			panic(err)
		}
		errBool = true
	}
	return errBool
}

func SearchForFiletype(folder, filetype string) (string, error) {
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Searching for "+filetype+" in "+folder+reset)
	matches, err := filepath.Glob(filepath.Join(folder, "*."+filetype))
	if err != nil {
		return "", err
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("no file matching *.%s in %s", filetype, folder)
	}
	return filepath.Base(matches[0]), nil
}

//Runs a command in the user's shell; hopefully used only to run system programs, however can be used to do anything
//
// Has command and then arguments:
//
// echo <--- command
//
// 'Hello World!' <--- argument
func RunShell(command string) error {
	cmd := exec.Command("sh", "-c", command)

	if err := cmd.Run(); err != nil {
		fmt.Println(bold+red+"↶(•̀ᴗ•́)Ψ :"+reset+yellow+" Your command failed!"+reset)
		return err
	}
	return nil
}

// Simply makes a directory, dont question why it exists.
func MakeDir(Dir string) bool {
	fmt.Println(bold+red+"↶(•̀ᴗ•́)Ψ :"+reset+yellow+" Making "+Dir+reset)
	if info, err := os.Stat(Dir); err == nil && info.IsDir() {
		return true
	}
	os.MkdirAll(Dir, 0755)
	return false
}

func MoveFile(src, dst string) error {
	if err := os.Rename(src, dst); err != nil {
		return fmt.Errorf("could not move %s to %s: %w", src, dst, err)
	}
	return nil
}

// Copies a file somewhere else
//
// File, Directory
func CopyFile(file string, dir string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return os.WriteFile(dir, data, 0755)
}

// Makes a file at the directory
//
// Returns an error
func MakeFile(fileName string) error {
	_, err := os.Create(fileName)
	return err
}

func Unzip(src, dest string) error {
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Unzipping "+src+" to "+dest+reset)
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		target := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(target, 0755)
			continue
		}
		os.MkdirAll(filepath.Dir(target), 0755)

		rc, err := f.Open()
		if err != nil {
			return err
		}

		if f.Mode()&os.ModeSymlink != 0 {
			link, _ := io.ReadAll(rc)
			os.Symlink(string(link), target)
		} else {
			out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode().Perm())
			if err != nil {
				rc.Close()
				return err
			}
			io.Copy(out, rc)
			out.Close()
		}
		rc.Close()
	}
	return nil
}

func FileExists(path string) (bool, error) {
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Checking path to "+path+reset)
	info, err := os.Stat(path)
	if err == nil {
		return !info.IsDir(), nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err // permissions or some other problem
}

// Writes a line for a given file
//
// Returns error
func WriteLn(fileName string, line int, text string) error{
	fmt.Println(bold+red+"↶(•̀ᴗ•́)Ψ :"+reset+yellow+" Writing "+text+" to "+fileName+reset)
	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")

	for len(lines) < line {
		lines = append(lines, "")
	}
	lines[line-1] = text
	return os.WriteFile(fileName, []byte(strings.Join(lines, "\n")), 0644)
}

// Reads the given line for the given file.ReadLn
//
// Returns data + success bool
func ReadLn(fileName string, line int) (string, bool){
	data, err := os.ReadFile(fileName)
	if err != nil {
		return "", false
	}
	lines := strings.Split(string(data), "\n")

	if line < 1 || line > len(lines) {
		return "", false
	}
	return lines[line-1], true
}

// Downloads a file from a given URL and puts it into a given directory with a certain name.
func DownloadFileURL(url string, fileName string) bool {

	data, err := http.Get(url)
	if err != nil {
		fmt.Println(bold+red+"↶(•̀ᴗ•́)Ψ :"+reset+yellow+" Failed to download file from "+url+reset)
		return false
	}
	defer data.Body.Close()

	fmt.Println(bold+red+"↶(•̀ᴗ•́)Ψ :"+reset+yellow+" Creating file "+fileName+reset)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println(bold+red+"↶(•̀ᴗ•́)Ψ :"+reset+yellow+" Failed to download file to "+fileName+reset)
		return false
	}
	defer file.Close()
	fmt.Println(bold+red+"↶(•̀ᴗ•́)Ψ :"+reset+green+" Copying file to "+fileName+reset)
	if _, err := io.Copy(file, data.Body); err != nil {
		fmt.Println(bold+red+"↶(•̀ᴗ•́)Ψ :"+reset+yellow+"An error occured during the download!"+reset)
		return false
	}

	return true
}

type clientVersion struct {
	Version             string `json:"version"`
	ClientVersionUpload string `json:"clientVersionUpload"`
}

// reads the roblox domain to get the version
func ReadDomain(url string) (version string, err error) {
	fmt.Println(bold+red+"↶(•̀ᴗ•́)Ψ :"+reset+yellow+" Reading data from "+url+reset)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var cv clientVersion
	if err := json.NewDecoder(resp.Body).Decode(&cv); err != nil {
		return "", err
	}
	if cv.ClientVersionUpload == "" {
		return "", errors.New("clientVersionUpload missing from response")
	}
	return cv.ClientVersionUpload, nil
}

type SoftwareInfo struct {
	Download    string
	Version     string
	Name        string
}

func DownloadSoftware(directory string, info SoftwareInfo) bool {
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Downloading "+info.Name+reset)
	if !DownloadFileURL(info.Download, directory) {return false}
	return true
}

// ReturnData returns a map of the user's data needed for setting up stuff
//
// returns map: HOME CURRENT SYSTYPE TEMPDIR CPUNUM USERNAME
func ReturnData() UserData {
	fmt.Println(bold+red+"↶(•̀ᴗ•́)Ψ :"+reset+yellow+" Fetching user Data!"+reset)
	var userData UserData
	var homeDir, _ = os.UserHomeDir()
	var currentDir, _ = os.Getwd()

	var systemType string = runtime.GOOS
	var tempFolder string = os.TempDir()
	var cpuNum string = strconv.Itoa(runtime.NumCPU())

	var userName string
	var applicationDir string
	switch runtime.GOOS {
		case "linux":
			applicationDir = filepath.Join(homeDir, ".local", "share")
			userName = os.Getenv("USER")
		case "windows":
			applicationDir = filepath.Join(os.Getenv("LOCALAPPDATA"))
			userName = os.Getenv("USERNAME")
		default: applicationDir = "UNKNOWN!"
	}
	userData.HmeDir = homeDir
	userData.CurDir = currentDir
	userData.SysType = systemType
	userData.TmpDir = tempFolder
	userData.CpuNum = cpuNum
	userData.Username = userName
	userData.AppDir = applicationDir
	return userData
}
