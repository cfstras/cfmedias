package db

import (
	"code.google.com/p/go.crypto/pbkdf2"
	"core"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"errrs"
	"fmt"
	"math/big"
	"regexp"
	"util"
)

const (
	// Settings for the stored hash, PBKDF2
	SaltSize   = 128
	Iterations = 10000
	KeySize    = 512

	// only characters easily typeable on a mobile
	AuthTokenChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// also, not too many of them
	AuthTokenSize = 16
)

func (db *DB) initLogin(c core.Core) {
	c.RegisterCommand(core.Command{
		[]string{"create_user"},
		"Creates a user in the database",
		map[string]string{
			"name":  "Username",
			"email": "E-Mail",
			"auth_level": fmt.Sprintf("User Rank: Guest(%d), User(%d), "+
				"Admin(%d), Root(%d)", core.AuthGuest, core.AuthUser,
				core.AuthAdmin, core.AuthRoot),
			"password": "Password"},
		core.AuthAdmin,
		func(ctx core.CommandContext) core.Result {
			args := ctx.Args
			var err error
			name, err := util.GetArg(args, "name", true, err)
			email, err := util.GetArg(args, "email", true, err)
			authLevelS, err := util.GetArg(args, "auth_level", true, err)
			password, err := util.GetArg(args, "password", true, err)

			authLevelI, err := util.CastUint(authLevelS, err)

			if err != nil {
				return core.Result{Status: core.StatusError, Error: err}
			}
			authLevel := core.AuthLevel(*authLevelI)

			if authLevel >= ctx.AuthLevel {
				return core.ResultByError(errrs.New("You cannot create a user" +
					" with that level!"))
			}

			user, err := db.CreateUser(*name, *email, authLevel, *password)
			if err == nil {
				return core.Result{Status: core.StatusOK, Results: []interface{}{user}}
			}
			return core.Result{Status: core.StatusError, Error: err}
		}})

	c.RegisterCommand(core.Command{
		[]string{"login"},
		"Logs in with user/password and returns the auth token",
		map[string]string{
			"name":     "Username",
			"password": "Password"},
		core.AuthGuest,
		func(ctx core.CommandContext) core.Result {
			args := ctx.Args
			var err error
			name, err := util.GetArg(args, "name", true, err)
			password, err := util.GetArg(args, "password", true, err)

			if err != nil {
				return core.ResultByError(err)
			}

			success, authToken, err := db.Login(*name, *password)
			if err == nil && success {
				return core.Result{Status: core.StatusOK, Results: []interface{}{
					map[string]string{"auth_token": authToken},
				}}
			}
			if err == nil {
				return core.Result{Status: core.StatusError,
					Error: errrs.New("Wrong username or password")}
			}
			return core.Result{Status: core.StatusError, Error: err}
		}})

	c.RegisterCommand(core.Command{
		[]string{"change_authtoken", "logout"},
		`Changes the authentication token of the user, thereby logging out all
clients`,
		map[string]string{
			"name": `(Optional) username of the user to log out. Leave empty to
use current user`,
		},
		core.AuthUser,
		func(ctx core.CommandContext) core.Result {
			args := ctx.Args
			var err error
			name, err := util.GetArg(args, "name", false, err)
			if err != nil {
				return core.ResultByError(err)
			}
			var user *User
			if name != nil {
				if ctx.AuthLevel < core.AuthAdmin {
					return core.ResultByError(core.ErrorNotAllowed)
				}
				user, err = db.GetUserByName(*name)
			} else {
				if ctx.UserId == nil {
					return core.ResultByError(core.ErrorNotLoggedIn)
				}
				user, err = db.GetUser(*ctx.UserId)
			}
			if user == nil {
				return core.ResultByError(core.ErrorUserNotFound)
			}
			_, err = db.ChangeAuthToken(user)
			return core.ResultByError(err)
		}})
}

func (db *DB) Login(name string, password string) (success bool,
	authToken string, err error) {
	user := User{}
	err = db.dbmap.SelectOne(&user,
		`select * from `+UserTable+` where name = ?`, name)
	if err != nil {
		return false, "", errrs.New(err.Error())
	}
	if user.Id == 0 {
		return false, "", nil
	}
	salt := user.Password[:SaltSize]
	expected := user.Password[SaltSize:]
	hashedPassword := HashPassword([]byte(password), salt)
	if 1 == subtle.ConstantTimeCompare(expected, hashedPassword) {
		//TODO add option to create new authtoken to logout all clients
		return true, user.AuthToken, nil
	}
	return false, "", nil
}

// Checks a given authentication token agains the database.
// On success, returns the permission level of the user and their ID
// On failure, returns (AuthGuest, nil, error)
func (db *DB) Authenticate(authtoken string) (core.AuthLevel, *uint64, error) {
	invalidErr := errrs.New("auth token invalid")
	user := User{}

	err := db.dbmap.SelectOne(&user,
		`select * from `+UserTable+` where auth_token = ?`, authtoken)
	if err != nil || user.Id == 0 { // not found
		return core.AuthGuest, nil, invalidErr
	}
	return user.AuthLevel, &user.Id, nil
}

func (db *DB) CreateUser(name string, email string, authLevel core.AuthLevel,
	password string) (*User, error) {

	if !IsSafe(name, TypeUsername) {
		return nil, errrs.New(`Username is invalid. Can contain a-z, 0-9,
underscore and dash, minimum 3 characters.`)
	}
	if !IsSafe(email, TypeEmail) {
		return nil, errrs.New("EMail is invalid.")
	}
	if !IsSafe(password, TypePassword) {
		return nil, errrs.New(`Password is invalid. Must be at least 6 and at
most 128 characters.`)
	}

	// check for unique
	num, err := db.dbmap.SelectInt(`select count(*) from `+UserTable+
		` where name = ?`, name)

	if err != nil {
		return nil, errrs.New("Error checking username: " + err.Error())
	}
	if num != 0 {
		return nil, errrs.New("User already exists!")
	}

	user := User{Name: name, AuthLevel: authLevel}

	// create authtoken
	user.AuthToken, err = db.makeAuthToken()
	if err != nil {
		return nil, err
	}

	// hash password
	user.Password = MakePassword([]byte(password))

	// store user
	err = db.dbmap.Insert(&user)

	if err != nil {
		return nil, errrs.New("Could not insert user: " + err.Error())
	}

	return &user, nil
}

// changes auth token and saves user
// returns: new auth token and/or error.
func (d *DB) ChangeAuthToken(user *User) (string, error) {
	var err error
	user.AuthToken, err = d.makeAuthToken()
	if err != nil {
		return "", err
	}
	d.dbmap.Update(user)
	return user.AuthToken, nil
}

func (d *DB) makeAuthToken() (string, error) {
	token := make([]rune, AuthTokenSize)
	for i := 0; i < AuthTokenSize; i++ {
		max := big.NewInt(int64(len(AuthTokenChars)))
		ind, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		token[i] = rune(AuthTokenChars[ind.Int64()])
	}
	return string(token), nil
}

type StringType uint

const (
	TypePassword StringType = iota
	TypeUsername
	TypeEmail
)

var (
	passwordRegex = regexp.MustCompile(`.{6,128}`)
	usernameRegex = regexp.MustCompile(`[a-z0-9_-]{3,64}`)
	emailRegex    = regexp.MustCompile(`[a-zA-Z0-9._-]+(\+[a-zA-Z0-9._-]+)?@[a-zA-Z0-9._-]+`)
	//TODO refine email regex (or use library)
)

func IsSafe(s string, t StringType) bool {
	switch t {
	case TypePassword:
		return passwordRegex.MatchString(s)
	case TypeUsername:
		return usernameRegex.MatchString(s)
	case TypeEmail:
		return emailRegex.MatchString(s)
	}
	return false
}

//TODO test MakePassword
func MakePassword(pass []byte) []byte {
	salt := make([]byte, SaltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return nil
	}

	key := HashPassword(pass, salt)
	if key == nil {
		return nil
	}
	hashed := make([]byte, KeySize+SaltSize)
	subtle.ConstantTimeCopy(1, hashed[:SaltSize], salt)
	subtle.ConstantTimeCopy(1, hashed[SaltSize:], key)

	return hashed
}

func HashPassword(pass, salt []byte) []byte {
	return pbkdf2.Key(pass, salt, Iterations, KeySize, sha256.New)
}
