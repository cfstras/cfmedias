package db

import (
	"code.google.com/p/go.crypto/pbkdf2"
	"core"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"errrs"
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
		core.AuthAdmin,
		func(ctx core.CommandContext) core.Result {
			args := ctx.Args
			var err error
			name, err := getArg(args, "name", true, err)
			email, err := getArg(args, "email", true, err)
			authLevelS, err := getArg(args, "auth_level", true, err)
			password, err := getArg(args, "password", true, err)

			authLevelI, err := castUint(authLevelS, err)

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
		core.AuthGuest,
		func(ctx core.CommandContext) core.Result {
			args := ctx.Args
			var err error
			name, err := getArg(args, "name", true, err)
			password, err := getArg(args, "password", true, err)

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

func (db *DB) CreateUser(name string, email string, authLevel core.AuthLevel,
	password string) (*User, error) {
	//TODO validate username format
	//TODO validate password format
	//TODO validate email format

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
