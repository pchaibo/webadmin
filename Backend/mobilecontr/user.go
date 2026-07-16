package mobilecontr

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"

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
	fmt.Println("Error:", message)
	c.JSON(code, gin.H{
		"status": 0,
		"error":  message,
	})
}

func md5Password(password string) string {
	sum := md5.Sum([]byte(password))
	return fmt.Sprintf("%x", sum)
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

func UserCreate(c *gin.Context) {
	var admin model.User
	if err := c.ShouldBindJSON(&admin); err != nil {
		errorResponse(c, 400, "Invalid request")
		return
	}

	admin.Username = strings.TrimSpace(admin.Username)
	admin.Email = strings.TrimSpace(admin.Email)
	admin.Password = strings.TrimSpace(admin.Password)
 	admin.Mobile = strings.TrimSpace(admin.Mobile)

	if admin.Username == "" {
		errorResponse(c, 400, "username is required")
		return
	}
	if admin.Email == "" {
		errorResponse(c, 400, "email is required")
		return
	}
	if admin.Password == "" {
		errorResponse(c, 400, "password is required")
		return
	}
 	if admin.Mobile == "" {
 		admin.Mobile = ""
 	}

	admin.Password = md5Password(admin.Password)

	if err := model.Db.Create(&admin).Error; err != nil {
		errorResponse(c, 500, "Failed to create user")
		return
	}
	successResponse(c, 201, 1, gin.H{"user": admin})
}

func UserUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		errorResponse(c, 400, "invalid user id")
		return
	}

	var admin model.User
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
		"username": admin.Username,
		"email":    admin.Email,
		"status":   admin.Status,
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

func UserLogin(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, "Invalid request")
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	if req.Password == "" {
		errorResponse(c, 400, "password is required")
		return
	}
	if req.Username == "" && req.Email == "" {
		errorResponse(c, 400, "username or email is required")
		return
	}

	req.Password = md5Password(req.Password)

	var admin model.User
	query := model.Db
	if req.Username != "" {
		query = query.Where("username = ?", req.Username)
	} else {
		query = query.Where("email = ?", req.Email)
	}

	if err := query.First(&admin).Error; err != nil {
		errorResponse(c, 401, "invalid username/email or password")
		return
	}

	if admin.Password != req.Password {
		errorResponse(c, 401, "invalid username/email or password")
		return
	}

	token, err := GenerateToken(admin.Id, admin.Username, admin.Email)
	if err != nil {
		errorResponse(c, 500, "Failed to generate token")
		return
	}

	successResponse(c, 200, 1, gin.H{
		"token":    token,
		"id":       admin.Id,
		"username": admin.Username,
		"email":    admin.Email,
	})
}
func UserInfo(c *gin.Context) {
	userId, _ := c.Get("user_id")
	uid, ok := userId.(int)
	if !ok || uid <= 0 {
		errorResponse(c, 401, "unauthorized")
		return
	}

	var user model.User
	if err := model.Db.First(&user, uid).Error; err != nil {
		errorResponse(c, 404, "user not found")
		return
	}

	successResponse(c, 200, 1, gin.H{
		"data": gin.H{
			"id":        user.Id,
			"username":  user.Username,
			"email":     user.Email,
			"usdt":      user.Usdt,
			"margin":    user.Margin,
			"bnaccess":  user.Bnaccess,
			"bnasecret": user.Bnasecret,
		},
	})
}

func UserUpdatePassword(c *gin.Context) {
	userId, _ := c.Get("user_id")
	uid, ok := userId.(int)
	if !ok || uid <= 0 {
		errorResponse(c, 401, "unauthorized")
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, "Invalid request")
		return
	}

	req.Password = strings.TrimSpace(req.Password)
	if req.Password == "" {
		errorResponse(c, 400, "password is required")
		return
	}

	updates := map[string]any{
		"password": md5Password(req.Password),
	}
	result := model.Db.Model(&model.User{}).Where("id = ?", uid).Updates(updates)
	if result.Error != nil {
		errorResponse(c, 500, "Failed to update password")
		return
	}

	successResponse(c, 200, 1, gin.H{"message": "密码更新成功"})
}

func UserUpdateBinance(c *gin.Context) {
	userId, _ := c.Get("user_id")
	uid, ok := userId.(int)
	if !ok || uid <= 0 {
		errorResponse(c, 401, "unauthorized")
		return
	}

	var req struct {
		Bnaccess  string  `json:"bnaccess"`
		Bnasecret string  `json:"bnasecret"`
		Margin    float64 `json:"margin"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, "Invalid request")
		return
	}

	req.Bnaccess = strings.TrimSpace(req.Bnaccess)
	req.Bnasecret = strings.TrimSpace(req.Bnasecret)

	updates := map[string]any{}
	updates["margin"] = req.Margin
	if req.Bnaccess != "" {
		updates["bnaccess"] = req.Bnaccess
	}
	if req.Bnasecret != "" {
		updates["bnasecret"] = req.Bnasecret
	}
	if len(updates) == 0 {
		errorResponse(c, 400, "bnaccess or bnasecret is required")
		return
	}

	result := model.Db.Model(&model.User{}).Where("id = ?", uid).Updates(updates)
	if result.Error != nil {
		errorResponse(c, 500, "Failed to update binance settings")
		return
	}

	successResponse(c, 200, 1, gin.H{"message": "币安设置更新成功"})
}
