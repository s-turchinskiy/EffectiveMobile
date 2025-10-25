package config

import (
	"go.uber.org/zap/zapcore"
	"strings"
)

type database struct {
	Host     string
	DBName   string
	Login    string
	Password string
	RawPath  string
}

func (d *database) MarshalLogObject(encoder zapcore.ObjectEncoder) error {

	encoder.AddString("Host", d.Host)
	encoder.AddString("DbName", d.DBName)
	encoder.AddString("Login", d.Login)
	encoder.AddString("Password", "********")
	return nil

}

func (d *database) String() string {
	/*return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
	d.Host, d.Login, d.Password, d.DBName)*/

	return d.RawPath
}

func (d *database) Set(s string) error {

	d.RawPath = s

	parsedStr := strings.ReplaceAll(s, "/", " ")
	hp := strings.Split(parsedStr, " ")
	parsedStr = hp[len(hp)-1]
	parsedStr = strings.ReplaceAll(s, "?", " ")
	hp = strings.Split(parsedStr, " ")
	d.DBName = hp[0]

	/*logger.Log.Debug("database path", "\""+s+"\"")

	s = strings.Replace(s, "://", " ", 1)
	s = strings.Replace(s, ":", " ", 1)
	s = strings.Replace(s, "@", " ", 1)
	s = strings.Replace(s, ":", " ", 1)
	s = strings.Replace(s, "/", " ", 1)
	s = strings.Replace(s, "?", " ", 1)

	hp := strings.Split(s, " ")
	if len(hp) < 6 {
		//return errors.New("need address in a form host=%s user=%s password=%s dbname=%s sslmode=disable")
		return errors.New("incorrect format database-dsn")
	}

	d.Host = hp[3]
	d.Login = hp[1]
	d.Password = hp[2]
	d.DBName = hp[5]*/

	return nil
}
