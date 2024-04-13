package utils

import (
	"math/big"
	"testing"
)

type TestStruct struct {
	Time *big.Int
}

func intToBigInt(val int) *big.Int {
	bigint := new(big.Int)
	bigint.SetInt64(int64(val))
	return bigint
}

func TestComparePtrFieldsDesc(t *testing.T) {
	tests := []struct {
		name string
		s1   TestStruct
		s2   TestStruct
		want bool
	}{
		{
			name: "Greater Value",
			s1:   TestStruct{Time: func() *big.Int { i := intToBigInt(5); return i }()},
			s2:   TestStruct{Time: func() *big.Int { i := intToBigInt(3); return i }()},
			want: true,
		},
		{
			name: "Lesser Value",
			s1:   TestStruct{Time: func() *big.Int { i := intToBigInt(3); return i }()},
			s2:   TestStruct{Time: func() *big.Int { i := intToBigInt(5); return i }()},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ComparePtrFieldsDesc(&tt.s1, &tt.s2); got != tt.want {
				t.Errorf("comparePtrFieldsDesc() = %v, want %v", got, tt.want)
			}
		})
	}
}
