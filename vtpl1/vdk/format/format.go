package format

import (
	"github.com/vtpl1/vdk/av/avutil"
	"github.com/vtpl1/vdk/format/rtsp"
)

func RegisterAll() {
	avutil.DefaultHandlers.Add(rtsp.Handler)
}
