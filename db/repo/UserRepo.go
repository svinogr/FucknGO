package repo

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	TABLE_NAME_USERS  = "users"
	COL_USER_ID       = "id"
	COL_USER_NAME     = "user_name"
	COL_USER_PASSWORD = "password"
	COL_USER_EMAIL    = "email"
	COL_USER_TYPE     = "type"
)

type TypeUser int

const (
	Admin = iota
	Client
	Shop
)

type UserModelRepo struct {
	Id       uint64
	Name     string
	Password string
	Email    string
	Type     TypeUser
}

type UserRepo struct {
	db *DataBase
}

func (u *UserRepo) CreateUser(user *UserModelRepo) (*UserModelRepo, error) {
	defer u.db.CloseDataBase()
	if err := u.db.Db.QueryRow("INSERT into "+TABLE_NAME_USERS+" ("+COL_USER_NAME+", "+COL_USER_PASSWORD+", "+COL_USER_EMAIL+","+COL_USER_TYPE+") VALUES ($1, $2, $3, $4) RETURNING "+COL_USER_ID,
		user.Name,
		user.Password,
		user.Email,
		user.Type).
		Scan(&user.Id); err != nil {
		return nil, err
	}

	return user, nil
}

/*func (u *UserRepo) openAndCloseDb() {
	u.db.OpenDataBase()
	//defer u.db.CloseDataBase()
}*/

func (u *UserRepo) UpdateUser(user *UserModelRepo) (*UserModelRepo, error) {
	defer u.db.CloseDataBase()
	err := u.db.Db.QueryRow("UPDATE "+TABLE_NAME_USERS+" set "+
		COL_USER_NAME+"=$1, "+
		COL_USER_PASSWORD+"=$2, "+
		COL_USER_EMAIL+"=$3 "+
		COL_USER_TYPE+"=$4 "+
		"WHERE "+COL_USER_ID+"=$5 returning id, user_name, password, email, type",
		user.Name,
		user.Password,
		user.Email,
		user.Type,
		user.Id).Scan(&user.Id, &user.Name, &user.Password, &user.Email, &user.Type)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepo) FindUserById(id uint64) (*UserModelRepo, error) {
	defer u.db.CloseDataBase()
	user := UserModelRepo{}
	//TODO переделать метод нафига еще один оьект user
	if err := u.db.Db.QueryRow("SELECT "+COL_USER_ID+", "+COL_USER_NAME+", "+COL_USER_PASSWORD+", "+COL_USER_EMAIL+", "+COL_USER_TYPE+" from "+TABLE_NAME_USERS+" where "+COL_USER_ID+"=$1",
		id).
		Scan(&user.Id, &user.Name, &user.Password, &user.Email, &user.Type); err != nil {

		return nil, err
	}

	return &user, nil
}

func (u *UserRepo) DeleteUser(user *UserModelRepo) (*UserModelRepo, error) {
	defer u.db.CloseDataBase()
	_, err := u.db.Db.Exec("DELETE from "+TABLE_NAME_USERS+" where "+COL_USER_ID+" = $1", user.Id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepo) FindUserByEmail(email string) (*UserModelRepo, error) {
	defer u.db.CloseDataBase()
	user := UserModelRepo{}

	if err := u.db.Db.QueryRow("SELECT "+COL_USER_ID+", "+COL_USER_NAME+", "+COL_USER_PASSWORD+", "+COL_USER_EMAIL+", "+COL_USER_TYPE+" from "+TABLE_NAME_USERS+" where "+COL_USER_EMAIL+"=$1",
		email).
		Scan(&user.Id, &user.Name, &user.Password, &user.Email, &user.Type); err != nil {

		return nil, err
	}

	return &user, nil
}

func (u *UserRepo) FindUserByName(name string) (*UserModelRepo, error) {
	defer u.db.CloseDataBase()
	user := UserModelRepo{}

	if err := u.db.Db.QueryRow("SELECT "+COL_USER_ID+", "+COL_USER_NAME+", "+COL_USER_PASSWORD+", "+COL_USER_EMAIL+", "+COL_USER_TYPE+" from "+TABLE_NAME_USERS+" where "+COL_USER_NAME+"=$1",
		name).
		Scan(&user.Id, &user.Name, &user.Password, &user.Email, &user.Type); err != nil {

		return nil, err
	}

	return &user, nil
}

func (u *UserRepo) FindAllUser() (*[]UserModelRepo, error) {
	defer u.db.CloseDataBase()
	userList := []UserModelRepo{}

	row, err := u.db.Db.Query("SELECT " + COL_USER_ID + ", " + COL_USER_NAME + ", " + COL_USER_PASSWORD + ", " + COL_USER_EMAIL + ", " + COL_USER_TYPE + " from " + TABLE_NAME_USERS)

	if err != nil {
		return &userList, err
	}

	for row.Next() {
		user := UserModelRepo{}
		err := row.Scan(&user.Id, &user.Name, &user.Password, &user.Email, &user.Type)

		if err != nil {
			return nil, err
		}
		userList = append(userList, user)
	}

	return &userList, nil
}

// validUser gets valid user by email and password
func (u *UserRepo) GetValidUser(user UserModelRepo) (*UserModelRepo, error) {
	defer u.db.CloseDataBase()
	uBemail, err := u.FindUserByEmail(user.Email)

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(uBemail.Password), []byte(user.Password))

	if err != nil {
		return nil, err
	}

	return uBemail, nil
}
