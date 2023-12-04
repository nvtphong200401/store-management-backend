package usecases

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"math"
	"net/mail"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-gomail/gomail"
	"github.com/nvtphong200401/store-management/pkg/handlers/models"
	"github.com/nvtphong200401/store-management/pkg/handlers/respository"
	"github.com/nvtphong200401/store-management/pkg/helpers"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCases struct {
	authRepo respository.AuthRepository
}

type Claims struct {
	UserID uint // Add additional claims as needed
	jwt.StandardClaims
}

func NewAuthUseCases(authRepo respository.AuthRepository) AuthUseCases {
	return AuthUseCases{
		authRepo: authRepo,
	}
}

func (a *AuthUseCases) Login(email string, password string) (int, gin.H) {
	user, err := a.authRepo.GetUserByEmail(email)
	if err != nil {
		return responseWithError(err)
	}
	err = comparePassword(string(user.Password), password)
	if err != nil {
		return responseWithError(err)
	}

	accessToken, err := generateToken(user.ID, helpers.ACCESS_TOKEN)
	if err != nil {
		return responseWithError(err)
	}
	refreshToken, err := generateToken(user.ID, helpers.REFRESH_TOKEN)
	if err != nil {
		return responseWithError(err)
	}

	return responseWithResult(map[string]string{
		helpers.ACCESS_TOKEN:  accessToken,
		helpers.REFRESH_TOKEN: refreshToken,
	})
}

func (a *AuthUseCases) RenewToken(refreshToken string) (int, gin.H) {
	claims, err := a.VerifyToken(refreshToken)
	if err != nil {
		return responseWithError(err)
	}
	if claims.Subject != helpers.REFRESH_TOKEN {
		return responseWithError(fmt.Errorf("invalid token type"))
	}
	newAccessToken, err := generateToken(claims.UserID, helpers.ACCESS_TOKEN)
	if err != nil {
		return responseWithError(err)
	}

	return responseWithResult(newAccessToken)
}

func (a *AuthUseCases) GetUserInfo(tokenString string) (*models.Employee, error) {
	claims, err := a.VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}
	return a.authRepo.GetUserByID(claims.UserID)
}

func (a *AuthUseCases) VerifyToken(tokenString string) (*Claims, error) {
	jwtKey := os.Getenv(helpers.JWT_KEY)
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

// Sign up use case
func (a *AuthUseCases) SignUp(email string, password string) (int, gin.H) {
	_, err := mail.ParseAddress(email)

	// check valid email pattern
	if err != nil {
		return responseWithError(fmt.Errorf("invalid email address"))
	}

	// create hashed password from raw password
	hashedPass, err := hashPassword(password)
	if err != nil {
		return responseWithError(err)
	}

	// generate verfication code
	verficationCode, err := generateVerificationCode()
	if err != nil {
		return responseWithError(err)
	}

	user := models.User{
		Email:            email,
		Password:         hashedPass,
		VerificationCode: verficationCode,
		Verified:         false,
	}

	// create user in database
	err = a.authRepo.SignUp(user)
	if err != nil {
		return responseWithError(err)
	}

	// send verification email
	err = sendVerificationEmail(email, verficationCode)

	if err != nil {
		return responseWithError(err)
	}

	return responseWithResult("User registered successfully. Check your email for verification instructions.")
}

func (a *AuthUseCases) RequestVerificationCode(user models.User) (int, gin.H) {
	const MaxRequests = 5

	window := 1 * time.Minute
	_, err := a.authRepo.IncrementRequestCounter(user.Email, window, MaxRequests)

	if err != nil {
		return responseWithError(err)
	}

	code, err := generateVerificationCode()
	if err != nil {
		return responseWithError(err)
	}

	err = sendVerificationEmail(user.Email, code)
	if err != nil {
		return responseWithError(err)
	}
	return responseWithResult("success")
}

func (a *AuthUseCases) VerifyCode(user models.User, code string) (int, gin.H) {
	expirationDuration := 15 * time.Minute
	maxVerficationAttempts := 5

	if time.Since(user.VerificationSentAt) > expirationDuration {
		return responseWithError(fmt.Errorf("verification code has expired"))
	}

	if user.VerificationAttempts >= maxVerficationAttempts {
		// Lock the account or apply a time delay before allowing the next attempt
		time.Sleep(calculateTimeDelay(user.VerificationAttempts))
		return responseWithError(fmt.Errorf("too many verification attempts, try again later"))
	}

	if user.VerificationCode != code {
		// Increment the number of failed attempts
		user.VerificationAttempts++
		err := a.authRepo.UpdateUser(user)
		if err != nil {
			return responseWithError(err)
		}

		return responseWithError(fmt.Errorf("invalid verification code"))
	}

	user.VerificationAttempts = 0
	user.Verified = true

	err := a.authRepo.UpdateUser(user)
	if err != nil {
		return responseWithError(err)
	}

	return responseWithResult("Verification success")
}

// Implement a time delay strategy, such as exponential backoff
// Adjust the parameters based on your security requirements
func calculateTimeDelay(attempts int) time.Duration {
	return time.Duration(math.Pow(2, float64(attempts))) * time.Second
}

func sendVerificationEmail(email, code string) error {
	// Retrieve email server details from environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	// Set up your email configuration (SMTP details)
	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Create the email message
	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Email Verification")
	m.SetBody("text/html", fmt.Sprintf("Your verification code is: %s", code))

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func generateVerificationCode() (string, error) {
	bytes := make([]byte, 6)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func comparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func generateToken(userID uint, tokenType string) (string, error) {
	var expiredTime int64
	jwtKey := os.Getenv(helpers.JWT_KEY)
	switch tokenType {
	case helpers.ACCESS_TOKEN:
		expiredTime = time.Now().Add(time.Hour * 24).Unix()
	case helpers.REFRESH_TOKEN:
		expiredTime = time.Now().Add(time.Hour * 24 * 7).Unix()
	}

	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredTime,
			IssuedAt:  time.Now().Unix(),
			Subject:   tokenType,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtKey))
}
