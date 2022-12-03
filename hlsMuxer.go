package main

import (
	"context"
	"sync"
)

//MuxerHLS struct
type MuxerHLS struct {
	mutex             sync.RWMutex
	UUID              string             //Current UUID
	MSN               int                //Current MSN
	FPS               int                //Current FPS
	MediaSequence     int                //Current MediaSequence
	CurrentFragmentID int                //Current fragment id
	CacheM3U8         string             //Current index cache
	CurrentSegment    *Segment           //Current segment link
	Segments          map[int]*Segment   //Current segments group
	FragmentCtx       context.Context    //chan 1-N
	FragmentCancel    context.CancelFunc //chan 1-N
}