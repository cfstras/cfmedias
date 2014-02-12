package db

import (
	"code.google.com/p/go.crypto/pbkdf2"
	"core"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"errrs"
	"regexp"
	"util"
)

const (
	SaltSize      = 128
	Iterations    = 10000
	AuthTokenSize = 128
	KeySize       = 512
)

func (db *DB) initLogin(c core.Core) {
	c.RegisterCommand(core.Command{
		[]string{"createuser"},
		"Creates a user in the database",
		map[string]string{
			"name":      "Username",
			"email":     "E-Mail",
			"authLevel": "User Rank: Guest(0), User(1), Admin(2), Root(3)",
			"password":  "Password"},
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
				return core.Result{Status: core.StatusError, Error: err}
			}

			success, authToken, err := db.Login(*name, *password)
			if err == nil && success {
				return core.Result{Status: core.StatusOK, Results: []interface{}{authToken}}
			}
			if err == nil {
				return core.Result{Status: core.StatusError,
					Error: errrs.New("Wrong username or password")}
			}
			return core.Result{Status: core.StatusError, Error: err}
		}})
}

func (db *DB) Login(name string, password string) (success bool, authToken []byte,
	err error) {
	user := User{}
	err = db.dbmap.SelectOne(&user, `select * from `+UserTable+` where name = ?`, name)
	if err != nil {
		return false, nil, errrs.New(err.Error())
	}
	if user.Id == 0 {
		return false, nil, nil
	}
	salt := user.Password[:SaltSize]
	expected := user.Password[SaltSize:]
	hashedPassword := HashPassword([]byte(password), salt)
	if 1 == subtle.ConstantTimeCompare(expected, hashedPassword) {
		//TODO create new auth token on login to logout old instances
		return true, user.AuthToken, nil
	}
	return false, nil, nil
}

func (db *DB) Authenticate(authtoken []byte) (core.AuthLevel, error) {
	invalidErr := errrs.New("auth token invalid")
	user := User{}

	err := db.dbmap.SelectOne(&user, `select * from `+UserTable+` where auth_token = ?`, authtoken)
	if err != nil {
		return core.AuthGuest, invalidErr
	}
	if user.Id == 0 {
		// not found
		return core.AuthGuest, invalidErr
	}
	return user.AuthLevel, nil
}

func (db *DB) CreateUser(name string, email string, authLevel core.AuthLevel,
	password string) (*User, error) {

	if !IsSafe(name, TypeUsername) {
		return nil, errrs.New("Username is invalid. Can contain a-z, 0-9, underscore and dash.")
	}
	if !IsSafe(email, TypeEmail) {
		return nil, errrs.New("EMail is invalid.")
	}
	if !IsSafe(password, TypePassword) {
		return nil, errrs.New("Password is invalid. Must be at least 6 and at most 128 characters.")
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
	user.AuthToken = make([]byte, AuthTokenSize)
	_, err = rand.Read(user.AuthToken)
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

type StringType uint

const (
	TypePassword StringType = iota
	TypeUsername
	TypeEmail
)

var (
	passwordRegex = regexp.MustCompile(`.{6,128}`)
	usernameRegex = regexp.MustCompile(`[a-z0-9_-]{4,64}`)
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
