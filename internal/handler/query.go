package handler

import (
	"github.com/bestkkii/saedori-api-server/internal/model"
	"github.com/bestkkii/saedori-api-server/pkg"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

/**
* 쿼리 스트링 있는 API 리스트
**/

func (h *Handler) GetInterestDetail(c *gin.Context) {
	category := c.DefaultQuery("category", "default_category")
	parsedCategories, length := pkg.ParseCategory(category)

	switch length {
	case 1:
		h.handleSingleCategory(c, parsedCategories[0])
	case 2:
		h.handleDoubleCategories(c, parsedCategories)
	default:
		h.handleAllCategories(c, parsedCategories)
	}
}

func (h *Handler) handleSingleCategory(c *gin.Context, category string) {
	switch category {
	case "music":
		h.GetMusicList(c)
	case "realtime-search":
		h.GetRealtimeSearchDetail(c)
	case "news":
		h.GetNewsDetails(c)
	}
}

func (h *Handler) handleDoubleCategories(c *gin.Context, categories []string) {
	switch {
	case categories[0] == "music" && categories[1] == "news":
		h.GetMusicAndNews(c)
	case categories[0] == "music" && categories[1] == "realtime-search":
		h.GetMusicAndRealtimeSearch(c)
	case categories[0] == "news" && categories[1] == "realtime-search":
		h.GetNewsAndRealtimeSearch(c)
	}
}

func (h *Handler) handleAllCategories(c *gin.Context, categories []string) {
	if categories[0] == "music" && categories[1] == "news" && categories[2] == "realtime-search" {
		h.GetAllCategories(c)
	}
}

func (h *Handler) GetMusicAndNews(c *gin.Context) {
	musics, err := h.dashboardService.GetMusicList()
	if err != nil {
		h.failedResponse(c, pkg.NewApiResponse("FAILED"))
		return
	}

	news, err := h.dashboardService.GetNewsDetails()
	if err != nil {
		h.failedResponse(c, pkg.NewApiResponse("FAILED"))
		return
	}

	h.okResponse(c, model.AllCategoriesResponse{
		ApiResponse: pkg.NewApiResponse("SUCCESS"),
		Musics:      musics,
		News:        []*model.News{news},
	})
}

func (h *Handler) GetMusicAndRealtimeSearch(c *gin.Context) {
	musics, err := h.dashboardService.GetMusicList()
	if err != nil {
		h.failedResponse(c, pkg.NewApiResponse("FAILED"))
		return
	}

	realtimeSearchDetail, err := h.dashboardService.GetRealtimeSearchDetailList()
	if err != nil {
		h.failedResponse(c, pkg.NewApiResponse("FAILED"))
		return
	}

	h.okResponse(c, model.AllCategoriesResponse{
		ApiResponse:                 pkg.NewApiResponse("SUCCESS"),
		Musics:                      musics,
		RealtimeSearchDetailWrapper: realtimeSearchDetail.RealtimeSearchDetailWrapper,
	})
}

func (h *Handler) GetNewsAndRealtimeSearch(c *gin.Context) {
	news, err := h.dashboardService.GetNewsDetails()
	if err != nil {
		h.failedResponse(c, pkg.NewApiResponse("FAILED"))
		return
	}

	realtimeSearchDetail, err := h.dashboardService.GetRealtimeSearchDetailList()
	if err != nil {
		h.failedResponse(c, pkg.NewApiResponse("FAILED"))
		return
	}

	h.okResponse(c, model.AllCategoriesResponse{
		ApiResponse:                 pkg.NewApiResponse("SUCCESS"),
		RealtimeSearchDetailWrapper: realtimeSearchDetail.RealtimeSearchDetailWrapper,
		News:                        []*model.News{news},
	})
}

func (h *Handler) GetAllCategories(c *gin.Context) {
	musics, err := h.dashboardService.GetMusicList()
	if err != nil {
		h.failedResponse(c, pkg.NewApiResponse("FAILED"))
		return
	}

	realtimeSearchDetail, err := h.dashboardService.GetRealtimeSearchDetailList()
	if err != nil {
		h.failedResponse(c, pkg.NewApiResponse("FAILED"))
		return
	}

	news, err := h.dashboardService.GetNewsDetails()
	if err != nil {
		h.failedResponse(c, pkg.NewApiResponse("FAILED"))
		return
	}

	h.okResponse(c, model.AllCategoriesResponse{
		ApiResponse:                 pkg.NewApiResponse("SUCCESS"),
		Musics:                      musics,
		RealtimeSearchDetailWrapper: realtimeSearchDetail.RealtimeSearchDetailWrapper,
		News:                        []*model.News{news},
	})
}

// GetDownloadData handles the download data request
func (h *Handler) GetDownloadData(c *gin.Context) {
	// Parse query parameters
	categoryStr := c.Query("category")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	// Split category string into array
	categories := strings.Split(categoryStr, ",")
	if len(categories) == 0 || categoryStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category parameter is required"})
		return
	}

	// Parse dates
	startDate, err := strconv.ParseInt(startDateStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date parameter"})
		return
	}

	endDate, err := strconv.ParseInt(endDateStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date parameter"})
		return
	}

	// Get data from service
	data, err := h.dashboardService.GetDownloadData(categories, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create response with ApiResponse
	response := &model.DownloadDataResponse{
		ApiResponse: pkg.NewApiResponse("Success"),
		Result:      data,
	}

	c.JSON(http.StatusOK, response)
}
