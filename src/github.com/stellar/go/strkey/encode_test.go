package strkey

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	cases := []struct {
		Name        string
		VersionByte VersionByte
		Payload     []byte
		Expected    string
	}{
		{
			Name:        "AccountID",
			VersionByte: VersionByteAccountID,
			Payload: []byte{
				0x36, 0x3e, 0xaa, 0x38, 0x67, 0x84, 0x1f, 0xba,
				0xd0, 0xf4, 0xed, 0x88, 0xc7, 0x79, 0xe4, 0xfe,
				0x66, 0xe5, 0x6a, 0x24, 0x70, 0xdc, 0x98, 0xc0,
				0xec, 0x9c, 0x07, 0x3d, 0x05, 0xc7, 0xb1, 0x03,
			},
			Expected: "GA3D5KRYM6CB7OWQ6TWYRR3Z4T7GNZLKERYNZGGA5SOAOPIFY6YQHES5",
		},
		{
			Name:        "Seed",
			VersionByte: VersionByteSeed,
			Payload: []byte{
				0x69, 0xa8, 0xc4, 0xcb, 0xb9, 0xf6, 0x4e, 0x8a,
				0x07, 0x98, 0xf6, 0xe1, 0xac, 0x65, 0xd0, 0x6c,
				0x31, 0x62, 0x92, 0x90, 0x56, 0xbc, 0xf4, 0xcd,
				0xb7, 0xd3, 0x73, 0x8d, 0x18, 0x55, 0xf3, 0x63,
			},
			Expected: "SBU2RRGLXH3E5CQHTD3ODLDF2BWDCYUSSBLLZ5GNW7JXHDIYKXZWHOKR",
		},
		{
			Name:        "HashTx",
			VersionByte: VersionByteHashTx,
			Payload: []byte{
				0x69, 0xa8, 0xc4, 0xcb, 0xb9, 0xf6, 0x4e, 0x8a,
				0x07, 0x98, 0xf6, 0xe1, 0xac, 0x65, 0xd0, 0x6c,
				0x31, 0x62, 0x92, 0x90, 0x56, 0xbc, 0xf4, 0xcd,
				0xb7, 0xd3, 0x73, 0x8d, 0x18, 0x55, 0xf3, 0x63,
			},
			Expected: "TBU2RRGLXH3E5CQHTD3ODLDF2BWDCYUSSBLLZ5GNW7JXHDIYKXZWHXL7",
		},
		{
			Name:        "HashX",
			VersionByte: VersionByteHashX,
			Payload: []byte{
				0x69, 0xa8, 0xc4, 0xcb, 0xb9, 0xf6, 0x4e, 0x8a,
				0x07, 0x98, 0xf6, 0xe1, 0xac, 0x65, 0xd0, 0x6c,
				0x31, 0x62, 0x92, 0x90, 0x56, 0xbc, 0xf4, 0xcd,
				0xb7, 0xd3, 0x73, 0x8d, 0x18, 0x55, 0xf3, 0x63,
			},
			Expected: "XBU2RRGLXH3E5CQHTD3ODLDF2BWDCYUSSBLLZ5GNW7JXHDIYKXZWGTOG",
		},
	}

	for _, kase := range cases {
		actual, err := Encode(kase.VersionByte, kase.Payload)
		if assert.NoError(t, err, "An error occured in case %s", kase.Name) {
			assert.Equal(t, kase.Expected, actual, "Output mismatch in case %s", kase.Name)
		}
	}

	// test bad version byte
	_, err := Encode(VersionByte(2), cases[0].Payload)
	assert.Error(t, err)
}
