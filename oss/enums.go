package oss

// The access control list (ACL) of the bucket
type BucketACLType string

// Enum values for BucketACLType
const (
	// Only the bucket owner can perform read and write operations on objects in the bucket.
	// Other users cannot access the objects in the bucket.
	BucketACLPrivate BucketACLType = "private"

	// Only the bucket owner can write data to objects in the bucket.
	// Other users, including anonymous users, can only read objects in the bucket.
	BucketACLPublicRead BucketACLType = "public-read"

	// All users, including anonymous users, can perform read and write operations on the bucket.
	BucketACLPublicReadWrite BucketACLType = "public-read-write"
)

// The storage class of the bucket
type StorageClassType string

// Enum values for StorageClassType
const (
	// Standard provides highly reliable, highly available,
	// and high-performance object storage for data that is frequently accessed.
	StorageClassStandard StorageClassType = "Standard"

	// IA provides highly durable storage at lower prices compared with Standard.
	// IA has a minimum billable size of 64 KB and a minimum billable storage duration of 30 days.
	StorageClassIA StorageClassType = "IA"

	// Archive provides high-durability storage at lower prices compared with Standard and IA.
	// Archive has a minimum billable size of 64 KB and a minimum billable storage duration of 60 days.
	StorageClassArchive StorageClassType = "Archive"

	// Cold Archive provides highly durable storage at lower prices compared with Archive.
	// Cold Archive has a minimum billable size of 64 KB and a minimum billable storage duration of 180 days.
	StorageClassColdArchive StorageClassType = "ColdArchive"

	// Deep Cold Archive provides highly durable storage at lower prices compared with Cold Archive.
	// Deep Cold Archive has a minimum billable size of 64 KB and a minimum billable storage duration of 180 days.
	StorageClassDeepColdArchive StorageClassType = "DeepColdArchive"
)

// The redundancy type of the bucket
type DataRedundancyType string

// Enum values for BucketACLType
const (
	// Locally redundant storage (LRS) stores copies of each object across different devices in the same zone.
	// This ensures data reliability and availability even if two storage devices are damaged at the same time.
	DataRedundancyLRS DataRedundancyType = "LRS"

	// Zone-redundant storage (ZRS) uses the multi-zone mechanism to distribute user data across
	// multiple zones in the same region. If one zone becomes unavailable, you can continue to access the data
	// that is stored in other zones.
	DataRedundancyZRS DataRedundancyType = "ZRS"
)

// The access control list (ACL) of the object
type ObjectACLType string

// Enum values for ObjectACLType
const (
	// Only the object owner is allowed to perform read and write operations on the object.
	// Other users cannot access the object.
	ObjectACLPrivate ObjectACLType = "private"

	// Only the object owner can write data to the object.
	// Other users, including anonymous users, can only read the object.
	ObjectACLPublicRead ObjectACLType = "public-read"

	// All users, including anonymous users, can perform read and write operations on the object.
	ObjectACLPublicReadWrite ObjectACLType = "public-read-write"

	// The ACL of the object is the same as that of the bucket in which the object is stored.
	ObjectACLDefault ObjectACLType = "default"
)

// OSS headers
const (
	HeaderOssPrefix                      string = "X-Oss-"
	HeaderOssMetaPrefix                         = "X-Oss-Meta-"
	HeaderOssACL                                = "X-Oss-Acl"
	HeaderOssObjectACL                          = "X-Oss-Object-Acl"
	HeaderOssSecurityToken                      = "X-Oss-Security-Token"
	HeaderOssServerSideEncryption               = "X-Oss-Server-Side-Encryption"
	HeaderOssServerSideEncryptionKeyID          = "X-Oss-Server-Side-Encryption-Key-Id"
	HeaderOssServerSideDataEncryption           = "X-Oss-Server-Side-Data-Encryption"
	HeaderOssSSECAlgorithm                      = "X-Oss-Server-Side-Encryption-Customer-Algorithm"
	HeaderOssSSECKey                            = "X-Oss-Server-Side-Encryption-Customer-Key"
	HeaderOssSSECKeyMd5                         = "X-Oss-Server-Side-Encryption-Customer-Key-MD5"
	HeaderOssCopySource                         = "X-Oss-Copy-Source"
	HeaderOssCopySourceRange                    = "X-Oss-Copy-Source-Range"
	HeaderOssCopySourceIfMatch                  = "X-Oss-Copy-Source-If-Match"
	HeaderOssCopySourceIfNoneMatch              = "X-Oss-Copy-Source-If-None-Match"
	HeaderOssCopySourceIfModifiedSince          = "X-Oss-Copy-Source-If-Modified-Since"
	HeaderOssCopySourceIfUnmodifiedSince        = "X-Oss-Copy-Source-If-Unmodified-Since"
	HeaderOssMetadataDirective                  = "X-Oss-Metadata-Directive"
	HeaderOssNextAppendPosition                 = "X-Oss-Next-Append-Position"
	HeaderOssRequestID                          = "X-Oss-Request-Id"
	HeaderOssCRC64                              = "X-Oss-Hash-Crc64ecma"
	HeaderOssSymlinkTarget                      = "X-Oss-Symlink-Target"
	HeaderOssStorageClass                       = "X-Oss-Storage-Class"
	HeaderOssCallback                           = "X-Oss-Callback"
	HeaderOssCallbackVar                        = "X-Oss-Callback-Var"
	HeaderOssRequester                          = "X-Oss-Request-Payer"
	HeaderOssTagging                            = "X-Oss-Tagging"
	HeaderOssTaggingDirective                   = "X-Oss-Tagging-Directive"
	HeaderOssTrafficLimit                       = "X-Oss-Traffic-Limit"
	HeaderOssForbidOverWrite                    = "X-Oss-Forbid-Overwrite"
	HeaderOssRangeBehavior                      = "X-Oss-Range-Behavior"
	HeaderOssAllowSameActionOverLap             = "X-Oss-Allow-Same-Action-Overlap"
	HeaderOssDate                               = "X-Oss-Date"
	HeaderOssContentSha256                      = "X-Oss-Content-Sha256"
	HeaderOssEC                                 = "X-Oss-Ec"
	HeaderOssERR                                = "X-Oss-Err"
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
	HTTPHeaderUserAgent                 = "User-Agent"
	HTTPHeaderIfModifiedSince           = "If-Modified-Since"
	HTTPHeaderIfUnmodifiedSince         = "If-Unmodified-Since"
	HTTPHeaderIfMatch                   = "If-Match"
	HTTPHeaderIfNoneMatch               = "If-None-Match"
)

type UrlStyleType int

const (
	UrlStyleVirtualHosted UrlStyleType = iota
	UrlStylePath
	UrlStyleCName
)

func (f UrlStyleType) String() string {
	switch f {
	default:
		return "virtual-hosted-style"
	case UrlStylePath:
		return "path-style"
	case UrlStyleCName:
		return "cname-style"
	}
}

type FeatureFlagsType int

const (
	// If the client time is different from server time by more than about 15 minutes,
	// the requests your application makes will be signed with the incorrect time, and the server will reject them.
	// The feature to help to identify this case, and SDK will correct for clock skew.
	FeatureCorrectClockSkew FeatureFlagsType = 1 << iota

	FeatureFlagsDefault = FeatureCorrectClockSkew
)

type SignatureVersionType int

const (
	SignatureVersionV1 SignatureVersionType = iota
)

func (f SignatureVersionType) String() string {
	switch f {
	default:
		return "OSS Signature Version 1"
	}
}

type AuthMethodType int

const (
	AuthMethodHeader AuthMethodType = iota
	AuthMethodQuery
)

func (f AuthMethodType) String() string {
	switch f {
	case AuthMethodQuery:
		return "authentication in query"
	default:
		return "authentication in header"
	}
}
