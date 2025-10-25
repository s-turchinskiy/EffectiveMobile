package postgresql

import (
	"github.com/s-turchinskiy/EffectiveMobile/cmd/subscriptions/config"
	"github.com/s-turchinskiy/EffectiveMobile/internal/common/errors"
)

func getRequest(fileName string) (string, error) {

	request, exist := config.PublicConfig.SQLRequests[fileName]
	if !exist {
		return "", commonerrors.NewErrorErrorSQLRequestNoExist(fileName)
	}

	return request, nil
}
