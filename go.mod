module github.com/almanac1631/scrubarr

go 1.24.0

toolchain go1.25.7

require (
	github.com/autobrr/go-rtorrent v1.12.0
	github.com/gdm85/go-libdeluge v0.6.0
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/gorilla/handlers v1.5.2
	github.com/knadh/koanf/parsers/toml/v2 v2.2.0
	github.com/knadh/koanf/providers/file v1.2.1
	github.com/knadh/koanf/v2 v2.3.2
	github.com/spf13/cobra v1.10.2
	golang.org/x/crypto v0.47.0
	golang.org/x/term v0.40.0
	golift.io/starr v1.3.0
)

require (
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/gdm85/go-rencode v0.1.8 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/knadh/koanf/maps v0.1.2 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
)

replace github.com/autobrr/go-rtorrent v1.12.0 => github.com/almanac1631/go-rtorrent v1.16.0
