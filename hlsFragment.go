package main

import (
	"time"

	"github.com/deepch/vdk/av"
)

//Fragment struct
type Fragment struct {
	Independent bool          //Fragment have i-frame (key frame)
	Finish      bool          //Fragment Ready
	Duration    time.Duration //Fragment Duration
	Packets     []*av.Packet  //Packet Slice
}