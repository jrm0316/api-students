package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/jrm0316/api-students/schemas"
	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (api *API) getStudents(c echo.Context) error {
	students, err := api.DB.GetStudents()
	if err != nil {
		return c.String(http.StatusNotFound, "Failed to get students")
	}

	return c.JSON(http.StatusOK, students)
}

func (api *API) createStudent(c echo.Context) error {
	studentReq := StudentRequest{}
	if err := c.Bind(&studentReq); err != nil {
		return err
	}

	if err := studentReq.Validate(); err != nil {
		log.Error().Err(err).Msgf("[api] error validating struct")
		return c.String(http.StatusBadRequest, "Error validating student")
	}

	student := schemas.Student{
		Name:   studentReq.Name,
		Email:  studentReq.Email,
		CPF:    studentReq.CPF,
		Age:    studentReq.Age,
		Active: *studentReq.Active,
	}

	if err := api.DB.AddStudent(student); err != nil {
		return c.String(http.StatusInternalServerError, "Error to create student")
	}

	return c.String(http.StatusOK, "Create student")
}

func (api *API) getStudent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get student ID")
	}

	student, err := api.DB.GetStudent(id)
	// não encontrar um student com ess id -> STATUS NOT FOUND (404)
	//ou pode ter algum problema para encontrar o student

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.String(http.StatusNotFound, "Student not found")
	}

	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to get student")
	}

	if err := api.DB.DeleteStudent(student); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to delete student")
	}

	return c.JSON(http.StatusOK, student)
}

//func (api *API) updateStudent(c echo.Context) error {
//	id, err := strconv.Atoi(c.Param("id"))
//	if err != nil {
//		return c.String(http.StatusInternalServerError, "Failed to get student")
//	}

//	receivedStudent := db.Student{}
//	if err := c.Bind(&receivedStudent); err != nil {
//		return err
//	}

//	updatingStudent, err := api.DB.GetStudent(id)
//	if errors.Is(err, gorm.ErrRecordNotFound) {
//		return c.String(http.StatusNotFound, "Student not found")
//	}

//	if err != nil {
//		return c.String(http.StatusInternalServerError, "Failed to get student")
//	}

//student := updateStudentInfo(receivedStudent, updatingStudent)

//	if err := api.DB.UpdateStudent(student); err != nil {
//		return c.String(http.StatusInternalServerError, "Failed")
//	}
//	return c.JSON(http.StatusOK, student)
//}

func (api *API) deleteStudent(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to delete student")
	}

	student, err := api.DB.GetStudent(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.String(http.StatusNotFound, "Student not found")
	}

	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed")
	}

	if err := api.DB.DeleteStudent(student); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to delete student")
	}

	return c.JSON(http.StatusOK, student)
}

//func updateStudentInfo(receivedStudent, student db.Student)
