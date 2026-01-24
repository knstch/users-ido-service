package users

import (
	"crypto/sha3"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/knstch/knstch-libs/auth"

	"users-service/internal/domain/dto"
	"users-service/internal/domain/enum"
)

func (s *ServiceImpl) mintJWT(userID uint64, role enum.Role) (dto.AccessTokens, error) {
	timeNow := time.Now()
	// Access token: signed JWT with user id and role.
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, auth.Claims{
		UserID: strconv.Itoa(int(userID)),
		Role:   role.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(timeNow.Add(30 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(timeNow),
		},
	})

	signedAccessToken, err := accessToken.SignedString([]byte(s.cfg.JwtSecret))
	if err != nil {
		return dto.AccessTokens{}, err
	}

	// Refresh token: opaque token (sha3 hash) stored server-side.
	rawRefreshToken := []byte(fmt.Sprintf("%s%d", signedAccessToken, time.Now().Unix()))
	hash := sha3.New256()
	_, err = hash.Write(rawRefreshToken)
	if err != nil {
		return dto.AccessTokens{}, err
	}
	refreshToken := hex.EncodeToString(hash.Sum(nil))

	return dto.AccessTokens{
		AccessToken:  signedAccessToken,
		RefreshToken: refreshToken,
	}, nil
}
