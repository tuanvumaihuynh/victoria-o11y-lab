package http

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/chi/v5"
)

func (s *Service) registerDocs(r chi.Router, api huma.API) error {
	specBytes, err := api.OpenAPI().YAML()
	if err != nil {
		return fmt.Errorf("get OpenAPI spec: %w", err)
	}
	template := createSwaggerHTML("/docs/openapi.yml")

	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(template)); err != nil {
			s.logger.ErrorContext(r.Context(), "Failed to write API docs template", slog.Any("error", err))
		}
	})

	r.Get("/docs/openapi.yml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(specBytes); err != nil {
			s.logger.ErrorContext(r.Context(), "Failed to write OpenAPI spec", slog.Any("error", err))
		}
	})

	return nil
}

// createSwaggerHTML returns the HTML template for Swagger API documentation
func createSwaggerHTML(specPath string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <meta name="description" content="SwaggerUI" />
  <title>Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.29.3/swagger-ui.css" />
  <link rel="icon" type="image/png" href="https://static1.smartbear.co/swagger/media/assets/swagger_fav.png" sizes="32x32" />
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5.29.3/swagger-ui-bundle.js" crossorigin></script>
<script>
  window.onload = () => {
    window.ui = SwaggerUIBundle({
      url: '%s',
      dom_id: '#swagger-ui',
      deepLinking: true,
	  showExtensions: true,
	  showCommonExtensions: true,
    });
  };
</script>
</body>
</html>
`, specPath)
}
