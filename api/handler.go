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

//getStudents godoc
//
// @Summary 	Get a list of students
// @Description Retrieve students details
// @Tags 		students
// @Accept 		json
// @Produce		json
// @Param 		register path 	int	false	"Registration"
// @Success 	200 {object} schemas.StudentResponse
// @Failure 	404
// Router 		/students [get]

func (api *API) getStudents(c echo.Context) error {
	students, err := api.DB.GetStudents()
	if err != nil {
		return c.String(http.StatusNotFound, "Failed to get students")
	}

	active := c.QueryParam("active")

	if active != "" {
		act, err := strconv.ParseBool(active)
		if err != nil {
			log.Error().Err(err).Msgf("[api] error to parse boolean")
			return c.String(http.StatusInternalServerError, "Failed to parse boolean")
		}
		students, err = api.DB.GetFilteredStudent(act)
	}

	listOfStudents := map[string][]schemas.StudentResponse{"students": schemas.NewResponse(students)}

	return c.JSON(http.StatusOK, listOfStudents)
}

// createStudent godoc
//
// @Summary 	Create student
// @Description Create student
// @Tags 		students
// @Accept 		json
// @Produce 	json
// @Success 	200 {object} schemas.StudentResponse
// @Failure 	400
// @Router 		/students [post]
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

	return c.JSON(http.StatusOK, student)
}

//getStudents godoc
//
// @Summary 	Get a list of students
// @Description Retrieve students details using ID
// @Tags 		students
// @Accept 		json
// @Produce		json
// @Success 	200 {object} schemas.StudentResponse
// @Failure 	404
// @Failure		500
// Router 		/students/{id} [get]

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

// updateStudent godoc
//
// @Summary		Update Student
// @Description	Update student details
// @Tagas 		students
// @Accept		json
// @Produce 	json
// @Success 	200 {object} schemas.StudentResponse
// @Failure 	404
// @Failure 	500
// @Router 		/students/{id}	[put]

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

// deleteStudent godoc
//
// @Summary		Delete Student
// @Description	Delete student details
// @Tags		students
// @Accept		json
// @Produce 	json
// @Success		200 {object schemas.StudentResponse}
// @Failure		404
// @Failure		500
// @Router		/students/{id}	[delete]

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
