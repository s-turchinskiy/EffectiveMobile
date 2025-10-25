package config

import (
	"os"
	path2 "path"
	"time"
)

type publicConfig struct {
	OutputPathsLogsDefault string
	EnvFilename            string
	HTTPServer             HTTPServerConfig
	SQLRequests            map[string]string
	ShutdownTimeout        time.Duration
}

type HTTPServerConfig struct {
	RetryStrategy []time.Duration
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
}

type AccrualServiceConfig struct {
	Interval       time.Duration
	TimeoutRequest time.Duration
	Address        string
	WorkerCount    int
	RetryStrategy  []time.Duration
}

var PublicConfig publicConfig

func InitializePublicConfig() (err error) {

	retryStrategy := []time.Duration{
		0,
		2 * time.Second,
		5 * time.Second}

	PublicConfig = publicConfig{
		OutputPathsLogsDefault: "server.log, stdout",
		EnvFilename:            "./cmd/subscriptions/.env",
		ShutdownTimeout:        15 * time.Second,
		HTTPServer: HTTPServerConfig{
			ReadTimeout:   5 * time.Second,
			WriteTimeout:  10 * time.Second,
			RetryStrategy: retryStrategy,
		},
	}

	PublicConfig.SQLRequests, err = getSQLRequests("./internal/repository/postgresql/sqlrequests/")
	return err

}

func getSQLRequests(dir string) (map[string]string, error) {

	files, err := os.ReadDir(dir)
	sqlRequests := make(map[string]string, len(files))

	if err != nil {
		return nil, err
	}

	for _, file := range files {
		filename := file.Name()
		path := path2.Join(dir, filename)
		file, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		sqlRequests[filename] = string(file)
	}

	return sqlRequests, nil
}
