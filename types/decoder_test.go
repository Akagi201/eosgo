package types_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
	"time"

	"github.com/Akagi201/eosgo/ecc"
	"github.com/Akagi201/eosgo/types"
	"github.com/stretchr/testify/assert"
)

func TestDecoder_Remaining(t *testing.T) {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint16(b, 1)
	binary.LittleEndian.PutUint16(b[2:], 2)

	d := types.NewDecoder(b)

	n, err := d.ReadUint16()
	assert.NoError(t, err)
	assert.Equal(t, uint16(1), n)
	assert.Equal(t, 2, d.Remaining())

	n, err = d.ReadUint16()
	assert.NoError(t, err)
	assert.Equal(t, uint16(2), n)
	assert.Equal(t, 0, d.Remaining())

}

func TestDecoder_Byte(t *testing.T) {
	buf := new(bytes.Buffer)
	enc := types.NewEncoder(buf)
	enc.WriteByte(0)
	enc.WriteByte(1)

	d := types.NewDecoder(buf.Bytes())

	n, err := d.ReadByte()
	assert.NoError(t, err)
	assert.Equal(t, byte(0), n)
	assert.Equal(t, 1, d.Remaining())

	n, err = d.ReadByte()
	assert.NoError(t, err)
	assert.Equal(t, byte(1), n)
	assert.Equal(t, 0, d.Remaining())

}

func TestDecoder_ByteArray(t *testing.T) {
	buf := new(bytes.Buffer)
	enc := types.NewEncoder(buf)
	enc.WriteByteArray([]byte{1, 2, 3})
	enc.WriteByteArray([]byte{4, 5, 6})

	d := types.NewDecoder(buf.Bytes())

	data, err := d.ReadByteArray()
	assert.NoError(t, err)
	assert.Equal(t, []byte{1, 2, 3}, data)
	assert.Equal(t, 4, d.Remaining())

	data, err = d.ReadByteArray()
	assert.Equal(t, []byte{4, 5, 6}, data)
	assert.Equal(t, 0, d.Remaining())

}

func TestDecoder_ByteArray_MissingData(t *testing.T) {
	buf := new(bytes.Buffer)
	enc := types.NewEncoder(buf)
	enc.WriteUVarInt(10)

	d := types.NewDecoder(buf.Bytes())

	_, err := d.ReadByteArray()
	assert.EqualError(t, err, "byte array: varlen=10, missing 10 bytes")

}

func TestDecoder_ByteArrayDataTooSmall(t *testing.T) {

	buf := new(bytes.Buffer)

	//to smalls
	d := types.NewDecoder(buf.Bytes())
	_, err := d.ReadByteArray()
	assert.Equal(t, types.ErrVarIntBufferSize, err)
}

func TestDecoder_Uint16(t *testing.T) {
	buf := new(bytes.Buffer)
	enc := types.NewEncoder(buf)
	enc.WriteUint16(uint16(99))
	enc.WriteUint16(uint16(100))

	d := types.NewDecoder(buf.Bytes())

	n, err := d.ReadUint16()
	assert.NoError(t, err)
	assert.Equal(t, uint16(99), n)
	assert.Equal(t, 2, d.Remaining())

	n, err = d.ReadUint16()
	assert.NoError(t, err)
	assert.Equal(t, uint16(100), n)
	assert.Equal(t, 0, d.Remaining())
}

func TestDecoder_int16(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := types.NewEncoder(buf)
	enc.WriteInt16(int16(-99))
	enc.WriteInt16(int16(100))

	d := types.NewDecoder(buf.Bytes())

	n, err := d.ReadInt16()
	assert.NoError(t, err)
	assert.Equal(t, int16(-99), n)
	assert.Equal(t, 2, d.Remaining())

	n, err = d.ReadInt16()
	assert.NoError(t, err)
	assert.Equal(t, int16(100), n)
	assert.Equal(t, 0, d.Remaining())
}

func TestDecoder_Uint32(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := types.NewEncoder(buf)
	enc.WriteUint32(uint32(342))
	enc.WriteUint32(uint32(100))

	d := types.NewDecoder(buf.Bytes())

	n, err := d.ReadUint32()
	assert.NoError(t, err)
	assert.Equal(t, uint32(342), n)
	assert.Equal(t, 4, d.Remaining())

	n, err = d.ReadUint32()
	assert.NoError(t, err)
	assert.Equal(t, uint32(100), n)
	assert.Equal(t, 0, d.Remaining())
}

func TestDecoder_Uint64(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := types.NewEncoder(buf)
	enc.WriteUint64(uint64(99))
	enc.WriteUint64(uint64(100))

	d := types.NewDecoder(buf.Bytes())

	n, err := d.ReadUint64()
	assert.NoError(t, err)
	assert.Equal(t, uint64(99), n)
	assert.Equal(t, 8, d.Remaining())

	n, err = d.ReadUint64()
	assert.NoError(t, err)
	assert.Equal(t, uint64(100), n)
	assert.Equal(t, 0, d.Remaining())
}

func TestDecoder_string(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := types.NewEncoder(buf)
	enc.WriteString("123")
	enc.WriteString("")
	enc.WriteString("abc")

	d := types.NewDecoder(buf.Bytes())

	s, err := d.ReadString()
	assert.NoError(t, err)
	assert.Equal(t, "123", s)
	assert.Equal(t, 5, d.Remaining())

	s, err = d.ReadString()
	assert.NoError(t, err)
	assert.Equal(t, "", s)
	assert.Equal(t, 4, d.Remaining())

	s, err = d.ReadString()
	assert.NoError(t, err)
	assert.Equal(t, "abc", s)
	assert.Equal(t, 0, d.Remaining())
}

func TestDecoder_SHA256Bytes(t *testing.T) {

	s := types.SHA256Bytes(bytes.Repeat([]byte{1}, 32))

	buf := new(bytes.Buffer)
	enc := types.NewEncoder(buf)
	enc.WriteSHA256Bytes(s)

	d := types.NewDecoder(buf.Bytes())

	rs, err := d.ReadSHA256Bytes()
	assert.NoError(t, err)

	assert.Equal(t, s, rs)
	assert.Equal(t, 0, d.Remaining())
}

func TestDecoder_Empty_SHA256Bytes(t *testing.T) {

	s := types.SHA256Bytes([]byte{})

	buf := new(bytes.Buffer)
	enc := types.NewEncoder(buf)
	enc.WriteSHA256Bytes(s)

	d := types.NewDecoder(buf.Bytes())

	s, err := d.ReadSHA256Bytes()
	assert.NoError(t, err)
	assert.Equal(t, s, types.SHA256Bytes(bytes.Repeat([]byte{0}, 32)))
	assert.Equal(t, 0, d.Remaining())
}

func TestDecoder_PublicKey(t *testing.T) {

	pk := ecc.PublicKey{Curve: ecc.CurveK1, Content: bytes.Repeat([]byte{1}, 33)}

	buf := new(bytes.Buffer)
	enc := types.NewEncoder(buf)
	assert.NoError(t, enc.WritePublicKey(pk))

	d := types.NewDecoder(buf.Bytes())

	rpk, err := d.ReadPublicKey()
	assert.NoError(t, err)

	assert.Equal(t, pk, rpk)
	assert.Equal(t, 0, d.Remaining())
}

func TestDecoder_Empty_PublicKey(t *testing.T) {

	pk := ecc.PublicKey{Curve: ecc.CurveK1, Content: []byte{}}

	buf := new(bytes.Buffer)
	enc := types.NewEncoder(buf)
	assert.Error(t, enc.WritePublicKey(pk))
}

func TestDecoder_Signature(t *testing.T) {

	sig := ecc.Signature{Curve: ecc.CurveK1, Content: bytes.Repeat([]byte{1}, 65)}

	buf := new(bytes.Buffer)
	enc := types.NewEncoder(buf)
	enc.WriteSignature(sig)

	d := types.NewDecoder(buf.Bytes())

	rsig, err := d.ReadSignature()
	assert.NoError(t, err)
	assert.Equal(t, sig, rsig)
	assert.Equal(t, 0, d.Remaining())
}

func TestDecoder_Empty_Signature(t *testing.T) {

	sig := ecc.Signature{Content: []byte{}}

	buf := new(bytes.Buffer)
	enc := types.NewEncoder(buf)
	assert.Error(t, enc.WriteSignature(sig))
}

func TestDecoder_Tstamp(t *testing.T) {

	ts := types.Tstamp{
		time.Unix(0, time.Now().UnixNano()),
	}

	buf := new(bytes.Buffer)
	enc := types.NewEncoder(buf)
	enc.WriteTstamp(ts)

	d := types.NewDecoder(buf.Bytes())

	rts, err := d.ReadTstamp()
	assert.NoError(t, err)
	assert.Equal(t, ts, rts)
	assert.Equal(t, 0, d.Remaining())
}

func TestDecoder_BlockTimestamp(t *testing.T) {

	ts := types.BlockTimestamp{
		time.Unix(time.Now().Unix(), 0),
	}

	buf := new(bytes.Buffer)
	enc := types.NewEncoder(buf)
	enc.WriteBlockTimestamp(ts)

	d := types.NewDecoder(buf.Bytes())

	rbt, err := d.ReadBlockTimestamp()
	assert.NoError(t, err)
	assert.Equal(t, ts, rbt)
	assert.Equal(t, 0, d.Remaining())
}

func TestDecoder_Time(t *testing.T) {

	time := time.Now()

	buf := new(bytes.Buffer)
	enc := types.NewEncoder(buf)
	enc.Encode(&time)

	fmt.Println(buf.Bytes())
}

type EncodeTestStruct struct {
	F1  string
	F2  int16
	F3  uint16
	F4  uint32
	F5  types.SHA256Bytes
	F6  []string
	F7  [2]string
	F8  map[string]string
	F9  ecc.PublicKey
	F10 ecc.Signature
	F11 byte
	F12 uint64
	F13 []byte
	F14 types.Tstamp
	F15 types.BlockTimestamp
	F16 types.Varuint32
	F17 bool
	F18 types.Asset
}

func TestDecoder_Encode(t *testing.T) {

	tstamp := types.Tstamp{Time: time.Unix(0, time.Now().UnixNano())}
	blockts := types.BlockTimestamp{time.Unix(time.Now().Unix(), 0)}
	s := &EncodeTestStruct{
		F1:  "abc",
		F2:  -75,
		F3:  99,
		F4:  999,
		F5:  bytes.Repeat([]byte{0}, 32),
		F6:  []string{"def", "789"},
		F7:  [2]string{"foo", "bar"},
		F8:  map[string]string{"foo": "bar", "hello": "you"},
		F9:  ecc.PublicKey{Curve: ecc.CurveK1, Content: bytes.Repeat([]byte{0}, 33)},
		F10: ecc.Signature{Curve: ecc.CurveK1, Content: bytes.Repeat([]byte{0}, 65)},
		F11: byte(1),
		F12: uint64(87),
		F13: []byte{1, 2, 3, 4, 5},
		F14: tstamp,
		F15: blockts,
		F16: types.Varuint32(999),
		F17: true,
		F18: types.NewEOSAsset(100000),
	}

	buf := new(bytes.Buffer)
	enc := types.NewEncoder(buf)
	enc.Encode(s)

	decoder := types.NewDecoder(buf.Bytes())
	err := decoder.Decode(s)
	assert.NoError(t, err)

	assert.Equal(t, "abc", s.F1)
	assert.Equal(t, int16(-75), s.F2)
	assert.Equal(t, uint16(99), s.F3)
	assert.Equal(t, uint32(999), s.F4)
	assert.Equal(t, types.SHA256Bytes(bytes.Repeat([]byte{0}, 32)), s.F5)
	assert.Equal(t, []string{"def", "789"}, s.F6)
	assert.Equal(t, [2]string{"foo", "bar"}, s.F7)
	assert.Equal(t, map[string]string{"foo": "bar", "hello": "you"}, s.F8)
	assert.Equal(t, ecc.PublicKey{Curve: ecc.CurveK1, Content: bytes.Repeat([]byte{0}, 33)}, s.F9)
	assert.Equal(t, ecc.Signature{Curve: ecc.CurveK1, Content: bytes.Repeat([]byte{0}, 65)}, s.F10)
	assert.Equal(t, byte(1), s.F11)
	assert.Equal(t, uint64(87), s.F12)
	assert.Equal(t, uint64(87), s.F12)
	assert.Equal(t, []byte{1, 2, 3, 4, 5}, s.F13)
	assert.Equal(t, tstamp, s.F14)
	assert.Equal(t, blockts, s.F15)
	assert.Equal(t, types.Varuint32(999), s.F16)
	assert.Equal(t, true, s.F17)
	assert.Equal(t, int64(100000), s.F18.Amount)
	assert.Equal(t, uint8(4), s.F18.Precision)
	assert.Equal(t, "EOS", s.F18.Symbol.Symbol)

}

func TestDecoder_Decode_No_Ptr(t *testing.T) {
	decoder := types.NewDecoder([]byte{})
	err := decoder.Decode(1)
	assert.EqualError(t, err, "decode, can only Decode to pointer type")
}

func TestDecoder_Decode_String_Err(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := types.NewEncoder(buf)
	enc.WriteUVarInt(10)

	decoder := types.NewDecoder(buf.Bytes())
	var s string
	err := decoder.Decode(&s)
	assert.EqualError(t, err, "byte array: varlen=10, missing 10 bytes")
}

func TestDecoder_Decode_Array(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := types.NewEncoder(buf)
	enc.Encode([3]byte{1, 2, 4})

	assert.Equal(t, []byte{1, 2, 4}, buf.Bytes())

	decoder := types.NewDecoder(buf.Bytes())
	var decoded [3]byte
	decoder.Decode(&decoded)
	assert.Equal(t, [3]byte{1, 2, 4}, decoded)

}

func TestDecoder_Decode_Slice_Err(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := types.NewEncoder(buf)

	decoder := types.NewDecoder(buf.Bytes())
	var s []string
	err := decoder.Decode(&s)
	assert.Equal(t, err, types.ErrVarIntBufferSize)

	enc.WriteUVarInt(1)
	decoder = types.NewDecoder(buf.Bytes())
	err = decoder.Decode(&s)
	assert.Equal(t, err, types.ErrVarIntBufferSize)

}

type structWithInvalidType struct {
	F1 time.Duration
}

func TestDecoder_Decode_Struct_Err(t *testing.T) {

	s := structWithInvalidType{}
	decoder := types.NewDecoder([]byte{})
	err := decoder.Decode(&s)
	assert.EqualError(t, err, "decode, unsupported type time.Duration")

}

func TestDecoder_Decode_Map_Err(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := types.NewEncoder(buf)

	decoder := types.NewDecoder(buf.Bytes())
	var m map[string]string
	err := decoder.Decode(&m)
	assert.Equal(t, err, types.ErrVarIntBufferSize)

	enc.WriteUVarInt(1)
	decoder = types.NewDecoder(buf.Bytes())
	err = decoder.Decode(&m)
	assert.Equal(t, err, types.ErrVarIntBufferSize)
}

func TestDecoder_Decode_Bad_Map(t *testing.T) {

	buf := new(bytes.Buffer)
	var m map[string]time.Duration
	enc := types.NewEncoder(buf)
	enc.WriteUVarInt(1)
	enc.WriteString("foo")
	enc.WriteString("bar")

	decoder := types.NewDecoder(buf.Bytes())
	err := decoder.Decode(&m)
	assert.EqualError(t, err, "decode, unsupported type time.Duration")

}

func TestEncoder_Encode_array_error(t *testing.T) {

	decoder := types.NewDecoder([]byte{1})

	toDecode := [1]time.Duration{}
	err := decoder.Decode(&toDecode)

	assert.EqualError(t, err, "decode, unsupported type time.Duration")

}

func TestEncoder_Decode_array_error(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := types.NewEncoder(buf)
	err := enc.Encode([1]time.Duration{time.Duration(0)})
	assert.EqualError(t, err, "Encode: unsupported type time.Duration")

}

func TestEncoder_Encode_slide_error(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := types.NewEncoder(buf)
	err := enc.Encode([]time.Duration{time.Duration(0)})
	assert.EqualError(t, err, "Encode: unsupported type time.Duration")

}
func TestEncoder_Encode_map_error(t *testing.T) {

	buf := new(bytes.Buffer)
	enc := types.NewEncoder(buf)
	err := enc.Encode(map[string]time.Duration{"key": time.Duration(0)})
	assert.EqualError(t, err, "Encode: unsupported type time.Duration")
	err = enc.Encode(map[time.Duration]string{time.Duration(0): "key"})
	assert.EqualError(t, err, "Encode: unsupported type time.Duration")

}

func TestEncoder_Encode_struct_error(t *testing.T) {

	s := struct {
		F time.Duration
	}{
		F: time.Duration(0),
	}

	buf := new(bytes.Buffer)
	enc := types.NewEncoder(buf)
	err := enc.Encode(&s)
	assert.EqualError(t, err, "Encode: unsupported type time.Duration")

}

type TagTestStruct struct {
	S1 string `eos:"-"`
	S2 string
}

func TestEncoder_Decode_struct_tag(t *testing.T) {
	var s TagTestStruct

	buf := new(bytes.Buffer)

	enc := types.NewEncoder(buf)
	enc.WriteString("123")

	d := types.NewDecoder(buf.Bytes())
	d.Decode(&s)
	assert.Equal(t, "", s.S1)
	assert.Equal(t, "123", s.S2)

}

func TestEncoder_Encode_struct_tag(t *testing.T) {

	s := &TagTestStruct{
		S1: "123",
		S2: "abc",
	}

	buf := new(bytes.Buffer)

	enc := types.NewEncoder(buf)
	enc.Encode(s)

	expected := []byte{0x3, 0x61, 0x62, 0x63}
	assert.Equal(t, expected, buf.Bytes())

}

func TestDecoder_readUint16_missing_data(t *testing.T) {

	_, err := types.NewDecoder([]byte{}).ReadByte()
	assert.EqualError(t, err, "byte required [1] byte, remaining [0]")

	_, err = types.NewDecoder([]byte{}).ReadUint16()
	assert.EqualError(t, err, "uint16 required [2] bytes, remaining [0]")

	_, err = types.NewDecoder([]byte{}).ReadUint32()
	assert.EqualError(t, err, "uint32 required [4] bytes, remaining [0]")

	_, err = types.NewDecoder([]byte{}).ReadUint64()
	assert.EqualError(t, err, "uint64 required [8] bytes, remaining [0]")

	_, err = types.NewDecoder([]byte{}).ReadSHA256Bytes()
	assert.EqualError(t, err, "sha256 required [32] bytes, remaining [0]")

	_, err = types.NewDecoder([]byte{}).ReadPublicKey()
	assert.EqualError(t, err, "publicKey required [34] bytes, remaining [0]")

	_, err = types.NewDecoder([]byte{}).ReadSignature()
	assert.EqualError(t, err, "signature required [66] bytes, remaining [0]")

	_, err = types.NewDecoder([]byte{}).ReadTstamp()
	assert.EqualError(t, err, "tstamp required [8] bytes, remaining [0]")

	_, err = types.NewDecoder([]byte{}).ReadBlockTimestamp()
	assert.EqualError(t, err, "blockTimestamp required [4] bytes, remaining [0]")
}
