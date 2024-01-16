package oss

import "context"

// NewDownloader creates a new Downloader instance to download objects.
func (c *Client) NewDownloader(optFns ...func(*DownloaderOptions)) *Downloader {
	return NewDownloader(c, optFns...)
}

// NewUploader creates a new Uploader instance to upload objects.
func (c *Client) NewUploader(optFns ...func(*UploaderOptions)) *Uploader {
	return NewUploader(c, optFns...)

}

// OpenFile opens the named file for reading.
func (c *Client) OpenFile(ctx context.Context, bucket string, key string, optFns ...func(*OpenOptions)) (*ReadOnlyFile, error) {
	return NewReadOnlyFile(ctx, c, bucket, key, optFns...)
}

// AppendFile opens or creates the named file for appending.
func (c *Client) AppendFile(ctx context.Context, bucket string, key string, optFns ...func(*AppendOptions)) (*AppendOnlyFile, error) {
	return NewAppendFile(ctx, c, bucket, key, optFns...)
}
