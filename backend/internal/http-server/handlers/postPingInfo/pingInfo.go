package postPingInfo

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"time"
)

type PingInfo struct {
	IPAddr     string    `json:"ipAddr"`
	PingTime   time.Time `json:"pingTime"`
	PacketLoss float64   `json:"packetLoss"`
}
type Req struct {
	Stats []PingInfo `json:"stats"`
}

type Resp struct {
	Stats string `json:"status"`
	Error string `json:"error,omitempty"`
}
type PingInfoSaver interface {
	SaveNewInfo(pingInfo []PingInfo) error
}

func New(log *slog.Logger, saver PingInfoSaver) echo.HandlerFunc {
	return func(c echo.Context) error {
		functionName := "http-server/handlers/postPingInfo"
		log = log.With(slog.String("funcName", functionName))
		var req []PingInfo
		err := json.NewDecoder(c.Request().Body).Decode(&req)
		if err != nil {
			log.Error("Failed to parse request body: %w", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to save pingInfo"})
		}

		log.Info("Request body decoded", slog.Any("request", req))

		//dateValue := req[0].PingTime.Format("2006-01-02 15:04:05")
		//pl := req[0].PacketLoss
		//fmt.Println("\n", dateValue, pl)

		err = saver.SaveNewInfo(req)
		if err != nil {
			log.Error("Failed to save in db: %w", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"status": "Error", "error": "Failed to save pingInfo"})
		}
		resp := Resp{
			Stats: "OK",
		}
		return c.JSON(200, resp)
	}
}
