package repo

import (
	"FucknGO/db"
	"FucknGO/internal/model/user"
)

const (
	TABLE_NAME   = "users"
	COL_ID       = "id"
	COL_NAME     = "user_name"
	COL_PASSWORD = "password"
	COL_EMAIL    = "email"
)

type UserRepo struct {
	Database *db.DataBase
}

func (u *UserRepo) CreateUser(user *user.UserModel) (*user.UserModel, error) {
	if err := u.Database.Db.QueryRow("INSERT into "+TABLE_NAME+" ("+COL_NAME+", "+COL_PASSWORD+", "+COL_EMAIL+") VALUES ($1, $2, $3) RETURNING "+COL_ID,
		user.Name,
		user.Password,
		user.Email).
		Scan(&user.Id); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepo) UpdateUser(user *user.UserModel) (*user.UserModel, error) {
	if err := u.Database.Db.QueryRow("INSERT into "+TABLE_NAME+" ("+COL_NAME+", "+COL_PASSWORD+", "+COL_EMAIL+")  VALUES ($1, $2, $3) where "+COL_ID+" = 4$",
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

	if err := u.Database.Db.QueryRow("SELECT "+COL_ID+","+COL_NAME+","+COL_NAME+","+COL_EMAIL+" from "+TABLE_NAME+" where "+COL_ID+" = 1$",
		id).
		Scan(&user.Name, &user.Password, &user.Email); err != nil {

		return nil, err
	}

	return &user, nil
}

func (u *UserRepo) DeleteUser(user *user.UserModel) (*user.UserModel, error) {
	if err := u.Database.Db.QueryRow("DELETE from "+TABLE_NAME+" where "+COL_ID+" = 1$", user.Id).
		Err(); err != nil {

		return nil, err
	}

	return user, nil
}
