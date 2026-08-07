package controller

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"webadmin/model"

	"github.com/gin-gonic/gin"
)

func successResponse(c *gin.Context, code int, status int, data gin.H) {
	if data == nil {
		data = gin.H{}
	}
	data["status"] = status
	c.JSON(code, data)
}

func errorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"status": 0,
		"error":  message,
	})
}

func md5Password(password string) string {
	sum := md5.Sum([]byte(password))
	return fmt.Sprintf("%x", sum)
}

type loginAttempt struct {
	failCount   int
	lockedUntil time.Time
}

var (
	loginAttempts   = make(map[string]*loginAttempt)
	loginAttemptsMu sync.RWMutex
)

func recordFailedLogin(username string) (int, time.Duration) {
	loginAttemptsMu.Lock()
	defer loginAttemptsMu.Unlock()

	attempt, exists := loginAttempts[username]
	if !exists {
		attempt = &loginAttempt{}
		loginAttempts[username] = attempt
	}

	attempt.failCount++

	if attempt.failCount >= 3 {
		attempt.lockedUntil = time.Now().Add(30 * time.Minute)
		return attempt.failCount, 30 * time.Minute
	}

	return attempt.failCount, 0
}

func clearLoginAttempts(username string) {
	loginAttemptsMu.Lock()
	delete(loginAttempts, username)
	loginAttemptsMu.Unlock()
}

func isLoginLocked(username string) (bool, time.Duration) {
	loginAttemptsMu.RLock()
	defer loginAttemptsMu.RUnlock()

	attempt, exists := loginAttempts[username]
	if !exists {
		return false, 0
	}

	remaining := time.Until(attempt.lockedUntil)
	if remaining > 0 {
		return true, remaining
	}

	return false, 0
}

func AdminList(c *gin.Context) {
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 20
	offset := (page - 1) * pageSize

	var admins []model.Admin
	var total int64

	if err := model.Db.Model(&model.Admin{}).Count(&total).Error; err != nil {
		errorResponse(c, 500, "Failed to count admins")
		return
	}

	if err := model.Db.Order("id asc").Limit(pageSize).Offset(offset).Find(&admins).Error; err != nil {
		errorResponse(c, 500, "Failed to retrieve admins")
		return
	}

	status := 0
	if len(admins) > 0 {
		status = 1
	}
	successResponse(c, 200, status, gin.H{
		"page":     page,
		"pagesize": pageSize,
		"total":    total,
		"data":     admins,
	})
}

func AdminCreate(c *gin.Context) {
	var admin model.Admin
	if err := c.ShouldBindJSON(&admin); err != nil {
		errorResponse(c, 400, "Invalid request")
		return
	}

	admin.Username = strings.TrimSpace(admin.Username)
	admin.Email = strings.TrimSpace(admin.Email)
	admin.Password = strings.TrimSpace(admin.Password)

	if admin.Username == "" {
		errorResponse(c, 400, "用户名不能为空")
		return
	}
	if admin.Email == "" {
		errorResponse(c, 400, "邮箱不能为空")
		return
	}
	if admin.Password == "" {
		errorResponse(c, 400, "密码不能为空")
		return
	}

	admin.Password = md5Password(admin.Password)

	if err := model.Db.Create(&admin).Error; err != nil {
		errorResponse(c, 500, "创建管理员失败")
		return
	}
	successResponse(c, 201, 1, gin.H{"admin": admin})
}

func AdminDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		errorResponse(c, 400, "无效的管理员ID")
		return
	}

	result := model.Db.Delete(&model.Admin{}, id)
	if result.Error != nil {
		errorResponse(c, 500, "删除管理员失败")
		return
	}

	if result.RowsAffected == 0 {
		errorResponse(c, 404, "管理员未找到")
		return
	}

	successResponse(c, 200, 1, gin.H{"message": "admin deleted"})
}

func AdminUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		errorResponse(c, 400, "invalid admin id")
		return
	}

	var admin model.Admin
	if err := c.ShouldBindJSON(&admin); err != nil {
		errorResponse(c, 400, "Invalid request")
		return
	}

	admin.Username = strings.TrimSpace(admin.Username)
	admin.Email = strings.TrimSpace(admin.Email)
	admin.Password = strings.TrimSpace(admin.Password)

	if admin.Username == "" {
		errorResponse(c, 400, "username is required")
		return
	}
	if admin.Email == "" {
		errorResponse(c, 400, "email is required")
		return
	}

	updates := map[string]any{
		"username":   admin.Username,
		"email":      admin.Email,
		"status":     admin.Status,
		"auth_group": admin.Auth_group,
	}
	if admin.Password != "" {
		updates["password"] = md5Password(admin.Password)
	}

	result := model.Db.Model(&model.Admin{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		errorResponse(c, 500, "Failed to update admin")
		return
	}
	if result.RowsAffected == 0 {
		errorResponse(c, 404, "admin not found")
		return
	}

	if err := model.Db.First(&admin, id).Error; err != nil {
		errorResponse(c, 500, "Failed to load admin")
		return
	}

	successResponse(c, 200, 1, gin.H{"admin": admin})
}

type loginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func AdminLogin(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, "Invalid request")
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	if req.Password == "" {
		errorResponse(c, 400, "密码不能为空")
		return
	}
	if req.Username == "" && req.Email == "" {
		errorResponse(c, 400, "用户名或邮箱不能为空")
		return
	}

	loginKey := req.Username
	if loginKey == "" {
		loginKey = req.Email
	}

	locked, remaining := isLoginLocked(loginKey)
	if locked {
		minutes := int(remaining.Minutes()) + 1
		errorResponse(c, 429, fmt.Sprintf("登录失败次数过多，请在%d分钟后重试", minutes))
		return
	}

	req.Password = md5Password(req.Password)

	var admin model.Admin
	query := model.Db
	if req.Username != "" {
		query = query.Where("username = ?", req.Username)
	} else {
		query = query.Where("email = ?", req.Email)
	}

	if err := query.First(&admin).Error; err != nil {
		recordFailedLogin(loginKey)
		errorResponse(c, 401, "用户名或密码不正确")
		return
	}

	if admin.Password != req.Password {
		recordFailedLogin(loginKey)
		errorResponse(c, 401, "用户名或密码不正确")
		return
	}

	clearLoginAttempts(loginKey)

	token, err := GenerateToken(admin.Id, admin.Username, admin.Email)
	if err != nil {
		errorResponse(c, 500, "生成令牌失败")
		return
	}

	successResponse(c, 200, 1, gin.H{
		"token":    token,
		"id":       admin.Id,
		"username": admin.Username,
		"email":    admin.Email,
	})
}
