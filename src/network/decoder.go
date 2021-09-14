package decoder

import "errors"

const (
	// MaxBodySize max proto body size
	MaxBodySize = int32(1 << 12)
)

var (
	// ErrProtoPackLen proto packet len error
	ErrProtoPackLen = errors.New("default server codec pack length error")
	// ErrProtoHeaderLen proto header len error
	ErrProtoHeaderLen = errors.New("default server codec header length error")
)

const (
	// size
	_packSize      = 4
	_headerSize    = 2
	_verSize       = 2
	_opSize        = 4
	_seqSize       = 4
	_heartSize     = 4
	_rawHeaderSize = _packSize + _headerSize + _verSize + _opSize + _seqSize
	_maxPackSize   = MaxBodySize + int32(_rawHeaderSize)
	// offset
	_packOffset   = 0
	_headerOffset = _packOffset + _packSize
	_verOffset    = _headerOffset + _headerSize
	_opOffset     = _verOffset + _verSize
	_seqOffset    = _opOffset + _opSize
	_heartOffset  = _seqOffset + _seqSize
)

// ReadTCP read a proto from TCP reader.
func ReadTCP(rr *Reader) (message *Message, err error) {
	var (
		bodyLen   int
		headerLen int16
		packLen   int32
		buf       []byte
	)
	if buf, err = rr.Pop(_rawHeaderSize); err != nil {
		return
	}
	message = &Message{}
	packLen = BigEndian.Int32(buf[_packOffset:_headerOffset])
	headerLen = BigEndian.Int16(buf[_headerOffset:_verOffset])
	message.Ver = int32(BigEndian.Int16(buf[_verOffset:_opOffset]))
	message.Op = BigEndian.Int32(buf[_opOffset:_seqOffset])
	message.Seq = BigEndian.Int32(buf[_seqOffset:])
	if packLen > _maxPackSize {
		return nil, ErrProtoPackLen
	}
	if headerLen != _rawHeaderSize {
		return nil, ErrProtoHeaderLen
	}
	if bodyLen = int(packLen - int32(headerLen)); bodyLen > 0 {
		message.Body, err = rr.Pop(bodyLen)
	} else {
		message.Body = nil
	}
	return message, err
}
