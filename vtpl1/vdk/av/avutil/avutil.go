package avutil

import (
	"io"

	"github.com/vtpl1/vdk/av"
)

type RegisterHandler struct {
	Ext           string
	ReaderDemuxer func(io.Reader) av.Demuxer
	WriterMuxer   func(io.Writer) av.Muxer
	UrlMuxer      func(string) (bool, av.MuxCloser, error)
	UrlDemuxer    func(string) (bool, av.DemuxCloser, error)
	UrlReader     func(string) (bool, io.ReadCloser, error)
	Probe         func([]byte) bool
	AudioEncoder  func(av.CodecType) (av.AudioEncoder, error)
	AudioDecoder  func(av.AudioCodecData) (av.AudioDecoder, error)
	ServerDemuxer func(string) (bool, av.DemuxCloser, error)
	ServerMuxer   func(string) (bool, av.MuxCloser, error)
	CodecTypes    []av.CodecType
}
type Handlers struct {
	handlers []RegisterHandler
}

var DefaultHandlers = &Handlers{}

func (self *Handlers) Add(fn func(*RegisterHandler)) {
	handler := &RegisterHandler{}
	fn(handler)
	self.handlers = append(self.handlers, *handler)
}
