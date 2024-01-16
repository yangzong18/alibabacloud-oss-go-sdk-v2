package oss

func (c *Client) NewDownloader(optFns ...func(*DownloaderOptions)) *Downloader {
	return NewDownloader(c, optFns...)
}

func (c *Client) NewUploader(optFns ...func(*UploaderOptions)) *Uploader {
	return NewUploader(c, optFns...)

}
