package repo

import (
	"FucknGO/db/user"
)

const (
	TABLE_NAME_TOKEN = "tokens"
	COL_ID_TOKEN     = "id"
	COL_TOKEN        = "token"
	COL_USER_ID      = "user_id"
)

type TokenRepo struct {
	Database *DataBase
}

func (t *TokenRepo) FindTokenByUserId(userId uint64) (*user.TokenModelRepo, error) {
	token := user.TokenModelRepo{}
	if err := t.Database.Db.QueryRow("SELECT "+COL_ID_TOKEN+", "+COL_TOKEN+", "+COL_ID_USER+" from "+TABLE_NAME_TOKEN+" where "+COL_USER_ID+" = $1", userId).
		Scan(token.Id, token.Token, token.UserId); err != nil {
		return nil, err
	}

	return &token, nil
}

func (t *TokenRepo) CreateToken(token *user.TokenModelRepo) (*user.TokenModelRepo, error) {
	if err := t.Database.Db.QueryRow("INSERT into "+TABLE_NAME_TOKEN+" ("+COL_TOKEN+", "+COL_ID_USER+") VALUES ($1, $2) RETURNING "+COL_ID_TOKEN,
		token.Token,
		token.UserId).
		Scan(&token.Id, &token.Token, &token.UserId); err != nil {
		return nil, err
	}

	return token, nil
}

func (t *TokenRepo) UpdateToken(token *user.TokenModelRepo) (*user.TokenModelRepo, error) {
	if err := t.Database.Db.QueryRow("INSERT into "+TABLE_NAME_TOKEN+" ("+COL_TOKEN+", "+COL_ID_USER+") VALUES ($1, $2) where "+COL_ID_TOKEN+" = 3$",
		token.Token,
		token.UserId,
		token.Id).Scan(&token.Id, &token.Token, &token.UserId); err != nil {
		return nil, err
	}

	return token, nil
}
