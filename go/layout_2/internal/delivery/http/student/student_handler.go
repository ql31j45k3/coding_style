package student

import (
	"fmt"
	"layout_2/internal/domain/entity"
	"layout_2/internal/domain/usecase"
	"layout_2/internal/domain/vo"
	"layout_2/internal/libs/response"
	"layout_2/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

type StudentHandlerCond struct {
	dig.In

	R *gin.Engine

	StudentUseCase usecase.StudentUseCase
}

func RegisterRouter(cond StudentHandlerCond) {
	router := studentRouter{
		StudentHandlerCond: cond,
	}

	routerGroup := cond.R.Group("/v2/student")
	routerGroup.POST("", router.create)
	routerGroup.PUT("/:id", router.updateID)
	routerGroup.GET("/:id", router.getID)
	routerGroup.GET("", router.get)
}

type studentRouter struct {
	StudentHandlerCond
}

func (sr *studentRouter) create(c *gin.Context) {
	var student entity.Student
	if err := response.ShouldBindJSON(c, &student); err != nil {
		response.NewReturnError(c, http.StatusBadRequest, response.HttpStatusBadRequest, err)
		return
	}

	rowID, err := sr.StudentUseCase.Create(student)
	if err != nil {
		if s, ok := response.FromError(err); ok {
			response.NewReturnError(c, http.StatusBadRequest, s.Code(), err)
			return
		}

		response.IsErrRecordNotFound(c, err)
		return
	}

	result := vo.ResponseStudentCreate{
		ID: fmt.Sprint(rowID),
	}

	c.JSON(http.StatusCreated, response.NewResponseSuccess(result))
}

func (sr *studentRouter) updateID(c *gin.Context) {
	cond := vo.StudentCond{}
	if err := cond.ParseID(c); err != nil {
		response.NewReturnError(c, http.StatusBadRequest, response.HttpStatusBadRequest, err)
		return
	}

	var student entity.Student
	if err := response.ShouldBindJSON(c, &student); err != nil {
		response.NewReturnError(c, http.StatusBadRequest, response.HttpStatusBadRequest, err)
		return
	}

	err := sr.StudentUseCase.UpdateID(cond, student)
	if err != nil {
		response.IsErrRecordNotFound(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (sr *studentRouter) getID(c *gin.Context) {
	cond := vo.StudentCond{}
	if err := cond.ParseID(c); err != nil {
		response.NewReturnError(c, http.StatusBadRequest, response.HttpStatusBadRequest, err)
		return
	}

	var responseStudent vo.ResponseStudent

	student, err := sr.StudentUseCase.GetID(cond)
	if err != nil {
		response.IsErrRecordNotFound(c, err)
		return
	}

	if err := utils.ConvSourceToData(&student, &responseStudent); err != nil {
		response.NewReturnError(c, http.StatusBadRequest, response.HttpStatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, response.NewResponseSuccess(responseStudent))
}

func (sr *studentRouter) get(c *gin.Context) {
	cond := vo.StudentCond{}
	if err := cond.ParseGet(c); err != nil {
		response.NewReturnError(c, http.StatusBadRequest, response.HttpStatusBadRequest, err)
		return
	}

	var responseStudents []vo.ResponseStudent

	students, count, err := sr.StudentUseCase.Get(cond)
	if err != nil {
		response.IsErrRecordNotFound(c, err)
		return
	}

	if err := utils.ConvSourceToData(&students, &responseStudents); err != nil {
		response.NewReturnError(c, http.StatusBadRequest, response.HttpStatusBadRequest, err)
		return
	}

	result := vo.ResponseStudentGet{
		Page:     cond.Pagination,
		Students: responseStudents,
	}

	result.Page.Count = count

	c.JSON(http.StatusOK, response.NewResponseSuccess(result))
}
