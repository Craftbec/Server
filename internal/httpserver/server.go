package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/Craftbec/Server/internal/storage"
	"github.com/gorilla/mux"
)

type ResponseSet struct {
	ErrorMsg string `json:"error_msg"`
}

type ResponseGet struct {
	ErrorMsg   string `json:"error_msg"`
	ReportInfo string `json:"report_info"`
}

type ResponseGetTime struct {
	ErrorMsg             string `json:"error_msg"`
	MaxObservationPeriod string `json:"Max_observation_period"`
}

type ReportSet struct {
	ReportInfo string `json:"report_info"`
}

func HTTPServer(ctx context.Context, storeData *storage.DB) error {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	r := mux.NewRouter()
	r.Use(loggingMiddleware(logger))

	r.HandleFunc("/api/v1/get_report/{report_id}", func(w http.ResponseWriter, r *http.Request) {
		getReportHandler(w, r, storeData)
	}).Methods("GET")

	r.HandleFunc("/api/v1/set_report", func(w http.ResponseWriter, r *http.Request) {
		setReportHandler(w, r, storeData)
	}).Methods("POST")

	r.HandleFunc("/api/v1/get_observation_time/{model_id}", func(w http.ResponseWriter, r *http.Request) {
		getObservationTimeHandler(w, r, storeData)
	}).Methods("GET")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler: r,
	}
	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(context.Background()); err != nil {
			logger.Error("Failed to shutdown server", zap.Error(err))
		}
	}()
	logger.Info("HTTP server started", zap.String("port", os.Getenv("PORT")))
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logger.Error("Server stopped with error", zap.Error(err))
		return err
	}
	return nil

}

func getReportHandler(w http.ResponseWriter, r *http.Request, storeData *storage.DB) {
	idStr := mux.Vars(r)["report_id"]
	id, err := strconv.Atoi(idStr)
	response := ResponseGet{}
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		response.ErrorMsg = err.Error()
		w.WriteHeader(http.StatusBadRequest)
	} else {
		res, err := storeData.GetReport(r.Context(), id)
		if err != nil {
			response.ErrorMsg = err.Error()
			w.WriteHeader(http.StatusNotFound)
		} else {
			response.ReportInfo = res
			w.WriteHeader(http.StatusOK)
		}
	}
	jsonResp, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write(jsonResp)
	}

}

func setReportHandler(w http.ResponseWriter, r *http.Request, storeData *storage.DB) {
	decoder := json.NewDecoder(r.Body)
	var report ReportSet
	response := ResponseSet{}
	w.Header().Set("Content-Type", "application/json")
	err := decoder.Decode(&report)
	if err != nil {
		response.ErrorMsg = err.Error()
		w.WriteHeader(http.StatusBadRequest)
	} else {
		if len(report.ReportInfo) == 0 {
			response.ErrorMsg = "Report_info is empty"
			w.WriteHeader(http.StatusBadRequest)
		} else {
			err = storeData.PostReport(r.Context(), report.ReportInfo)
			if err != nil {
				response.ErrorMsg = err.Error()
				w.WriteHeader(http.StatusNotImplemented)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		}
	}
	jsonResp, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write(jsonResp)
	}

}

func getObservationTimeHandler(w http.ResponseWriter, r *http.Request, storeData *storage.DB) {
	idStr := mux.Vars(r)["model_id"]
	id, err := strconv.Atoi(idStr)
	response := ResponseGetTime{}
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		response.ErrorMsg = err.Error()
		w.WriteHeader(http.StatusBadRequest)
	} else {
		res, err := storeData.GetObservationTime(r.Context(), id)
		if err != nil {
			response.ErrorMsg = err.Error()
			w.WriteHeader(http.StatusNotFound)
		} else {
			response.MaxObservationPeriod = strconv.Itoa(res)
			w.WriteHeader(http.StatusOK)
		}
	}
	jsonResp, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write(jsonResp)
	}

}

func loggingMiddleware(logger *zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("Request received",
				zap.String("request_time", time.Now().Format("2006-01-02 15:04:05")),
				zap.String("remote_address", r.RemoteAddr),
				zap.String("endpoint", r.URL.Path),
			)
			next.ServeHTTP(w, r)
		})
	}
}
