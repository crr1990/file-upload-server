package Controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"
	"time"
)

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
	FileName   string
}

type DeleteFileData struct {
	SavePath string
}

func TestInsert(c *gin.Context) {
	var BathPathInfo = viper.GetString("path")
	year := time.Now().Year()
	month := int(time.Now().Month())
	day := time.Now().Day()
	fmt.Println(year, month, day)
	c.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": BathPathInfo,
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

	pathInfo := viper.GetString("path") + "/" + savePath
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
		err := os.MkdirAll(pathInfo+"/"+data.Guid, 0777)
		if err != nil {
			log.Print(err)
			return
		}
	}
	err = c.SaveUploadedFile(file, pathInfo+"/"+data.Guid+"/"+data.Name+data.Stuffix)
	if err != nil {
		log.Print(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    viper.GetString("host") + savePath + "/" + data.Guid + "/" + data.Name + data.Stuffix,
	})
}
func DeleteFile(c *gin.Context) {
	isremove := false //删除文件是否成功
	var d DeleteFileData
	c.BindJSON(&d)

	//删除文件
	cuowu := os.RemoveAll(viper.GetString("path") + "/" + d.SavePath + "/")
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
		err := DoneMergeFile(p.Identifier, p.FileName, viper.GetString("path")+"/"+p.SavePath)
		if err != nil {
			log.Println("MergeFileErr", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": msg,
	})
}

func DoneMergeFile(guid string, fileName string, pathInfo string) error {
	log.Println(pathInfo + "/" + guid)
	log.Println(PathExists(pathInfo + "/" + guid))
	ok, err := PathExists(pathInfo + "/" + guid)
	if err == nil && ok {
		var data []int
		data = make([]int, 0)
		fileInfo, _ := ioutil.ReadDir(pathInfo + "/" + guid)
		for _, val := range fileInfo {
			d, _ := strconv.Atoi(val.Name())
			data = append(data, d)
		}

		sort.Ints(data)
		f, err := os.OpenFile(pathInfo+"/"+fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
		if err != nil {
			log.Println("DoneMergeFileErr", err)
			return err
		}

		for _, val := range data {
			contents, err := ioutil.ReadFile(pathInfo + "/" + guid + "/" + strconv.Itoa(val))
			if err != nil {
				log.Println("DoneMergeFileErr2", err)
				continue
			}
			_, err = f.Write(contents)
			if err != nil {
				log.Println("DoneMergeFileErr3", err)
				continue
			}
			err = os.Remove(pathInfo + "/" + guid + "/" + strconv.Itoa(val))
			if err != nil {
				log.Println("DoneMergeFileErr4", err)
				continue
			}
		}
		err = os.Remove(pathInfo + "/" + guid)
		if err != nil {
			log.Println("DoneMergeFileErr5", err)
			return err
		}

		defer f.Close()
	} else {
		log.Println("DoneMergeFileFail", err, ok)
		return err
	}

	return nil
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
