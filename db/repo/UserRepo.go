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
	db *DataBase
}

func (u *UserRepo) CreateUser(user *user.UserModelRepo) (*user.UserModelRepo, error) {
	//defer u.db.CloseDataBase()
	if err := u.db.Db.QueryRow("INSERT into "+TABLE_NAME_USERS+" ("+COL_NAME+", "+COL_PASSWORD+", "+COL_EMAIL+") VALUES ($1, $2, $3) RETURNING "+COL_ID_USER,
		user.Name,
		user.Password,
		user.Email).
		Scan(&user.Id); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepo) openAndCloseDb() {
	u.db.OpenDataBase()
	//defer u.db.CloseDataBase()
}

func (u *UserRepo) UpdateUser(user *user.UserModelRepo) (*user.UserModelRepo, error) {
	err := u.db.Db.QueryRow("UPDATE "+TABLE_NAME_USERS+" set "+
		COL_NAME+"= $1, "+
		COL_PASSWORD+"= $2, "+
		COL_EMAIL+"= $3"+
		" WHERE "+COL_ID_USER+"=$4 returning id, user_name, password, email",
		user.Name,
		user.Password,
		user.Email,
		user.Id).Scan(&user.Id, &user.Name, &user.Password, &user.Email)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepo) FindUserById(id uint64) (*user.UserModelRepo, error) {
	user := user.UserModelRepo{}

	if err := u.db.Db.QueryRow("SELECT "+COL_ID_USER+", "+COL_NAME+", "+COL_PASSWORD+", "+COL_EMAIL+" from "+TABLE_NAME_USERS+" where "+COL_ID_USER+"=$1",
		id).
		Scan(&user.Id, &user.Name, &user.Password, &user.Email); err != nil {

		return nil, err
	}

	return &user, nil
}

func (u *UserRepo) DeleteUser(user *user.UserModelRepo) (*user.UserModelRepo, error) {
	_, err := u.db.Db.Exec("DELETE from "+TABLE_NAME_USERS+" where "+COL_ID_USER+" = $1", user.Id)

	if err != nil {
		return nil, err
	}

	return user, nil
}
