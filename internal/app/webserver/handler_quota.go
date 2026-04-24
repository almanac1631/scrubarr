package webserver

import "net/http"

func (handler *handler) handleDiskQuotaEndpoint(writer http.ResponseWriter, request *http.Request) {
	logger := getRequestLogger(request)
	diskQuota, err := handler.quotaService.GetDiskQuota()
	if err != nil {
		logger.Error("could not get disk quota", "err", err)
		http.Error(writer, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := handler.ExecuteSubTemplate(writer, "disk_quota.gohtml", "disk_quota", diskQuota); err != nil {
		logger.Error(err.Error())
	}
}
