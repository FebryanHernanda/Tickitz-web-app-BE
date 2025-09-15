package utils

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadFile(ctx *gin.Context, formField, uploadPath, prefix, folderPath string) (string, error) {
	file, err := ctx.FormFile(formField)
	if err != nil {
		return "", nil
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_images_%d%s", prefix, time.Now().UnixNano(), ext)
	fullPath := filepath.Join(uploadPath, filename)

	if err := ctx.SaveUploadedFile(file, fullPath); err != nil {
		return "", err
	}

	return fmt.Sprintf("/%s/%s", folderPath, filename), nil
}
