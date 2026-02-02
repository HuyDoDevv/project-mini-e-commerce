package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"project-mini-e-commerce/pkg/logger"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *CustomResponseWriter) Write(data []byte) (n int, err error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

func LoggerMiddleware(httpLogger *zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		contentType := ctx.GetHeader("Content-Type")
		requestBody := make(map[string]any)
		var formFiles []map[string]any

		if strings.HasPrefix(contentType, "multipart/form-data") {
			if err := ctx.Request.ParseMultipartForm(32 << 20); err == nil && ctx.Request.MultipartForm != nil {
				// value
				for key, vals := range ctx.Request.MultipartForm.Value {
					if len(vals) == 1 {
						requestBody[key] = vals[0]
					} else {
						requestBody[key] = vals
					}
				}
				// file
				for field, files := range ctx.Request.MultipartForm.File {
					for _, f := range files {
						formFiles = append(formFiles, map[string]any{
							"field":        field,
							"file_name":    f.Filename,
							"size":         formatFileSize(f.Size),
							"content-type": f.Header.Get("Content-Type"),
						})
					}
				}
				if len(formFiles) > 0 {
					requestBody["form_file"] = formFiles
				}
			}
		} else {
			bodyBytes, err := io.ReadAll(ctx.Request.Body)
			if err != nil {
				httpLogger.Error().Err(err).Msg("Field to read request body")
			}
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			if strings.HasPrefix(contentType, "application/json") {
				_ = json.Unmarshal(bodyBytes, &requestBody)
			} else {
				values, _ := url.ParseQuery(string(bodyBytes))
				for key, val := range values {
					if len(val) == 1 {
						requestBody[key] = val[0]
					} else {
						requestBody[key] = val
					}
				}
			}
		}

		customWriter := &CustomResponseWriter{
			ResponseWriter: ctx.Writer,
			body:           bytes.NewBufferString(""),
		}
		ctx.Writer = customWriter

		ctx.Next()
		duration := time.Since(start)

		responseContentType := ctx.Writer.Header().Get("Content-Type")
		responseBodyRaw := customWriter.body.String()
		var responseBodyParse interface{}

		if strings.HasPrefix(responseContentType, "image/") {
			responseBodyParse = "[BINARY IMAGE]"
		} else if strings.HasPrefix(responseContentType, "application/json") ||
			strings.HasPrefix(responseContentType, "{") ||
			strings.HasPrefix(responseContentType, "}") {
			if err := json.Unmarshal([]byte(responseBodyRaw), &responseBodyParse); err != nil {
				responseBodyParse = responseBodyRaw
			}
		} else {
			responseBodyParse = responseBodyRaw
		}

		logEvent := httpLogger.Info()

		logEvent.
			Str("trace_id", logger.GetTraceId(ctx.Request.Context())).
			Str("method", ctx.Request.Method).
			Str("path", ctx.Request.URL.Path).
			Str("query", ctx.Request.URL.RawQuery).
			Str("client_ip", ctx.ClientIP()).
			Str("user_agent", ctx.Request.UserAgent()).
			Str("refer", ctx.Request.Referer()).
			Str("protocol", ctx.Request.Proto).
			Str("host", ctx.Request.Host).
			Str("remote_addr", ctx.Request.RemoteAddr).
			Str("request_uri", ctx.Request.RequestURI).
			Int64("content_length", ctx.Request.ContentLength).
			Interface("header", ctx.Request.Header).
			Interface("request_body", requestBody).
			Interface("response_body", responseBodyParse).
			Int("status_code", ctx.Writer.Status()).
			Int64("duration_ms", duration.Milliseconds()).
			Msg("Logger https")
	}
}

func formatFileSize(size int64) string {
	switch {
	case size >= 1<<20:
		return fmt.Sprintf("%.2f MB", float64(size)/(1<<20))
	case size >= 1<<10:
		return fmt.Sprintf("%.2f KB", float64(size)/(1<<20))
	default:
		return fmt.Sprintf("%d B", size)
	}
}
