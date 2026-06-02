package security

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Ajulll22/payment-ai-assistant/pkg/jutil"
	"github.com/cristalhq/jwt/v4"
)

type UserClaims struct {
	jwt.RegisteredClaims
	LoginID string `json:"login_id"`
	UserID  int    `json:"jti"`
	// Email           string `json:"email"`
}

func RefreshJWT(token string, appKey string) (string, error) {
	key := []byte(appKey)

	// create a Verifier (HMAC in this example)
	verifier, err := jwt.NewVerifierHS(jwt.HS256, key)
	if err != nil {
		return "", err
	}

	// parse and verify a token
	tokenBytes := []byte(token)
	newToken, _ := jwt.Parse(tokenBytes, verifier)

	// or just verify it's signature
	err = verifier.Verify(newToken)
	if err != nil {
		return "", err
	}

	// get REGISTERED claims
	var newClaims UserClaims
	err = json.Unmarshal(newToken.Claims(), &newClaims)
	if err != nil {
		return "", err
	}

	iss := newClaims.Issuer
	sub := newClaims.Subject
	login_id := newClaims.LoginID
	jti := newClaims.UserID

	now := time.Now()

	signer, err := jwt.NewSignerHS(jwt.HS256, key)
	if err != nil {
		return "", err
	}
	builder := jwt.NewBuilder(signer)

	// create claims (you can create your own, see: Example_BuildUserClaims)
	claims := &UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    iss,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(3600 * time.Second)),
			Subject:   sub,
		},
		LoginID: login_id,
		UserID:  jti,
	}

	jwtToken, err := builder.Build(claims)
	if err != nil {
		return "", err
	}

	resToken := string(jwtToken.Bytes())

	return resToken, nil

}

func VerifyJWS(bodyRequest *jutil.OrderedMap2, postedToken, appKey string) error {
	b, _ := bodyRequest.MarshalJSON()

	key := []byte(appKey)
	signer, err := jwt.NewSignerHS(jwt.HS256, key)
	if err != nil {
		return err
	}
	builder := jwt.NewBuilder(signer)

	token, err := builder.Build(b)
	if err != nil {
		return err
	}

	strToken := string(token.Bytes())
	// get new signature
	s := strings.Split(strToken, ".")
	newSignedToken := s[2]

	// get posted signature
	s = strings.Split(postedToken, ".")
	if len(s) < 3 {
		return fmt.Errorf("Invalid JWS")
	}
	postedSignature := s[2]

	if postedSignature != newSignedToken {
		return errors.New("Invalid signature")
	}

	return nil

}

func GenerateJWT(bodyRequest map[string]interface{}, appName, appKey string) (string, UserClaims, error) {
	user_id := bodyRequest["user_id"].(int)
	login_id := bodyRequest["login_id"].(string)
	user_name := bodyRequest["user_name"].(string)

	now := time.Now()

	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    appName + "/login",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(3600 * time.Second)),
			Subject:   user_name,
		},
		LoginID: login_id,
		UserID:  user_id,
	}

	key := []byte(appKey)
	signer, err := jwt.NewSignerHS(jwt.HS256, key)
	if err != nil {
		return "", claims, err
	}
	builder := jwt.NewBuilder(signer)

	token, err := builder.Build(&claims)
	if err != nil {
		return "", claims, err
	}

	return string(token.Bytes()), claims, nil

}

func VerifyJWT(token string, appName, appKey string) (map[string]interface{}, error) {
	key := []byte(appKey)

	// create a Verifier (HMAC in this example)
	verifier, err := jwt.NewVerifierHS(jwt.HS256, key)
	if err != nil {
		return nil, err
	}

	// parse and verify a token
	tokenBytes := []byte(token)
	newToken, err := jwt.Parse(tokenBytes, verifier)
	if err != nil {
		return nil, err
	}

	// or just verify it's signature
	err = verifier.Verify(newToken)
	if err != nil {
		return nil, err
	}

	// get REGISTERED claims
	var newClaims UserClaims
	err = json.Unmarshal(newToken.Claims(), &newClaims)
	if err != nil {
		return nil, err
	}

	payload := map[string]interface{}{}
	payload["sub"] = newClaims.Subject
	payload["iss"] = newClaims.Issuer
	payload["login_id"] = newClaims.LoginID
	payload["jti"] = newClaims.UserID

	// check token expired time
	expiredStatus := newClaims.IsValidExpiresAt(time.Now())
	if !expiredStatus {
		return nil, fmt.Errorf("Token expired")
	}

	// check token issuer
	issuerStatus := newClaims.IsIssuer(appName + "/login")

	if !issuerStatus {
		return nil, fmt.Errorf("Token issuer invalid")
	}

	return payload, nil

}
