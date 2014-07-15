package db

func (d *DB) GetUserByName(name string) (*User, error) {
	var user User
	err := d.dbmap.SelectOne(&user,
		`select * from `+UserTable+` where name = ?`, name)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *DB) GetUser(uid uint64) (*User, error) {
	u, err := d.dbmap.Get(User{}, uid)
	if err != nil {
		return nil, err
	}
	user := u.(User)
	return &user, nil
}
