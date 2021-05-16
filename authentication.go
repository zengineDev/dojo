package dojo

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gofrs/uuid"
	"github.com/zengineDev/dojo/helpers"
	"golang.org/x/crypto/argon2"
	"strings"
)

var (
	// ErrInvalidHash in returned by ComparePasswordAndHash if the provided
	// hash isn't in the expected format.
	ErrInvalidHash = errors.New("argon2id: hash is not in the correct format")

	// ErrIncompatibleVersion in returned by ComparePasswordAndHash if the
	// provided hash was created using a different version of Argon2.
	ErrIncompatibleVersion = errors.New("argon2id: incompatible version of argon2")
)

var DefaultConfigs = &PasswordConfig{
	Memory:      64 * 1024,
	Iterations:  1,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

type PasswordConfig struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

const authUserSessionKey = "auth_user"
const oauthStateSessionKey = "oauth_state"

type AuthUserType string

const (
	GuestUserType AuthUserType = "guest"
	UserUserType  AuthUserType = "user"
)

type Authenticable interface {
	GetAuthType() AuthUserType
	GetAuthID() uuid.UUID
	GetAuthData() interface{}
}

type AuthUser struct {
	ID   uuid.UUID
	Data interface{}
}

func (u *AuthUser) GetAuthType() AuthUserType {
	if u.ID == uuid.Nil {
		return GuestUserType
	}
	return UserUserType
}

func (u *AuthUser) GetAuthID() uuid.UUID {
	return u.ID
}

func (u *AuthUser) GetAuthData() interface{} {
	return u.Data
}

type Authentication struct {
	app *Application
}

func NewAuthentication(app *Application) *Authentication {
	gob.Register(AuthUser{})
	return &Authentication{app: app}
}

func (auth *Authentication) GetAuthUser(ctx Context) AuthUser {
	session := auth.app.getSession(ctx.Request(), ctx.Response())
	sessionData := session.Get(authUserSessionKey)
	if sessionData == nil {
		return AuthUser{
			ID:   uuid.Nil,
			Data: nil,
		}
	}
	return sessionData.(AuthUser)
}

func (auth *Authentication) Login(ctx Context, user Authenticable) error {
	session := auth.app.getSession(ctx.Request(), ctx.Response())
	session.Set(authUserSessionKey, AuthUser{
		ID:   user.GetAuthID(),
		Data: user.GetAuthData(),
	})
	return session.Save()
}

func (auth *Authentication) GetAuthorizationUri(ctx Context) string {
	cfg := auth.app.Configuration.Auth
	state := helpers.RandomString(16)
	session := auth.app.getSession(ctx.Request(), ctx.Response())
	session.Set(oauthStateSessionKey, state)
	return fmt.Sprintf("%s?response_type=%s&client_id=%sredirect_uri=%s&scope=%s&state=%s",
		fmt.Sprintf("%s/auhtorize", cfg.Endpoint),
		"code",
		cfg.ClientID,
		fmt.Sprintf("%s%s", auth.app.Configuration.App.Domain, cfg.RedirectPath),
		strings.Join(cfg.Scopes, "+"),
		state,
	)
}

type OAuthResult struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_at"`
	Scope        string `json:"scope"`
}

type ExchangeAuthorisationCodeRequest struct {
	GrantType    string
	ClientId     string
	ClientSecret string
	RedirectUri  string
	Code         string
}

func (auth Authentication) CompareOAuthState(ctx Context, state string) error {
	session := auth.app.getSession(ctx.Request(), ctx.Response())
	sessionState := session.Get(oauthStateSessionKey)
	if fmt.Sprintf("%s", sessionState) != state {
		return errors.New("oauth state dont match")
	}
	return nil
}

func (auth *Authentication) ExchangeAuthorisationCode(ctx Context, authorisationCode string) (OAuthResult, error) {
	var result OAuthResult
	cfg := auth.app.Configuration.Auth
	body := ExchangeAuthorisationCodeRequest{
		GrantType:    "authorisation_code",
		ClientId:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectUri:  "",
		Code:         authorisationCode,
	}
	client := resty.New()
	resp, err := client.R().
		EnableTrace().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&result).
		Post(fmt.Sprintf("%s/token", cfg.Endpoint))

	if err != nil {
		return result, err
	}

	if resp.Error() != nil {
		return result, fmt.Errorf("%s", resp.Error())
	}

	return result, nil
}

func (auth Authentication) GeneratePasswordHash(c *PasswordConfig, password string) (string, error) {
	salt, err := generateRandomBytes(c.SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, c.Iterations, c.Memory, c.Parallelism, c.KeyLength)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	full := fmt.Sprintf(format, argon2.Version, c.Memory, c.Iterations, c.Parallelism, b64Salt, b64Hash)
	return full, nil
}

func (auth Authentication) ComparePasswordAndHash(password, hash string) (match bool, err error) {
	match, _, err = checkHash(password, hash)
	return match, err
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func checkHash(password, hash string) (match bool, params *PasswordConfig, err error) {
	params, salt, key, err := decodeHash(hash)
	if err != nil {
		return false, nil, err
	}

	otherKey := argon2.IDKey([]byte(password), salt, params.Iterations, params.Memory, params.Parallelism, params.KeyLength)

	keyLen := int32(len(key))
	otherKeyLen := int32(len(otherKey))

	if subtle.ConstantTimeEq(keyLen, otherKeyLen) == 0 {
		return false, params, nil
	}
	if subtle.ConstantTimeCompare(key, otherKey) == 1 {
		return true, params, nil
	}
	return false, params, nil
}

func decodeHash(hash string) (params *PasswordConfig, salt, key []byte, err error) {
	vals := strings.Split(hash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	params = &PasswordConfig{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &params.Memory, &params.Iterations, &params.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	params.SaltLength = uint32(len(salt))

	key, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	params.KeyLength = uint32(len(key))

	return params, salt, key, nil
}
