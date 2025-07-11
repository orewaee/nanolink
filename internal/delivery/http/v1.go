package delivery

import "net/http"

func (c *RestController) muxV1() http.Handler {
	v1 := http.NewServeMux()

	return http.StripPrefix("/v1", v1)
}
