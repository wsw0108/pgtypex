package pgtypex

import (
	"bytes"

	wktparser "github.com/Succo/wktToOrb"
	"github.com/jackc/pgtype"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkb"
	"github.com/paulmach/orb/encoding/wkt"
)

type Geometry struct {
	Geom orb.Geometry
}

func NewGeometry(geom orb.Geometry) *Geometry {
	return &Geometry{geom}
}

func (Geometry) PreferredParamFormat() int16 {
	return pgtype.BinaryFormatCode
}

func (g *Geometry) EncodeText(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	// FIXME: use buf
	// w := bytes.NewBuffer(buf)
	wktString := wkt.MarshalString(g.Geom)
	return []byte(wktString), nil
}

func (g *Geometry) EncodeBinary(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	w := bytes.NewBuffer(buf)
	enc := wkb.NewEncoder(w)
	if err := enc.Encode(g.Geom); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (Geometry) PreferredResultFormat() int16 {
	return pgtype.BinaryFormatCode
}

func (g *Geometry) DecodeText(ci *pgtype.ConnInfo, src []byte) (err error) {
	wktString := string(src)
	geom, err := wktparser.Scan(wktString)
	if err != nil {
		return
	}
	g.Geom = geom
	return
}

func (g *Geometry) DecodeBinary(ci *pgtype.ConnInfo, src []byte) (err error) {
	s := wkb.Scanner(nil)
	if err = s.Scan(src); err != nil {
		return
	}
	g.Geom = s.Geometry
	return
}
