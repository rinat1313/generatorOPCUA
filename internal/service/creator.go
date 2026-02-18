package service

import "os"

func CreateFile(path string, body string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(body)
	if err != nil {
		return err
	}
	return nil
}
