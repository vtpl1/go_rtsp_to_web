package format

import (
	"github.com/vtpl1/vdk/av/avutil"
	"github.com/vtpl1/vdk/format/aac"
	"github.com/vtpl1/vdk/format/flv"
	"github.com/vtpl1/vdk/format/mp4"
	"github.com/vtpl1/vdk/format/rtmp"
	"github.com/vtpl1/vdk/format/rtsp"
	"github.com/vtpl1/vdk/format/ts"
)

func RegisterAll() {
	avutil.DefaultHandlers.Add(mp4.Handler)
	avutil.DefaultHandlers.Add(ts.Handler)
	avutil.DefaultHandlers.Add(rtmp.Handler)
	avutil.DefaultHandlers.Add(rtsp.Handler)
	avutil.DefaultHandlers.Add(flv.Handler)
	avutil.DefaultHandlers.Add(aac.Handler)
}
