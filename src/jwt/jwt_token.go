package jwt

import (
	"strconv"
	"time"

	"git.xenonstack.com/util/drive-portal/config"
	"git.xenonstack.com/util/drive-portal/src/types"
	"gopkg.in/dgrijalva/jwt-go.v3"
)

// defining structure for user defined jwt claims
type JWTClaims struct {
	Claim map[string]interface{} `json:"claim"`
	jwt.StandardClaims
}

// function for creating new signed token
func NewToken(claim map[string]interface{}, minutes int64) map[string]interface{} {
	mapd := make(map[string]interface{})
	claim["startTime"] = time.Now().Unix()
	// setting userdefined claims in adition to standard claims
	claims := JWTClaims{
		claim,
		jwt.StandardClaims{
			// setting expirattion time of token
			ExpiresAt: time.Now().Add(time.Second * time.Duration(minutes)).Unix(),
		},
	}

	// creating token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// signed token and signing with private key
	tokenStr, err := token.SignedString([]byte(config.TestPortalKey))
	if err != nil {
		return mapd
	}

	err = saveToken(tokenStr, claim["startTime"].(int64), claims.StandardClaims.ExpiresAt, claim)

	//set token str and expire time in map
	mapd["token"] = tokenStr
	mapd["start"] = claim["start"]
	mapd["expire"] = time.Unix(claims.StandardClaims.ExpiresAt, 0).Add(time.Duration(-120) * time.Second).Format(time.RFC3339)
	return mapd
}

func saveToken(token string, start, expire int64, claim map[string]interface{}) error {
	email, _ := claim["email"]
	drive, _ := claim["drive"]
	driveId, _ := strconv.Atoi(drive.(string))

	db := config.DB

	err := db.Create(&types.UserSession{
		Email:     email.(string),
		DriveId:   driveId,
		Token:     token,
		Expire:    expire - 120,
		Start:     start,
		TimeTaken: int(expire - start - 120),
	}).Error
	return err
}
