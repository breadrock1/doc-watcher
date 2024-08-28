package server

import (
	"bytes"
	"encoding/json"
	"log"
	"mime/multipart"

	"github.com/labstack/echo/v4"
)

func (s *Service) CreateStorageGroup() error {
	group := s.server.Group("/storage")

	group.GET("/buckets", s.GetBuckets)
	group.PUT("/bucket", s.CreateBucket)
	group.DELETE("/:bucket", s.RemoveBucket)

	group.POST("/:bucket/file/copy", s.CopyFile)
	group.POST("/:bucket/file/move", s.MoveFile)
	group.POST("/:bucket/file/upload", s.UploadFile)
	group.POST("/:bucket/file/download", s.DownloadFile)
	group.POST("/:bucket/file/remove", s.RemoveFile)

	group.POST("/:bucket/files", s.GetListFiles)

	return nil
}

// GetBuckets
// @Summary Get watched bucket list
// @Description Get watched bucket list
// @ID get-buckets
// @Tags storage
// @Produce  json
// @Success 200 {array} string "Ok"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /storage/buckets [get]
func (s *Service) GetBuckets(c echo.Context) error {
	watcherDirs := s.watcher.Watcher.GetBuckets()
	return c.JSON(200, watcherDirs)
}

// CreateBucket
// @Summary Create new bucket into storage
// @Description Create new bucket into storage
// @ID create-bucket
// @Tags storage
// @Accept  json
// @Produce json
// @Param bucket path string true "Bucket name to create"
// @Param jsonQuery body BucketNameForm true "Bucket name to create"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /storage/bucket [put]
func (s *Service) CreateBucket(c echo.Context) error {
	jsonForm := &BucketNameForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		respErr := createStatusResponse(400, err.Error())
		return c.JSON(400, respErr)
	}

	if err := s.watcher.Watcher.CreateBucket(jsonForm.BucketName); err != nil {
		respErr := createStatusResponse(208, err.Error())
		return c.JSON(208, respErr)
	}

	return c.JSON(200, createStatusResponse(200, "Ok"))
}

// RemoveBucket
// @Summary Remove bucket from storage
// @Description Remove bucket from storage
// @ID remove-bucket
// @Tags storage
// @Produce  json
// @Param bucket path string true "Bucket name to remove"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /storage/{bucket} [delete]
func (s *Service) RemoveBucket(c echo.Context) error {
	bucketName := c.Param("bucket")
	if err := s.watcher.Watcher.RemoveBucket(bucketName); err != nil {
		respErr := createStatusResponse(400, err.Error())
		return c.JSON(400, respErr)
	}

	return c.JSON(200, createStatusResponse(200, "Ok"))
}

// CopyFile
// @Summary Copy file to another location into bucket
// @Description Copy file to another location into bucket
// @ID copy-file
// @Tags storage
// @Accept  json
// @Produce json
// @Param bucket path string true "Bucket name of src file"
// @Param jsonQuery body CopyFileForm true "Params to copy file"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /storage/{bucket}/file/copy [post]
func (s *Service) CopyFile(c echo.Context) error {
	bucketName := c.Param("bucket")

	jsonForm := &CopyFileForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		respErr := createStatusResponse(400, err.Error())
		return c.JSON(400, respErr)
	}

	copyErr := s.watcher.Watcher.CopyFile(bucketName, jsonForm.SrcPath, jsonForm.DstPath)
	if copyErr != nil {
		respErr := createStatusResponse(400, copyErr.Error())
		return c.JSON(400, respErr)
	}

	return c.JSON(200, createStatusResponse(200, "Ok"))
}

// MoveFile
// @Summary Move file to another location into bucket
// @Description Move file to another location into bucket
// @ID move-file
// @Tags storage
// @Accept  json
// @Produce json
// @Param bucket path string true "Bucket name of src file"
// @Param jsonQuery body CopyFileForm true "Params to move file"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /storage/{bucket}/file/move [post]
func (s *Service) MoveFile(c echo.Context) error {
	bucketName := c.Param("bucket")

	jsonForm := &CopyFileForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		respErr := createStatusResponse(400, err.Error())
		return c.JSON(400, respErr)
	}

	moveErr := s.watcher.Watcher.CopyFile(bucketName, jsonForm.SrcPath, jsonForm.DstPath)
	if moveErr != nil {
		respErr := createStatusResponse(400, moveErr.Error())
		return c.JSON(400, respErr)
	}

	return c.JSON(200, createStatusResponse(200, "Ok"))
}

// UploadFile
// @Summary Upload files to storage
// @Description Upload files to storage
// @ID upload-files
// @Tags storage
// @Accept  multipart/form
// @Produce  json
// @Param bucket path string true "Bucket name to upload files"
// @Param files formData file true "Files multipart form"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /storage/{bucket}/file/upload [post]
func (s *Service) UploadFile(c echo.Context) error {
	var uploadErr error
	var fileData bytes.Buffer
	var multipartForm *multipart.Form

	if multipartForm, uploadErr = c.MultipartForm(); uploadErr != nil {
		respErr := createStatusResponse(400, uploadErr.Error())
		return c.JSON(400, respErr)
	}

	bucketName := c.Param("bucket")
	if multipartForm.File["files"] == nil {
		return c.JSON(400, createStatusResponse(400, "Empty File body"))
	}

	for _, fileForm := range multipartForm.File["files"] {
		fileName := fileForm.Filename
		fileHandler, uploadErr := fileForm.Open()
		if uploadErr != nil {
			log.Println(uploadErr)
			continue
		}
		defer func(fileHandler multipart.File) {
			err := fileHandler.Close()
			if err != nil {
				log.Println("failed to close file handler: ", err)
				return
			}
		}(fileHandler)

		_, uploadErr = fileData.ReadFrom(fileHandler)
		if uploadErr != nil {
			log.Println(uploadErr)
			continue
		}
		defer fileData.Reset()

		uploadErr = s.watcher.Watcher.UploadFile(bucketName, fileName, fileData)
		if uploadErr != nil {
			log.Println(uploadErr)
			continue
		}
	}

	return c.JSON(200, createStatusResponse(200, "Ok"))
}

// DownloadFile
// @Summary Download file from storage
// @Description Download file from storage
// @ID download-file
// @Tags storage
// @Accept  json
// @Produce json
// @Param bucket path string true "Bucket name to download file"
// @Param jsonQuery body DownloadFile true "Parameters to download file"
// @Success 200 {file} io.Writer "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /storage/{bucket}/file/download [post]
func (s *Service) DownloadFile(c echo.Context) error {
	bucketName := c.Param("bucket")

	jsonForm := &DownloadFile{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		respErr := createStatusResponse(400, err.Error())
		return c.JSON(400, respErr)
	}

	fileData, err := s.watcher.Watcher.DownloadFile(bucketName, jsonForm.FileName)
	if err != nil {
		respErr := createStatusResponse(400, err.Error())
		return c.JSON(400, respErr)
	}
	defer fileData.Reset()

	return c.Blob(200, echo.MIMEMultipartForm, fileData.Bytes())
}

// RemoveFile
// @Summary Remove file from storage
// @Description Remove file from storage
// @ID remove-file
// @Tags storage
// @Produce  json
// @Param bucket path string true "Bucket name to remove file"
// @Param jsonQuery body RemoveFile true "Parameters to remove file"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /storage/{bucket}/file/remove [post]
func (s *Service) RemoveFile(c echo.Context) error {
	bucketName := c.Param("bucket")

	jsonForm := &RemoveFile{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		respErr := createStatusResponse(400, err.Error())
		return c.JSON(400, respErr)
	}

	return s.watcher.Watcher.RemoveFile(bucketName, jsonForm.FileName)
}

// GetListFiles
// @Summary Get files list into bucket
// @Description Get files list into bucket
// @ID get-list-files
// @Tags storage
// @Accept  json
// @Produce json
// @Param bucket path string true "Bucket name to get list files"
// @Param jsonQuery body ListFilesForm true "Parameters to get list files"
// @Success 200 {object} ResponseForm "Ok"
// @Failure	400 {object} BadRequestForm "Bad Request message"
// @Failure	503 {object} ServerErrorForm "Server does not available"
// @Router /storage/{bucket}/files [post]
func (s *Service) GetListFiles(c echo.Context) error {
	bucketName := c.Param("bucket")

	jsonForm := &ListFilesForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		respErr := createStatusResponse(400, err.Error())
		return c.JSON(400, respErr)
	}

	listObjects := s.watcher.Watcher.GetListFiles(bucketName, jsonForm.DirectoryName)
	return c.JSON(200, listObjects)
}
