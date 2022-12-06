module github.com/vtpl1/go_rtsp_to_web

go 1.19

require (
	github.com/sirupsen/logrus v1.9.0
	github.com/vtpl1/vdk v0.0.0-00010101000000-000000000000
)

require golang.org/x/sys v0.3.0 // indirect

replace github.com/vtpl1/vdk => ./vtpl1/vdk
