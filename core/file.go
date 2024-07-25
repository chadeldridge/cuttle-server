package core

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"syscall"

	"gopkg.in/yaml.v3"
)

var (
	tester func(string) error
	reader func(string) ([]byte, error)
)

func init() {
	tester = AssertReadable
	reader = os.ReadFile
}

func SetTester(t func(string) error) {
	tester = t
}

func SetReader(r func(string) ([]byte, error)) {
	reader = r
}

func ParseYAML[T any](file string, obj *T) error {
	data, err := reader(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, obj)
	if err != nil {
		return err
	}

	return nil
}

// FindFiles returns the location of the config file if found. Returns "" and an error if not.
func FindFiles(fileNames ...string) (string, error) {
	if len(fileNames) == 0 {
		return "", fmt.Errorf("no file names provided")
	}

	// Add the user's home directory.
	hd, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	file, err := FindInDir(hd, fileNames...)
	if err == nil {
		return file, nil
	}

	file, err = FindInDir(hd+"/.config/cuttle", fileNames...)
	if err == nil {
		return file, nil
	}

	// Add the current working directory. Assume the filename contains the app name so we do
	// not load in another app's config by mistake.
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	file, err = FindInDir(wd, fileNames...)
	if err == nil {
		return file, nil
	}

	return "", fmt.Errorf(FileNotFound)
}

// FindInDir returns the contents of the first file found in the directory. dir should not contain
// a trailing "/". If no file is found, return FileNotFound.
func FindInDir(dir string, fileNames ...string) (string, error) {
	switch dir {
	case "":
		break
	case "/":
		break
	default:
		dir = dir + "/"
	}

	if len(fileNames) == 0 || (len(fileNames) == 1 && fileNames[0] == "") {
		return "", fmt.Errorf("no file names provided")
	}

	// Check default locations for the file.
	for _, name := range fileNames {
		err := tester(dir + name)
		if err == nil {
			return dir + name, nil
		}
	}

	return "", fmt.Errorf(FileNotFound)
}

func ReadFile(file string) ([]byte, error) {
	if file == "" {
		return nil, fmt.Errorf("read: no file provided")
	}

	data, err := reader(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// CheckReadability returns nil if file is readable.
func AssertReadable(file string) error {
	if file == "" {
		return fmt.Errorf("read: no file provided")
	}

	s, err := os.Stat(file)
	if err != nil {
		return err
	}

	if s.IsDir() {
		return fmt.Errorf("read %s: is a directory", file)
	}

	err = HasReadPerm(s)
	if err != nil {
		return err
	}

	return nil
}

// HasReadPerm returns nil if the app has read permission to the file.
func HasReadPerm(info fs.FileInfo) error {
	if info == nil {
		return fmt.Errorf("read: no info provided")
	}

	u, err := user.Current()
	if err != nil {
		return err
	}

	if info.Mode().Perm()&0o004 == 0o004 {
		return nil
	}

	fileUid := fmt.Sprint(info.Sys().(*syscall.Stat_t).Uid)
	if u.Uid == fileUid {
		if info.Mode().Perm()&0o400 == 0o400 {
			return nil
		}
	}

	groups, err := u.GroupIds()
	if err != nil {
		return err
	}

	fileGid := fmt.Sprint(info.Sys().(*syscall.Stat_t).Gid)
	for _, group := range groups {
		if group == fileGid {
			if info.Mode().Perm()&0o040 == 0o040 {
				return nil
			}
		}
	}

	return fmt.Errorf("read: permission denied")
}
