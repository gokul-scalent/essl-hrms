package auth

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/scalent.io/scalent-hrms/entity"
	"github.com/scalent.io/scalent-hrms/entity/commonConstants"
	"github.com/scalent.io/scalent-hrms/pkg/log"

	"github.com/casbin/casbin/v2"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"

	mailoraContext "github.com/scalent.io/scalent-hrms/pkg/context"
	"github.com/scalent.io/scalent-hrms/pkg/errors"

	"github.com/google/uuid"
)

const (
	TOKEN_PREFIX = `scalent-hrms:`
)

type GetSessionRequest struct {
	Token  string
	Method string
	URI    string
}

type CreateSessionRequest struct {
	UserID       int
	CustomerID   int
	TechnicianID int
	Role         string
	FirstName    string
	LastName     string
	TTL          int
}

type UpdateSessionRequest struct {
	UserID int
	Role   string
	Token  string
	TTL    uint64
}

type GetSessionResponse struct {
	UserID int
	Role   string
	TTL    int
}

type CreateSessionResponse struct {
	Token string
}

type AuthImpl struct {
	rdb      *redis.Client
	enforcer *casbin.Enforcer
	sqlDB    *sqlx.DB
}

func NewAuthImpl(rdb *redis.Client, sqlDB *sqlx.DB, enforcer *casbin.Enforcer) (*AuthImpl, error) {
	return &AuthImpl{
		rdb:      rdb,
		enforcer: enforcer,
		sqlDB:    sqlDB,
	}, nil
}

func (authImpl *AuthImpl) CreateSession(ctx context.Context, request CreateSessionRequest) (CreateSessionResponse, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("internal>auth: CreateSession started", reqID)

	response := CreateSessionResponse{}
	token := uuid.New().String() + uuid.New().String()

	encodedSession, err := json.Marshal(&request)
	if err != nil {
		log.Error(err.Error(), reqID)
		return response, errors.ResponseInternalServerError(errors.COULD_NOT_CREATE_SESSION)
	}

	// store the session in top level key value
	_, err = authImpl.rdb.Set(ctx, TOKEN_PREFIX+token, encodedSession, time.Minute*time.Duration(request.TTL)).Result()
	if err != nil {
		log.Error(err.Error(), reqID)
		return response, errors.ResponseInternalServerError(errors.COULD_NOT_CREATE_SESSION)
	}

	response.Token = token
	log.Info("internal>auth: CreateSession completed", reqID)
	return response, nil
}

func (authImpl *AuthImpl) GetSession(ctx context.Context, request GetSessionRequest) (entity.Session, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("internal>auth: GetSession started", reqID)

	sessionEntity := entity.Session{}

	sessionByte, err := authImpl.rdb.Get(ctx, TOKEN_PREFIX+request.Token).Bytes()
	if err != nil {
		log.Error(err.Error(), reqID)
		return sessionEntity, errors.ResponseLoginTimeoutError(errors.SESSION_ERROR)
	}

	authImpl.rdb.Expire(ctx, TOKEN_PREFIX+request.Token, time.Minute*time.Duration(commonConstants.SESSION_TTL))

	err = json.Unmarshal(sessionByte, &sessionEntity)
	if err != nil {
		log.Error(err.Error(), reqID)
		return sessionEntity, errors.ResponseLoginTimeoutError(errors.SESSION_ERROR)
	}

	log.Info(
		"Role="+sessionEntity.Role+
			", URI="+request.URI+
			", Method="+request.Method,
		reqID,
	)
	aclResponse, aclErr := authImpl.enforcer.Enforce(sessionEntity.Role, request.URI, request.Method)
	if aclErr != nil || !aclResponse {
		log.Info("internal>auth: "+"authorization failed while user :"+strconv.Itoa(sessionEntity.UserID)+" with Role :"+sessionEntity.Role+" trying to access:"+request.URI+" with method: "+request.Method, reqID)
		return sessionEntity, errors.ResponseUnauthorizedError(errors.MSG_ACCESS_DENIED)
	}

	log.Info("internal>auth: GetSession completed", reqID)
	return sessionEntity, nil
}

func (authImpl *AuthImpl) DeleteSession(ctx context.Context) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("internal>auth: DeleteSession started", reqID)

	token, err := mailoraContext.GetTokenFromContext(ctx)
	if err != nil {
		log.Error(err.Error(), reqID)
		return errors.ResponseInternalServerError(errors.INTERNAL_SERVER_ERROR)
	}

	err = authImpl.rdb.Del(ctx, TOKEN_PREFIX+token).Err()
	if err != nil {
		log.Error(err.Error(), reqID)
		return errors.ResponseInternalServerError(errors.DELETE_SESSION_ERR)
	}

	log.Info("internal>auth: DeleteSession completed", reqID)
	return nil
}

func (authImpl *AuthImpl) UpdateSession(ctx context.Context, request UpdateSessionRequest) (CreateSessionResponse, errors.Response) {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("internal>auth: CreateSession started", reqID)

	response := CreateSessionResponse{}

	encodedSession, err := json.Marshal(&request)
	if err != nil {
		log.Error(err.Error(), reqID)
		return response, errors.ResponseInternalServerError(errors.COULD_NOT_CREATE_SESSION)
	}

	// store the session in top level key value
	_, err = authImpl.rdb.Set(ctx, TOKEN_PREFIX+request.Token, encodedSession, time.Minute*time.Duration(request.TTL)).Result()
	if err != nil {
		log.Error(err.Error(), reqID)
		return response, errors.ResponseInternalServerError(errors.COULD_NOT_CREATE_SESSION)
	}

	log.Info("internal>auth: CreateSession completed", reqID)
	return response, nil
}

// write a function to remove the prevoious session token from redis
func (authImpl *AuthImpl) RemovePreviousSessionToken(ctx context.Context, token string) errors.Response {
	reqID, _ := mailoraContext.GetRequestIDFromContext(ctx)
	log.Info("internal>auth: RemovePreviousSessionToken started", reqID)

	err := authImpl.rdb.Del(ctx, TOKEN_PREFIX+token).Err()
	if err != nil {
		log.Error(err.Error(), reqID)
		return errors.ResponseInternalServerError(errors.DELETE_SESSION_ERR)
	}

	log.Info("internal>auth: RemovePreviousSessionToken completed", reqID)
	return nil
}
