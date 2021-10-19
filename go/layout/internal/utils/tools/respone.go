package tools

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// BindJSON 取得 Request body 資料，JSON 資料轉行為 struct
func BindJSON(c *gin.Context, trans ut.Translator, obj interface{}) error {
	if err := c.BindJSON(obj); err != nil {
		var errs []string
		if _, ok := err.(validator.ValidationErrors); ok {
			for _, err2 := range err.(validator.ValidationErrors) {
				errs = append(errs, err2.Translate(trans))
			}
		} else {
			errs = append(errs, err.Error())
		}

		c.JSON(http.StatusBadRequest, NewResponseBasicError(errs))
		return err
	}

	return nil
}

// NewReturnError 設定回傳錯誤，統一回傳錯誤格式
func NewReturnError(c *gin.Context, code int, err error) {
	messages := []string{err.Error()}
	c.JSON(code, NewResponseBasicError(messages))
}

// IsErrRecordNotFound 驗證 SQL 語法執行但查無資料情況，調整 http status
func IsErrRecordNotFound(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		NewReturnError(c, http.StatusNotFound, err)
	} else {
		NewReturnError(c, http.StatusInternalServerError, err)
	}
}

type ResponseBasic struct {
	_ struct{}

	Code     int         `json:"code"`
	Messages []string    `json:"messages"`
	Data     interface{} `json:"data"`
}

func NewResponseBasicSuccess(data interface{}) ResponseBasic {
	return ResponseBasic{
		Code:     1,
		Messages: []string{"success"},
		Data:     data,
	}
}

func NewResponseBasicError(messages []string) ResponseBasic {
	return ResponseBasic{
		Code:     0,
		Messages: messages,
	}
}

// Model 對外回傳基礎欄位
type Model struct {
	_ struct{}

	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Pagination 查詢分頁欄位
type Pagination struct {
	_ struct{}

	PageIndex int `json:"page_index"`
	PageSize  int `json:"page_size"`
}

func (p *Pagination) GetOffset() int {
	if p.PageIndex == DefaultNotAssignInt {
		return 0
	}

	return (p.PageIndex - 1) * p.PageSize
}

func (p *Pagination) GetRowCount() int {
	if p.PageSize == DefaultNotAssignInt || p.PageSize == 0 {
		return 25
	}

	return p.PageSize
}
