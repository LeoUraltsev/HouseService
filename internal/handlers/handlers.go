package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/LeoUraltsev/HouseService/internal/gen"
	mv "github.com/LeoUraltsev/HouseService/internal/middleware"
	"github.com/LeoUraltsev/HouseService/internal/models"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validate = validator.New()

type HouseService interface {
	HouseCreate(ctx context.Context, house models.House) (*models.House, error)
	HouseByID(ctx context.Context, id int) (*models.House, error)
}

type FlatService interface {
	FlatCreate(ctx context.Context, flat models.Flat) (*models.Flat, error)
	FlatUpdate(ctx context.Context, flatId int, newStatus models.Status) (*models.Flat, error)
}

type AuthService interface {
	Login(ctx context.Context, user models.User) (string, error)
	Register(ctx context.Context, user models.User) (*uuid.UUID, error)
	DummyLogin(ctx context.Context, userType models.UserType) (string, error)
}

type Handler struct {
	HouseService HouseService
	FlatService  FlatService
	AuthService  AuthService

	log *slog.Logger
}

type RegResponse struct {
	UserID gen.UserId `json:"user_id"`
}

type LoginResponse struct {
	Token gen.Token `json:"token"`
}

func New(
	houseService HouseService,
	flatService FlatService,
	authService AuthService,
	log *slog.Logger,
) *Handler {
	return &Handler{
		HouseService: houseService,
		FlatService:  flatService,
		AuthService:  authService,
		log:          log,
	}
}

// GetDummyLogin implements gen.ServerInterface.
func (h *Handler) GetDummyLogin(w http.ResponseWriter, r *http.Request, params gen.GetDummyLoginParams) {
	const op = "handlers.GetDummyLogin"

	reqID := middleware.GetReqID(r.Context())
	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", reqID),
		slog.String("type", string(params.UserType)),
	)

	log.Info("attempting getting token")

	userType := params.UserType
	if !(userType == gen.Client || userType == gen.Moderator) {
		log.Warn("params invalide")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := h.AuthService.DummyLogin(context.Background(), models.UserType(params.UserType))
	if err != nil {
		log.Error("internal error", slog.String("err", err.Error()))
		respError(w, r, "что-то пошло не так", http.StatusInternalServerError)
		return
	}

	log.Info("success geting token")
	render.JSON(w, r, LoginResponse{
		Token: token,
	})
}

// GetHouseId implements gen.ServerInterface.
func (h *Handler) GetHouseId(w http.ResponseWriter, r *http.Request, id gen.HouseId) {
	const op = "handlers.GetHouseId"
	reqID := middleware.GetReqID(r.Context())
	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", reqID),
		slog.Int("id", id),
	)

	log.Info("attempting getting house")

	house, err := h.HouseService.HouseByID(context.Background(), id)
	if err != nil {
		log.Error("internal error", slog.String("err", err.Error()))
		respError(w, r, "что-то пошло не так", http.StatusInternalServerError)
		return
	}

	log.Info("success getting house")

	render.JSON(w, r, gen.House{
		Address:   house.Address,
		CreatedAt: &house.CreatedAt,
		Developer: house.Developer,
		Id:        house.UID,
		UpdateAt:  &house.LastFlatAddAt,
		Year:      gen.Year(house.Year),
	})
}

// PostFlatCreate implements gen.ServerInterface.
func (h *Handler) PostFlatCreate(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.PostFlatCreate"
	reqID := middleware.GetReqID(r.Context())
	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", reqID),
	)

	log.Info("attempting create flat")

	var f gen.PostFlatCreateJSONRequestBody
	if err := render.DecodeJSON(r.Body, &f); err != nil {
		log.Warn("incorrect json body", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Debug("success decode json", slog.Any("req", f))

	flat, err := h.FlatService.FlatCreate(context.Background(), models.Flat{
		HouseID: f.HouseId,
		Price:   uint(f.Price),
		Rooms:   uint(*f.Rooms),
		Status:  models.Created,
	})
	if err != nil {
		log.Error("create flat", slog.String("err", err.Error()))

		render.Status(r, http.StatusInternalServerError)
		respError(w, r, "что-то пошло не так", http.StatusInternalServerError)
		return
	}

	log.Info("success adding flat")

	render.JSON(w, r, gen.Flat{
		HouseId: gen.HouseId(flat.HouseID),
		Id:      gen.FlatId(flat.ID),
		Price:   gen.Price(flat.Price),
		Rooms:   gen.Rooms(flat.Rooms),
		Status:  gen.Status(flat.Status),
	})

}

// PostFlatUpdate implements gen.ServerInterface.
func (h *Handler) PostFlatUpdate(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.PostFlatUpdate"

	reqID := middleware.GetReqID(r.Context())

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", reqID),
	)

	log.Info("attempting update status flat")

	userType, ok := r.Context().Value(mv.UserTypeContextKey).(models.UserType)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if userType != models.Moderator {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var f gen.PostFlatUpdateJSONRequestBody
	if err := render.DecodeJSON(r.Body, &f); err != nil {
		log.Warn("incorrect json body", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := validateFlatUpdate(gen.PostFlatUpdateJSONBody(f))
	if err != nil {
		log.Warn("incorrect json", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	flat, err := h.FlatService.FlatUpdate(context.Background(), f.Id, models.Status(*f.Status))
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		respError(w, r, "что-то пошло не так", http.StatusInternalServerError)
		return
	}

	log.Info("success update status flat")

	render.JSON(w, r, gen.Flat{
		HouseId: gen.HouseId(flat.HouseID),
		Id:      gen.FlatId(flat.ID),
		Price:   gen.Price(flat.Price),
		Rooms:   gen.Rooms(flat.Rooms),
		Status:  gen.Status(flat.Status),
	})

}

// PostHouseCreate implements gen.ServerInterface.
func (h *Handler) PostHouseCreate(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.PostHouseCreate"
	var req gen.PostHouseCreateJSONRequestBody

	reqID := middleware.GetReqID(r.Context())

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", reqID),
	)

	log.Info("attempting create house")

	userType := r.Context().Value(mv.UserTypeContextKey).(models.UserType)
	if userType != models.Moderator {
		log.Warn("unautorized", slog.String("user_type", string(userType)))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Warn("invilide json", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := validateHouse(gen.PostHouseCreateJSONBody(req))
	if err != nil {
		log.Warn("validation params", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	house, err := h.HouseService.HouseCreate(context.Background(), models.House{
		Address:   req.Address,
		Year:      uint(req.Year),
		Developer: req.Developer,
	})

	if err != nil {
		log.Error("internal error", slog.String("err", err.Error()))
		respError(w, r, "что-то пошло не так", http.StatusInternalServerError)
		return
	}

	log.Info("create hause", slog.Int("hause_id", house.UID))

	render.JSON(w, r, gen.House{
		Address:   house.Address,
		CreatedAt: &house.CreatedAt,
		Developer: house.Developer,
		Id:        house.UID,
		UpdateAt:  &house.LastFlatAddAt,
		Year:      gen.Year(house.Year),
	})
}

// PostHouseIdSubscribe implements gen.ServerInterface.
func (h *Handler) PostHouseIdSubscribe(w http.ResponseWriter, r *http.Request, id gen.HouseId) {
	panic("unimplemented")
}

// PostLogin implements gen.ServerInterface.
func (h *Handler) PostLogin(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.PostLogin"
	var req gen.PostLoginJSONRequestBody

	reqID := middleware.GetReqID(r.Context())
	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", reqID),
	)

	log.Info("attempting login")

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("decode json", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := validateLogin(gen.PostLoginJSONBody(req))
	if err != nil {
		log.Warn("failed validate", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Debug("success decode json", slog.String("id", req.Id.String()))

	token, err := h.AuthService.Login(context.Background(), models.User{
		ID:           *req.Id,
		PasswordHash: *req.Password,
	})
	if err != nil {
		log.Error("login", slog.String("err", err.Error()))
		respError(w, r, "что-то пошло не так", http.StatusInternalServerError)
		return
	}
	log.Debug("success get token", slog.String("token", token))
	log.Info("success login", slog.String("user_id", req.Id.String()))

	render.JSON(w, r, &LoginResponse{
		Token: token,
	})
}

// PostRegister implements gen.ServerInterface.
func (h *Handler) PostRegister(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.PostRegister"

	reqID := middleware.GetReqID(r.Context())
	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", reqID),
	)

	log.Info("attempting registration")

	var req gen.PostRegisterJSONRequestBody

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Warn("incorect json", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := validateRegister(gen.PostRegisterJSONBody(req))
	if err != nil {
		log.Warn("validate failed", slog.String("err", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := h.AuthService.Register(context.Background(), models.User{
		Email:        string(*req.Email),
		PasswordHash: *req.Password,
		UserType:     models.UserType(*req.UserType),
	})

	if err != nil {
		log.Error("internal error", slog.String("err", err.Error()))
		respError(w, r, "что-то пошло не так", http.StatusInternalServerError)
		return
	}

	log.Info("success registation", slog.String("user_id", id.String()))
	render.JSON(w, r, &RegResponse{UserID: *id})
}

func respError(w http.ResponseWriter, r *http.Request, msg string, statusCode int) {
	requestID := middleware.GetReqID(r.Context())
	render.Status(r, statusCode)
	render.JSON(w, r, gen.N5xx{
		Code:      &statusCode,
		Message:   msg,
		RequestId: &requestID,
	})
}

func validateHouse(house gen.PostHouseCreateJSONBody) error {

	rules := map[string]string{
		"Address": "min=10",
		"Year":    "min=1900",
	}

	validate.RegisterStructValidationMapRules(rules, gen.PostHouseCreateJSONBody{})
	err := validate.Struct(house)
	return err
}

func validateRegister(user gen.PostRegisterJSONBody) error {

	rules := map[string]string{
		"Email":    "required,email",
		"UserType": "eq=client|eq=moderator",
		"Password": "min=6",
	}

	validate.RegisterStructValidationMapRules(rules, gen.PostRegisterJSONBody{})
	err := validate.Struct(user)

	return err
}

func validateLogin(user gen.PostLoginJSONBody) error {
	rules := map[string]string{
		"Id":       "uuid",
		"Password": "min=6",
	}

	validate.RegisterStructValidationMapRules(rules, gen.PostRegisterJSONBody{})
	err := validate.Struct(user)

	return err
}

func validateFlatUpdate(user gen.PostFlatUpdateJSONBody) error {
	rules := map[string]string{
		"Id":     "uuid",
		"Status": "eq=approved|eq=created|eq=declined|eq=on moderation",
	}

	validate.RegisterStructValidationMapRules(rules, gen.PostRegisterJSONBody{})
	err := validate.Struct(user)

	return err
}
