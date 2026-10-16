package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func ClearDir(dir string) {
	fmt.Println(bold+green+"/ᐠ˵- ⩊ -˵マ: "+reset+blue+"Cleaning "+dir+reset)
	os.RemoveAll(dir)
	os.MkdirAll(dir, 0755)
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

func CopyFile(file string, filepath string, permission int32) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath, data, os.FileMode(permission))
}

func MakeFile(fileName string) error {
	file, err := os.Create(fileName)
	if err == nil {
		defer file.Close()
	}
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
		if !filepath.IsLocal(f.Name) {return fmt.Errorf("Illegal path in zip: %s", f.Name)}

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
	return false, err
}

func WriteLn(fileName string, line int, text string) error{
	if line < 1 {
		return fmt.Errorf("Could not write line, line written is less than 1!")
	}
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
