package utils

import (
	"os"
	"path/filepath"
)

func SaveFileToReportFolder(fileName string, data []byte) error {
	// Pastikan folder report/ ada
	err := os.MkdirAll("report", os.ModePerm)
	if err != nil {
		return err
	}

	// Path lengkap
	fullPath := filepath.Join("report", fileName)

	// Simpan file
	return os.WriteFile(fullPath, data, 0644)
}
