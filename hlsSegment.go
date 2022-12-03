package main

import "time"

//Segment struct
type Segment struct {
	FPS               int               //Current fps
	CurrentFragment   *Fragment         //CurrentFragment link
	CurrentFragmentID int               //CurrentFragment ID
	Finish            bool              //Segment Ready
	Duration          time.Duration     //Segment Duration
	Time              time.Time         //Realtime EXT-X-PROGRAM-DATE-TIME
	Fragment          map[int]*Fragment //Fragment map
}