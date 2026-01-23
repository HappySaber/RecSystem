package tests

import (
	"log"
	ssov1 "sso-microservice/internal/pb/sso.v1"
	"sso-microservice/tests/suite"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyAppID = 0
	appID      = 1
	appSecret  = "testsecret"

	passDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomPassword()
	log.Println("Registering email:", email)

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
		Name:     gofakeit.Name(),
		Surname:  gofakeit.LastName(),
		Role:     "user",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: pass,
		AppId:    appID,
	})
	require.NoError(t, err)

	loginTime := time.Now()

	token := respLogin.GetAccessToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, respReg.GetUserId(), (claims["uid"].(string)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))

	const deltaSeconds = 5

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func randomPassword() string {
	return gofakeit.Password(true, true, true, true, true, passDefaultLen)
}

func TestRegisterLogin_DuplicationRegistration(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomPassword()

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
		Name:     gofakeit.Name(),
		Surname:  gofakeit.LastName(),
		Role:     "user",
	})

	require.NoError(t, err)
	require.NotEmpty(t, respReg.GetUserId())

	respReg, err = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
		Name:     gofakeit.Name(),
		Surname:  gofakeit.LastName(),
		Role:     "user",
	})

	require.Error(t, err)
	require.Empty(t, respReg.GetUserId())
	assert.ErrorContains(t, err, "user already exists")
}

func TestRegister_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		testName    string
		email       string
		password    string
		name        string
		surname     string
		role        string
		expectedErr string
	}{
		{
			testName:    "Register with empty password",
			email:       gofakeit.Email(),
			password:    "",
			name:        gofakeit.Name(),
			surname:     gofakeit.LastName(),
			role:        "user",
			expectedErr: "wrong arguments",
		},
		{
			testName:    "Register with empty email",
			email:       "",
			password:    randomPassword(),
			name:        gofakeit.Name(),
			surname:     gofakeit.LastName(),
			role:        "user",
			expectedErr: "wrong arguments",
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    tt.email,
				Password: tt.password,
				Name:     tt.name,
				Surname:  tt.surname,
				Role:     tt.role,
			})
			require.Error(t, err)
			assert.ErrorContains(t, err, tt.expectedErr)
		})
	}
}

func TestLogin_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		appID       int32
		expectedErr string
	}{
		{
			name:        "Login with Empty Password",
			email:       gofakeit.Email(),
			password:    "",
			appID:       appID,
			expectedErr: "wrong arguments",
		},
		{
			name:        "Login with Empty Email",
			email:       "",
			password:    randomPassword(),
			appID:       appID,
			expectedErr: "wrong arguments",
		},
		{
			name:        "Login with Both Empty Email and Password",
			email:       "",
			password:    "",
			appID:       appID,
			expectedErr: "wrong arguments",
		},
		{
			name:        "Login with Non-Matching Password",
			email:       gofakeit.Email(),
			password:    randomPassword(),
			appID:       appID,
			expectedErr: "invalid credentials",
		},
		{
			name:        "Login without AppID",
			email:       gofakeit.Email(),
			password:    randomPassword(),
			appID:       emptyAppID,
			expectedErr: "invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    gofakeit.Email(),
				Password: randomPassword(),
				Name:     gofakeit.Name(),
				Surname:  gofakeit.LastName(),
				Role:     "user",
			})
			require.NoError(t, err)

			_, err = st.AuthClient.Login(ctx, &ssov1.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
				AppId:    tt.appID,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}
