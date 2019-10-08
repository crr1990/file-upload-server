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

var PathInfo = "/my/uploads"

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
	name := c.PostForm("name")
	index := c.PostForm("chunk")
	guid := c.PostForm("guid")
	isLast := c.PostForm("is_last")

	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "a bad request")
		return
	}

	filename := file.Filename
	fileStuffix := path.Ext(filename)
	if index == "" {
		index = "0"
	}

	data := UploadNames{name, index, file, fileStuffix, guid}

	if ok, _ := PathExists(PathInfo + "/" + data.Guid); ok {

	} else {
		os.MkdirAll(PathInfo+"/"+data.Guid, 0777)
	}
	c.SaveUploadedFile(file, PathInfo+"/"+data.Guid+"/"+data.Ids+data.Stuffix)

	if isLast == "1" {
		DoneMergeFile(guid, fileStuffix)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"message": "success",
	})
}

func DoneMergeFile(guid string, suffix string) {
	if ok, _ := PathExists(PathInfo + "/" + guid); ok {
		var data []string
		data = make([]string, 0)
		fileInfo, _ := ioutil.ReadDir(PathInfo + "/" + guid)
		for _, val := range fileInfo {
			data = append(data, val.Name())
		}

		sort.Strings(data)
		f, _ := os.OpenFile(PathInfo+"/"+guid+suffix, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
		for _, val := range data {
			contents, _ := ioutil.ReadFile(PathInfo + "/" + guid + "/" + val)
			f.Write(contents)
			os.Remove(PathInfo + "/" + guid + "/" + val)
		}
		os.Remove(PathInfo + "/" + guid)
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
