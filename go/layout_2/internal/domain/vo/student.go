package vo

import (
	"layout_2/internal/libs/response"
	"layout_2/internal/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type StudentCond struct {
	response.Pagination

	ID uint

	Name string

	Gender int
	Status int
}

func (cond *StudentCond) ParseID(c *gin.Context) error {
	var err error
	idStr := c.Param("id")

	cond.ID, err = cond.getID(idStr)
	if err != nil {
		return err
	}

	return nil
}

func (cond *StudentCond) getID(idStr string) (uint, error) {
	if utils.IsEmpty(idStr) {
		return 0, nil
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

func (cond *StudentCond) ParseGet(c *gin.Context) error {
	var err error
	idStr := c.Query("id")

	cond.ID, err = cond.getID(idStr)
	if err != nil {
		return err
	}

	cond.PageIndex, err = utils.AtoiAndDefaultNotAssignInt(c.Query("page_index"))
	if err != nil {
		return err
	}

	cond.PageSize, err = utils.AtoiAndDefaultNotAssignInt(c.Query("page_size"))
	if err != nil {
		return err
	}

	cond.Name = strings.TrimSpace(c.Query("name"))

	cond.Gender, err = utils.AtoiAndDefaultNotAssignInt(c.Query("gender"))
	if err != nil {
		return err
	}

	cond.Status, err = utils.AtoiAndDefaultNotAssignInt(c.Query("status"))
	if err != nil {
		return err
	}

	return nil
}
