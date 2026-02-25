package handlers

import (
	"net/http"
	"strings"

	"example.com/student-api/models"
	"github.com/gin-gonic/gin"

	"example.com/student-api/services"
)

type StudentHandler struct {
	Service *services.StudentService
}

func (h *StudentHandler) GetStudents(c *gin.Context) {
	students, err := h.Service.GetStudents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if students == nil {
		students = []models.Student{}
	}

	c.JSON(http.StatusOK, students)
}

func (h *StudentHandler) GetStudentByID(c *gin.Context) {
	id := c.Param("id")
	student, err := h.Service.GetStudentByID(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, student)
}

func (h *StudentHandler) CreateStudent(c *gin.Context) {
	var student models.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	if validationErrors := validateStudentInput(student); len(validationErrors) > 0 {
		errMsg := strings.Join(validationErrors, ", ")
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	if err := h.Service.CreateStudent(student); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create student"})
		return
	}

	c.JSON(http.StatusCreated, student)
}

func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	id := c.Param("id")
	var student models.Student

	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	student.Id = id

	if validationErrors := validateStudentInput(student); len(validationErrors) > 0 {
		errMsg := strings.Join(validationErrors, ", ")
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	err := h.Service.UpdateStudent(id, student)
	if err != nil {
		if err.Error() == "student not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student"})
		return
	}

	c.JSON(http.StatusOK, student)
}

func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	id := c.Param("id")

	err := h.Service.DeleteStudent(id)
	if err != nil {
		if err.Error() == "student not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete student"})
		return
	}

	c.Status(http.StatusNoContent)
}

func validateStudentInput(student models.Student) []string {
	var errs []string

	if student.Id == "" {
		errs = append(errs, "id must not be empty")
	}
	if student.Name == "" {
		errs = append(errs, "name must not be empty")
	}
	if student.GPA < 0.00 || student.GPA > 4.00 {
		errs = append(errs, "gpa must be between 0.00 and 4.00")
	}

	return errs
}
