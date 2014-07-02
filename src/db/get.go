package db

func (d *DB) GetUserByName(name string) (*User, error) {
	var user User
	err := d.db.Where("name = ?", name).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *DB) GetUser(uid uint64) (*User, error) {
	var u User
	err := d.db.First(u, uid).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}
