// Package gen provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package gen

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// Defines values for Status.
const (
	Approved     Status = "approved"
	Created      Status = "created"
	Declined     Status = "declined"
	OnModeration Status = "on moderation"
)

// Defines values for UserType.
const (
	Client    UserType = "client"
	Moderator UserType = "moderator"
)

// Address Адрес дома
type Address = string

// Date Дата + время
type Date = time.Time

// Developer Застройщик
type Developer = string

// Email Email пользователя
type Email = openapi_types.Email

// Flat Квартира
type Flat struct {
	// HouseId Идентификатор дома
	HouseId HouseId `json:"house_id"`

	// Id Идентификатор квартиры
	Id FlatId `json:"id"`

	// Price Цена квартиры в у.е.
	Price Price `json:"price"`

	// Rooms Количество комнат в квартире
	Rooms Rooms `json:"rooms"`

	// Status Статус квартиры
	Status Status `json:"status"`
}

// FlatId Идентификатор квартиры
type FlatId = int

// House Дом
type House struct {
	// Address Адрес дома
	Address Address `json:"address"`

	// CreatedAt Дата + время
	CreatedAt *Date `json:"created_at,omitempty"`

	// Developer Застройщик
	Developer *Developer `json:"developer"`

	// Id Идентификатор дома
	Id HouseId `json:"id"`

	// UpdateAt Дата + время
	UpdateAt *Date `json:"update_at,omitempty"`

	// Year Год постройки дома
	Year Year `json:"year"`
}

// HouseId Идентификатор дома
type HouseId = int

// Password Пароль пользователя
type Password = string

// Price Цена квартиры в у.е.
type Price = int

// Rooms Количество комнат в квартире
type Rooms = int

// Status Статус квартиры
type Status string

// Token Авторизационный токен
type Token = string

// UserId Идентификатор пользователя
type UserId = openapi_types.UUID

// UserType Тип пользователя
type UserType string

// Year Год постройки дома
type Year = int

// N5xx defines model for 5xx.
type N5xx struct {
	// Code Код ошибки. Предназначен для классификации проблем и более быстрого решения проблем.
	Code *int `json:"code,omitempty"`

	// Message Описание ошибки
	Message string `json:"message"`

	// RequestId Идентификатор запроса. Предназначен для более быстрого поиска проблем.
	RequestId *string `json:"request_id,omitempty"`
}

// GetDummyLoginParams defines parameters for GetDummyLogin.
type GetDummyLoginParams struct {
	UserType UserType `form:"user_type" json:"user_type"`
}

// PostFlatCreateJSONBody defines parameters for PostFlatCreate.
type PostFlatCreateJSONBody struct {
	// HouseId Идентификатор дома
	HouseId HouseId `json:"house_id"`

	// Price Цена квартиры в у.е.
	Price Price `json:"price"`

	// Rooms Количество комнат в квартире
	Rooms *Rooms `json:"rooms,omitempty"`
}

// PostFlatUpdateJSONBody defines parameters for PostFlatUpdate.
type PostFlatUpdateJSONBody struct {
	// Id Идентификатор квартиры
	Id FlatId `json:"id"`

	// Status Статус квартиры
	Status *Status `json:"status,omitempty"`
}

// PostHouseCreateJSONBody defines parameters for PostHouseCreate.
type PostHouseCreateJSONBody struct {
	// Address Адрес дома
	Address Address `json:"address"`

	// Developer Застройщик
	Developer *Developer `json:"developer"`

	// Year Год постройки дома
	Year Year `json:"year"`
}

// PostHouseIdSubscribeJSONBody defines parameters for PostHouseIdSubscribe.
type PostHouseIdSubscribeJSONBody struct {
	// Email Email пользователя
	Email Email `json:"email"`
}

// PostLoginJSONBody defines parameters for PostLogin.
type PostLoginJSONBody struct {
	// Id Идентификатор пользователя
	Id *UserId `json:"id,omitempty"`

	// Password Пароль пользователя
	Password *Password `json:"password,omitempty"`
}

// PostRegisterJSONBody defines parameters for PostRegister.
type PostRegisterJSONBody struct {
	// Email Email пользователя
	Email *Email `json:"email,omitempty"`

	// Password Пароль пользователя
	Password *Password `json:"password,omitempty"`

	// UserType Тип пользователя
	UserType *UserType `json:"user_type,omitempty"`
}

// PostFlatCreateJSONRequestBody defines body for PostFlatCreate for application/json ContentType.
type PostFlatCreateJSONRequestBody PostFlatCreateJSONBody

// PostFlatUpdateJSONRequestBody defines body for PostFlatUpdate for application/json ContentType.
type PostFlatUpdateJSONRequestBody PostFlatUpdateJSONBody

// PostHouseCreateJSONRequestBody defines body for PostHouseCreate for application/json ContentType.
type PostHouseCreateJSONRequestBody PostHouseCreateJSONBody

// PostHouseIdSubscribeJSONRequestBody defines body for PostHouseIdSubscribe for application/json ContentType.
type PostHouseIdSubscribeJSONRequestBody PostHouseIdSubscribeJSONBody

// PostLoginJSONRequestBody defines body for PostLogin for application/json ContentType.
type PostLoginJSONRequestBody PostLoginJSONBody

// PostRegisterJSONRequestBody defines body for PostRegister for application/json ContentType.
type PostRegisterJSONRequestBody PostRegisterJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /dummyLogin)
	GetDummyLogin(w http.ResponseWriter, r *http.Request, params GetDummyLoginParams)

	// (POST /flat/create)
	PostFlatCreate(w http.ResponseWriter, r *http.Request)

	// (POST /flat/update)
	PostFlatUpdate(w http.ResponseWriter, r *http.Request)

	// (POST /house/create)
	PostHouseCreate(w http.ResponseWriter, r *http.Request)

	// (GET /house/{id})
	GetHouseId(w http.ResponseWriter, r *http.Request, id HouseId)

	// (POST /house/{id}/subscribe)
	PostHouseIdSubscribe(w http.ResponseWriter, r *http.Request, id HouseId)

	// (POST /login)
	PostLogin(w http.ResponseWriter, r *http.Request)

	// (POST /register)
	PostRegister(w http.ResponseWriter, r *http.Request)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// (GET /dummyLogin)
func (_ Unimplemented) GetDummyLogin(w http.ResponseWriter, r *http.Request, params GetDummyLoginParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (POST /flat/create)
func (_ Unimplemented) PostFlatCreate(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (POST /flat/update)
func (_ Unimplemented) PostFlatUpdate(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (POST /house/create)
func (_ Unimplemented) PostHouseCreate(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (GET /house/{id})
func (_ Unimplemented) GetHouseId(w http.ResponseWriter, r *http.Request, id HouseId) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (POST /house/{id}/subscribe)
func (_ Unimplemented) PostHouseIdSubscribe(w http.ResponseWriter, r *http.Request, id HouseId) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (POST /login)
func (_ Unimplemented) PostLogin(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// (POST /register)
func (_ Unimplemented) PostRegister(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// GetDummyLogin operation middleware
func (siw *ServerInterfaceWrapper) GetDummyLogin(w http.ResponseWriter, r *http.Request) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetDummyLoginParams

	// ------------- Required query parameter "user_type" -------------

	if paramValue := r.URL.Query().Get("user_type"); paramValue != "" {

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "user_type"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "user_type", r.URL.Query(), &params.UserType)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "user_type", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetDummyLogin(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// PostFlatCreate operation middleware
func (siw *ServerInterfaceWrapper) PostFlatCreate(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostFlatCreate(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// PostFlatUpdate operation middleware
func (siw *ServerInterfaceWrapper) PostFlatUpdate(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostFlatUpdate(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// PostHouseCreate operation middleware
func (siw *ServerInterfaceWrapper) PostHouseCreate(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostHouseCreate(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetHouseId operation middleware
func (siw *ServerInterfaceWrapper) GetHouseId(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "id" -------------
	var id HouseId

	err = runtime.BindStyledParameterWithOptions("simple", "id", chi.URLParam(r, "id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetHouseId(w, r, id)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// PostHouseIdSubscribe operation middleware
func (siw *ServerInterfaceWrapper) PostHouseIdSubscribe(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "id" -------------
	var id HouseId

	err = runtime.BindStyledParameterWithOptions("simple", "id", chi.URLParam(r, "id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostHouseIdSubscribe(w, r, id)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// PostLogin operation middleware
func (siw *ServerInterfaceWrapper) PostLogin(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostLogin(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// PostRegister operation middleware
func (siw *ServerInterfaceWrapper) PostRegister(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostRegister(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/dummyLogin", wrapper.GetDummyLogin)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/flat/create", wrapper.PostFlatCreate)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/flat/update", wrapper.PostFlatUpdate)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/house/create", wrapper.PostHouseCreate)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/house/{id}", wrapper.GetHouseId)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/house/{id}/subscribe", wrapper.PostHouseIdSubscribe)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/login", wrapper.PostLogin)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/register", wrapper.PostRegister)
	})

	return r
}
