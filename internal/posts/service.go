package posts

import (
	"context"
	"fmt"
	"time"
)

// Service -> Contract
type Service interface {
	Post(ctx context.Context, req PostRequest, userId uint) (PostResponse, error)
}

// Structure
type postsService struct {
	repo Repository
}

// NewService - Constructor
func NewService(repo Repository) Service {
	return &postsService{
		repo: repo,
	}
}

func (s *postsService) Post(ctx context.Context, req PostRequest, userId uint) (PostResponse, error) {
	// Build entity post
	post := &Post{
		Title:     req.Title,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}

	err := s.repo.Create(ctx, post)
	if err != nil {
		return PostResponse{}, fmt.Errorf("error on create new post: %w", err)
	}

	return PostResponse{
		Id:        post.Id,
		UserId:    post.UserId,
		Title:     post.Title,
		Content:   post.Content,
		ImageURL:  post.ImageURL,
		CreatedAt: post.CreatedAt,
	}, nil
}
