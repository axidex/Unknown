package archive

import (
	"mime/multipart"
)

const MbMul = 1024 * 1024

type Manager interface {
	Validate(file *multipart.FileHeader) (*Info, error)
	Extract(archivePath string, extractPath string) error
	ExtractStraight(fileHeader *multipart.FileHeader, extractPath string) error
}
