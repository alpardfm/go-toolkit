package codes

import (
	"math"

	"github.com/alpardfm/go-toolkit/language"
	"github.com/alpardfm/go-toolkit/operator"
)

type Code uint32

type AppMessage map[Code]Message

type DisplayMessage struct {
	StatusCode int    `json:"statusCode"`
	Title      string `json:"title"`
	Body       string `json:"body"`
}

const (
	NoCode Code = math.MaxUint32
)

// Success code
const (
	CodeSuccess = Code(iota + 10)

	// reserving codes up to 999 for success codes
)

// Common errors
const (
	CodeInvalidValue = Code(iota + 1000)
	CodeContextDeadlineExceeded
	CodeContextCanceled
	CodeInternalServerError
	CodeServerUnavailable
	CodeNotImplemented
	CodeBadRequest
	CodeNotFound
	CodeConflict
	CodeUnauthorized
	CodeTooManyRequest
	CodeMarshal
	CodeUnmarshal
)

// SQL errors
const (
	CodeSQL = Code(iota + 1300)
	CodeSQLInit
	CodeSQLBuilder
	CodeSQLTxBegin
	CodeSQLTxCommit
	CodeSQLTxRollback
	CodeSQLTxExec
	CodeSQLPrepareStmt
	CodeSQLRead
	CodeSQLRowScan
	CodeSQLRecordDoesNotExist
	CodeSQLUniqueConstraint
	CodeSQLConflict
	CodeSQLNoRowsAffected
)

// third party/client errors
const (
	CodeClient = Code(iota + 1500)
	CodeClientMarshal
	CodeClientUnmarshal
	CodeClientErrorOnRequest
	CodeClientErrorOnReadBody
)

// auth errors
const (
	CodeAuth = Code(iota + 1700)
	CodeAuthRefreshTokenExpired
	CodeAuthAccessTokenExpired
	CodeAuthFailure
	CodeAuthInvalidToken
	CodeForbidden
)

// JSON encoding errors
const (
	CodeJSONSchema = Code(iota + 1900)
	CodeJSONSchemaInvalid
	CodeJSONSchemaNotFound
	CodeJSONStructInvalid
	CodeJSONRawInvalid
	CodeJSONValidationError
	CodeJSONMarshalError
	CodeJSONUnmarshalError
)

// storage errors
const (
	CodeStorage = Code(iota + 2000)
	CodeStorageNoFile
	CodeStorageGenerateURLFailure
	CodeStorageReadFileFailure
	CodeStorageNoClient
)

// nosql errors
const (
	CodeNoSQL = Code(iota + 2200)
	CodeNoSQLInit
	CodeNoSQLRead
	CodeNoSQLClose
	CodeNoSQLDecode
	CodeNoSQLUpdate
	CodeNoSQLInsert
)

// jwt token errors
const (
	CodeJWTInvalidMethod = Code(iota + 2400)
	CodeJWTParseWithClaimsError
	CodeJWTInvalidClaimsType
	CodeJWTSignedStringError
)

// gql errors
const (
	CodeGQLInvalidValue = Code(iota + 2800)
	CodeGQLBuilder
)

// argon2 hash errors
const (
	CodeArgon2InvalidEncodedHash = Code(iota + 3000)
	CodeArgon2EncodeHashError
	CodeArgon2DecodeHashError
	CodeArgon2IncompatibleVersion
)

// aes 256 gcm errors
const (
	CodeAES256GCMOpenError = Code(iota + 3200)
)

// smtp errors
const (
	CodeSMTPError = Code(iota + 3400)
	CodeSMTPBadRequest
	CodeSMTPRequestTimeout
)

// hash errors
const (
	CodeBcryptEncodeHashError = Code(iota + 3600)
	CodeBcryptCompareHashError
)

// s3 error
const (
	CodeS3SessionError = Code(iota + 4000)
)

const (
	// Go Queue codes
	CodeQueueEmpty = Code(iota + 5000)
	CodeQueueFull
	// Other codes
)

const (
	// Go string template codes
	CodeStrTemplateStart = Code(iota + 6000)
	CodeStrTemplateInvalidFormat
	CodeStrTemplateExecuteErr
	CodeStrTemplateEnd
	// Other codes
)

// Error messages only
var ErrorMessages = AppMessage{
	CodeInvalidValue:            ErrMsgBadRequest,
	CodeContextDeadlineExceeded: ErrMsgContextTimeout,
	CodeContextCanceled:         ErrMsgContextTimeout,
	CodeInternalServerError:     ErrMsgInternalServerError,
	CodeServerUnavailable:       ErrMsgServiceUnavailable,
	CodeNotImplemented:          ErrMsgNotImplemented,
	CodeBadRequest:              ErrMsgBadRequest,
	CodeNotFound:                ErrMsgNotFound,
	CodeConflict:                ErrMsgConflict,
	CodeUnauthorized:            ErrMsgUnauthorized,
	CodeTooManyRequest:          ErrMsgTooManyRequest,
	CodeMarshal:                 ErrMsgBadRequest,
	CodeUnmarshal:               ErrMsgBadRequest,
	CodeJSONMarshalError:        ErrMsgBadRequest,
	CodeJSONUnmarshalError:      ErrMsgBadRequest,

	CodeSQL:                   ErrMsgInternalServerError,
	CodeSQLInit:               ErrMsgInternalServerError,
	CodeSQLBuilder:            ErrMsgInternalServerError,
	CodeSQLTxBegin:            ErrMsgInternalServerError,
	CodeSQLTxCommit:           ErrMsgInternalServerError,
	CodeSQLTxRollback:         ErrMsgInternalServerError,
	CodeSQLTxExec:             ErrMsgInternalServerError,
	CodeSQLPrepareStmt:        ErrMsgInternalServerError,
	CodeSQLRead:               ErrMsgInternalServerError,
	CodeSQLRowScan:            ErrMsgInternalServerError,
	CodeSQLRecordDoesNotExist: ErrMsgNotFound,
	CodeSQLUniqueConstraint:   ErrMsgConflict,
	CodeSQLConflict:           ErrMsgConflict,
	CodeSQLNoRowsAffected:     ErrMsgNotFound,

	CodeClientMarshal:         ErrMsgInternalServerError,
	CodeClientUnmarshal:       ErrMsgInternalServerError,
	CodeClientErrorOnRequest:  ErrMsgInternalServerError,
	CodeClientErrorOnReadBody: ErrMsgInternalServerError,

	CodeAuth:                    ErrMsgUnauthorized,
	CodeAuthRefreshTokenExpired: ErrMsgRefreshTokenExpired,
	CodeAuthAccessTokenExpired:  ErrMsgAccessTokenExpired,
	CodeAuthFailure:             ErrMsgUnauthorized,
	CodeAuthInvalidToken:        ErrMsgInvalidToken,
	CodeForbidden:               ErrMsgForbidden,
}

// Successful messages only
var ApplicationMessages = AppMessage{
	CodeSuccess: MsgSuccessDefault,
}

func Compile(code Code, lang string) DisplayMessage {
	if appMsg, ok := ApplicationMessages[code]; ok {
		return DisplayMessage{
			StatusCode: appMsg.StatusCode,
			Title:      operator.TernaryString(lang == language.Indonesian, appMsg.TitleID, appMsg.TitleEN),
			Body:       operator.TernaryString(lang == language.Indonesian, appMsg.BodyID, appMsg.BodyEN),
		}
	}

	return DisplayMessage{
		StatusCode: MsgSuccessDefault.StatusCode,
		Title:      operator.TernaryString(lang == language.Indonesian, MsgSuccessDefault.TitleID, MsgSuccessDefault.TitleEN),
		Body:       operator.TernaryString(lang == language.Indonesian, MsgSuccessDefault.BodyID, MsgSuccessDefault.BodyEN),
	}
}
