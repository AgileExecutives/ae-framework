package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/AgileExecutives/ae-framework/serverbase/pkg/core"
	"github.com/gin-gonic/gin"
)

type StaticHandlers struct {
	logger core.Logger
	repo   StaticRepo
}

func NewStaticHandlers(logger core.Logger, repo StaticRepo) *StaticHandlers {
	return &StaticHandlers{logger: logger, repo: repo}
}

func (h *StaticHandlers) ListStaticJSON(c *gin.Context) {
	entries, err := h.repo.ListFiles(context.Background())
	if err != nil {
		h.logger.Error("Failed to list static JSON files", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read directory"})
		return
	}
	var files []string
	for _, name := range entries {
		if strings.HasSuffix(strings.ToLower(name), ".json") {
			filename := strings.TrimSuffix(name, ".json")
			files = append(files, filename)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"available_files": files,
		"base_url":        "/api/v1/static/",
		"example_usage":   "GET /api/v1/static/{filename}",
		"security_note":   "Only JSON files from statics/json/ directory are accessible",
		"restrictions":    "Filenames must be alphanumeric with hyphens/underscores only",
		"note":            "Drop any .json file in ./statics/json/ directory to make it available",
	})
}

func (h *StaticHandlers) ServeStaticJSON(c *gin.Context) {
	fileName := c.Param("filename")
	if fileName == "" || len(fileName) > 100 {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	if strings.HasPrefix(fileName, ".") || strings.Contains(fileName, "\x00") {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	for _, char := range fileName {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '-' || char == '_') {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file name"})
			return
		}
	}
	expectedFilename := fileName + ".json"
	entries, lerr := h.repo.ListFiles(context.Background())
	if lerr != nil {
		h.logger.Error("Failed to list static JSON files for lookup", "error", lerr)
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	found := false
	for _, e := range entries {
		if e == expectedFilename {
			found = true
			break
		}
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	data, err := h.repo.ReadFile(context.Background(), expectedFilename)
	if err != nil {
		h.logger.Error("Failed to read JSON file", "file", expectedFilename, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	c.Header("Content-Type", "application/json")
	c.Data(http.StatusOK, "application/json", data)
}
