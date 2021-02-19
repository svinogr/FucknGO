package repo

import (
	"FucknGO/db/user"
)

const (
	TABLE_NAME_USERS = "users"
	COL_ID_USER      = "id"
	COL_NAME         = "user_name"
	COL_PASSWORD     = "password"
	COL_EMAIL        = "email"
)

type UserRepo struct {
	Database *DataBase
}

func (u *UserRepo) CreateUser(user *user.UserModel) (*user.UserModel, error) {
	if err := u.Database.Db.QueryRow("INSERT into "+TABLE_NAME_USERS+" ("+COL_NAME+", "+COL_PASSWORD+", "+COL_EMAIL+") VALUES ($1, $2, $3) RETURNING "+COL_ID_USER,
		user.Name,
		user.Password,
		user.Email).
		Scan(&user.Id); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepo) UpdateUser(user *user.UserModel) (*user.UserModel, error) {
	if err := u.Database.Db.QueryRow("INSERT into "+TABLE_NAME_USERS+" ("+COL_NAME+", "+COL_PASSWORD+", "+COL_EMAIL+")  VALUES ($1, $2, $3) where "+COL_ID_USER+" = 4$",
		user.Name,
		user.Password,
		user.Email,
		user.Id).Scan(&user.Name, &user.Password, &user.Email); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepo) FindUserById(id uint64) (*user.UserModel, error) {
	user := user.UserModel{}

	if err := u.Database.Db.QueryRow("SELECT "+COL_ID_USER+","+COL_NAME+","+COL_NAME+","+COL_EMAIL+" from "+TABLE_NAME_USERS+" where "+COL_ID_USER+" = 1$",
		id).
		Scan(&user.Name, &user.Password, &user.Email); err != nil {

		return nil, err
	}

	return &user, nil
}

func (u *UserRepo) DeleteUser(user *user.UserModel) (*user.UserModel, error) {
	if err := u.Database.Db.QueryRow("DELETE from "+TABLE_NAME_USERS+" where "+COL_ID_USER+" = 1$", user.Id).
		Err(); err != nil {

		return nil, err
	}

	return user, nil
}
