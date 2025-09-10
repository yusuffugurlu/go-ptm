package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/yusuffugurlu/go-project/config/logger"
	"github.com/yusuffugurlu/go-project/internal/cache"
	"github.com/yusuffugurlu/go-project/internal/dtos"
	"github.com/yusuffugurlu/go-project/internal/models"
	"github.com/yusuffugurlu/go-project/internal/repositories"
	appErrors "github.com/yusuffugurlu/go-project/pkg/errors"
	"gorm.io/gorm"
)

type UserService interface {
	GetAllUsers() ([]models.User, error)
	GetUserById(id int) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) (*models.User, error)
	UpdateUser(id int, user *dtos.UpdateUserRequest) (*models.User, error)
	DeleteUser(userId int) error
}

type userService struct {
	userRepo     repositories.UserRepository
	logService   AuditLogService
	cacheService *cache.CacheService
}

func NewUserService(userRepository repositories.UserRepository, logService AuditLogService) UserService {
	return &userService{
		userRepo:     userRepository,
		logService:   logService,
		cacheService: nil,
	}
}

func NewUserServiceWithCache(userRepository repositories.UserRepository, logService AuditLogService, cacheService *cache.CacheService) UserService {
	return &userService{
		userRepo:     userRepository,
		logService:   logService,
		cacheService: cacheService,
	}
}

func (u *userService) SetCacheService(cacheService *cache.CacheService) {
	u.cacheService = cacheService
}

func (u *userService) CreateUser(user *models.User) (*models.User, error) {
	_, err := u.userRepo.GetByEmail(user.Email)
	if err == nil {
		return nil, appErrors.NewConflict(nil, fmt.Sprintf("email %s is already in use", user.Email))
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err := user.HashPassword(); err != nil {
		return nil, appErrors.NewInternalServerError(err)
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}

	if err := u.logService.CreateAuditLog(int(user.Id), "user", "create", fmt.Sprintf("user %d created", user.Id)); err != nil {
		return nil, err
	}

	if u.cacheService != nil {
		ctx := context.Background()
		u.invalidateUserCache(ctx, int(user.Id), user.Email)
	}

	return user, nil
}

func (u *userService) DeleteUser(userId int) error {
	user, err := u.userRepo.GetById(userId)
	if err != nil {
		return err
	}

	if err := u.userRepo.Delete(userId); err != nil {
		return err
	}

	if err := u.logService.CreateAuditLog(userId, "user", "delete", fmt.Sprintf("user %d deleted", userId)); err != nil {
		return err
	}

	if u.cacheService != nil {
		ctx := context.Background()
		u.invalidateUserCache(ctx, userId, user.Email)
	}

	return nil
}

func (u *userService) GetAllUsers() ([]models.User, error) {
	if u.cacheService != nil {
		ctx := context.Background()
		cacheKey := "users:all"

		var users []models.User
		if err := u.cacheService.GetJSON(ctx, cacheKey, &users); err == nil {
			logger.Log.Debug("Users retrieved from cache")
			return users, nil
		}
	}

	users, err := u.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	if u.cacheService != nil {
		ctx := context.Background()
		cacheKey := "users:all"
		if err := u.cacheService.SetJSON(ctx, cacheKey, users, 5*time.Minute); err != nil {
			logger.Log.Error("Failed to cache users", err)
		}
	}

	return users, nil
}

func (u *userService) GetUserByEmail(email string) (*models.User, error) {
	if u.cacheService != nil {
		ctx := context.Background()
		cacheKey := u.cacheService.GenerateCacheKey("user:email", email)

		var user models.User
		if err := u.cacheService.GetJSON(ctx, cacheKey, &user); err == nil {
			logger.Log.Debug("User retrieved from cache by email", "email", email)
			return &user, nil
		}
	}

	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	if u.cacheService != nil {
		ctx := context.Background()
		cacheKey := u.cacheService.GenerateCacheKey("user:email", email)
		if err := u.cacheService.SetJSON(ctx, cacheKey, user, 10*time.Minute); err != nil {
			logger.Log.Error("Failed to cache user by email", err)
		}
	}

	return user, nil
}

func (u *userService) GetUserById(id int) (*models.User, error) {
	if u.cacheService != nil {
		ctx := context.Background()
		cacheKey := u.cacheService.GenerateCacheKey("user", fmt.Sprintf("%d", id))

		var user models.User
		if err := u.cacheService.GetJSON(ctx, cacheKey, &user); err == nil {
			logger.Log.Debug("User retrieved from cache", "id", id)
			return &user, nil
		}
	}

	user, err := u.userRepo.GetById(id)
	if err != nil {
		return nil, err
	}

	if u.cacheService != nil {
		ctx := context.Background()
		cacheKey := u.cacheService.GenerateCacheKey("user", fmt.Sprintf("%d", id))
		if err := u.cacheService.SetJSON(ctx, cacheKey, user, 10*time.Minute); err != nil {
			logger.Log.Error("Failed to cache user", err)
		}
	}

	return user, nil
}

func (u *userService) UpdateUser(id int, userData *dtos.UpdateUserRequest) (*models.User, error) {
	existingUser, err := u.userRepo.GetById(id)
	if err != nil {
		return nil, err
	}

	if userData.Email != "" && userData.Email != existingUser.Email {
		user, err := u.userRepo.GetByEmail(userData.Email)
		if err == nil && user != nil {
			return nil, appErrors.NewConflict(nil, fmt.Sprintf("email %s is already in use", userData.Email))
		} else if !errors.Is(err, gorm.ErrRecordNotFound) && !appErrors.IsAppError(err) {
			return nil, appErrors.NewDatabaseError(err, "failed to check email uniqueness")
		}
	}

	if userData.Username != "" {
		existingUser.Username = userData.Username
	}

	if userData.Email != "" {
		existingUser.Email = userData.Email
	}

	if userData.Password != "" {
		existingUser.PasswordHash = userData.Password
		if err := existingUser.HashPassword(); err != nil {
			return nil, appErrors.NewInternalServerError(err)
		}
	}

	if err := u.userRepo.Update(existingUser); err != nil {
		return nil, err
	}

	if err := u.logService.CreateAuditLog(id, "user", "update", fmt.Sprintf("user %d updated", id)); err != nil {
		return nil, err
	}

	if u.cacheService != nil {
		ctx := context.Background()
		u.invalidateUserCache(ctx, id, existingUser.Email)

		if userData.Email != "" && userData.Email != existingUser.Email {
			u.invalidateUserCache(ctx, id, userData.Email)
		}
	}

	return existingUser, nil
}

func (u *userService) invalidateUserCache(ctx context.Context, userID int, email string) {
	userCacheKey := u.cacheService.GenerateCacheKey("user", fmt.Sprintf("%d", userID))
	u.cacheService.Delete(ctx, userCacheKey)

	emailCacheKey := u.cacheService.GenerateCacheKey("user:email", email)
	u.cacheService.Delete(ctx, emailCacheKey)

	allUsersCacheKey := "users:all"
	u.cacheService.Delete(ctx, allUsersCacheKey)

	logger.Log.Debug("User cache invalidated", "userID", userID, "email", email)
}
