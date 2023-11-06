package oss

import (
	"net/url"
)

func isValidEndpoint(endpoint *url.URL) bool {
	return (endpoint != nil)
}

func isValidBucketName(bucketName *string) bool {
	if bucketName == nil {
		return false
	}

	nameLen := len(*bucketName)
	if nameLen < 3 || nameLen > 63 {
		return false
	}

	if (*bucketName)[0] == '-' || (*bucketName)[nameLen-1] == '-' {
		return false
	}

	for _, v := range *bucketName {
		if !(('a' <= v && v <= 'z') || ('0' <= v && v <= '9') || v == '-') {
			return false
		}
	}
	return true
}

func isValidObjectName(objectName *string) bool {
	if objectName == nil || len(*objectName) == 0 {
		return false
	}
	return true
}

var supportedMethod = map[string]interface{}{
	"GET":    nil,
	"PUT":    nil,
	"POST":   nil,
	"DELETE": nil,
	"OPTION": nil,
}

func isValidMethod(method string) bool {
	if _, ok := supportedMethod[method]; ok {
		return true
	}
	return false
}
