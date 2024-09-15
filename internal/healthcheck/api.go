package healthcheck

import routing "github.com/go-ozzo/ozzo-routing/v2"

func RegisterHealthCheckHandler(r *routing.Router) {
	r.To("GET", "/healthcheck", getHealthCheckHandler())
}

func getHealthCheckHandler() routing.Handler {
	return func(c *routing.Context) error {
		return c.Write("OK")
	}
}
