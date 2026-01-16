package users_test

import (
	"context"

	"github.com/knstch/knstch-libs/svcerrs"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	googleconn "users-service/internal/connector/google"
	"users-service/internal/domain/dto"
	"users-service/internal/users"
	"users-service/internal/users/filters"
	"users-service/internal/users/modles"
	"users-service/testhelper"
)

func (s *UsersServiceTestSuite) TestAuthViaGoogle_StoresStateAndReturnsLoginURL() {
	t := require.New(s.T())
	ctx := context.Background()

	loginURL, err := s.svc.AuthViaGoogle(ctx, testhelper.TestReturnHome)
	t.NoError(err)
	t.NotEmpty(loginURL)

	state := mustExtractState(t, loginURL)
	st := mustDecodeOAuthState(t, state)
	t.NotEmpty(st.CSRF)
	t.Equal(testhelper.TestReturnHome, st.Return)

	val, err := s.rdb.Get(ctx, testhelper.TestStateKeyPrefix+st.CSRF).Result()
	t.NoError(err)
	t.Equal(testhelper.TestReturnHome, val)
}

func (s *UsersServiceTestSuite) TestAuthViaGoogle_RejectsUnsafeReturnURL() {
	t := require.New(s.T())
	ctx := context.Background()

	_, err := s.svc.AuthViaGoogle(ctx, "//evil.com")
	t.Error(err)

	_, err = s.svc.AuthViaGoogle(ctx, "https://evil.com/hack")
	t.Error(err)
}

func (s *UsersServiceTestSuite) TestCompleteLogin_CreatesUserAndTokens_AndReturnsReturnURL() {
	t := require.New(s.T())
	ctx := context.Background()

	loginURL, err := s.svc.AuthViaGoogle(ctx, testhelper.TestReturnFeed)
	t.NoError(err)
	state := mustExtractState(t, loginURL)

	idToken := makeFakeIDToken(t, map[string]any{
		"sub":     testhelper.TestGoogleSub,
		"email":   testhelper.TestEmail,
		"name":    testhelper.TestNameTwoParts,
		"picture": testhelper.TestPicture,
	})

	s.googleMock.
		On("ExchangeCodeToToken", mock.Anything, mock.MatchedBy(func(req googleconn.ExchangeCodeToTokenRequest) bool {
			return req.Code == testhelper.TestCode &&
				req.GoogleClientID == s.cfg.GoogleAPI.GoogleClientID &&
				req.ClientSecret == s.cfg.GoogleAPI.GoogleOAuthClientSecret &&
				req.RedirectURI == s.cfg.GoogleAPI.GoogleRedirectURI
		})).
		Return(&googleconn.ExchangeCodeToTokenResponse{
			IDToken: idToken,
		}, nil).
		Once()

	tokens, returnURL, err := s.svc.CompleteLogin(ctx, state, testhelper.TestCode)
	t.NoError(err)
	t.Equal(testhelper.TestReturnFeed, returnURL)
	t.NotEmpty(tokens.AccessToken)
	t.NotEmpty(tokens.RefreshToken)

	// token pair is persisted
	var tokenRow modles.AccessToken
	tokenTx := s.db.
		WithContext(ctx).
		Model(&modles.AccessToken{}).
		Scopes((&filters.AccessTokenFilter{RefreshToken: tokens.RefreshToken}).ToScope())
	err = tokenTx.First(&tokenRow).Error
	t.NoError(err)
	t.Equal(tokens.AccessToken, tokenRow.AccessToken)
	t.Equal(tokens.RefreshToken, tokenRow.RefreshToken)

	// user row exists in DB
	var userRow modles.User
	userTx := s.db.
		WithContext(ctx).
		Model(&modles.User{}).
		Where("email = ?", testhelper.TestEmail)
	err = userTx.First(&userRow).Error
	t.NoError(err)
	t.Equal(testhelper.TestGoogleSub, userRow.GoogleSub)
	t.Equal(testhelper.TestEmail, userRow.Email)
	t.Equal("John", userRow.FirstName)
	t.Equal("Doe", userRow.LastName)
	t.Equal(testhelper.TestPicture, userRow.ProfilePic)

	s.googleMock.AssertExpectations(s.T())
}

func (s *UsersServiceTestSuite) TestCompleteLogin_FailsWhenStateMissingInRedis() {
	t := require.New(s.T())
	ctx := context.Background()

	// syntactically valid state, but redis key missing
	state := mustEncodeOAuthState(t, users.OAuthState{CSRF: "csrf-missing", Return: testhelper.TestReturnFeed})
	_, _, err := s.svc.CompleteLogin(ctx, state, testhelper.TestCode)
	t.Error(err)
}

func (s *UsersServiceTestSuite) TestCompleteLogin_FailsWhenStateMismatch_DoesNotDeleteRedisKey() {
	t := require.New(s.T())
	ctx := context.Background()

	loginURL, err := s.svc.AuthViaGoogle(ctx, testhelper.TestReturnFeed)
	t.NoError(err)

	origState := mustExtractState(t, loginURL)
	st := mustDecodeOAuthState(t, origState)

	// tamper "return" in state, keep CSRF
	tampered := mustEncodeOAuthState(t, users.OAuthState{CSRF: st.CSRF, Return: "/other"})
	_, _, err = s.svc.CompleteLogin(ctx, tampered, testhelper.TestCode)
	t.Error(err)
	t.ErrorIs(err, svcerrs.ErrInvalidData)

	// key must still exist
	_, err = s.rdb.Get(ctx, testhelper.TestStateKeyPrefix+st.CSRF).Result()
	t.NoError(err)
}

func (s *UsersServiceTestSuite) TestCompleteLogin_DeletesStateOnSuccess() {
	t := require.New(s.T())
	ctx := context.Background()

	loginURL, err := s.svc.AuthViaGoogle(ctx, testhelper.TestReturnFeed)
	t.NoError(err)
	state := mustExtractState(t, loginURL)

	st := mustDecodeOAuthState(t, state)

	idToken := makeFakeIDToken(t, map[string]any{
		"sub":     testhelper.TestGoogleSub,
		"email":   testhelper.TestEmail,
		"name":    testhelper.TestNameTwoParts,
		"picture": testhelper.TestPicture,
	})
	s.googleMock.
		On("ExchangeCodeToToken", mock.Anything, mock.Anything).
		Return(&googleconn.ExchangeCodeToTokenResponse{IDToken: idToken}, nil).
		Once()

	_, _, err = s.svc.CompleteLogin(ctx, state, testhelper.TestCode)
	t.NoError(err)

	_, err = s.rdb.Get(ctx, testhelper.TestStateKeyPrefix+st.CSRF).Result()
	t.Error(err) // should be redis.Nil, but enough to assert it's gone
}

func (s *UsersServiceTestSuite) TestCompleteLogin_NameOnePart_DefaultsLovely() {
	t := require.New(s.T())
	ctx := context.Background()

	loginURL, err := s.svc.AuthViaGoogle(ctx, testhelper.TestReturnFeed)
	t.NoError(err)
	state := mustExtractState(t, loginURL)

	idToken := makeFakeIDToken(t, map[string]any{
		"sub":     testhelper.TestGoogleSub,
		"email":   testhelper.TestEmail,
		"name":    testhelper.TestNameOnePart,
		"picture": testhelper.TestPicture,
	})
	s.googleMock.
		On("ExchangeCodeToToken", mock.Anything, mock.Anything).
		Return(&googleconn.ExchangeCodeToTokenResponse{IDToken: idToken}, nil).
		Once()

	_, _, err = s.svc.CompleteLogin(ctx, state, testhelper.TestCode)
	t.NoError(err)

	var row modles.User
	err = s.db.WithContext(ctx).
		Model(&modles.User{}).
		Where("email = ?", testhelper.TestEmail).
		First(&row).Error
	t.NoError(err)
	t.Equal("Lovely", row.FirstName)
	t.Equal(testhelper.TestNameOnePart, row.LastName)
}

func (s *UsersServiceTestSuite) TestCompleteLogin_ExistingUser_UpdatesMetadata() {
	t := require.New(s.T())
	ctx := context.Background()

	// Seed user in DB using the repo model.
	seed := &modles.User{
		GoogleSub:  testhelper.TestGoogleSub,
		Email:      testhelper.TestEmail,
		FirstName:  "Old",
		LastName:   "Name",
		ProfilePic: testhelper.TestPicture,
	}
	err := s.db.
		WithContext(ctx).
		Model(&modles.User{}).
		Create(seed).Error
	t.NoError(err)
	t.NotZero(seed.ID)

	loginURL, err := s.svc.AuthViaGoogle(ctx, testhelper.TestReturnFeed)
	t.NoError(err)
	state := mustExtractState(t, loginURL)

	idToken := makeFakeIDToken(t, map[string]any{
		"sub":     testhelper.TestGoogleSub,
		"email":   testhelper.TestEmail,
		"name":    testhelper.TestUpdatedName,
		"picture": testhelper.TestUpdatedPicture,
	})
	s.googleMock.
		On("ExchangeCodeToToken", mock.Anything, mock.Anything).
		Return(&googleconn.ExchangeCodeToTokenResponse{IDToken: idToken}, nil).
		Once()

	_, _, err = s.svc.CompleteLogin(ctx, state, testhelper.TestCode)
	t.NoError(err)

	// still only one user and metadata updated
	var count int64
	err = s.db.
		WithContext(ctx).
		Model(&modles.User{}).
		Where("email = ?", testhelper.TestEmail).
		Count(&count).Error
	t.NoError(err)
	t.Equal(int64(1), count)

	var row modles.User
	err = s.db.
		WithContext(ctx).
		Model(&modles.User{}).
		Where("email = ?", testhelper.TestEmail).
		First(&row).Error
	t.NoError(err)
	t.Equal("Jane", row.FirstName)
	t.Equal("Roe", row.LastName)
	t.Equal(testhelper.TestUpdatedPicture, row.ProfilePic)
}

func (s *UsersServiceTestSuite) TestGetUser_ByEmail() {
	t := require.New(s.T())
	ctx := context.Background()

	loginURL, err := s.svc.AuthViaGoogle(ctx, testhelper.TestReturnFeed)
	t.NoError(err)
	state := mustExtractState(t, loginURL)

	idToken := makeFakeIDToken(t, map[string]any{
		"sub":     testhelper.TestGoogleSub,
		"email":   testhelper.TestEmail,
		"name":    testhelper.TestNameTwoParts,
		"picture": testhelper.TestPicture,
	})
	s.googleMock.
		On("ExchangeCodeToToken", mock.Anything, mock.Anything).
		Return(&googleconn.ExchangeCodeToTokenResponse{IDToken: idToken}, nil).
		Once()

	_, _, err = s.svc.CompleteLogin(ctx, state, testhelper.TestCode)
	t.NoError(err)

	u, err := s.svc.GetUser(ctx, dto.GetUser{Email: testhelper.TestEmail})
	t.NoError(err)
	t.Equal(testhelper.TestEmail, u.Email)
	t.NotZero(u.ID)
}
