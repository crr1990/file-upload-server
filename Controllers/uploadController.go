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
	"log"
)

var BathPathInfo = "D:/CUSTOMER"

type UploadNames struct {
	Name    string
	Ids     string
	FIle    *multipart.FileHeader
	Stuffix string
	Guid    string
}

type s_data struct {
	Identifier string
	Suffix     string
	SavePath   string
}

type DeleteFileData struct {
	SavePath   string
}

func TestInsert(c *gin.Context) {
	year := time.Now().Year()
	month := int(time.Now().Month())
	day := time.Now().Day()
	fmt.Println(year, month, day)
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "success",
	})

}

func Upload(c *gin.Context) {
	chunkNumber := c.PostForm("chunkNumber")
	identifier := c.PostForm("identifier")
	name := c.PostForm("filename")
	savePath := c.PostForm("savePath")
	if savePath == "" {
		c.String(http.StatusBadRequest, "参数不能为空")
		return
	}

	pathInfo := BathPathInfo + "/"+savePath
	log.Println(pathInfo)
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "a bad request")
		return
	}

	filename := file.Filename
	fileStuffix := path.Ext(filename)

	data := UploadNames{name, chunkNumber, file, fileStuffix, identifier}
	log.Println("pathInfo:" + pathInfo)
	if ok, _ := PathExists(pathInfo + "/" + data.Guid); ok {

	} else {
		os.MkdirAll(pathInfo+"/"+data.Guid, 0777)
	}
	c.SaveUploadedFile(file, pathInfo+"/"+data.Guid+"/"+data.Ids+data.Stuffix)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
	})
}
func DeleteFile(c *gin.Context) {
	isremove := false //删除文件是否成功
	var d  DeleteFileData
	c.BindJSON(&d)


	//删除文件
	cuowu := os.RemoveAll(BathPathInfo + "/" + d.SavePath + "/")
	if cuowu != nil {
		//如果删除失败则输出 file remove Error!
		fmt.Println("file remove Error!")
		//输出错误详细信息
		fmt.Printf("%s", cuowu)
	} else {
		//如果删除成功则输出 file remove OK!
		fmt.Print("file remove OK!")
		isremove = true
	}
	//返回结果
	c.JSON(http.StatusOK, gin.H{
		"success": isremove,
		"code":    0,
	})
}

func MergeFile(c *gin.Context) {
	var p s_data
	c.BindJSON(&p)

	var code int
	var msg string

	if p.Identifier == "" {
		code = 1001
		msg = "Identifier is null."
	} else if p.Suffix == "" {
		code = 1001
		msg = "Suffix is null."
	} else if p.SavePath == "" {
		code = 1001
		msg = "SavePath is null."
	} else {
		code = 0
		DoneMergeFile(p.Identifier, p.Suffix, BathPathInfo+"/"+p.SavePath)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"message":         msg,
	})
}

func DoneMergeFile(guid string, suffix string, pathInfo string) {
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
