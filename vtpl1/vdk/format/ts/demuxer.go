package ts

import (
	"bufio"

	"github.com/vtpl1/vdk/av"
	"github.com/vtpl1/vdk/format/ts/tsio"
)

type Demuxer struct {
	r *bufio.Reader

	pkts []av.Packet

	pat     *tsio.PAT
	pmt     *tsio.PMT
	streams []*Stream
	tshdr   []byte
	AnnexB  bool
	stage   int
}