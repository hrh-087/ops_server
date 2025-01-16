package utils

import "os"

func CreateFile(filePath, fileName string, text string) (string, error) {

	// 检测目录是否存在
	err := os.MkdirAll(filePath, 0755)
	if err != nil {
		return "", err
	}

	fullPath := filePath + "/" + fileName

	err = os.WriteFile(fullPath, []byte(text), 0644)
	if err != nil {
		return "", err
	}

	return fullPath, nil
}
