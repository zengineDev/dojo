package dojo

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Application struct {
	Configuration DefaultConfiguration
	Middleware    *MiddlewareStack `json:"-"`
	router        *mux.Router
}

// New creates a new instance of Application
func New(conf DefaultConfiguration) *Application {

	app := &Application{
		Configuration: conf,
		router:        mux.NewRouter(),
		Middleware:    newMiddlewareStack(),
	}

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
		//return err
	}

	err = <-shutdownError
	if err != nil {
		//return err
	}

}
