package helpers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"mime"
	"os"
	"strconv"
	"strings"
)

func ParseFloat(val string) float64 {
	f, _ := strconv.ParseFloat(val, 64)
	return f
}

func ParseInt(val string) int {
	i, _ := strconv.Atoi(val)
	return i
}

type TblFile struct {
	Name string `gorm:"size:255" json:"name"`
	Ext  string `gorm:"size:255" json:"ext"`
	Size int    `gorm:"size:255" json:"size"`
}

func GetDataFile64(base64Str string) (TblFile, error) {

	splitData := strings.SplitN(base64Str, ",", 2)
	if len(splitData) != 2 {
		return TblFile{}, errors.New("invalid base64 format")
	}

	metaData := splitData[0]
	data := splitData[1]

	// Deteksi MIME type
	mimeType := strings.TrimSuffix(strings.TrimPrefix(metaData, "data:"), ";base64")
	exts, err := getStandardImageExt(mimeType)
	if err != nil || len(exts) == 0 {
		return TblFile{}, errors.New("cannot determine file extension")
	}

	// Decode base64 ke []byte
	fileBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return TblFile{}, err
	}

	// Hitung ukuran file
	sizeInBytes := len(fileBytes)

	return TblFile{
		Size: sizeInBytes,
		Ext:  exts,
	}, nil
}

func Base64ToFile(base64str, filename string) (*os.File, error) {
	// Split base64 prefix
	parts := strings.Split(base64str, ",")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid base64 format")
	}

	// Decode
	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	// Buat file sementara
	tmpFile, err := os.CreateTemp("", filename)
	if err != nil {
		return nil, err
	}

	// Tulis data ke file
	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		return nil, err
	}

	// Reset cursor ke awal
	if _, err := tmpFile.Seek(0, io.SeekStart); err != nil {
		tmpFile.Close()
		return nil, err
	}

	return tmpFile, nil
}
func getStandardImageExt(mimeType string) (string, error) {
	switch mimeType {
	case "image/jpeg":
		return ".jpeg", nil
	case "image/png":
		return ".png", nil
	case "image/gif":
		return ".gif", nil
	case "image/webp":
		return ".webp", nil
	case "video/mp4":
		return ".mp4", nil
	default:
		exts, err := mime.ExtensionsByType(mimeType)
		if err != nil || len(exts) == 0 {
			return "", fmt.Errorf("extension invalid")
		}
		return exts[0], nil
	}
}
