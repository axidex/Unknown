package archive

import (
	"archive/tar"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type TarManager struct {
	maxZipSize int
}

func CreateTarManager(maxZipSize int) Manager {
	return &TarManager{maxZipSize: maxZipSize}
}

func (m *TarManager) ExtractStraight(fileHeader *multipart.FileHeader, extractPath string) error {
	// Open the file from the multipart.FileHeader
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create a tar.Reader from the file
	reader := tar.NewReader(file)

	// Allocate a buffer for copying files
	buffer := make([]byte, 32*1024) // 32 KB buffer

	for {
		header, err := reader.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		path := filepath.Join(extractPath, header.Name)
		info := header.FileInfo()

		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			continue
		}

		if err = os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create parent directories: %w", err)
		}

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
		if err != nil {
			return fmt.Errorf("failed to create file on disk: %w", err)
		}

		_, err = io.CopyBuffer(outFile, reader, buffer)
		outFile.Close()
		if err != nil {
			return fmt.Errorf("failed to copy file content: %w", err)
		}
	}

	return nil
}

func (m *TarManager) Validate(file *multipart.FileHeader) (*Info, error) {
	fileHandle, err := file.Open()
	defer fileHandle.Close()

	if err != nil {
		return nil, err
	}
	_, err = fileHandle.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	tarReader := tar.NewReader(fileHandle)

	archiveInfo, err := m.getArchiveInfo(tarReader, MbMul*m.maxZipSize)
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

func (m *TarManager) calculateHash(reader *tar.Reader) ([]byte, error) {
	// Compute the Hash of the file content
	fileHasher := sha256.New()

	if _, err := io.Copy(fileHasher, reader); err != nil {
		return nil, errors.New(fmt.Sprintf("Can't read file to calculate Hash - %v", err))
	}

	return fileHasher.Sum(nil), nil
}

func (m *TarManager) getArchiveInfo(tarReader *tar.Reader, maxSize int) (*archive, error) {
	size := 0
	hasher := sha256.New()

	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		size += int(header.Size)
		if size > maxSize {
			return nil, errors.New(fmt.Sprintf("max Size exceeded - %s", ConvertSize(int64(maxSize))))
		}

		hash, err := m.calculateHash(tarReader)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("can't calculate Hash - %v", err))
		}

		_, err = hasher.Write(hash)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("can't write Hash - %v", err))
		}
	}
	// Compute the final Hash of the archive
	archiveHash := hasher.Sum(nil)

	return &archive{
		Size: size,
		Hash: fmt.Sprintf("%x", archiveHash),
	}, nil
}

func (m *TarManager) Extract(archivePath string, extractPath string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := tar.NewReader(file)

	buffer := make([]byte, 32*1024)

	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		path := filepath.Join(extractPath, header.Name)
		info := header.FileInfo()

		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		if err = os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
		if err != nil {
			return err
		}

		_, err = io.CopyBuffer(outFile, reader, buffer)
		outFile.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
