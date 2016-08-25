package store

import "github.com/ugorji/go/codec"

var (
	bh codec.BincHandle
	mh codec.MsgpackHandle
	ch codec.CborHandle
	jh codec.JsonHandle
)

var defCodec = &jh

func (s *store) Encode(obj interface{}) ([]byte, error) {
	var (
		buf []byte
		err error
		enc *codec.Encoder
	)

	enc = codec.NewEncoderBytes(&buf, s.codec)

	if err = enc.Encode(obj); err != nil {
		return nil, err
	}

	return buf, nil
}

func (s *store) Decode(buf []byte, obj interface{}) error {
	var (
		err error
		dec *codec.Decoder
	)

	dec = codec.NewDecoderBytes(buf, s.codec)
	if err = dec.Decode(obj); err != nil {
		return err
	}

	return nil
}
