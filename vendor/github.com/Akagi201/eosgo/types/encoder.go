package types

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/Akagi201/eosgo/ecc"
)

// --------------------------------------------------------------
// Encoder implements the EOS packing, similar to FC_BUFFER
// --------------------------------------------------------------
type Encoder struct {
	output io.Writer
	Order  binary.ByteOrder
	count  int
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		output: w,
		Order:  binary.LittleEndian,
		count:  0,
	}
}

func (e *Encoder) WriteName(name Name) error {
	val, err := StringToName(string(name))
	if err != nil {
		return fmt.Errorf("WriteName: %s", err)
	}
	return e.WriteUint64(val)
}

func (e *Encoder) Encode(v interface{}) (err error) {
	switch cv := v.(type) {
	case Name:
		return e.WriteName(cv)
	case AccountName:
		name := Name(cv)
		return e.WriteName(name)
	case PermissionName:
		name := Name(cv)
		return e.WriteName(name)
	case ActionName:
		name := Name(cv)
		return e.WriteName(name)
	case TableName:
		name := Name(cv)
		return e.WriteName(name)
	case ScopeName:
		name := Name(cv)
		return e.WriteName(name)
	case string:
		return e.WriteString(cv)
	case TransactionStatus:
		return e.WriteByte(uint8(cv))
	case IDListMode:
		return e.WriteByte(byte(cv))
	case byte:
		return e.WriteByte(cv)
	case int8:
		return e.WriteByte(byte(cv))
	case int16:
		return e.WriteInt16(cv)
	case uint16:
		return e.WriteUint16(cv)
	case uint32:
		return e.WriteUint32(cv)
	case uint64:
		return e.WriteUint64(cv)
	case Varuint32:
		return e.WriteUVarInt(int(cv))
	case bool:
		return e.WriteBool(cv)
	case Bool:
		return e.WriteBool(bool(cv))
	case JSONTime:
		return e.writeJSONTime(cv)
	case HexBytes:
		return e.WriteByteArray(cv)
	case []byte:
		return e.WriteByteArray(cv)
	case SHA256Bytes:
		return e.WriteSHA256Bytes(cv)
	case ecc.PublicKey:
		return e.WritePublicKey(cv)
	case ecc.Signature:
		return e.WriteSignature(cv)
	case Tstamp:
		return e.WriteTstamp(cv)
	case BlockTimestamp:
		return e.WriteBlockTimestamp(cv)
	case CurrencyName:
		return e.writeCurrencyName(cv)
	case Asset:
		return e.writeAsset(cv)
		// case *OptionalProducerSchedule:
		// 	isPresent := cv != nil
		// 	e.WriteBool(isPresent)
		// 	if isPresent {

		// 	}
	case ActionData:
		println("ActionData")
		return e.writeActionData(cv)
	case *ActionData:
		println("*ActionData")
		return e.writeActionData(*cv)
	case *P2PMessageEnvelope:
		return e.WriteBlockP2PMessageEnvelope(*cv)
	default:

		rv := reflect.Indirect(reflect.ValueOf(v))
		t := rv.Type()

		switch t.Kind() {

		case reflect.Array:
			l := t.Len()
			//prefix = append(prefix, "     ")
			println(fmt.Sprintf("Encode: array [%T] of length: %d", v, l))

			for i := 0; i < l; i++ {
				if err = e.Encode(rv.Index(i).Interface()); err != nil {
					return
				}
			}
			//prefix = prefix[:len(prefix)-1]
		case reflect.Slice:
			l := rv.Len()
			if err = e.WriteUVarInt(l); err != nil {
				return
			}
			//prefix = append(prefix, "     ")
			println(fmt.Sprintf("Encode: slice [%T] of length: %d", v, l))

			for i := 0; i < l; i++ {
				if err = e.Encode(rv.Index(i).Interface()); err != nil {
					return
				}
			}
			//prefix = prefix[:len(prefix)-1]
			//case reflect.Ptr:
			//	println("*************************************************")
			//	println("*************************************************")
			//	println(fmt.Sprintf("PTR [%T]", v))
			//	println("*************************************************")
			//	println("*************************************************")
		case reflect.Struct:
			l := rv.NumField()
			println(fmt.Sprintf("Encode: struct [%T] with %d field.", v, l))
			//prefix = append(prefix, "     ")

			for i := 0; i < l; i++ {
				field := t.Field(i)
				println(fmt.Sprintf("field -> %s", field.Name))
				//fmt.Println(fmt.Sprintf("field -> %s", field.Name))

				tag := field.Tag.Get("eos")
				if tag == "-" {
					continue
				}

				if v := rv.Field(i); t.Field(i).Name != "_" {
					if v.CanInterface() {
						isPresent := true
						if tag == "optional" {
							isPresent = !v.IsNil()
							e.WriteBool(isPresent)
						}

						//fmt.Printf("IS PRESENT: %T %#v\n", iface, iface, isPresent)

						if isPresent {
							if err = e.Encode(v.Interface()); err != nil {
								return
							}
						}
					}
				}
			}
			//prefix = prefix[:len(prefix)-1]

		case reflect.Map:
			l := rv.Len()
			if err = e.WriteUVarInt(l); err != nil {
				return
			}
			println(fmt.Sprintf("Map [%T] of length: %d", v, l))
			for _, key := range rv.MapKeys() {
				value := rv.MapIndex(key)
				if err = e.Encode(key.Interface()); err != nil {
					return err
				}
				if err = e.Encode(value.Interface()); err != nil {
					return err
				}
			}
		default:
			return errors.New("Encode: unsupported type " + t.String())
		}
	}

	return
}

func (e *Encoder) toWriter(bytes []byte) (err error) {

	e.count += len(bytes)
	println(fmt.Sprintf("    Appending : [%s] pos [%d]", hex.EncodeToString(bytes), e.count))
	_, err = e.output.Write(bytes)
	return
}

func (e *Encoder) WriteByteArray(b []byte) error {
	println(fmt.Sprintf("writing byte array of len [%d]", len(b)))
	if err := e.WriteUVarInt(len(b)); err != nil {
		return err
	}
	return e.toWriter(b)
}

func (e *Encoder) WriteUVarInt(v int) (err error) {
	buf := make([]byte, 8)
	l := binary.PutUvarint(buf, uint64(v))
	return e.toWriter(buf[:l])
}

func (e *Encoder) WriteByte(b byte) (err error) {
	return e.toWriter([]byte{b})
}

func (e *Encoder) WriteBool(b bool) (err error) {
	var out byte
	if b {
		out = 1
	}
	return e.WriteByte(out)
}

func (e *Encoder) WriteUint16(i uint16) (err error) {
	buf := make([]byte, TypeSize.UInt16)
	binary.LittleEndian.PutUint16(buf, i)
	return e.toWriter(buf)
}

func (e *Encoder) WriteInt16(i int16) (err error) {
	return e.WriteUint16(uint16(i))
}

func (e *Encoder) WriteUint32(i uint32) (err error) {
	buf := make([]byte, TypeSize.UInt32)
	binary.LittleEndian.PutUint32(buf, i)
	return e.toWriter(buf)

}

func (e *Encoder) WriteUint64(i uint64) (err error) {
	buf := make([]byte, TypeSize.UInt64)
	binary.LittleEndian.PutUint64(buf, i)
	return e.toWriter(buf)

}

func (e *Encoder) WriteString(s string) (err error) {
	return e.WriteByteArray([]byte(s))
}

func (e *Encoder) WriteSHA256Bytes(s SHA256Bytes) error {
	if len(s) == 0 {
		return e.toWriter(bytes.Repeat([]byte{0}, TypeSize.SHA256Bytes))
	}
	return e.toWriter(s)
}

func (e *Encoder) WritePublicKey(pk ecc.PublicKey) (err error) {
	if len(pk.Content) != 33 {
		return fmt.Errorf("public key %q should be 33 bytes, was %d", hex.EncodeToString(pk.Content), len(pk.Content))
	}

	if err = e.WriteByte(byte(pk.Curve)); err != nil {
		return err
	}

	return e.toWriter(pk.Content)
}

func (e *Encoder) WriteSignature(s ecc.Signature) (err error) {
	if len(s.Content) != 65 {
		return fmt.Errorf("signature should be 65 bytes, was %d", len(s.Content))
	}

	if err = e.WriteByte(byte(s.Curve)); err != nil {
		return
	}

	return e.toWriter(s.Content) // should write 65 bytes
}

func (e *Encoder) WriteTstamp(t Tstamp) (err error) {
	n := uint64(t.UnixNano())
	return e.WriteUint64(n)
}

func (e *Encoder) WriteBlockTimestamp(bt BlockTimestamp) (err error) {
	n := uint32(bt.Unix() - 946684800)
	return e.WriteUint32(n)
}

func (e *Encoder) writeCurrencyName(curreny CurrencyName) (err error) {
	out := make([]byte, 7, 7)
	copy(out, []byte(curreny))

	return e.toWriter(out)
}

func (e *Encoder) writeAsset(asset Asset) (err error) {

	e.WriteUint64(uint64(asset.Amount))
	e.WriteByte(asset.Precision)

	symbol := make([]byte, 7, 7)

	copy(symbol[:], []byte(asset.Symbol.Symbol))
	return e.toWriter(symbol)
}

func (e *Encoder) writeJSONTime(time JSONTime) (err error) {
	return e.WriteUint32(uint32(time.Unix()))
}

func (e *Encoder) WriteBlockP2PMessageEnvelope(envelope P2PMessageEnvelope) (err error) {

	println("WriteBlockP2PMessageEnvelope")

	if envelope.P2PMessage != nil {
		buf := new(bytes.Buffer)
		subEncoder := NewEncoder(buf)
		err = subEncoder.Encode(envelope.P2PMessage)
		if err != nil {
			err = fmt.Errorf("p2p message, %s", err)
			return
		}
		envelope.Payload = buf.Bytes()
	}

	messageLen := uint32(len(envelope.Payload) + 1)
	println(fmt.Sprintf("Message length: %d", messageLen))
	err = e.WriteUint32(messageLen)
	if err == nil {
		err = e.WriteByte(byte(envelope.Type))

		if err == nil {
			return e.toWriter(envelope.Payload)
		}
	}
	return
}

func (e *Encoder) writeActionData(actionData ActionData) (err error) {
	if actionData.Data != nil {
		//if reflect.TypeOf(actionData.Data) == reflect.TypeOf(&ActionData{}) {
		//	log.Fatal("pas cool")
		//}

		println(fmt.Sprintf("entering action data, %T", actionData))
		var d interface{}
		d = actionData.Data
		if reflect.TypeOf(d).Kind() == reflect.Ptr {
			d = reflect.ValueOf(actionData.Data).Elem().Interface()
		}

		if reflect.TypeOf(d).Kind() == reflect.String { //todo : this is a very bad ack ......

			data, err := hex.DecodeString(d.(string))
			if err != nil {
				return fmt.Errorf("ack, %s", err)
			}
			e.WriteByteArray(data)
			return nil

		}

		println(fmt.Sprintf("encoding action data, %T", d))
		raw, err := MarshalBinary(d)
		if err != nil {
			return err
		}
		println(fmt.Sprintf("writing action data, %T", d))
		return e.WriteByteArray(raw)
	}

	return e.WriteByteArray(actionData.HexData)
}

func MarshalBinary(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := NewEncoder(buf)
	err := encoder.Encode(v)
	return buf.Bytes(), err
}
