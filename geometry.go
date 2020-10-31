package pgtypex

import (
	"bytes"
	"database/sql/driver"
	"fmt"

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

func (dst *Geometry) Set(src interface{}) error {
	if src == nil {
		return nil
	}

	switch value := src.(type) {
	case string:
		geom, err := wktparser.Scan(value)
		if err != nil {
			return err
		}
		dst.Geom = geom
	case []byte:
		geom, err := wkb.Unmarshal(value)
		if err != nil {
			return err
		}
		dst.Geom = geom
	case orb.Geometry:
		// TODO: clone?
		dst.Geom = value
	default:
		return fmt.Errorf("cannot convert %v to Geometry", value)
	}

	return nil
}

func (dst Geometry) Get() interface{} {
	return dst.Geom
}

func (src *Geometry) AssignTo(dst interface{}) error {
	return fmt.Errorf("cannot assign %v to %T", src, dst)
}

func (Geometry) PreferredParamFormat() int16 {
	return pgtype.BinaryFormatCode
}

func (src Geometry) EncodeText(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	// FIXME: use buf
	// w := bytes.NewBuffer(buf)
	wktString := wkt.MarshalString(src.Geom)
	return []byte(wktString), nil
}

func (src Geometry) EncodeBinary(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	w := bytes.NewBuffer(buf)
	enc := wkb.NewEncoder(w)
	if err := enc.Encode(src.Geom); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (Geometry) PreferredResultFormat() int16 {
	return pgtype.BinaryFormatCode
}

func (dst *Geometry) DecodeText(ci *pgtype.ConnInfo, src []byte) error {
	if src == nil {
		return nil
	}

	wktString := string(src)
	geom, err := wktparser.Scan(wktString)
	if err != nil {
		return err
	}

	dst.Geom = geom
	return nil
}

func (dst *Geometry) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
	if src == nil {
		return nil
	}

	geom, err := wkb.Unmarshal(src)
	if err != nil {
		return err
	}

	dst.Geom = geom
	return nil
}

// Scan implements the database/sql Scanner interface.
func (dst *Geometry) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	switch src := src.(type) {
	case string:
		geom, err := wktparser.Scan(src)
		if err != nil {
			return err
		}
		dst.Geom = geom
		return nil
	case []byte:
		srcCopy := make([]byte, len(src))
		copy(srcCopy, src)
		s := wkb.Scanner(nil)
		err := s.Scan(srcCopy)
		if err != nil {
			return err
		}
		dst.Geom = s.Geometry
		return nil
	}

	return fmt.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src Geometry) Value() (driver.Value, error) {
	if src.Geom == nil {
		return nil, nil
	}
	return wkb.Marshal(src.Geom)
}
