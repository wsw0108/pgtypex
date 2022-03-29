package pgtypex

import (
	"bytes"
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgtype"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkb"
	"github.com/twpayne/go-geom/encoding/wkt"
)

var byteOrder = wkb.NDR

type Geometry struct {
	Geom geom.T
}

func NewGeometry(geom geom.T) *Geometry {
	return &Geometry{geom}
}

func (dst *Geometry) Set(src interface{}) error {
	if src == nil {
		return nil
	}

	switch value := src.(type) {
	case string:
		g, err := wkt.Unmarshal(value)
		if err != nil {
			return err
		}
		dst.Geom = g
	case []byte:
		geom, err := wkb.Unmarshal(value)
		if err != nil {
			return err
		}
		dst.Geom = geom
	case geom.T:
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
	wktString, err := wkt.Marshal(src.Geom)
	return []byte(wktString), err
}

func (src Geometry) EncodeBinary(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	w := bytes.NewBuffer(buf)
	if err := wkb.Write(w, byteOrder, src.Geom); err != nil {
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
	geom, err := wkt.Unmarshal(wktString)
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
		geom, err := wkt.Unmarshal(src)
		if err != nil {
			return err
		}
		dst.Geom = geom
		return nil
	case []byte:
		srcCopy := make([]byte, len(src))
		copy(srcCopy, src)
		g, err := wkb.Unmarshal(srcCopy)
		if err != nil {
			return err
		}
		dst.Geom = g
		return nil
	}

	return fmt.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src Geometry) Value() (driver.Value, error) {
	if src.Geom == nil {
		return nil, nil
	}
	return wkb.Marshal(src.Geom, byteOrder)
}
