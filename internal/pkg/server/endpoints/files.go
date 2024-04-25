package endpoints

import (
	watcher2 "doc-notifier/internal/pkg/watcher"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

func DownloadFile(c echo.Context) error {
	filePath := c.FormValue("file_path")
	return c.File(filePath)
}

func UploadFile(c echo.Context) error {
	fileHandle, err := c.FormFile("file_name")
	if err != nil {
		return err
	}

	srcStream, err := fileHandle.Open()
	if err != nil {
		return err
	}
	defer func() { _ = srcStream.Close() }()

	filePath := path.Join("./upload/", fileHandle.Filename)
	dstStream, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() { _ = dstStream.Close() }()

	if _, err = io.Copy(dstStream, srcStream); err != nil {
		return err
	}

	time.Sleep(5 * time.Second)
	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	for key, value := range watcher.Ocr.Ocr.GetProcessingJobs() {
		if value.Document.DocumentName == fileHandle.Filename {
			return c.JSON(200, returnStatusResponse(200, key))
		}
	}

	return c.JSON(403, returnStatusResponse(403, "Failed while uploading"))
}

const UploadHTMLForm = `
<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Single file upload</title>
</head>
<body>
<h1>Upload single file</h1>

<form action="/file/upload" method="post" enctype="multipart/form-data">
    Files: <input type="file" name="file_name"><br><br>
    <input type="submit" value="Submit">
</form>
</body>
</html>
`

func UploadFileForm(c echo.Context) error {
	return c.HTML(http.StatusOK, UploadHTMLForm)
}

type MoveFileForm struct {
	FilePath   string `json:"file_path"`
	TargetPath string `json:"target_path"`
}

func MoveFile(c echo.Context) error {
	jsonForm := &MoveFileForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		return c.JSON(403, returnStatusResponse(403, err.Error()))
	}

	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	if err := watcher.Reader.MoveFileTo(jsonForm.FilePath, jsonForm.TargetPath); err != nil {
		return c.JSON(403, returnStatusResponse(403, err.Error()))
	}

	return c.JSON(200, returnStatusResponse(200, "Ok"))
}

type RemoveFileForm struct {
	FilePath string `json:"file_path"`
}

func RemoveFile(c echo.Context) error {
	jsonForm := &RemoveFileForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		return c.JSON(403, returnStatusResponse(403, err.Error()))
	}

	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	documents := watcher.Reader.ParseCaughtFiles(jsonForm.FilePath)
	res, err := os.ReadFile(jsonForm.FilePath)
	if err != nil {
		return c.JSON(403, returnStatusResponse(403, err.Error()))
	}
	watcher.Reader.ComputeMd5HashByData(documents[0], res)

	targetURL := fmt.Sprintf("%s/document/%s/%s", watcher.Searcher.Address, documents[0].BucketUUID, documents[0].DocumentMD5)
	response, err := http.Get(targetURL)
	if err != nil || response.StatusCode != 200 {
		return c.JSON(403, returnStatusResponse(403, err.Error()))
	}

	if err = os.RemoveAll(jsonForm.FilePath); err != nil {
		return c.JSON(403, returnStatusResponse(403, err.Error()))
	}

	return c.JSON(200, returnStatusResponse(200, "Ok"))
}

type ParsingJobForm struct {
	JobId string `json:"job_id"`
}

func GetUploadingStatus(c echo.Context) error {
	jsonForm := &ParsingJobForm{}
	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(jsonForm); err != nil {
		return c.JSON(403, returnStatusResponse(403, err.Error()))
	}

	watcher := c.Get("Watcher").(*watcher2.NotifyWatcher)
	jobRes := watcher.Ocr.Ocr.GetProcessingJob(jsonForm.JobId)
	return c.JSON(200, jobRes)
}
