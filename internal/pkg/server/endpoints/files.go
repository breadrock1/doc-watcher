package endpoints

import (
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
	"path"
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

	return nil
}

const UploadHtmlForm = `
<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <title>Single file upload</title>
</head>
<body>
<h1>Upload single file</h1>

<form action="/upload" method="post" enctype="multipart/form-data">
    Files: <input type="file" name="file_name"><br><br>
    <input type="submit" value="Submit">
</form>
</body>
</html>
`

func UploadFileForm(c echo.Context) error {
	return c.HTML(http.StatusOK, UploadHtmlForm)
}
