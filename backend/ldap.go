package main

import (
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"errors"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

type User struct{
	Username string		`json:"cid"`
	Nick string			`json:"nick"`
	MemberOfPRIT bool	`json:"prit"`
	jwt.StandardClaims
}

func login_ldap(username string, pass string) (User, error){
	conn, err := ldap.DialURL("ldap://kamino.chalmers.it")
	if err != nil {
		return User{}, errors.New("Could not connect to ldap")
	}
	defer conn.Close()

	err = conn.Bind(fmt.Sprintf("uid=%s,ou=people,dc=chalmers,dc=it", username), pass) //Password for chalmers.it
	if err != nil {
		return User{}, errors.New("Username and password did not match")
	}

	user := User{Username: username}
	user.MemberOfPRIT = memberOfFkit(conn, username, "digit")

	s, er := conn.Search(ldap.NewSearchRequest(
		"ou=people,dc=chalmers,dc=it",
		ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, false,
		fmt.Sprintf("(uid=%s)", username),
		[]string{"nickname"},
		nil))

	if er != nil {
		return user, errors.New("Could not find nickname of user")
	}

	user.Nick = s.Entries[0].GetAttributeValue("nickname")

	return user, nil
}

func search(conn *ldap.Conn, dc string, filter string, attributes []string) (*ldap.SearchResult, error){
	return conn.Search(ldap.NewSearchRequest(
		dc,
		ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, false,
		filter,
		attributes,
		nil))
}

func memberOfFkit(conn *ldap.Conn, cid string, committee string) bool {
	members := membersOf(conn, membersOf(conn, []string{fmt.Sprintf("cn=%s,ou=%s,ou=fkit,ou=groups,dc=chalmers,dc=it", committee, committee)}))

	for i := 0; i < len(members); i++ {
		if strings.Contains(members[i], cid) {
			return true
		}
	}

	return false
}

func membersOf(conn *ldap.Conn, groups []string) []string {
	members := []string{}

	for _, group := range groups {
		nameAndDir := strings.SplitN(group, ",", 2)

		r, err := search(conn, nameAndDir[1], fmt.Sprintf("(%s)", nameAndDir[0]), []string{"member"})
		if err == nil {
			members = append(members, r.Entries[0].GetAttributeValues("member")...)
		}
	}

	return members
}