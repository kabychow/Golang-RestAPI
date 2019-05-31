package models

import (
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"os"
	"regexp"
	"strings"
	"time"
)

type Token struct {
	UserId uint
	jwt.StandardClaims
}

type Account struct {
	ID uint `gorm:"primary_key"`
	Email string `gorm:"type:varchar(255);unique;not null"`
	Password string `json:",omitempty",gorm:"type:varchar(255);not null"`
	Token string `gorm:"-"`
	UpdatedAt time.Time
	CreatedAt time.Time
	DeletedAt *time.Time `json:",omitempty",sql:"index"`
}

func (account *Account) Validate() int {
	if regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`).MatchString(account.Email) &&
		len(account.Password) >= 6 && len(account.Password) <= 32 {
		var temp Account
		GetDB().Where(&Account{Email: account.Email}).First(&temp)
		account.Email = strings.TrimSpace(account.Email)
		if temp.Email == "" {
			return 0
		}
		return 409
	}
	return 422
}

func (account *Account) Create() (int, *Account) {
	if code := account.Validate(); code > 0 {
		return code, nil
	}
	password, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(password)
	GetDB().Create(account)
	if account.ID > 0 {
		GenerateToken(account)
		return 201, account
	}
	return 500, nil
}

func Login(email, password string) (int, *Account) {
	var account Account
	GetDB().Where(&Account{Email: email}).First(&account)
	if bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password)) == nil {
		GenerateToken(&account)
		return 200, &account
	}
	return 401, nil
}

func GenerateToken(account *Account) {
	account.Password = ""
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	account.Token, _ = token.SignedString([]byte(os.Getenv("token_password")))
}