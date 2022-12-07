package main

import "github.com/vtpl1/go_rtsp_to_web/utils"

func (obj *StorageST) StreamChannelRunAll() {
	utils.Logger.Info("Here")
	for k, v := range obj.Streams {
		for ks, vs := range v.Channels {
			if !vs.OnDemand {
				vs.runLock = true
				// go StreamServerRunStreamDo(k, ks)
				v.Channels[ks] = vs
				obj.Streams[k] = v
			}
		}
	}
}
