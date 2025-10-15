package user

import (
	"fmt"
	"log"
	"math/rand"
	"auto-forge/internal/models"
	"auto-forge/pkg/cache"
	"auto-forge/pkg/common"
	"auto-forge/pkg/database"
	"auto-forge/pkg/email"
	"auto-forge/pkg/errors"
	"auto-forge/pkg/utils"
	"time"
)

var userService *UserService


type UserService struct {

}


func InitUserService() {
	userService = &UserService{}
}


func GetUserService() *UserService {
	return userService
}


func Login(account, password string) (map[string]interface{}, string, time.Time, error) {
	db := database.GetDB()
	if db == nil {
		return nil, "", time.Time{}, errors.New(errors.CodeDBConnectionFailed, "数据库连接失败")
	}


	var userRow struct {
		ID        string `db:"id"`
		Username  string `db:"username"`
		Password  string `db:"password"`
		Email     string `db:"email"`
		Avatar    string `db:"avatar"`
		Bio       string `db:"bio"`
		Status    int    `db:"status"`
		Role      int    `db:"role"`
		CreatedAt string `db:"created_at"`
		UpdatedAt string `db:"updated_at"`
	}

	err := db.Raw("SELECT id, username, password, email, avatar, bio, status, role, created_at, updated_at FROM user WHERE username = ? OR email = ? LIMIT 1", account, account).Scan(&userRow).Error
	if err != nil {
		return nil, "", time.Time{}, errors.New(errors.CodeQueryFailed, "数据库查询失败")
	}

	if userRow.ID == "" {
		return nil, "", time.Time{}, errors.New(errors.CodeUserNotFound, "用户不存在")
	}


	if !utils.ComparePasswords(userRow.Password, password) {
		return nil, "", time.Time{}, errors.New(errors.CodeWrongPassword, "密码错误")
	}


	if userRow.Status != common.UserStatusNormal {
		return nil, "", time.Time{}, errors.New(errors.CodeUserDisabled, "账号已被禁用")
	}


	userID, err := common.ParseUUID(userRow.ID)
	if err != nil {
		return nil, "", time.Time{}, errors.New(errors.CodeInternal, "用户ID格式错误")
	}


	token, err := common.GenerateToken(userID, userRow.Username, userRow.Role)
	if err != nil {
		return nil, "", time.Time{}, errors.New(errors.CodeInternal, "生成token失败")
	}


	expiresAt := time.Now().Add(time.Duration(24) * time.Hour)


	userInfo := map[string]interface{}{
		"id":       userRow.ID,
		"username": userRow.Username,
		"email":    userRow.Email,
		"avatar":   userRow.Avatar,
		"bio":      userRow.Bio,
		"role":     userRow.Role,
		"status":   userRow.Status,
	}

	return userInfo, token, expiresAt, nil
}


func FindUsers() ([]models.User, error) {
	db := database.GetDB()
	var users []models.User
	result := db.Find(&users)
	return users, result.Error
}


func FindUserByID(id string) (*models.User, error) {
	db := database.GetDB()
	var user models.User
	result := db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}


func FindUserByEmail(email string) (*models.User, error) {
	db := database.GetDB()
	var user models.User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}


func generateVerificationCode(email string, codeType string) string {

	rand.Seed(time.Now().UnixNano())


	code := fmt.Sprintf("%06d", rand.Intn(1000000))


	key := fmt.Sprintf("%s:%s:code", email, codeType)
	err := cache.GetCache().Set(key, code, 5*time.Minute)
	if err != nil {
		log.Printf("存储验证码到缓存失败: %v", err)
		return ""
	}

	return code
}


func SendRegistrationCode(email string) error {

	_, err := FindUserByEmail(email)
	if err == nil {
		return errors.New(errors.CodeEmailExists, "该邮箱已被注册")
	}


	code := generateVerificationCode(email, common.CodeTypeRegister)
	if code == "" {
		return errors.New(errors.CodeInternal, "生成验证码失败")
	}


	if err := sendVerificationEmail(email, code, common.CodeTypeRegister); err != nil {
		return fmt.Errorf("发送验证码失败: %v", err)
	}

	return nil
}


func SendResetPasswordCode(email string) error {

	_, err := FindUserByEmail(email)
	if err != nil {
		return errors.New(errors.CodeUserNotFound, "该邮箱尚未注册")
	}


	code := generateVerificationCode(email, common.CodeTypeResetPassword)
	if code == "" {
		return errors.New(errors.CodeInternal, "生成验证码失败")
	}


	if err := sendVerificationEmail(email, code, common.CodeTypeResetPassword); err != nil {
		return fmt.Errorf("发送验证码失败: %v", err)
	}

	return nil
}


func SendChangeEmailCode(userID, newEmail string) error {

	existingUser, err := FindUserByEmail(newEmail)
	if err == nil && existingUser.ID.String() != userID {
		return errors.New(errors.CodeEmailExists, "该邮箱已被其他用户使用")
	}


	code := generateVerificationCode(newEmail, common.CodeTypeChangeEmail)
	if code == "" {
		return errors.New(errors.CodeInternal, "生成验证码失败")
	}


	if err := sendVerificationEmail(newEmail, code, common.CodeTypeChangeEmail); err != nil {
		return fmt.Errorf("发送验证码失败: %v", err)
	}

	return nil
}


func ValidateCode(email, code, codeType string) bool {
	key := fmt.Sprintf("%s:%s:code", email, codeType)
	cachedCode, err := cache.GetCache().Get(key)
	if err != nil {
		log.Printf("获取验证码失败: %v", err)
		return false
	}


	if code == cachedCode {

		_ = cache.GetCache().Del(key)
		return true
	}

	return false
}


func RegisterUser(username, email, password, code string) error {
	db := database.GetDB()


	if !ValidateCode(email, code, common.CodeTypeRegister) {
		return errors.New(errors.CodeInvalidVerifyCode, "验证码无效或已过期")
	}


	var count int64
	db.Model(&models.User{}).Where("username = ?", username).Count(&count)
	if count > 0 {
		return errors.New(errors.CodeUserExists, "用户名已存在")
	}


	db.Model(&models.User{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		return errors.New(errors.CodeEmailExists, "邮箱已被注册")
	}


	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return errors.New(errors.CodeInternal, "密码加密失败")
	}


	var totalUserCount int64
	db.Model(&models.User{}).Count(&totalUserCount)

	role := common.UserRoleUser
	if totalUserCount == 0 {
		role = common.UserRoleSuperAdmin
		log.Printf("第一个注册用户 %s 已设置为超级管理员", username)
	}

	user := models.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
		Status:   common.UserStatusNormal,
		Role:     role,
	}

	if err := db.Create(&user).Error; err != nil {
		return errors.New(errors.CodeInternal, "创建用户失败")
	}

	return nil
}


func sendVerificationEmail(emailAddr string, code string, codeType string) error {

	if !email.IsMailEnabled() {
		return errors.New(errors.CodeEmailServiceError, "邮件服务不可用，请联系管理员")
	}

	var subject string
	if codeType == common.CodeTypeRegister {
		subject = "注册验证码"
	} else if codeType == common.CodeTypeChangeEmail {
		subject = "修改邮箱验证码"
	} else {
		subject = "重置密码验证码"
	}


	err := email.SendMail(emailAddr, subject, fmt.Sprintf("您的验证码是: %s，5分钟内有效。", code))
	if err != nil {
		return err
	}

	return nil
}


func ResetPassword(email, code, newPassword string) error {
	db := database.GetDB()


	if !ValidateCode(email, code, common.CodeTypeResetPassword) {
		return errors.New(errors.CodeInvalidVerifyCode, "验证码无效或已过期")
	}


	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New(errors.CodeInternal, "密码加密失败")
	}

	result := db.Model(&models.User{}).Where("email = ?", email).Update("password", hashedPassword)
	if result.Error != nil {
		return errors.New(errors.CodeInternal, "更新密码失败")
	}

	if result.RowsAffected == 0 {
		return errors.New(errors.CodeUserNotFound, "未找到用户")
	}

	return nil
}


func GetUserInfo(userID string) (map[string]interface{}, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New(errors.CodeDBConnectionFailed, "数据库连接失败")
	}


	var userRow struct {
		ID        string `db:"id"`
		Username  string `db:"username"`
		Email     string `db:"email"`
		Avatar    string `db:"avatar"`
		Bio       string `db:"bio"`
		Status    int    `db:"status"`
		Role      int    `db:"role"`
		CreatedAt string `db:"created_at"`
		UpdatedAt string `db:"updated_at"`
	}

	err := db.Raw("SELECT id, username, email, avatar, bio, status, role, created_at, updated_at FROM user WHERE id = ? LIMIT 1", userID).Scan(&userRow).Error
	if err != nil {
		return nil, errors.New(errors.CodeQueryFailed, "数据库查询失败")
	}

	if userRow.ID == "" {
		return nil, errors.New(errors.CodeUserNotFound, "用户不存在")
	}


	userInfo := map[string]interface{}{
		"id":         userRow.ID,
		"username":   userRow.Username,
		"email":      userRow.Email,
		"avatar":     userRow.Avatar,
		"bio":        userRow.Bio,
		"role":       userRow.Role,
		"status":     userRow.Status,
		"created_at": userRow.CreatedAt,
		"updated_at": userRow.UpdatedAt,
	}

	return userInfo, nil
}


func UpdateProfile(userID, username, email, avatar, code string) (map[string]interface{}, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New(errors.CodeDBConnectionFailed, "数据库连接失败")
	}


	updateData := make(map[string]interface{})
	if username != "" {

		var count int64
		db.Model(&models.User{}).Where("username = ? AND id != ?", username, userID).Count(&count)
		if count > 0 {
			return nil, errors.New(errors.CodeUserExists, "用户名已被使用")
		}
		updateData["username"] = username
	}
	if email != "" {

		if code != "" {

			if !ValidateCode(email, code, common.CodeTypeChangeEmail) {
				return nil, errors.New(errors.CodeInvalidVerifyCode, "验证码无效或已过期")
			}
		}


		var count int64
		db.Model(&models.User{}).Where("email = ? AND id != ?", email, userID).Count(&count)
		if count > 0 {
			return nil, errors.New(errors.CodeEmailExists, "邮箱已被使用")
		}
		updateData["email"] = email
	}
	if avatar != "" {
		updateData["avatar"] = avatar
	}

	if len(updateData) == 0 {
		return nil, errors.New(errors.CodeInvalidParameter, "没有需要更新的数据")
	}


	result := db.Model(&models.User{}).Where("id = ?", userID).Updates(updateData)
	if result.Error != nil {
		return nil, errors.New(errors.CodeInternal, "更新用户信息失败")
	}

	if result.RowsAffected == 0 {
		return nil, errors.New(errors.CodeUserNotFound, "用户不存在")
	}


	return GetUserInfo(userID)
}


func ChangePassword(userID, oldPassword, newPassword string) error {
	db := database.GetDB()
	if db == nil {
		return errors.New(errors.CodeDBConnectionFailed, "数据库连接失败")
	}


	var userRow struct {
		Password string `db:"password"`
	}

	err := db.Raw("SELECT password FROM user WHERE id = ? LIMIT 1", userID).Scan(&userRow).Error
	if err != nil {
		return errors.New(errors.CodeQueryFailed, "数据库查询失败")
	}

	if userRow.Password == "" {
		return errors.New(errors.CodeUserNotFound, "用户不存在")
	}


	if !utils.ComparePasswords(userRow.Password, oldPassword) {
		return errors.New(errors.CodeWrongPassword, "原密码错误")
	}


	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New(errors.CodeInternal, "密码加密失败")
	}


	result := db.Model(&models.User{}).Where("id = ?", userID).Update("password", hashedPassword)
	if result.Error != nil {
		return errors.New(errors.CodeInternal, "更新密码失败")
	}

	if result.RowsAffected == 0 {
		return errors.New(errors.CodeUserNotFound, "用户不存在")
	}

	return nil
}
