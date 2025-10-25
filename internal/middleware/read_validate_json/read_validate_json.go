package readvalidatejson

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/mailru/easyjson"
	commonerrors "github.com/s-turchinskiy/EffectiveMobile/internal/common/errors"
	"github.com/s-turchinskiy/EffectiveMobile/internal/middleware/logger"
	"github.com/s-turchinskiy/EffectiveMobile/internal/models"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type KeyUnmarshalJSONType string

const (
	UnmarshalJSON KeyUnmarshalJSONType = "unmarshalJson"
)

func Middleware(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		if !(r.RequestURI == "/api/subscriptions/sum" ||
			r.RequestURI == "/api/subscription/create" ||
			r.RequestURI == "/api/subscription/update") {
			next.ServeHTTP(w, r)
			return
		}

		defer r.Body.Close()
		bodyByte, err := io.ReadAll(r.Body)
		if err != nil {
			commonerrors.ErrorGettingData(w, err)
			return
		}

		var req any
		switch r.RequestURI {

		case "/api/subscriptions/sum":
			var reqTyped models.SumSubscriptionsJSON
			err = easyjson.Unmarshal(bodyByte, &reqTyped)
			req = reqTyped
		case "/api/subscription/create":
			var reqTyped models.CreateSubscriptionJSON
			err = easyjson.Unmarshal(bodyByte, &reqTyped)
			req = reqTyped
		case "/api/subscription/update":
			var reqTyped models.UpdateSubscriptionJSON
			err = easyjson.Unmarshal(bodyByte, &reqTyped)
			req = reqTyped
		}

		if err != nil {

			logger.Log.Info("cannot decode request JSON body", zap.Error(commonerrors.WrapError(err)), zap.String("body", string(bodyByte)))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = validator.New().Struct(req)
		if err != nil {
			logger.Log.Info(zap.Error(err))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UnmarshalJSON, req)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
