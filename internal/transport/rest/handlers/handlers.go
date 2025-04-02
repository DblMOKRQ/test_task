package handlers

import (
	"net/http"
	"strconv"

	_ "github.com/DblMOKRQ/test_task/docs"
	"github.com/DblMOKRQ/test_task/internal/entity"
	"github.com/DblMOKRQ/test_task/internal/service"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	service *service.Service
}

func NewHandlers(service *service.Service) *Handlers {
	return &Handlers{service: service}
}

// AddUser добавляет нового пользователя
// @Summary Добавить пользователя
// @Description Создает новую запись пользователя в системе
// @Tags Users
// @Accept json
// @Produce json
// @Param user body entity.User true "Данные пользователя"
// @Success 200 {object} map[string]interface{} "Успешное добавление"
// @Failure 400 {object} map[string]interface{} "Неверные данные"
// @Failure 500 {object} map[string]interface{} "Ошибка сервера"
// @Router /add [post]
func (h *Handlers) AddUser(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if user.Name == "" || user.Surname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and Surname are required"})
		return
	}
	id, err := h.service.AddUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User added successfully", "id": id})
}

// GetUsersByNationality возвращает пользователей по национальности
// @Summary Получить по национальности
// @Description Возвращает список пользователей с указанной национальностью
// @Tags Users
// @Produce json
// @Param nationality path string true "Код национальности (2 заглавные буквы)"
// @Success 200 {array} entity.FullUser
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /get/nationality/{nationality} [get]
func (h *Handlers) GetUsersByNationality(c *gin.Context) {
	nationality := c.Param("nationality")
	users, err := h.service.GetUsersByNationality(nationality)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUsersByGender возвращает пользователей по полу
// @Summary Получить по полу
// @Description Возвращает список пользователей с указанным полом
// @Tags Users
// @Produce json
// @Param gender path string true "Пол (male/female)"
// @Success 200 {array} entity.FullUser
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /get/gender/{gender} [get]
func (h *Handlers) GetUsersByGender(c *gin.Context) {
	gender := c.Param("gender")
	users, err := h.service.GetUsersByGender(gender)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUsersByAge возвращает пользователей по возрасту
// @Summary Получить по возрасту
// @Description Возвращает список пользователей указанного возраста
// @Tags Users
// @Produce json
// @Param age path int true "Возраст (целое число)"
// @Success 200 {array} entity.FullUser
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /get/age/{age} [get]
func (h *Handlers) GetUsersByAge(c *gin.Context) {
	age, _ := strconv.Atoi(c.Param("age"))
	users, err := h.service.GetUsersByAge(age)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUsersByName возвращает пользователей по имени
// @Summary Получить по имени
// @Description Возвращает список пользователей с указанным именем
// @Tags Users
// @Produce json
// @Param name path string true "Имя пользователя"
// @Success 200 {array} entity.FullUser
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /get/name/{name} [get]
func (h *Handlers) GetUsersByName(c *gin.Context) {
	name := c.Param("name")
	users, err := h.service.GetUsersByName(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetAllUsers возвращает всех пользователей
// @Summary Получить всех пользователей
// @Description Возвращает полный список пользователей в системе
// @Tags Users
// @Produce json
// @Success 200 {array} entity.FullUser
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /get/all [get]
func (h *Handlers) GetAllUsers(c *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUserByID возвращает пользователя по ID
// @Summary Получить по ID
// @Description Возвращает полную информацию о пользователе по его идентификатору
// @Tags Users
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} entity.FullUser
// @Failure 400 {object} map[string]interface{} "Неверный формат ID"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /get/id/{id} [get]
func (h *Handlers) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.service.GetUserByID(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser обновляет данные пользователя
// @Summary Обновить пользователя
// @Description Обновляет информацию о пользователе по его ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Param user body entity.User true "Новые данные пользователя"
// @Success 200 {object} map[string]interface{} "Успешное обновление"
// @Failure 400 {object} map[string]interface{} "Неверные данные"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /update/{id} [put]
func (h *Handlers) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.UpdateUser(&user, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "id": userId})

}

// DeleteUser удаляет пользователя
// @Summary Удалить пользователя
// @Description Удаляет пользователя по его ID
// @Tags Users
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} map[string]interface{} "Успешное удаление"
// @Failure 400 {object} map[string]interface{} "Неверный формат ID"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /delete/{id} [delete]
func (h *Handlers) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.DeleteUser(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})

}
