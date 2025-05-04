package middleware

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/color"
)

type CustomLogger struct {
	// Add fields if needed
}

func NewCustomLogger() *CustomLogger {
	return &CustomLogger{}
}

func (l *CustomLogger) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			latency := time.Since(start)

			// Get status color based on status code
			statusColor := color.White
			switch {
			case res.Status >= 500:
				statusColor = color.Red
			case res.Status >= 400:
				statusColor = color.Yellow
			case res.Status >= 300:
				statusColor = color.Cyan
			default:
				statusColor = color.Green
			}

			// Format the log output with colored components
			logMessage := fmt.Sprintf("\n[%s] %s | %s | %s%s | %s | %s | %s | %d bytes\n",
				time.Now().Format("2006-01-02 15:04:05"),
				statusColor(fmt.Sprint(res.Status)),
				color.Blue(req.Method),
				color.Magenta(req.Host),
				color.Blue(req.RequestURI),
				color.Yellow(latency.String()),
				color.Cyan(req.UserAgent()),
				c.RealIP(),
				res.Size,
			)

			fmt.Print(logMessage)
			return nil
		}
	}
}
