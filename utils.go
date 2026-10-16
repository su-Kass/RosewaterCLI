package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

type UserData struct {
	HmeDir	     string	// user home dir
	CurDir	     string	// current dir
	TmpDir	     string	// temp dir
	AppDir	     string	// application dir
	SysType	     string	// reports system type
	CpuNum	     int	// reports amount of cpu cores
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


func ErrHandle(err error, exit bool) error {
	if err != nil {
		if exit {
			panic(err)
		}
	}
	return err
}

func RunShell(command string) error {
	cmd := exec.Command("sh", "-c", command)

	if err := cmd.Run(); err != nil {
		var errout = fmt.Errorf("You had an error: %s", err)
		fmt.Println(bold+red+"↶(•̀ᴗ•́)Ψ :"+reset+yellow+" %w"+reset, errout)
		return err
	}
	return nil
}

func DownloadFileURL(url string, fileName string) error {
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.ResponseHeaderTimeout = 30 * time.Second
	client := &http.Client{Transport: tr}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(bold + red + "↶(•̀ᴗ•́)Ψ :" + reset + yellow + " Failed to download file from " + url + reset)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(bold + red + "↶(•̀ᴗ•́)Ψ :" + reset + yellow + " Failed to download file from " + url + reset)
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	tmp := fileName + ".part"
	file, err := os.Create(tmp)
	if err != nil {
		fmt.Println(bold + red + "↶(•̀ᴗ•́)Ψ :" + reset + yellow + " Failed to create " + tmp + reset)
		return err
	}

	n, err := io.Copy(file, resp.Body)
	if cerr := file.Close(); err == nil {
		err = cerr
	}
	if err == nil && resp.ContentLength > 0 && n != resp.ContentLength {
		err = fmt.Errorf("truncated download: got %d of %d bytes", n, resp.ContentLength)
	}
	if err != nil {
		os.Remove(tmp)
		fmt.Println(bold + red + "↶(•̀ᴗ•́)Ψ :" + reset + yellow + " An error occured during the download!" + reset)
		return err
	}

	return os.Rename(tmp, fileName)
}

const assetRepo string = "https://raw.githubusercontent.com/su-Kass/RosewaterCLI/main/"
type Asset struct {
	Directory string
	Asset 	  string
}

func DownloadAsset(directory string, asset Asset) string {
	var url string = filepath.Join(assetRepo, asset.Directory, asset.Asset)
	var fileName = filepath.Join(directory, asset.Asset)
	DownloadFileURL(url, fileName)
	return fileName
}

type clientVersion struct {
	Version             string `json:"version"`
	ClientVersionUpload string `json:"clientVersionUpload"`
}

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
	if DownloadFileURL(info.Download, directory) != nil {return false}
	return true
}

func ReturnData() UserData {
	fmt.Println(bold+red+"↶(•̀ᴗ•́)Ψ :"+reset+yellow+" Fetching user Data!"+reset)
	var userData UserData
	var homeDir, _ = os.UserHomeDir()
	var currentDir, _ = os.Getwd()

	var systemType string = runtime.GOOS
	var tempFolder string = os.TempDir()
	var cpuNum int = runtime.NumCPU()

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
