package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid" // <-- YENİ: jti (JWT ID) için eklendi
)

// ErrInvalidToken, token doğrulama hataları için standart bir hata.
var ErrInvalidToken = errors.New("geçersiz veya süresi dolmuş token")

// TokenManager, JWT oluşturma ve doğrulama işlemlerini yönetir.
type TokenManager struct {
	secretKey            []byte
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

// --- YENİ: Token Tipi Sabitleri ---
// "access", "refresh" gibi sihirli dizgeleri (magic strings) önler.
type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

// --- GÜNCELLENDİ: CustomClaims ---
// UserID artık 'sub' (Subject) standardı içinde taşınacak.
type CustomClaims struct {
	Type TokenType `json:"type"` // 'access' veya 'refresh'
	jwt.RegisteredClaims
}

// --- GÜNCELLENDİ: NewTokenManager ---
// Yeni bir TokenManager örneği oluşturur.
func NewTokenManager(secret string, accessDuration, refreshDuration time.Duration) (*TokenManager, error) {
	// HS256 (SHA-256), 256 bitlik (32 byte) bir anahtar bekler.
	// Güvenlik için minimum bir uzunluk zorlaması ekliyoruz.
	if len(secret) < 32 {
		return nil, errors.New("JWT secret key, güvenlik için en az 32 byte olmalıdır")
	}

	if accessDuration <= 0 || refreshDuration <= 0 {
		return nil, errors.New("token süreleri pozitif bir değer olmalıdır")
	}

	return &TokenManager{
		secretKey:            []byte(secret),
		accessTokenDuration:  accessDuration,
		refreshTokenDuration: refreshDuration,
	}, nil
}

// generateToken, verilen claim'ler ile yeni bir token imzalar.
func (tm *TokenManager) generateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(tm.secretKey)
}

// --- GÜNCELLENDİ: GenerateTokens ---
// Bir kullanıcı ID'si için standart claim'leri (sub, jti, iat vb.) içeren
// yeni bir access ve refresh token çifti oluşturur.
func (tm *TokenManager) GenerateTokens(userID string) (string, string, error) {
	now := time.Now()

	// Access Token Claims
	accessClaims := CustomClaims{
		Type: TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			// 'sub' (Subject) standardını UserID için kullanıyoruz
			Subject: userID,
			// 'jti' (JWT ID) token'ı benzersiz kılar, iptal (revocation) için kullanılır
			ID: uuid.NewString(),
			// 'iss' (Issuer) token'ı kimin oluşturduğu
			Issuer: "my-auth-service",
			// 'aud' (Audience) token'ın kimin/hangi servis için olduğu
			Audience: jwt.ClaimStrings{"my-app-client"},
			// 'iat' (Issued At) ne zaman oluşturulduğu
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(tm.accessTokenDuration)),
		},
	}
	accessTokenString, err := tm.generateToken(accessClaims)
	if err != nil {
		return "", "", fmt.Errorf("access token imzalanamadı: %w", err)
	}

	// Refresh Token Claims
	refreshClaims := CustomClaims{
		Type: TokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: userID,
			ID:      uuid.NewString(), // Refresh token için de benzersiz ID
			Issuer:  "my-auth-service",
			// Refresh token'ın hedef kitlesi *sadece* auth servisidir
			Audience:  jwt.ClaimStrings{"my-auth-service"},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(tm.refreshTokenDuration)),
		},
	}
	refreshTokenString, err := tm.generateToken(refreshClaims)
	if err != nil {
		return "", "", fmt.Errorf("refresh token imzalanamadı: %w", err)
	}

	return accessTokenString, refreshTokenString, nil
}

// --- GÜNCELLENDİ: ValidateToken ---
// Bir token string'ini doğrular ve içindeki claims'i döndürür.
// Hata yönetimini daha spesifik hale getirir.
func (tm *TokenManager) ValidateToken(tokenString string) (*CustomClaims, error) {

	// Token'ı CustomClaims yapımıza göre ayrıştırıyoruz.
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// İmzalama yöntemini (alg) doğruluyoruz
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("beklenmedik imzalama algoritması: %v", token.Header["alg"])
		}
		// Gizli anahtarımızı döndürüyoruz
		return tm.secretKey, nil
	})

	// Hata Kontrolü:
	// Hata Kontrolü:
	if err != nil {
		// Hatanın süresinin dolmasından mı kaynaklandığını kontrol et
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrInvalidToken // Süresi dolmuş
		}

		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, ErrInvalidToken // İmza geçersiz
		}

		// Diğer tüm JWT ile ilgili hatalar için (örn: format bozuk, henüz geçerli değil vb.)
		// yine standart hatamızı dönelim. İstemciye detay vermemek en güvenlisidir.
		return nil, ErrInvalidToken
	}

	// Token'ı ve claims'i al
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
