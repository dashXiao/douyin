// Shared status and error codes for the API gateway and business services.
package errno

const (
	// API gateway success response.
	StatusSuccessCode = 0
	StatusSuccessMsg  = "ok"

	// RPC success response.
	SuccessCode = 10000
	SuccessMsg  = "ok"

	// Error codes.
	ServiceErrorCode           = 10001 // Unknown microservice error
	ParamErrorCode             = 10002 // Invalid parameters
	AuthorizationFailedErrCode = 10003 // Authorization failed
	UnexpectedTypeErrorCode    = 10004 // Unexpected action type
	NotImplementErrorCode      = 10005 // Not implemented
	LikeAlreadyExistErrorCode  = 10007 // Video already liked
	LikeNoExistErrorCode       = 10008 // Video not liked
	FileUploadErrorCode        = 10009 // File upload failed
	FollowYourselfErrorCode    = 10010 // Cannot follow yourself
	NotFollowErrorCode         = 10011 // Not following this user
	AlreadyFollowErrorCode     = 10012 // Already following this user
	CharacterLimitErrorCode    = 10013 // Invalid message length
)
