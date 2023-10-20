package usecases

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func responseWithResult(data interface{}) (int, gin.H) {
	return http.StatusOK, gin.H{
		"result": data,
	}
}

func responseWithError(err error) (int, gin.H) {
	return http.StatusInternalServerError, gin.H{
		"error": err.Error(),
	}
}

func responseWithPagination(totalItems int, limit int, page int, data interface{}) (int, gin.H) {
	return http.StatusOK, gin.H{
		"totalItems":  totalItems,
		"totalPages":  int(int(totalItems)/limit) + 1,
		"currentPage": page,
		"data":        data,
	}
}
