package repo

const (
	TABLE_NAME_TOKEN  = "tokens"
	COL_ID_TOKEN      = "id"
	COL_TOKEN         = "token"
	COL_TOKEN_USER_ID = "user_id"
)

type TokenModelRepo struct {
	Id     uint64
	Token  string
	UserId uint64
}

type TokenRepo struct {
	Database *DataBase
}

func (t *TokenRepo) FindTokenByUserId(userId uint64) (*TokenModelRepo, error) {
	token := TokenModelRepo{}

	if err := t.Database.Db.QueryRow("SELECT "+COL_ID_TOKEN+", "+
		COL_TOKEN+", "+COL_TOKEN_USER_ID+
		" from "+TABLE_NAME_TOKEN+
		" where "+COL_TOKEN_USER_ID+" = $1", userId).
		Scan(&token.Id, &token.Token, &token.UserId); err != nil {
		return nil, err
	}

	return &token, nil
}

func (t *TokenRepo) DeleteTokenByUserId(userId uint64) (int64, error) {
	result, err := t.Database.Db.Exec("DELETE from "+TABLE_NAME_TOKEN+" where "+COL_TOKEN_USER_ID+" = $1", userId)

	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return affected, nil
}

func (t *TokenRepo) CreateToken(token *TokenModelRepo) (*TokenModelRepo, error) {
	if err := t.Database.Db.QueryRow("INSERT into "+TABLE_NAME_TOKEN+" ("+COL_TOKEN+", "+COL_TOKEN_USER_ID+") "+
		"VALUES ($1, $2) RETURNING "+COL_ID_TOKEN,
		token.Token,
		token.UserId).
		Scan(&token.Id); err != nil {
		return nil, err
	}

	return token, nil
}

func (t *TokenRepo) UpdateToken(token *TokenModelRepo) (*TokenModelRepo, error) {
	if err := t.Database.Db.QueryRow("UPDATE "+TABLE_NAME_TOKEN+
		" set "+COL_TOKEN+
		"=$1 where "+COL_TOKEN_USER_ID+
		"=$2 returning id, token, user_id",
		token.Token,
		token.UserId).Scan(&token.Id, &token.Token, &token.UserId); err != nil {
		return nil, err
	}

	return token, nil
}

func (t *TokenRepo) DeleteToken(token *TokenModelRepo) (*TokenModelRepo, error) {
	_, err := t.Database.Db.Exec("DELETE from "+TABLE_NAME_TOKEN+" where "+COL_TOKEN_USER_ID+" = $1", token.Id)

	if err != nil {
		return nil, err
	}

	return token, nil
}
