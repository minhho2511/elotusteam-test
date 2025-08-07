package transports

import (
	"context"
	"errors"
	"fmt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/minhho2511/elotusteam-test/cfg"
	"github.com/minhho2511/elotusteam-test/internal/kit/endpoints"
	"github.com/minhho2511/elotusteam-test/internal/kit/services"
	"github.com/minhho2511/elotusteam-test/internal/middleware"
	"github.com/minhho2511/elotusteam-test/internal/transforms"
	"github.com/minhho2511/elotusteam-test/pkgs/clog"
	"github.com/minhho2511/elotusteam-test/utils"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var TempDir = "/tmp"

func generateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("upload_%d%s", timestamp, ext)
}

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies/load balancers)
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		// Take the first IP if multiple are present
		ips := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return xRealIP
	}

	// Fall back to RemoteAddr
	return strings.Split(r.RemoteAddr, ":")[0]
}

func decodeUpload(c cfg.Config, logger clog.Logger) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (rqst interface{}, err error) {
		size, err := c.GetMaxFileSize()
		if err != nil {
			return nil, err
		}
		err = r.ParseMultipartForm(size)
		if err != nil {
			return nil, err
		}
		file, fileHeader, err := r.FormFile("data")
		if err != nil {
			return nil, err
		}
		defer file.Close()

		if fileHeader.Size > size {
			return nil, errors.New("FILE_TOO_LARGE")
		}

		contentType := fileHeader.Header.Get("Content-Type")
		if contentType == "" {
			// Detect content type from file extension
			contentType = mime.TypeByExtension(filepath.Ext(fileHeader.Filename))
		}

		fileName := generateFileName(fileHeader.Filename)
		filePath := filepath.Join(TempDir, fileName)

		tempFile, err := os.Create(filePath)
		if err != nil {
			logger.Error(err)
			return nil, errors.New("CREATE_FILE_ERROR")
		}
		defer tempFile.Close()

		_, err = io.Copy(tempFile, file)
		if err != nil {
			log.Printf("Failed to copy file: %v", err)
			os.Remove(filePath) // Clean up
			logger.Error(err)
			return nil, errors.New("COPY_FILE_ERROR")
		}

		userAgent := r.Header.Get("User-Agent")
		ipAddress := getClientIP(r)
		referer := r.Header.Get("Referer")

		req := transforms.FileReq{
			FileName:     fileName,
			OriginalName: fileHeader.Filename,
			FilePath:     filePath,
			ContentType:  contentType,
			FileSize:     fileHeader.Size,
			UserAgent:    userAgent,
			IPAddress:    ipAddress,
			Referer:      referer,
		}
		return req, nil
	}
}

func FileHttpHandler(fileSvc services.FileSvc, logger clog.Logger, c cfg.Config, jwtSvc *utils.JWTService) http.Handler {
	pr := mux.NewRouter()
	pr.Use(middleware.Authenticate(jwtSvc))

	file := endpoints.NewFileEndpoint(fileSvc)

	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(logger),
		httptransport.ServerErrorEncoder(utils.EncodeError),
	}

	pr.Methods("POST").Path("/file/upload").Handler(httptransport.NewServer(
		file.Upload(),
		decodeUpload(c, logger),
		utils.EncodeResponseHTTP,
		options...,
	))

	return pr
}
