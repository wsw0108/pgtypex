package pgtypex

import (
	"bytes"

	"github.com/jackc/pgtype"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkb"
)

type Geometry struct {
	Geom orb.Geometry
}

func NewGeometry(geom orb.Geometry) *Geometry {
	return &Geometry{geom}
}

func (g *Geometry) EncodeBinary(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	w := bytes.NewBuffer(buf)
	enc := wkb.NewEncoder(w)
	if err := enc.Encode(g.Geom); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (g *Geometry) DecodeBinary(ci *pgtype.ConnInfo, src []byte) (err error) {
	s := wkb.Scanner(nil)
	if err = s.Scan(src); err != nil {
		return
	}
	g.Geom = s.Geometry
	return
}
