package mjpeg

import "github.com/vtpl1/vdk/av"

type CodecData struct {
}

func (d CodecData) Type() av.CodecType {
	return av.MJPEG
}
