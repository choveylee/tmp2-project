package lib

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/choveylee/terror"
	"github.com/choveylee/thttp"
	"github.com/choveylee/tlog"

	constant "dev.choveylee.top/knowledge-base-backend/internal/const"
)

func ParseTikaFile(ctx context.Context, fileName string, mimeType string, fileBytes []byte) (string, *terror.Terror) {
	tikaURL := tikaEndpoint + "/tika"

	requestOption := thttp.NewRequestOption().WithHeader("Accept", "text/plain")
	if strings.TrimSpace(mimeType) != "" {
		requestOption.WithContentType(mimeType)
	}

	response, err := tikaClient.Put(ctx, tikaURL, requestOption, fileBytes)
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msgf("Parse tika file (tika endpoint: %s, tika url: %s, file name: %s, mime type: %s, file size: %d) err (put request %v)",
			tikaEndpoint, tikaURL, fileName, mimeType, len(fileBytes), err)

		errx := terror.NewTerror(ctx, err, constant.ErrorCodeDocumentParseFailed, errMsg)

		return "", errx
	}

	statusCode, body, err := response.ToString()
	if err != nil {
		errMsg := tlog.E(ctx).Err(err).Msgf("Parse tika file (tika endpoint: %s, tika url: %s, file name: %s, mime type: %s, file size: %d, status code: %d) err (read response body %v)",
			tikaEndpoint, tikaURL, fileName, mimeType, len(fileBytes), statusCode, err)

		errx := terror.NewTerror(ctx, err, constant.ErrorCodeDocumentParseFailed, errMsg)

		return "", errx
	}

	if statusCode >= http.StatusBadRequest {
		errMsg := tlog.E(ctx).Msgf("Parse tika file (tika endpoint: %s, tika url: %s, file name: %s, mime type: %s, file size: %d, status code: %d, response len: %d) err (bad response status)",
			tikaEndpoint, tikaURL, fileName, mimeType, len(fileBytes), statusCode, len(body))

		errx := terror.NewTerror(ctx, errors.New("tika bad response status"), constant.ErrorCodeDocumentParseFailed, errMsg)

		return "", errx
	}

	return body, nil
}
