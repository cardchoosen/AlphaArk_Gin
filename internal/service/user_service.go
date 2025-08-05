package service

import (
	"errors"
	"time"

	"github.com/cardchoosen/AlphaArk_Gin/internal/models"
	"github.com/cardchoosen/AlphaArk_Gin/internal/repository"
)

// UserService 用户服务接口
type UserService interface {
	GetAllUsers() ([]models.UserResponse, error)
	GetUserByID(id uint) (*models.UserResponse, error)
	CreateUser(req models.UserCreateRequest) (*models.UserResponse, error)
	UpdateUser(id uint, req models.UserUpdateRequest) (*models.UserResponse, error)
	DeleteUser(id uint) error
}

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// GetAllUsers 获取所有用户
func (s *userService) GetAllUsers() ([]models.UserResponse, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var responses []models.UserResponse
	for _, user := range users {
		responses = append(responses, s.toUserResponse(user))
	}

	return responses, nil
}

// GetUserByID 根据ID获取用户
func (s *userService) GetUserByID(id uint) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := s.toUserResponse(*user)
	return &response, nil
}

// CreateUser 创建用户
func (s *userService) CreateUser(req models.UserCreateRequest) (*models.UserResponse, error) {
	// 检查用户名是否已存在
	existingUser, _ := s.userRepo.FindByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// 检查邮箱是否已存在
	existingUser, _ = s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// 创建用户
	user := models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password, // TODO: 加密密码
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdUser, err := s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	response := s.toUserResponse(*createdUser)
	return &response, nil
}

// UpdateUser 更新用户
func (s *userService) UpdateUser(id uint, req models.UserUpdateRequest) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	user.UpdatedAt = time.Now()

	updatedUser, err := s.userRepo.Update(*user)
	if err != nil {
		return nil, err
	}

	response := s.toUserResponse(*updatedUser)
	return &response, nil
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}

// toUserResponse 转换为用户响应
func (s *userService) toUserResponse(user models.User) models.UserResponse {
	return models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
} 