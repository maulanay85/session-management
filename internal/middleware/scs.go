package middleware

import (
	"context"
	"net/http"
	"scs-session/internal/config"
	"scs-session/internal/usecase"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gin-gonic/gin"
)

func LoadAndSave(sessionManager *scs.SessionManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// var newToken string
		var maxAge int
		cookie, err := c.Cookie(sessionManager.Cookie.Name)
		if err != nil {
			cookie = ""
		}
		ctx, err := sessionManager.Load(c.Request.Context(), cookie)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load session"})
			c.Abort()
			return
		}
		c.Request = c.Request.WithContext(ctx)

		if c.Request.URL.Path == "/favicon.ico" {
			// Allow the request to pass without further middleware processing
			return
		}

		c.Next()

		if c.Request.MultipartForm != nil {
			c.Request.MultipartForm.RemoveAll()
		}

		switch sessionManager.Status(ctx) {
		case scs.Modified:
			t, expiry, err := sessionManager.Commit(ctx)
			cookie = t
			if expiry.IsZero() {
				maxAge = -1
			} else if sessionManager.Cookie.Persist {
				maxAge = int(time.Until(expiry).Seconds() + 1)
			}
			if err != nil {
				c.Abort()
				return
			}

		case scs.Destroyed:
			cookie = ""
			maxAge = -1

		case scs.Unmodified:
			if cookie != "" {
				maxAge = int(time.Until(time.Now().Add(sessionManager.Lifetime).Add(1 * time.Second)).Seconds())
			} else {
				maxAge = -1
			}
		}

		// if newToken != "" {
		cr := http.Cookie{
			Name:     sessionManager.Cookie.Name,
			Value:    cookie,
			MaxAge:   maxAge,
			Path:     sessionManager.Cookie.Path,
			Domain:   sessionManager.Cookie.Domain,
			Secure:   sessionManager.Cookie.Secure,
			HttpOnly: sessionManager.Cookie.HttpOnly,
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(c.Writer, &cr)
		// }

		response, exist := c.Get("data")
		if exist {
			c.JSON(c.Writer.Status(), response)
			c.Done()
		}
	}
}

func SessionMiddleware(config config.Config, sessionManager *scs.SessionManager, sessionUsecase usecase.SessionUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("token")
		if err != nil {
			sessionManager.Destroy(c.Request.Context())
			c.Status(401)
			c.Abort()
			return
		}
		session, err := sessionUsecase.GetByToken(c.Request.Context(), cookie)
		if err != nil {
			c.Status(500)
			c.Abort()
			return
		}
		ctxKey := "id"
		ctx := context.WithValue(c.Request.Context(), ctxKey, session.UserID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()

		d, err := sessionUsecase.Validate(c.Request.Context(), cookie)
		if err != nil {
			c.Status(500)
			c.Abort()
			return
		}

		cr := &http.Cookie{
			Name:     "token",
			Value:    d.Token,
			MaxAge:   int(time.Until(d.ExpiredAt)) + 1,
			Domain:   "localhost",
			Path:     "/",
			Secure:   false,
			HttpOnly: false,
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(c.Writer, cr)
	}
}
