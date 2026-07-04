package controller

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"
	"time"

	"webadmin/model"

	"github.com/gin-gonic/gin"
)

func UserList(c *gin.Context) {
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 20
	offset := (page - 1) * pageSize

	var users []model.User
	var total int64

	if err := model.Db.Model(&model.User{}).Count(&total).Error; err != nil {
		errorResponse(c, 500, "Failed to count users")
		return
	}

	if err := model.Db.Order("id asc").Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
		errorResponse(c, 500, "Failed to retrieve users")
		return
	}

	status := 0
	if len(users) > 0 {
		status = 1
	}
	successResponse(c, 200, status, gin.H{
		"page":     page,
		"pagesize": pageSize,
		"total":    total,
		"data":     users,
	})
}

func UserCreate(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		errorResponse(c, 400, "Invalid request")
		return
	}

	user.Username = strings.TrimSpace(user.Username)
	user.Email = strings.TrimSpace(user.Email)
	user.Password = strings.TrimSpace(user.Password)
	user.Bnaccess = strings.TrimSpace(user.Bnaccess)
	user.Bnasecret = strings.TrimSpace(user.Bnasecret)

	if user.Username == "" {
		errorResponse(c, 400, "username is required")
		return
	}
	if user.Password == "" {
		errorResponse(c, 400, "password is required")
		return
	}

	// MD5 encrypt password
	sum := md5.Sum([]byte(user.Password))
	user.Password = fmt.Sprintf("%x", sum)

	user.Addtime = int(time.Now().Unix())

	if err := model.Db.Create(&user).Error; err != nil {
		errorResponse(c, 500, "Failed to create user")
		return
	}
	successResponse(c, 201, 1, gin.H{"user": user})
}

func UserUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		errorResponse(c, 400, "invalid user id")
		return
	}

	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, "Invalid request")
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
	req.Bnaccess = strings.TrimSpace(req.Bnaccess)
	req.Bnasecret = strings.TrimSpace(req.Bnasecret)

	if req.Username == "" {
		errorResponse(c, 400, "username is required")
		return
	}

	updates := map[string]any{
		"username":  req.Username,
		"email":     req.Email,
		"usdt":      req.Usdt,
		"status":    req.Status,
		"bnaccess":  req.Bnaccess,
		"bnasecret": req.Bnasecret,
	}
	if req.Password != "" {
		sum := md5.Sum([]byte(req.Password))
		updates["password"] = fmt.Sprintf("%x", sum)
	}

	result := model.Db.Model(&model.User{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		errorResponse(c, 500, "Failed to update user")
		return
	}
	if result.RowsAffected == 0 {
		errorResponse(c, 404, "user not found")
		return
	}

	var user model.User
	if err := model.Db.First(&user, id).Error; err != nil {
		errorResponse(c, 500, "Failed to load user")
		return
	}

	successResponse(c, 200, 1, gin.H{"user": user})
}

func UserDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		errorResponse(c, 400, "invalid user id")
		return
	}

	result := model.Db.Delete(&model.User{}, id)
	if result.Error != nil {
		errorResponse(c, 500, "Failed to delete user")
		return
	}

	if result.RowsAffected == 0 {
		errorResponse(c, 404, "user not found")
		return
	}

	successResponse(c, 200, 1, gin.H{"message": "user deleted"})
}
