package oss

import "os"

const (
	MaxUploadParts int32 = 10000

	MinUploadPartSize int64 = 5 * 1024 * 1024

	DefaultUploadPartSize = MinUploadPartSize

	DefaultUploadParallel = 3

	DefaultDownloadPartSize = MinUploadPartSize

	DefaultDownloadParallel = 3

	FilePermMode = os.FileMode(0664) // File permission

	TempFileSuffix = ".temp" // Temp file suffix
)
