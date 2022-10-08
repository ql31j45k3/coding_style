package http

import (
	"fmt"
	"layout_2/internal/domain/student"
	"layout_2/internal/libs/response"
	"layout_2/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

type StudentHandlerCond struct {
	dig.In

	R *gin.Engine

	StudentUseCase student.StudentUseCase
}

func RegisterRouter(cond StudentHandlerCond) {
	router := studentRouter{
		studentUseCase: cond.StudentUseCase,
	}

	routerGroup := cond.R.Group("/v1/student")
	routerGroup.POST("", router.create)
	routerGroup.PUT("/:id", router.updateID)
	routerGroup.GET("/:id", router.getID)
	routerGroup.GET("", router.get)
}

type studentRouter struct {
	studentUseCase student.StudentUseCase
}

type responseStudentCreate struct {
	ID string `json:"id"`
}

func (sr *studentRouter) create(c *gin.Context) {
	var student student.Student
	if err := response.BindJSON(c, &student); err != nil {
		response.NewReturnError(c, http.StatusBadRequest, response.HttpStatusBadRequest, err)
		return
	}

	rowID, err := sr.studentUseCase.Create(student)
	if err != nil {
		response.IsErrRecordNotFound(c, err)
		return
	}

	result := responseStudentCreate{
		ID: fmt.Sprint(rowID),
	}

	c.JSON(http.StatusCreated, response.NewResponseSuccess(result))
}

func (sr *studentRouter) updateID(c *gin.Context) {
	cond := student.StudentCond{}
	if err := cond.ParseID(c); err != nil {
		response.NewReturnError(c, http.StatusBadRequest, response.HttpStatusBadRequest, err)
		return
	}

	var student student.Student
	if err := response.BindJSON(c, &student); err != nil {
		response.NewReturnError(c, http.StatusBadRequest, response.HttpStatusBadRequest, err)
		return
	}

	err := sr.studentUseCase.UpdateID(cond, student)
	if err != nil {
		response.IsErrRecordNotFound(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

type responseStudent struct {
	Name string `json:"name"`

	Gender int `json:"gender"`
	Status int `json:"status"`
}

func (sr *studentRouter) getID(c *gin.Context) {
	cond := student.StudentCond{}
	if err := cond.ParseID(c); err != nil {
		response.NewReturnError(c, http.StatusBadRequest, response.HttpStatusBadRequest, err)
		return
	}

	var responseStudent responseStudent

	student, err := sr.studentUseCase.GetID(cond)
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

type responseStudentGet struct {
	Page response.Pagination `json:"page"`

	Students []responseStudent `json:"students"`
}

func (sr *studentRouter) get(c *gin.Context) {
	cond := student.StudentCond{}
	if err := cond.ParseGet(c); err != nil {
		response.NewReturnError(c, http.StatusBadRequest, response.HttpStatusBadRequest, err)
		return
	}

	var responseStudents []responseStudent

	students, count, err := sr.studentUseCase.Get(cond)
	if err != nil {
		response.IsErrRecordNotFound(c, err)
		return
	}

	if err := utils.ConvSourceToData(&students, &responseStudents); err != nil {
		response.NewReturnError(c, http.StatusBadRequest, response.HttpStatusBadRequest, err)
		return
	}

	result := responseStudentGet{
		Page:     cond.Pagination,
		Students: responseStudents,
	}

	result.Page.Count = count

	c.JSON(http.StatusOK, response.NewResponseSuccess(result))
}
