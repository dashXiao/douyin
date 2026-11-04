// Shared errors used by the business services.

package errno

var (
	// Success
	Success = NewErrNo(SuccessCode, "Success")

	ServiceError             = NewErrNo(ServiceErrorCode, "service could not start")
	ServiceInternalError     = NewErrNo(ServiceErrorCode, "internal service error")
	ParamError               = NewErrNo(ParamErrorCode, "invalid parameters")
	AuthorizationFailedError = NewErrNo(AuthorizationFailedErrCode, "authorization failed")

	// User
	UserExistedError = NewErrNo(ParamErrorCode, "user already exists")

	// Comment
	UnexpectedTypeError = NewErrNo(UnexpectedTypeErrorCode, "unexpected action type")
	NotImplementError   = NewErrNo(NotImplementErrorCode, "not implemented")

	// Favorite
	LikeNoExistError      = NewErrNo(LikeNoExistErrorCode, "you have not liked this video")
	LikeAlreadyExistError = NewErrNo(LikeAlreadyExistErrorCode, "you have already liked this video")

	// Video
	FileUploadError = NewErrNo(FileUploadErrorCode, "file upload failed")
	NotVideoFile    = NewErrNo(FileUploadError.ErrorCode, "uploaded file is not a video")

	// Follow
	FollowYourselfError = NewErrNo(FollowYourselfErrorCode, "you cannot follow yourself")
	NotFollowError      = NewErrNo(NotFollowErrorCode, "you are not following this user")
	AlreadyFollowError  = NewErrNo(AlreadyFollowErrorCode, "you are already following this user")
	UserNotFoundError   = NewErrNo(ParamErrorCode, "user not found")

	// Chat
	CharacterLimitError = NewErrNo(CharacterLimitErrorCode, "invalid message length")
)
