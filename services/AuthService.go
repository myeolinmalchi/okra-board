package services

import (
	"errors"
	"okra_board2/models"
	"okra_board2/repositories"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type AuthService interface {

    // 관리자 id를 통해 새로운 Access Token, Refresh Token 쌍을 발급하고, db에 저장한다.
    CreateTokenPair(string)             (auth *models.AdminAuth, err error)

    // 이미 존재 하는 토큰 쌍의 uuid와 관리자 id 정보를 가지고 Access Token을 새롭게 발급한다.
    // 재발급된 Access Token은 db상에서 업데이트된다.
    CreateAccessToken(uuid, id string)   (at string, err error)

    // Refresh Token에서 추출한 uuid로 db상의 토큰 쌍을 검색하여
    // 주어진 토큰 쌍과 일치하는지 검증한다.
    // 검증에 실패할 경우 error를 반환한다.
    VerifyTokenPair(at, rt string)      (uuid string, err error)

    // Access Token의 유효성을 검증하고, claim과 error를 반환한다.
    VerifyAccessToken(string)           (claims jwt.MapClaims, err error)

    // Refresh Token의 유효성을 검증하고, claim과 error를 반환한다.
    VerifyRefreshToken(string)          (claims jwt.MapClaims, err error)

    // Access Token, Refresh Token 쌍을 db에서 삭제한다.
    DeleteTokenPair(uuid string)        (err error)
}

type AuthServiceImpl struct {
    authRepo repositories.AuthRepository
    adminService AdminService
}

func NewAuthServiceImpl(
    authRepo repositories.AuthRepository,
    adminService AdminService,
) AuthService {
    return &AuthServiceImpl{ 
        authRepo: authRepo,
        adminService: adminService,
    }
}

func (s *AuthServiceImpl) CreateTokenPair(id string) (*models.AdminAuth, error) {
    var err error

    adminAuth := &models.AdminAuth{
        UUID: uuid.NewString(),
        AdminID: id,
    }

    admin, err := s.adminService.GetAdmin(id)
    if err != nil { return nil, err }

    atClaims := jwt.MapClaims{}
    atClaims["authorized"] = true
    atClaims["uuid"] = adminAuth.UUID
    atClaims["id"] = id
    atClaims["name"] = admin.Name
    atClaims["exp"] = time.Now().Add(time.Second * 15).Unix()
    
    at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
    adminAuth.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
    if err != nil {
        return nil, err
    }

    rtClaims := jwt.MapClaims{}
    rtClaims["uuid"] = adminAuth.UUID
    rtClaims["id"] = id
    rtClaims["name"] = admin.Name
    rtClaims["exp"] = time.Now().Add(time.Hour * 24 * 3).Unix()
    rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
    adminAuth.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
    if err != nil {
        return nil, err
    }
    err = s.authRepo.InsertAdminAuth(adminAuth)
    if err != nil {
        return nil, err
    }
    return adminAuth, nil
}

func (s *AuthServiceImpl) CreateAccessToken(uuid, id string) (string, error) {
    atClaims := jwt.MapClaims{}

    admin, err := s.adminService.GetAdmin(id)
    if err != nil { return "", err }

    atClaims["authorized"] = true
    atClaims["uuid"] = uuid
    atClaims["id"] = id
    atClaims["name"] = admin.Name
    atClaims["exp"] = time.Now().Add(time.Second * 15).Unix()
    at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
    token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
    if err != nil {
        return "", err
    }
    if err := s.authRepo.UpdateAccessToken(uuid, token); err != nil {
        return "", err
    }
    return token, nil
}

func (s *AuthServiceImpl) VerifyTokenPair(at, rt string) (string, error) {
    rtClaims, err := s.VerifyRefreshToken(rt)
    if err != nil {
        return "", err
    }
    uuid := rtClaims["uuid"].(string)
    adminAuth := &models.AdminAuth{}
    adminAuth, err = s.authRepo.GetAdminAuth(uuid)
    if err != nil {
        return "", err
    }

    if adminAuth.AccessToken == at && adminAuth.RefreshToken == rt {
        return uuid, nil
    } else {
        return "", errors.New("Invalid Token Pair.")
    }
}

func (s *AuthServiceImpl) VerifyAccessToken(token string) (jwt.MapClaims, error) {
    claims := jwt.MapClaims{}
    verifying := func(token *jwt.Token) (interface{}, error) {
        if token.Method != jwt.SigningMethodHS256 {
            return nil, errors.New("Unexpected Signing Method")
        }
        return []byte(os.Getenv("ACCESS_SECRET")), nil
    }
    _, err := jwt.ParseWithClaims(token, &claims, verifying)
    return claims, err
}

func (s *AuthServiceImpl) VerifyRefreshToken(token string) (jwt.MapClaims, error) {
    claims := jwt.MapClaims{}
    verifying := func(token *jwt.Token) (interface{}, error) {
        if token.Method != jwt.SigningMethodHS256 {
            return nil, errors.New("Unexpected Signing Method")
        }
        return []byte(os.Getenv("REFRESH_SECRET")), nil
    }
    _, err := jwt.ParseWithClaims(token, &claims, verifying)
    return claims, err
}

func (s *AuthServiceImpl) DeleteTokenPair(uuid string) error {
    return s.authRepo.DeleteAdminAuth(uuid)
}
