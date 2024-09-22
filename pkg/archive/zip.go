package archive

import (
	"archive/zip"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type ZipManager struct {
	maxZipSize int
}

func CreateZipManager(maxZipSize int) Manager {
	return &ZipManager{
		maxZipSize,
	}
}

func (m *ZipManager) ExtractStraight(fileHeader *multipart.FileHeader, extractPath string) error {
	// Open the file from the multipart.FileHeader
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create a zip.Reader from the file
	zipReader, err := zip.NewReader(file, fileHeader.Size)
	if err != nil {
		return fmt.Errorf("failed to create zip reader: %w", err)
	}

	// Allocate a buffer for copying files
	buffer := make([]byte, 32*1024) // 32 KB buffer

	for _, f := range zipReader.File {
		err := func() error {
			rc, err := f.Open()
			if err != nil {
				return fmt.Errorf("failed to open zip file: %w", err)
			}
			defer rc.Close()

			path := filepath.Join(extractPath, f.Name)
			if f.FileInfo().IsDir() {
				return os.MkdirAll(path, f.Mode())
			}

			err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}

			outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return fmt.Errorf("failed to create file on disk: %w", err)
			}
			defer outFile.Close()

			_, err = io.CopyBuffer(outFile, rc, buffer)
			if err != nil {
				return fmt.Errorf("failed to copy file content: %w", err)
			}

			return nil
		}()
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *ZipManager) Validate(file *multipart.FileHeader) (*Info, error) {

	fileHandle, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer fileHandle.Close()

	zipReader, err := zip.NewReader(fileHandle, file.Size)
	if err != nil {
		//fileHandle.Close()
		return nil, err
	}

	archiveInfo, err := m.getArchiveInfo(zipReader, MbMul*m.maxZipSize)
	if err != nil {
		//fileHandle.Close()
		return nil, err
	}
	if archiveInfo.Size < 0 {
		//fileHandle.Close()
		return nil, errors.New(fmt.Sprintf("Content in zip Is Too Large | %s, maxSize - %s", file.Filename, ConvertSize(int64(MbMul*m.maxZipSize))))
	}

	_, err = fileHandle.Seek(0, 0)
	if err != nil {
		//fileHandle.Close()
		return nil, err
	}

	return &Info{
		Archive: archiveInfo,
	}, nil
}

func (m *ZipManager) calculateHash(file *zip.File) ([]byte, error) {
	rc, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	// Compute the Hash of the file content
	fileHasher := sha256.New()
	if _, err := io.Copy(fileHasher, rc); err != nil {
		return nil, err
	}

	return fileHasher.Sum(nil), nil
}

func (m *ZipManager) getArchiveInfo(zipReader *zip.Reader, maxSize int) (*archive, error) {
	size := 0
	hasher := sha256.New()

	for _, file := range zipReader.File {
		size += int(file.FileInfo().Size())
		if size > maxSize {
			return nil, errors.New(fmt.Sprintf("max Size exceeded - %s", ConvertSize(int64(maxSize))))
		}

		hash, err := m.calculateHash(file)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("can't calculate Hash - %v", err))
		}

		_, err = hasher.Write(hash)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("can't write Hash - %v", err))
		}
	}

	return &archive{
		Size: size,
		Hash: fmt.Sprintf("%x", hasher.Sum(nil)),
	}, nil
}

func (m *ZipManager) Extract(archivePath string, extractPath string) error {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer r.Close()

	// Allocate a buffer for copying files
	buffer := make([]byte, 32*1024) // 32 KB buffer

	for _, f := range r.File {
		err := func() error {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			path := filepath.Join(extractPath, f.Name)
			if f.FileInfo().IsDir() {
				return os.MkdirAll(path, f.Mode())
			}

			err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
			if err != nil {
				return err
			}

			outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.CopyBuffer(outFile, rc, buffer)
			return err
		}()
		if err != nil {
			return err
		}
	}
	return nil
}
