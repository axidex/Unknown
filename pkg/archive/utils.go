package archive

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

func GetExtension(archivePath string, archiveCfg map[string][]string) (string, error) {
	ext := filepath.Ext(archivePath)
	ext = strings.ToLower(ext)

	for curExt, synonyms := range archiveCfg {
		for _, synonym := range synonyms {
			if ext == synonym {
				return curExt, nil
			}
		}
	}

	return "", errors.New(fmt.Sprintf("no extension found for %s | %s", archivePath, ext))
}

func SaveFile(fileHeader *multipart.FileHeader, dstPath string) error {
	// Open the uploaded file
	srcFile, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy the contents from the uploaded file to the destination file
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	_, err = srcFile.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	return nil
}

func CopyFile(srcPath, dstPath string) error {
	// Open the source file
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy the contents from the source file to the destination file
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func ConvertSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
