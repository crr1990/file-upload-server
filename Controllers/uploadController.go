package Controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"sort"
	"time"
)

var BathPathInfo = "/my/uploads"

type UploadNames struct {
	Name    string
	Ids     string
	FIle    *multipart.FileHeader
	Stuffix string
	Guid    string
}

func TestInsert(c *gin.Context) {
	year := time.Now().Year()
	month := int(time.Now().Month())
	day := time.Now().Day()
	fmt.Println(year, month, day)
	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"message": "success",
	})

}

func Upload(c *gin.Context) {
	// new
	chunkNumber := c.PostForm("chunkNumber")
	totalChunks := c.PostForm("totalChunks")
	//chunkSize := c.PostForm("chunkSize")
	//chunkNumber := c.PostForm("currentChunkSize")
	//chunkNumber := c.PostForm("totalSize")
	identifier := c.PostForm("identifier")
	name := c.PostForm("filename")
	savePath := c.PostForm("savePath")
	pathInfo := BathPathInfo + "/"+ savePath

	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "a bad request")
		return
	}

	filename := file.Filename
	fileStuffix := path.Ext(filename)


	data := UploadNames{name, chunkNumber, file, fileStuffix, identifier}

	if ok, _ := PathExists(pathInfo + "/" + data.Guid); ok {

	} else {
		os.MkdirAll(pathInfo+"/"+data.Guid, 0777)
	}
	c.SaveUploadedFile(file, pathInfo+"/"+data.Guid+"/"+data.Ids+data.Stuffix)

	if totalChunks == chunkNumber {
		DoneMergeFile(identifier, fileStuffix,pathInfo )
	}

	c.JSON(http.StatusOK, gin.H{
		"successStatuses": 200,
		"message": "success",
	})
}

func DoneMergeFile(guid string, suffix string,pathInfo string) {
	if ok, _ := PathExists(pathInfo + "/" + guid); ok {
		var data []string
		data = make([]string, 0)
		fileInfo, _ := ioutil.ReadDir(pathInfo + "/" + guid)
		for _, val := range fileInfo {
			data = append(data, val.Name())
		}

		sort.Strings(data)
		f, _ := os.OpenFile(pathInfo+"/"+guid+suffix, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
		for _, val := range data {
			contents, _ := ioutil.ReadFile(pathInfo + "/" + guid + "/" + val)
			f.Write(contents)
			os.Remove(pathInfo + "/" + guid + "/" + val)
		}
		os.Remove(pathInfo + "/" + guid)
		defer f.Close()
	}
}

/**
*判断文件夹是否存在
 */
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
