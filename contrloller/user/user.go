package user

import (
	"github.com/revianto/goTemplate/conn"
	model "github.com/revianto/goTemplate/model"
	mystruct "github.com/revianto/goTemplate/struct/user"

	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

var (
	errNotExist        = errors.New("Data are not exist")
	errInvalidID       = errors.New("Invalid ID")
	errInvalidParam    = errors.New("Invalid parameter")
	errInvalidInt      = errors.New("Invalid integer")
	errInvalidBody     = errors.New("Invalid request body")
	errGetFailed       = errors.New("Error at getting data")
	errInsertionFailed = errors.New("Error in the data insertion")
	errUpdationFailed  = errors.New("Error in the data updation")
)

const dataCollection = "user"

var db = conn.GetDB()

// FilterParamenter to filter param get
type FilterParamenter struct {
	ID     bson.ObjectId `form:"id"`
	Page   string        `form:"page"`
	SortBy string        `form:"sortBy"`
}

// GetData to get data
func GetData(c *gin.Context) {

	structData := make([]mystruct.GetStruct, 0)
	quary := bson.M{}
	var ft FilterParamenter

	err := c.ShouldBind(&ft)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": model.ErrorHandle(errInvalidParam.Error(), err.Error())})
		return
	}

	/*
		- page -> pagination default 1
		- sortBY -> sort data berdasarkan filed (- untuk desc) ex sortBy : -name
		- selectField -> memilih field yang akan di tampilkan sesuai denga struct
	*/
	page := ft.Page
	sortBy := ft.SortBy
	selectField := bson.M{}
	// selectField1 := bson.M{"nama": 1}

	rowCount, data, err := model.GetData(sortBy, page, selectField, quary, structData, dataCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": model.ErrorHandle(errGetFailed.Error(), err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"offset": 10, "page": 1, "row_count": rowCount, "data": &data})
}

// AddData to add data
func AddData(c *gin.Context) {
	structData := mystruct.PostStruct{}
	err := c.BindJSON(&structData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": model.ErrorHandle(errInvalidBody.Error(), err.Error())})
		return
	}
	if structData.IsAktif == 0 {
		structData.IsAktif = 1
	}
	data, err := model.InsertData(structData, dataCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": model.ErrorHandle(errInsertionFailed.Error(), err.Error())})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": &data})

}

// UpdateData to update data
func UpdateData(c *gin.Context) {
	structData := mystruct.PutStruct{}
	err := c.BindJSON(&structData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": model.ErrorHandle(errInvalidBody.Error(), err.Error())})
		return
	}

	selectData := mystruct.IDStruct{}
	selectData.ID = bson.ObjectIdHex(c.Query("id"))

	data, err := model.UpdateData(selectData, structData, dataCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": model.ErrorHandle(errUpdationFailed.Error(), err.Error())})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": &data})
}
