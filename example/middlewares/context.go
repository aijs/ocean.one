package middlewares

import (
	"net/http"

	"cloud.google.com/go/spanner"
	"github.com/MixinNetwork/ocean.one/example/durable"
	"github.com/MixinNetwork/ocean.one/example/session"
	"github.com/newrelic/go-agent"
	"github.com/unrolled/render"
)

func Context(handler http.Handler, spannerClient *spanner.Client, limiter *durable.Limiter, render *render.Render) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var nT newrelic.Transaction = nil
		if v, ok := w.(newrelic.Transaction); ok {
			nT = v
		}
		db := durable.WrapDatabase(spannerClient, nT)
		ctx := session.WithRequest(r.Context(), r)
		ctx = session.WithDatabase(ctx, db)
		ctx = session.WithLimiter(ctx, limiter)
		ctx = session.WithRender(ctx, render)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
