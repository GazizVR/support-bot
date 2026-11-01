package bot

import (
	"context"

	"github.com/mymmrac/telego"
)

type Route struct {
	Match  func(telego.Update) bool
	Handle func(ctx context.Context, upd telego.Update) error
}

type Router struct {
	routes []Route
}

func NewRouter(routes []Route) *Router {
	return &Router{
		routes: routes,
	}
}

func (r *Router) Route(
	ctx context.Context,
	update telego.Update,
) error {
	for _, r := range r.routes {
		if r.Match(update) {
			if err := r.Handle(ctx, update); err != nil {
				return err
			}
		}
	}
	return nil
}
