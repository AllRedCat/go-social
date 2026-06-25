package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// O ideal é que essa chave venha de uma variável de ambiente, ex: os.Getenv("JWT_SECRET")
var jwtSecretKey = []byte("sua_chave_super_secreta_aqui_mudar_em_prod") 

// Contract
type Service interface {
	Register(ctx context.Context, req RegisterRequest) (UserResponse, error)
	Login(ctx context.Context, req LoginRequest) (string, error)
	UpdateAvatar(ctx context.Context, userID uint, avatarURL string) error
	UpdateUser(ctx context.Context, user *User) error
	SoftDelete(ctx context.Context, userID uint) error
}

// Structure
type authService struct {
	repo Repository
}

// Constructor
func NewService(repo Repository) Service {
	return &authService{
		repo: repo,
	}
}

// Implementations

func (s *authService) Register(ctx context.Context, req RegisterRequest) (UserResponse, error) {
	// 1. Criptografar a senha usando bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return UserResponse{}, fmt.Errorf("falha ao processar a senha: %w", err)
	}

	// 2. Montar a entidade User
	user := &User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	// 3. Salvar no banco via repositório
	err = s.repo.Create(ctx, user)
	if err != nil {
		// Geralmente falha por causa de UNIQUE constraint (email já existe)
		return UserResponse{}, fmt.Errorf("erro ao registrar usuário (email já pode estar em uso)")
	}

	// 4. Retornar o DTO limpo
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *authService) Login(ctx context.Context, req LoginRequest) (string, error) {
	// 1. Buscar o usuário pelo email
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		// Para segurança, sempre retornamos a mesma mensagem genérica se errar email ou senha
		return "", fmt.Errorf("credenciais inválidas") 
	}

	// 2. Comparar a senha fornecida com o hash salvo
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", fmt.Errorf("credenciais inválidas")
	}

	// 3. Gerar o token JWT (JSON Web Token)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID, // Subject: ID do usuário
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Expira em 24h
	})

	// 4. Assinar o token
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", fmt.Errorf("erro ao gerar token de sessão: %w", err)
	}

	return tokenString, nil
}

func (s *authService) UpdateAvatar(ctx context.Context, userID uint, avatarURL string) error {
	return s.repo.UpdateAvatar(ctx, userID, avatarURL)
}

func (s *authService) UpdateUser(ctx context.Context, user *User) error {
	return s.repo.UpdateUser(ctx, user)
}

func (s *authService) SoftDelete(ctx context.Context, userID uint) error {
	return s.repo.SoftDelete(ctx, userID)
}
