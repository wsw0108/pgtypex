package pgtypex

import (
	"reflect"
	"testing"

	"github.com/jackc/pgtype"
	"github.com/paulmach/orb"
)

func TestGeometry_EncodeBinary(t *testing.T) {
	type fields struct {
		geom orb.Geometry
	}
	type args struct {
		ci  *pgtype.ConnInfo
		buf []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:   "encode return updated buf",
			fields: fields{geom: orb.Point{113.68328500000001, 31.257848300000003}},
			args:   args{nil, make([]byte, 0, 32)},
			want:   []byte{1, 1, 0, 0, 0, 60, 54, 2, 241, 186, 107, 92, 64, 71, 212, 159, 88, 2, 66, 63, 64},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Geometry{
				Geom: tt.fields.geom,
			}
			got, err := g.EncodeBinary(tt.args.ci, tt.args.buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeBinary() got = %v, want %v", got, tt.want)
			}
		})
	}
}
