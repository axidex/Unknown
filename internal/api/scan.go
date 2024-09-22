package api

import (
	"context"
	"github.com/axidex/Unknown/pkg/archive"
	"github.com/axidex/Unknown/pkg/shell"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"os"
	"path"
	"time"
)

// Scan
// @Summary Start a new scan
// @Description Starts a new Secret Scanning task
// @Tags scans
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Archive file to scan"
// @Success 201 {object} map[string]interface{} "Task created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/v1/scan [post]
func (app *App) Scan(c *gin.Context) {

	scannerName := "gitleaks"

	timeoutCtx, cancel := context.WithTimeout(c, time.Second*time.Duration(app.config.Server.Deadlines.SS))
	defer cancel()

	// Unique task ID
	taskId := uuid.New().String()

	// Getting form file
	file, err := c.FormFile("file")
	if err != nil {
		app.logger.Errorf("Failed getting FormFile | %v", err)
		ResponseJSON(c, ErrFileNotProvided)
		return
	}

	// Checking archive size
	convertedArchiveSize := archive.ConvertSize(file.Size)
	if int(file.Size) > archive.MbMul*app.config.Archive.MaxSize {
		app.logger.Errorf("Archive size too large: %s", convertedArchiveSize)
		ResponseJSON(c, ErrUploadedFileTooLarge)
		return
	}
	app.logger.Infof("Archive Size: %s", convertedArchiveSize)

	// Getting extension and validating archive
	extension, err := archive.GetExtension(file.Filename, app.config.Archive.Extensions)
	if err != nil {
		app.logger.Errorf("Failed getting file extension | %s", err)
		ResponseJSON(c, ErrArchiveType)
		return
	}
	app.logger.Infof("File Extension: %s | %s", extension, taskId)
	archiveManager := app.archiveManagers[extension]

	archiveInfo, err := archiveManager.Validate(file)
	if err != nil {
		app.logger.Errorf("Failed validating archive | %s", err)
		ResponseJSON(c, ErrArchiveValidation)
		return
	}
	app.logger.Infof("Archive unpacked size - %s", archive.ConvertSize(int64(archiveInfo.Archive.Size)))

	app.logger.Infof("Start unpacking... | %s", taskId)
	tmpScanPath := path.Join(app.config.Workdir, taskId)
	err = archiveManager.ExtractStraight(file, tmpScanPath)
	if err != nil {
		app.logger.Errorf("Failed to extract archive | %s", err)
		ResponseJSON(c, ErrArchiveSaving)
		return
	}

	app.logger.Infof("Unpacking finished | %s", taskId)

	gitleaksService := app.services[scannerName]
	reportPath, err := gitleaksService.Scan(timeoutCtx, shell.CreateNewScan(tmpScanPath, taskId, scannerName))
	if err != nil {
		app.logger.Errorf("Failed scanning | %s", err)
		ResponseJSON(c, ErrScanFailed)
		return
	}

	app.logger.Infof("Report file - %s", reportPath)

	reportBytes, err := os.ReadFile(reportPath)
	if err != nil {
		app.logger.Errorf("Failed reading file | %s", err)
		ResponseJSON(c, ErrFailedToReadReport)
		return
	}

	gitleaksParser := app.parsers[scannerName]
	findings, err := gitleaksParser.Parse(reportBytes, taskId)
	if err != nil {
		app.logger.Errorf("Failed parsing file | %s", err)
		ResponseJSON(c, ErrFailedToReadReport)
		return
	}

	app.logger.Infof("Findings - %d", len(findings))

	// Success response
	c.JSON(http.StatusCreated, findings)
	app.logger.Infof("Response sent to User | %s", taskId)
	return
}
