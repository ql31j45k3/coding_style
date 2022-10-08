package response

import (
	"errors"
	"layout_2/internal/utils"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BindJSON 取得 Request body 資料，JSON 資料轉行為 struct
func BindJSON(c *gin.Context, obj interface{}) error {
	if err := c.BindJSON(obj); err != nil {
		var messages string
		if _, ok := err.(validator.ValidationErrors); ok {
			for _, err2 := range err.(validator.ValidationErrors) {
				messages += err2.Translate(translator) + ", "
			}

			// 過濾最後字尾 ", "
			messages = messages[:len(messages)-2]

		} else {
			messages = err.Error()
		}

		return errors.New(messages)
	}

	return nil
}

// NewReturnError 設定回傳錯誤，統一回傳錯誤格式
func NewReturnError(c *gin.Context, code int, httpStatus HttpStatus, err error) {
	messages := err.Error()
	c.JSON(code, NewResponseError(httpStatus, messages))
}

func NewErrorByBasic(c *gin.Context, code int, response ResponseBasic) {
	c.JSON(code, response)
}

// IsErrRecordNotFound 驗證 SQL 語法執行但查無資料情況，調整 http status
func IsErrRecordNotFound(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		NewReturnError(c, http.StatusNotFound, HttpStatusNotFound, err)
	} else {
		NewReturnError(c, http.StatusInternalServerError, HttpStatusInternalServerError, err)
	}
}

type ResponseBasic struct {
	Code     HttpStatus  `json:"code"`
	Messages string      `json:"messages"`
	Data     interface{} `json:"data"`
}

func NewResponseSuccess(data interface{}) ResponseBasic {
	return ResponseBasic{
		Code:     HttpStatusOk,
		Messages: "success",
		Data:     data,
	}
}

func NewResponseError(httpStatus HttpStatus, messages string) ResponseBasic {
	return ResponseBasic{
		Code:     httpStatus,
		Messages: messages,
		Data:     []string{},
	}
}

// Pagination 查詢分頁欄位
type Pagination struct {
	PageIndex int `json:"page_index"`
	PageSize  int `json:"page_size"`

	Count int64 `json:"count"`
}

func (p *Pagination) GetOffset() int {
	if p.PageIndex == utils.DefaultNotAssignInt || p.PageIndex == 0 {
		return 0
	}

	return (p.PageIndex - 1) * p.PageSize
}

func (p *Pagination) GetRowCount() int {
	if p.PageSize == utils.DefaultNotAssignInt || p.PageSize == 0 {
		return 25
	}

	return p.PageSize
}
