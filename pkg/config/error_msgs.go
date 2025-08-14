package config

import "errors"

var (
	ErrReadConfig  = errors.New("config: cannot read config file")
	ErrParseConfig = errors.New("config: cannot parse yaml")
)
