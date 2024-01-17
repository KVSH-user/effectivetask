package identifier

import (
	resp "effectivetask/internal/lib/api/response"
	"effectivetask/internal/lib/getdata"
	"effectivetask/internal/storage/psql"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Request struct {
	Id         int    `json:"id,omitempty"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
	Country    string `json:"country,omitempty"`
	Gender     string `json:"gender,omitempty"`
	Age        int    `json:"age,omitempty"`
}

type Response struct {
	resp.Response
}

type PersonSaver interface {
	SavePeople(name string, surname string, patronymic string, country string, gender string, age int) error
}

type PersonDeliter interface {
	DelPeople(id int) error
}

type People struct {
	Id         int    `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
	Country    string `json:"country,omitempty"`
	Gender     string `json:"gender,omitempty"`
	Age        int    `json:"age,omitempty"`
}

type SearchId interface {
	SearchId(id int) (People, error)
}

type SearchName interface {
	SearchName(name string, surname string) (People, error)
}

type EditById interface {
	Edit(id int, name string, surname string, patronymic string, country string, gender string, age int) error
}

func New(log *slog.Logger, personSaver PersonSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.handlers.person.identifier.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", err)

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		country := req.Country
		if country == "" {
			country = getdata.GetCountry(req.Name)
		}

		gender := req.Gender
		if gender == "" {
			gender = getdata.GetSex(req.Name)
		}

		age := req.Age
		if age == 0 {
			age = getdata.GetAge(req.Name)
		}

		err = personSaver.SavePeople(req.Name, req.Surname, req.Patronymic, country, gender, age)
		if err != nil {
			log.Info("failed to add people: ", err)

			render.JSON(w, r, resp.Error("failed to add people"))

			return
		}

		log.Info("people added")

		render.JSON(w, r, Response{resp.OK()})
	}
}

func Del(log *slog.Logger, personDeliter PersonDeliter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.handlers.person.identifier.Del"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", err)

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		err = personDeliter.DelPeople(req.Id)
		if err != nil {
			log.Info("failed to delete people: ", err)

			render.JSON(w, r, resp.Error("failed to delete people"))

			return
		}

		log.Info("people deleted")

		render.JSON(w, r, Response{resp.OK()})
	}
}

func SearchById(log *slog.Logger, searchId *psql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.handlers.person.identifier.SearchById"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", err)

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		respo, err := searchId.SearchId(req.Id)
		if err != nil {
			log.Info("failed to search people: ", err)

			render.JSON(w, r, resp.Error("failed to search people"))

			return
		}

		log.Info("people searched")

		render.JSON(w, r, respo)
	}
}

func SearchByName(log *slog.Logger, searchName *psql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.handlers.person.identifier.SearchByName"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", err)

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		respo, err := searchName.SearchName(req.Name, req.Surname)
		if err != nil {
			log.Info("failed to search people: ", err)

			render.JSON(w, r, resp.Error("failed to search people"))

			return
		}

		log.Info("people searched")

		render.JSON(w, r, respo)
	}
}

func Edit(log *slog.Logger, editById EditById) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.handlers.person.identifier.Edit"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", err)

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		err = editById.Edit(req.Id, req.Name, req.Surname, req.Patronymic, req.Country, req.Gender, req.Age)
		if err != nil {
			log.Info("failed to edit people: ", err)

			render.JSON(w, r, resp.Error("failed to edit people"))

			return
		}

		log.Info("people edited")

		render.JSON(w, r, Response{resp.OK()})
	}
}
