package dojo

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Application struct {
	Configuration      DefaultConfiguration
	MiddlewareRegistry *MiddlewareRegistry `json:"-"`
	router             *mux.Router
	Logger             *logrus.Logger
	root               *Application
	SessionStore       *sessions.CookieStore
	Auth               *Authentication
}

// New creates a new instance of Application
func New(conf DefaultConfiguration) *Application {

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logger.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	logger.SetLevel(logrus.DebugLevel)

	cookieStore := sessions.NewCookieStore([]byte(conf.Session.Secret))
	cookieStore.Options.HttpOnly = true
	if conf.App.Environment == "production" {
		cookieStore.Options.Secure = true
	}

	app := &Application{
		Configuration:      conf,
		router:             mux.NewRouter(),
		MiddlewareRegistry: NewMiddlewareRegistry(),
		Logger:             logger,
		SessionStore:       cookieStore,
	}

	app.Auth = NewAuthentication(app)

	return app
}

func (app *Application) Serve() {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Configuration.App.Port),
		Handler:      app.router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		fmt.Printf("signal %s", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		// return err
		app.Logger.Error(err)
	}

	err = <-shutdownError
	if err != nil {
		app.Logger.Error(err)
	}

}
