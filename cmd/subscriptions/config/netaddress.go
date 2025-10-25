package config

import (
	"errors"
	"go.uber.org/zap/zapcore"
	"strconv"
	"strings"
)

type netAddress struct {
	Host string
	Port int
}

func (a *netAddress) String() string {
	return a.Host + ":" + strconv.Itoa(a.Port)
}

func (a *netAddress) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("need address in a form host:port")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	a.Host = hp[0]
	a.Port = port
	return nil
}

func (a *netAddress) MarshalLogObject(encoder zapcore.ObjectEncoder) error {

	encoder.AddString("Host", a.Host)
	encoder.AddInt("Port", a.Port)
	return nil
}
