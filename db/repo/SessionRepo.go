package repo

import (
	"time"
)

const (
	TABLE_REFRESH_SESSIONS  = "refresh_sessions"
	COL_ID_REFRESH_SESSIONS = "id"
	COL_USER_ID_REF         = "user_id"
	COL_REFRESH_TOKEN       = "refresh_token"
	COL_USERAGENT           = "user_agent"
	COL_FINGERPRINT         = "fingerprint"
	COL_IP                  = "ip"
	COL_EXPIREDIN           = "expires_in"
	COL_CREATED_AT          = "created_at"
)

type SessionModelRepo struct {
	Id           uint64
	UserId       uint64
	RefreshToken string
	UserAgent    string
	Fingerprint  string
	Ip           string
	ExpireIn     int64
	CreatedAt    time.Time
}

type SessionRepo struct {
	Database *DataBase
}

/*func (t *SessionRepo) DeleteSessionByUserId(userId uint64) (int64, error) {
	result, err := t.Database.Db.Exec("DELETE from "+TABLE_REFRESH_SESSIONS+"where "+COL_USER_ID+" = $1", userId)

	if err != nil {
		return 0, err
	}
	affected, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return affected, nil
}*/

func (t *SessionRepo) CreateSession(session *SessionModelRepo) (*SessionModelRepo, error) {
	if err := t.Database.Db.QueryRow("INSERT into "+TABLE_REFRESH_SESSIONS+" ("+
		COL_USER_ID_REF+
		", "+COL_REFRESH_TOKEN+
		", "+COL_USERAGENT+
		", "+COL_FINGERPRINT+
		", "+COL_IP+
		", "+COL_EXPIREDIN+") "+
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING "+COL_ID_REFRESH_SESSIONS+", "+COL_CREATED_AT,
		session.UserId,
		session.RefreshToken,
		session.UserAgent,
		session.Fingerprint,
		session.Id,
		session.ExpireIn).
		Scan(&session.Id, &session.CreatedAt); err != nil {
		return nil, err
	}

	return session, nil
}

/*
func (t *SessionRepo) UpdateSession(session *SessionModelRepo) (*SessionModelRepo, error) {
	if err := t.Database.Db.QueryRow("UPDATE "+TABLE_NAME_TOKEN+
		" set "+
		COL_REFRESH_TOKEN+" = $1, "+
		COL_USERAGENT+" = $2, "+
		COL_FINGERPRINT+" = $3, "+
		COL_IP+" = $4, "+
		COL_EXPIREDIN+" = $5, "+
		COL_CREATED_AT+" = $6 "+
		"where "+COL_ID_REFRESH_SESSIONS+
		"=$7 returning id",
		session.RefreshToken,
		session.UserAgent,
		session.Fingerprint,
		session.Ip,
		session.ExpireIn,
		session.CreatedAt).Scan(&session); err != nil {
		return nil, err
	}

	return session, nil
}
*/
