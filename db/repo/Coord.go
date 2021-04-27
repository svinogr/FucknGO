package repo

const (
	TABLE_NAME_COORD = "coordRepo"
	COL_COORD_ID     = "id"
	COL_COORD_X      = "coord_x"
	COL_COORD_Y      = "coord_y"
)

const radius float64 = 20

type CoordModelRepo struct {
	Id     uint64
	CoordX float64
	CoordY float64
}

type CoordRepo struct {
	db *DataBase
}

func (s *CoordRepo) Create(coord *CoordModelRepo) (*CoordModelRepo, error) {
	if err := s.db.Db.QueryRow("INSERT into "+TABLE_NAME_COORD+
		" ("+COL_COORD_X+", "+
		COL_COORD_Y+") VALUES ($1, $2, $3, $4) RETURNING "+COL_COORD_ID,
		coord.CoordX,
		coord.CoordY).
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
		coord.CoordX,
		coord.CoordY).
		Scan(&coord.Id, &coord.CoordX, &coord.CoordY)

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

func (s *CoordRepo) FindById(coord *CoordModelRepo) (*CoordModelRepo, error) {
	if err := s.db.Db.QueryRow("SELECT * from "+TABLE_NAME_COORD+" where "+COL_COORD_ID+"=$1",
		coord.Id).
		Scan(&coord.Id, &coord.CoordX, &coord.CoordY); err != nil {

		return nil, err
	}

	return coord, nil
}

func (s *CoordRepo) FindAll() (*[]CoordModelRepo, error) {
	coordList := []CoordModelRepo{}

	row, err := s.db.Db.Query("SELECT * from " + TABLE_NAME_COORD)

	if err != nil {
		return &coordList, err
	}

	for row.Next() {
		coord := CoordModelRepo{}
		err := row.Scan(&coord.Id, &coord.CoordX, &coord.CoordY)

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

	xH := coord.CoordX + radius
	xL := coord.CoordX - radius
	yH := coord.CoordY + radius
	yL := coord.CoordX + radius

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
		err := row.Scan(&coord.Id, &coord.CoordX, &coord.CoordY)

		if err != nil {
			return nil, err
		}
		coordList = append(coordList, coord)
	}

	return &coordList, nil
}
