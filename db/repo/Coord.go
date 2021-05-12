package repo

const (
	TABLE_NAME_COORD = "coord"
	COL_COORD_ID     = "id"
	COL_COORD_X      = "coord_x"
	COL_COORD_Y      = "coord_y"
)

const radius float64 = 20

type CoordModelRepo struct {
	Id       uint64
	CoordLat float64
	CoordLng float64
}

type CoordRepo struct {
	db *DataBase
}

func (s *CoordRepo) Create(coord *CoordModelRepo) (*CoordModelRepo, error) {
	if err := s.db.Db.QueryRow("INSERT into "+TABLE_NAME_COORD+
		" ("+COL_COORD_X+", "+
		COL_COORD_Y+") VALUES ($1, $2) RETURNING "+COL_COORD_ID,
		coord.CoordLat,
		coord.CoordLng).
		Scan(&coord.Id); err != nil {
		return nil, err
	}

	return coord, nil
}

func (s *CoordRepo) Update(coord *CoordModelRepo) (*CoordModelRepo, error) {
	err := s.db.Db.QueryRow("UPDATE "+TABLE_NAME_COORD+" set "+
		COL_COORD_X+"= $1, "+
		COL_COORD_Y+"= $2, "+
		" WHERE "+COL_COORD_ID+"=$3 returning "+COL_COORD_ID+", "+
		COL_COORD_X+", "+
		COL_COORD_Y,
		coord.CoordLat,
		coord.CoordLng).
		Scan(&coord.Id, &coord.CoordLat, &coord.CoordLng)

	if err != nil {
		return nil, err
	}

	return coord, nil
}

func (s *CoordRepo) Delete(coord *CoordModelRepo) (*CoordModelRepo, error) {
	_, err := s.db.Db.Exec("DELETE from "+TABLE_NAME_COORD+" where "+COL_COORD_ID+" = $1", coord.Id)

	if err != nil {
		return nil, err
	}

	return coord, nil
}

func (s *CoordRepo) FindById(id uint64) (*CoordModelRepo, error) {
	coord := CoordModelRepo{}

	if err := s.db.Db.QueryRow("SELECT * from "+TABLE_NAME_COORD+" where "+COL_COORD_ID+"=$1",
		id).
		Scan(&coord.Id, &coord.CoordLat, &coord.CoordLng); err != nil {

		return nil, err
	}

	return &coord, nil
}

func (s *CoordRepo) FindAll() (*[]CoordModelRepo, error) {
	coordList := []CoordModelRepo{}

	row, err := s.db.Db.Query("SELECT * from " + TABLE_NAME_COORD)

	if err != nil {
		return &coordList, err
	}

	for row.Next() {
		coord := CoordModelRepo{}
		err := row.Scan(&coord.Id, &coord.CoordLat, &coord.CoordLng)

		if err != nil {
			return nil, err
		}
		coordList = append(coordList, coord)
	}

	return &coordList, nil
}

// FindInRadius finds all coords in radius
func (s *CoordRepo) FindInRadius(coord *CoordModelRepo) (*[]CoordModelRepo, error) {
	coordList := []CoordModelRepo{}

	xH := coord.CoordLat + radius
	xL := coord.CoordLat - radius
	yH := coord.CoordLng + radius
	yL := coord.CoordLat + radius

	row, err := s.db.Db.Query("SELECT * from "+TABLE_NAME_COORD+" where "+
		COL_COORD_X+"< $1"+" && "+
		COL_COORD_X+"> $2"+" && "+
		COL_COORD_Y+"< $3"+" && "+
		COL_COORD_Y+"> $4", xH, xL, yH, yL)

	if err != nil {
		return &coordList, err
	}

	for row.Next() {
		coord := CoordModelRepo{}
		err := row.Scan(&coord.Id, &coord.CoordLat, &coord.CoordLng)

		if err != nil {
			return nil, err
		}
		coordList = append(coordList, coord)
	}

	return &coordList, nil
}
