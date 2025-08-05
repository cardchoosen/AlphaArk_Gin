package tests

import (
	"testing"

	"github.com/cardchoosen/AlphaArk_Gin/internal/models"
)

func TestUserModel(t *testing.T) {
	user := models.User{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "password123",
		FirstName: "Test",
		LastName:  "User",
		IsActive:  true,
	}

	if user.Username != "testuser" {
		t.Errorf("Expected username to be 'testuser', got %s", user.Username)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Expected email to be 'test@example.com', got %s", user.Email)
	}
}

func TestUserCreateRequest(t *testing.T) {
	req := models.UserCreateRequest{
		Username:  "newuser",
		Email:     "new@example.com",
		Password:  "newpassword",
		FirstName: "New",
		LastName:  "User",
	}

	if req.Username != "newuser" {
		t.Errorf("Expected username to be 'newuser', got %s", req.Username)
	}
}
