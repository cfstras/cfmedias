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

			user, err := db.CreateUser(*name, authLevel, *password)
			if err == nil {
				return core.Result{Status: core.StatusOK, Results: []interface{}{user}}
			} else {
				return core.Result{Status: core.StatusError, Error: err}
			}
		}})
}

func (db *DB) CreateUser(name string, authLevel core.AuthLevel, password string) (
	*User, error) {
	//TODO validate username format
	//TODO validate password format

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
	user.Password = HashPassword(password)

	// store user
	err = db.dbmap.Insert(&user)

	if err != nil {
		return nil, errrs.New("Could not insert user: " + err.Error())
	}

	return &user, nil
}

//TODO test HashPassword
func HashPassword(pass string) []byte {
	salt := make([]byte, SaltSize)
	_, err := rand.Read(salt)
	if err != nil {
		return nil
	}

	key := pbkdf2.Key([]byte(pass), salt, Iterations, KeySize, sha256.New)
	if key == nil {
		return nil
	}
	hashed := make([]byte, KeySize+SaltSize)
	subtle.ConstantTimeCopy(1, hashed[:SaltSize], salt)
	subtle.ConstantTimeCopy(1, hashed[SaltSize:], key)

	return hashed
}
