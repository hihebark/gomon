package engine

import (
	"io/ioutil"
	"os"
)

func ReadFile(path string) (string, error) {
	if _, err := os.Stat(path); os.IsExist(err) {
		return "", err
	}
	file, err := os.OpenFile(path, os.O_RDONLY, 0555)
	defer file.Close()
	if err != nil {
		return "", err
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func FileExists(fileName string) bool {
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
