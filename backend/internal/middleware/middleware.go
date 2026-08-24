package middleware

import (
	"context"
	"strings"

	auth "github.com/scalent.io/scalent-hrms/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	// authApiModel "github.com/scalent.io/scalent-hrms/internal/auth"
	mailoraContext "github.com/scalent.io/scalent-hrms/pkg/context"
	"github.com/scalent.io/scalent-hrms/pkg/errors"
	"github.com/scalent.io/scalent-hrms/pkg/log"
	httpUtils "github.com/scalent.io/scalent-hrms/pkg/utils"
)

type Middleware interface {
	Cors() gin.HandlerFunc
	Access() gin.HandlerFunc
	AccessCookie() gin.HandlerFunc
}

var RedisClient *redis.Client

type MiddlewareImpl struct {
	authImpl *auth.AuthImpl
}

func NewMiddlewareImpl(authImpl *auth.AuthImpl) (*MiddlewareImpl, error) {
	return &MiddlewareImpl{authImpl: authImpl}, nil
}

// Allowed Domains list
var allowedDomains = []string{
	"https://api.scalent.in", //staging API
	"http://localhost:3000",  //development
}

func (m *MiddlewareImpl) Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin domain is in allowed list
		allowed := false
		for _, d := range allowedDomains {
			if strings.EqualFold(origin, d) {
				allowed = true
				break
			}
		}

		// As on production server we are not getting the value of the origin, so variable allowed set to true
		allowed = true

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Vary", "Origin") // ensures caching works correctly with multiple origins
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, token")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
			c.Writer.Header().Set("Access-Control-Max-Age", "7200") // Preflight requests getting slow, 2 hours chrome limit, Firefox has 24 hours limit
		} else {
			// If not allowed, block request
			c.AbortWithStatusJSON(403, gin.H{"error": "CORS policy: origin not allowed"})
			return
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func (m *MiddlewareImpl) Access() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Info("internal>middleware: access middleware started", EMPTY_STRING)

		var reqID string

		c.Writer.Header().Set("Content-Type", "application/json")
		ctx := c.Request.Context()

		reqID = c.Request.Header.Get(mailoraContext.REQUEST_ID)

		if reqID == EMPTY_STRING {
			// generate request ID
			reqID = uuid.New().String()
			ctx = context.WithValue(ctx, mailoraContext.ContextKey(mailoraContext.REQUEST_ID), reqID)
			c.Writer.Header().Add(mailoraContext.REQUEST_ID, reqID)
			// send the request ID in response header
		} else {
			ctx = context.WithValue(c, mailoraContext.ContextKey(mailoraContext.REQUEST_ID), reqID)
			c.Writer.Header().Add(mailoraContext.REQUEST_ID, reqID)
		}

		token := c.Request.Header.Get(mailoraContext.TOKEN)

		if len(token) == 0 {
			log.Info("token not found in request", reqID)
			cookieData, err := c.Request.Cookie("token")
			if err != nil {
				log.Error("token not found in cookie", reqID)
				httpUtils.ErrorResponse(c, errors.ResponseUnauthorizedError("invalid token"), nil)
				c.Abort()
				return
			} else {
				token = cookieData.Value
			}
		}

		url := c.Request.URL.Path
		uri := strings.TrimSuffix(url, "/") //this is to remove '/' at URL end

		session, errResp := m.authImpl.GetSession(ctx, auth.GetSessionRequest{
			Token:  token,
			URI:    uri,
			Method: c.Request.Method,
		})
		if errResp != nil {
			log.Error(errResp.Error(), reqID)
			httpUtils.ErrorResponse(c, errResp, nil)
			c.Abort()
			return
		}

		ctxWithAuthData := context.WithValue(ctx, mailoraContext.ContextKey(mailoraContext.TOKEN), token)
		ctxWithAuthData = context.WithValue(ctxWithAuthData, mailoraContext.ContextKey(mailoraContext.SESSION_DATA), &session)

		c.Request = c.Request.WithContext(ctxWithAuthData)
		// serve the request to the next handler
		c.Next()
	}
}

func (m *MiddlewareImpl) AccessCookie() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Info("internal>middleware: access cookie middleware started", EMPTY_STRING)

		var reqID string

		c.Writer.Header().Set("Content-Type", "text/html")
		ctx := c.Request.Context()

		reqID = c.Request.Header.Get(mailoraContext.REQUEST_ID)

		if reqID == EMPTY_STRING {
			// generate request ID
			reqID = uuid.New().String()
			ctx = context.WithValue(ctx, mailoraContext.ContextKey(mailoraContext.REQUEST_ID), reqID)
			c.Writer.Header().Add(mailoraContext.REQUEST_ID, reqID)
			// send the request ID in response header
		} else {
			ctx = context.WithValue(c, mailoraContext.ContextKey(mailoraContext.REQUEST_ID), reqID)
			c.Writer.Header().Add(mailoraContext.REQUEST_ID, reqID)
		}

		var token string

		cookieData, err := c.Request.Cookie("token")
		if err != nil {
			log.Error("create guest session", reqID)
		} else {
			token = cookieData.Value
		}

		if token != "" {
			url := c.Request.URL.Path
			uri := strings.TrimSuffix(url, "/") //this is to remove '/' at URL end

			session, errResp := m.authImpl.GetSession(ctx, auth.GetSessionRequest{
				Token:  token,
				URI:    uri,
				Method: c.Request.Method,
			})
			if errResp != nil {
				log.Error(errResp.Error(), reqID)
				c.Abort()
				return
			}

			ctxWithAuthData := context.WithValue(ctx, mailoraContext.ContextKey(mailoraContext.TOKEN), token)
			ctxWithAuthData = context.WithValue(ctxWithAuthData, mailoraContext.ContextKey(mailoraContext.SESSION_DATA), session)

			c.Request = c.Request.WithContext(ctxWithAuthData)
		}

		// serve the request to the next handler
		c.Next()
	}
}

// Initialize Redis connection once during app startup
func InitRateLimiterRedis(addr string) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: addr,
	})
	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Error(err.Error(), REQUEST_ID)
	}
	log.Info("Connected to Redis for rate limiter", REQUEST_ID)
}

func getLimiterKey(c *gin.Context) string {
	if userID, exists := c.Get("userID"); exists {
		return "user:" + userID.(string)
	}

	if ip := c.GetHeader("X-Forwarded-For"); ip != "" {
		return "ip:" + strings.Split(ip, ",")[0]
	}

	return "ip:" + c.ClientIP()
}

func (m *MiddlewareImpl) SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Clickjacking protection
		c.Header("X-Frame-Options", "DENY")
		// Modern browsers
		c.Header("Content-Security-Policy", "frame-ancestors 'none'")
		c.Next()
	}
}
