module github.com/vtpl1/go_rtsp_to_web

go 1.19

require (
	github.com/go-chi/chi/v5 v5.0.8
	github.com/imdario/mergo v0.3.13
	github.com/vtpl1/vdk v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.24.0
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

require (
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

replace github.com/vtpl1/vdk => ./vtpl1/vdk
