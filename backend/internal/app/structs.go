package app

import "github.com/dgrijalva/jwt-go"

type TokenResponse struct {
	Access_token string `json:"access_token"`
	Token_type   string `json:"token_type"`
	Expires_in   string `json:"expires_in"`
	Scope        string `json:"scope"`
	Nick         string `json:"nick"`
	Uid          string `json:"uid"`
	ISS          string `json:"iss"`
	JTI          string `json:"jti"`
}

type TokenClaims struct {
	Nick     string   `json:"nick"`
	Uid      string   `json:"uid"`
	UserName string   `json:"user_name"`
	Scope    []string `json:"scope"`
	ISS      string   `json:"iss"`
	Exp      int      `json:"exp"`
	JTI      string   `json:"jti"`
	ClientId string   `json:"client_id"`
	jwt.StandardClaims
}

type SvEn struct {
	Sv string `json:"sv"`
	En string `json:"en"`
}

type Group struct {
	Id              string `json:"id"`
	BecomesActive   int    `json:"becomesActive"`
	BecomesInactive int    `json:"becomesInactive"`
	Description     SvEn   `json:"description"`
	Email           string `json:"email"`
	Function        SvEn   `json:"function"`
	Name            string `json:"name"`
	PrettyName      string `json:"prettyName"`
	AvatarURL       string `json:"avatarURL"`
	Active          bool   `json:"active"`
}

type SuperGroup struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	PrettyName string `json:"prettyName"`
	Type       string `json:"type"`
	Email      string `json:"email"`
}

type Relationship struct {
	SuperGroup SuperGroup `json:"superGroup"`
	Group      Group      `json:"group"`
}

type User struct {
	Id                    string         `json:"id"`
	Cid                   string         `json:"cid"`
	Nick                  string         `json:"nick"`
	FirstName             string         `json:"firstName"`
	LastName              string         `json:"lastName"`
	Email                 string         `json:"email"`
	Phone                 string         `json:"phone"`
	Language              string         `json:"language"`
	AvatarUrl             string         `json:"avatarUrl"`
	Gdpr                  bool           `json:"gdpr"`
	UserAgreement         bool           `json:"userAgreement"`
	AccountLocked         bool           `json:"accountLocked"`
	AcceptanceYear        int            `json:"acceptanceYear"`
	Authorities           []interface{}  `json:"authorities"`
	Activated             bool           `json:"activated"`
	Username              string         `json:"username"`
	Enabled               bool           `json:"enabled"`
	AccountNonLocked      bool           `json:"accountNonLocked"`
	AccountNonExpired     bool           `json:"accountNonExpired"`
	CredentialsNonExpired bool           `json:"credentialsNonExpired"`
	WebsiteURLs           string         `json:"websiteURLs"`
	Relationships         []Relationship `json:"relationships"`
}