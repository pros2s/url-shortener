package save

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"url-shortener/internal/lib/random"
	"url-shortener/internal/lib/response"
	"url-shortener/internal/lib/sl"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

const aliasLength = 6

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

type UrlSaver interface {
	SaveToUrl(url string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver UrlSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn := "internal.http.handlers.save.New"

		log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// decode request
		var req Request

		if err := render.DecodeJSON(r.Body, req); err != nil {
			errorMessage := "Failed to decode request body"

			log.Error(errorMessage, sl.AttrByErr(err))
			render.JSON(w, r, response.Error(errorMessage))

			return
		}

		log.Info("Request body has decoded", slog.Any("request", req))

		// validation
		if err := validator.New().Struct(req); err != nil {

			log.Error("Validation error", sl.AttrByErr(err))

			errs := err.(validator.ValidationErrors)
			render.JSON(w, r, response.ValidateErrors(errs))

			return
		}

		// random alias
		alias := req.Alias
		if alias == "" {
			alias = random.RandomAlias(aliasLength)
		}

		// save
		id, err := urlSaver.SaveToUrl(req.URL, alias)
		// url exists
		if errors.Is(err, storage.ErrUrlExists) {
			errMessage := fmt.Sprintf("Url with alias \"%s\" already exists", alias)

			log.Info(errMessage, sl.AttrByErr(err))
			render.JSON(w, r, response.Error(errMessage))

			return
		}

		// url error
		if err != nil {
			errMessage := fmt.Sprintf("Error with saving url by alias \"%s\"", alias)

			log.Error(errMessage, sl.AttrByErr(err))
			render.JSON(w, r, response.Error(errMessage))

			return
		}

		// success response
		log.Info(fmt.Sprintf("Url by alias \"%s\" with id %d has saved", alias, id))
		render.JSON(w, r, Response{Response: response.Ok(), Alias: alias})
	}
}
