package types

// ACLType bucket/object ACL
type ACLType string

const (
	// ACLPrivate definition : private read and write
	ACLPrivate ACLType = "private"

	// ACLPublicRead definition : public read and private write
	ACLPublicRead ACLType = "public-read"

	// ACLPublicReadWrite definition : public read and public write
	ACLPublicReadWrite ACLType = "public-read-write"

	// ACLDefault Object. It's only applicable for object.
	ACLDefault ACLType = "default"
)

// bucket versioning status
type VersioningStatus string

const (
	// Versioning Status definition: Enabled
	VersionEnabled VersioningStatus = "Enabled"

	// Versioning Status definition: Suspended
	VersionSuspended VersioningStatus = "Suspended"
)

// MetadataDirectiveType specifying whether use the metadata of source object when copying object.
type MetadataDirectiveType string

const (
	// MetaCopy the target object's metadata is copied from the source one
	MetaCopy MetadataDirectiveType = "COPY"

	// MetaReplace the target object's metadata is created as part of the copy request (not same as the source one)
	MetaReplace MetadataDirectiveType = "REPLACE"
)

// TaggingDirectiveType specifying whether use the tagging of source object when copying object.
type TaggingDirectiveType string

const (
	// TaggingCopy the target object's tagging is copied from the source one
	TaggingCopy TaggingDirectiveType = "COPY"

	// TaggingReplace the target object's tagging is created as part of the copy request (not same as the source one)
	TaggingReplace TaggingDirectiveType = "REPLACE"
)

// AlgorithmType specifying the server side encryption algorithm name
type AlgorithmType string

const (
	KMSAlgorithm AlgorithmType = "KMS"
	AESAlgorithm AlgorithmType = "AES256"
	SM4Algorithm AlgorithmType = "SM4"
)

// StorageClassType bucket storage type
type StorageClassType string

const (
	// StorageStandard standard
	StorageStandard StorageClassType = "Standard"

	// StorageIA infrequent access
	StorageIA StorageClassType = "IA"

	// StorageArchive archive
	StorageArchive StorageClassType = "Archive"

	// StorageColdArchive cold archive
	StorageColdArchive StorageClassType = "ColdArchive"

	// StorageDeepColdArchive deep cold archive
	StorageDeepColdArchive StorageClassType = "DeepColdArchive"
)

// RedundancyType bucket data Redundancy type
type DataRedundancyType string

const (
	// RedundancyLRS Local redundancy, default value
	RedundancyLRS DataRedundancyType = "LRS"

	// RedundancyZRS Same city redundancy
	RedundancyZRS DataRedundancyType = "ZRS"
)

// ObjecthashFuncType
type ObjecthashFuncType string

const (
	HashFuncSha1   ObjecthashFuncType = "SHA-1"
	HashFuncSha256 ObjecthashFuncType = "SHA-256"
)

// PayerType the type of request payer
type PayerType string

const (
	// Requester the requester who send the request
	Requester PayerType = "Requester"

	// BucketOwner the requester who send the request
	BucketOwner PayerType = "BucketOwner"
)

// RestoreMode the restore mode for coldArchive object
type RestoreMode string

const (
	//RestoreExpedited object will be restored in 1 hour
	RestoreExpedited RestoreMode = "Expedited"

	//RestoreStandard object will be restored in 2-5 hours
	RestoreStandard RestoreMode = "Standard"

	//RestoreBulk object will be restored in 5-10 hours
	RestoreBulk RestoreMode = "Bulk"
)

// HTTPMethod HTTP request method
type HTTPMethod string

const (
	// HTTPGet HTTP GET
	HTTPGet HTTPMethod = "GET"

	// HTTPPut HTTP PUT
	HTTPPut HTTPMethod = "PUT"

	// HTTPHead HTTP HEAD
	HTTPHead HTTPMethod = "HEAD"

	// HTTPPost HTTP POST
	HTTPPost HTTPMethod = "POST"

	// HTTPDelete HTTP DELETE
	HTTPDelete HTTPMethod = "DELETE"
)

// HTTP headers
const (
	HTTPHeaderAcceptEncoding     string = "Accept-Encoding"
	HTTPHeaderAuthorization             = "Authorization"
	HTTPHeaderCacheControl              = "Cache-Control"
	HTTPHeaderContentDisposition        = "Content-Disposition"
	HTTPHeaderContentEncoding           = "Content-Encoding"
	HTTPHeaderContentLength             = "Content-Length"
	HTTPHeaderContentMD5                = "Content-MD5"
	HTTPHeaderContentType               = "Content-Type"
	HTTPHeaderContentLanguage           = "Content-Language"
	HTTPHeaderDate                      = "Date"
	HTTPHeaderEtag                      = "ETag"
	HTTPHeaderExpires                   = "Expires"
	HTTPHeaderHost                      = "Host"
	HTTPHeaderLastModified              = "Last-Modified"
	HTTPHeaderRange                     = "Range"
	HTTPHeaderLocation                  = "Location"
	HTTPHeaderOrigin                    = "Origin"
	HTTPHeaderServer                    = "Server"
	HTTPHeaderUserAgent                 = "User-Agent"
	HTTPHeaderIfModifiedSince           = "If-Modified-Since"
	HTTPHeaderIfUnmodifiedSince         = "If-Unmodified-Since"
	HTTPHeaderIfMatch                   = "If-Match"
	HTTPHeaderIfNoneMatch               = "If-None-Match"
	HTTPHeaderACReqMethod               = "Access-Control-Request-Method"
	HTTPHeaderACReqHeaders              = "Access-Control-Request-Headers"

	HTTPHeaderOssACL                         = "X-Oss-Acl"
	HTTPHeaderOssMetaPrefix                  = "X-Oss-Meta-"
	HTTPHeaderOssObjectACL                   = "X-Oss-Object-Acl"
	HTTPHeaderOssSecurityToken               = "X-Oss-Security-Token"
	HTTPHeaderOssServerSideEncryption        = "X-Oss-Server-Side-Encryption"
	HTTPHeaderOssServerSideEncryptionKeyID   = "X-Oss-Server-Side-Encryption-Key-Id"
	HTTPHeaderOssServerSideDataEncryption    = "X-Oss-Server-Side-Data-Encryption"
	HTTPHeaderSSECAlgorithm                  = "X-Oss-Server-Side-Encryption-Customer-Algorithm"
	HTTPHeaderSSECKey                        = "X-Oss-Server-Side-Encryption-Customer-Key"
	HTTPHeaderSSECKeyMd5                     = "X-Oss-Server-Side-Encryption-Customer-Key-MD5"
	HTTPHeaderOssCopySource                  = "X-Oss-Copy-Source"
	HTTPHeaderOssCopySourceRange             = "X-Oss-Copy-Source-Range"
	HTTPHeaderOssCopySourceIfMatch           = "X-Oss-Copy-Source-If-Match"
	HTTPHeaderOssCopySourceIfNoneMatch       = "X-Oss-Copy-Source-If-None-Match"
	HTTPHeaderOssCopySourceIfModifiedSince   = "X-Oss-Copy-Source-If-Modified-Since"
	HTTPHeaderOssCopySourceIfUnmodifiedSince = "X-Oss-Copy-Source-If-Unmodified-Since"
	HTTPHeaderOssMetadataDirective           = "X-Oss-Metadata-Directive"
	HTTPHeaderOssNextAppendPosition          = "X-Oss-Next-Append-Position"
	HTTPHeaderOssRequestID                   = "X-Oss-Request-Id"
	HTTPHeaderOssCRC64                       = "X-Oss-Hash-Crc64ecma"
	HTTPHeaderOssSymlinkTarget               = "X-Oss-Symlink-Target"
	HTTPHeaderOssStorageClass                = "X-Oss-Storage-Class"
	HTTPHeaderOssCallback                    = "X-Oss-Callback"
	HTTPHeaderOssCallbackVar                 = "X-Oss-Callback-Var"
	HTTPHeaderOssRequester                   = "X-Oss-Request-Payer"
	HTTPHeaderOssTagging                     = "X-Oss-Tagging"
	HTTPHeaderOssTaggingDirective            = "X-Oss-Tagging-Directive"
	HTTPHeaderOssTrafficLimit                = "X-Oss-Traffic-Limit"
	HTTPHeaderOssForbidOverWrite             = "X-Oss-Forbid-Overwrite"
	HTTPHeaderOssRangeBehavior               = "X-Oss-Range-Behavior"
	HTTPHeaderOssTaskID                      = "X-Oss-Task-Id"
	HTTPHeaderOssHashCtx                     = "X-Oss-Hash-Ctx"
	HTTPHeaderOssMd5Ctx                      = "X-Oss-Md5-Ctx"
	HTTPHeaderAllowSameActionOverLap         = "X-Oss-Allow-Same-Action-Overlap"
	HttpHeaderOssDate                        = "X-Oss-Date"
	HttpHeaderOssContentSha256               = "X-Oss-Content-Sha256"
	HttpHeaderOssNotification                = "X-Oss-Notification"
	HTTPHeaderOssEc                          = "X-Oss-Ec"
	HTTPHeaderOssErr                         = "X-Oss-Err"
)

// AuthVersion the version of auth
type AuthVersionType string

const (
	// AuthV1 v1
	AuthV1 AuthVersionType = "v1"
	// AuthV2 v2
	AuthV2 AuthVersionType = "v2"
	// AuthV4 v4
	AuthV4 AuthVersionType = "v4"
)

type UrlRequestStyle int

const (
	VirtualHostedStyle UrlRequestStyle = iota
	PathStyle
	CNameStyle
)

func (f UrlRequestStyle) String() string {
	switch f {
	case VirtualHostedStyle:
		return "virtual-hosted–style"
	case PathStyle:
		return "path-style"
	case CNameStyle:
		return "cname-style"
	default:
		return ""
	}
}
