package ipp

const (
	StatusCupsInvalid                         uint8 = -1
	StatusOk                                  uint8 = 0x0000
	StatusOkIgnoredOrSubstituted              uint8 = 0x0001
	StatusOkConflicting                       uint8 = 0x0002
	statusOkIgnoredSubscriptions              uint8 = 0x0003
	statusOkIgnoredNotifications              uint8 = 0x0004
	statusOkTooManyEvents                     uint8 = 0x0005
	statusOkButCancelSubscription             uint8 = 0x0006
	statusOkEventsComplete                    uint8 = 0x0007
	statusRedirectionOtherSite                uint8 = 0x0200
	statusCupsSeeOther                        uint8 = 0x0280
	statusErrorBadRequest                     uint8 = 0x0400
	StatusErrorForbidden                      uint8 = 0x0401
	StatusErrorNotAuthenticated               uint8 = 0x0402
	StatusErrorNotAuthorized                  uint8 = 0x0403
	StatusErrorNotPossible                    uint8 = 0x0404
	StatusErrorTimeout                        uint8 = 0x0405
	StatusErrorNotFound                       uint8 = 0x0406
	StatusErrorGone                           uint8 = 0x0407
	StatusErrorRequestEntity                  uint8 = 0x0408
	StatusErrorRequestValue                   uint8 = 0x0409
	StatusErrorDocumentFormatNotSupported     uint8 = 0x040a
	StatusErrorAttributesOrValues             uint8 = 0x040b
	StatusErrorUriScheme                      uint8 = 0x040c
	StatusErrorCharset                        uint8 = 0x040d
	StatusErrorConflicting                    uint8 = 0x040e
	StatusErrorCompressionError               uint8 = 0x040f
	StatusErrorDocumentFormatError            uint8 = 0x0410
	StatusErrorDocumentAccess                 uint8 = 0x0411
	StatusErrorAttributesNotSettable          uint8 = 0x0412
	StatusErrorIgnoredAllSubscriptions        uint8 = 0x0413
	StatusErrorTooManySubscriptions           uint8 = 0x0414
	StatusErrorIgnoredAllNotifications        uint8 = 0x0415
	StatusErrorPrintSupportFileNotFound       uint8 = 0x0416
	StatusErrorDocumentPassword               uint8 = 0x0417
	StatusErrorDocumentPermission             uint8 = 0x0418
	StatusErrorDocumentSecurity               uint8 = 0x0419
	StatusErrorDocumentUnprintable            uint8 = 0x041a
	StatusErrorAccountInfoNeeded              uint8 = 0x041b
	StatusErrorAccountClosed                  uint8 = 0x041c
	StatusErrorAccountLimitReached            uint8 = 0x041d
	StatusErrorAccountAuthorizationFailed     uint8 = 0x041e
	StatusErrorNotFetchable                   uint8 = 0x041f
	StatusErrorCupsAccountInfoNeeded          uint8 = 0x049C
	StatusErrorCupsAccountClosed              uint8 = 0x049d
	StatusErrorCupsAccountLimitReached        uint8 = 0x049e
	StatusErrorCupsAccountAuthorizationFailed uint8 = 0x049f
	StatusErrorInternal                       uint8 = 0x0500
	StatusErrorOperationNotSupported          uint8 = 0x0501
	StatusErrorServiceUnavailable             uint8 = 0x0502
	StatusErrorVersionNotSupported            uint8 = 0x0503
	StatusErrorDevice                         uint8 = 0x0504
	StatusErrorTemporary                      uint8 = 0x0505
	StatusErrorNotAcceptingJobs               uint8 = 0x0506
	StatusErrorBusy                           uint8 = 0x0507
	StatusErrorJobCanceled                    uint8 = 0x0508
	StatusErrorMultipleJobsNotSupported       uint8 = 0x0509
	StatusErrorPrinterIsDeactivated           uint8 = 0x050a
	StatusErrorTooManyJobs                    uint8 = 0x050b
	StatusErrorTooManyDocuments               uint8 = 0x050c
	StatusErrorCupsAuthenticationCanceled     uint8 = 0x1000
	StatusErrorCupsPki                        uint8 = 0x1001
	StatusErrorCupsUpgradeRequired            uint8 = 0x1002
)

const (
	OperationCupsInvalid                     int16 = -0x0001
	OperationCupsNone                        int16 = 0x0000
	OperationPrintJob                        int16 = 0x0002
	OperationPrintUri                        int16 = 0x0003
	OperationValidateJob                     int16 = 0x0004
	OperationCreateJob                       int16 = 0x0005
	OperationSendDocument                    int16 = 0x0006
	OperationSendUri                         int16 = 0x0007
	OperationCancelJob                       int16 = 0x0008
	OperationGetJobAttributes                int16 = 0x0009
	OperationGetJobs                         int16 = 0x000a
	OperationGetPrinterAttributes            int16 = 0x000b
	OperationHoldJob                         int16 = 0x000c
	OperationReleaseJob                      int16 = 0x000d
	OperationRestartJob                      int16 = 0x000e
	OperationPausePrinter                    int16 = 0x0010
	OperationResumePrinter                   int16 = 0x0011
	OperationPurgeJobs                       int16 = 0x0012
	OperationSetPrinterAttributes            int16 = 0x0013
	OperationSetJobAttributes                int16 = 0x0014
	OperationGetPrinterSupportedValues       int16 = 0x0015
	OperationCreatePrinterSubscriptions      int16 = 0x0016
	OperationCreateJobSubscriptions          int16 = 0x0017
	OperationGetSubscriptionAttributes       int16 = 0x0018
	OperationGetSubscriptions                int16 = 0x0019
	OperationRenewSubscription               int16 = 0x001a
	OperationCancelSubscription              int16 = 0x001b
	OperationGetNotifications                int16 = 0x001c
	OperationSendNotifications               int16 = 0x001d
	OperationGetResourceAttributes           int16 = 0x001e
	OperationGetResourceData                 int16 = 0x001f
	OperationGetResources                    int16 = 0x0020
	OperationGetPrintSupportFiles            int16 = 0x0021
	OperationEnablePrinter                   int16 = 0x0022
	OperationDisablePrinter                  int16 = 0x0023
	OperationPausePrinterAfterCurrentJob     int16 = 0x0024
	OperationHoldNewJobs                     int16 = 0x0025
	OperationReleaseHeldNewJobs              int16 = 0x0026
	OperationDeactivatePrinter               int16 = 0x0027
	OperationActivatePrinter                 int16 = 0x0028
	OperationRestartPrinter                  int16 = 0x0029
	OperationShutdownPrinter                 int16 = 0x002a
	OperationStartupPrinter                  int16 = 0x002b
	OperationReprocessJob                    int16 = 0x002c
	OperationCancelCurrentJob                int16 = 0x002d
	OperationSuspendCurrentJob               int16 = 0x002e
	OperationResumeJob                       int16 = 0x002f
	OperationOperationPromoteJob             int16 = 0x0030
	OperationScheduleJobAfter                int16 = 0x0031
	OperationCancelDocument                  int16 = 0x0033
	OperationGetDocumentAttributes           int16 = 0x0034
	OperationGetDocuments                    int16 = 0x0035
	OperationDeleteDocument                  int16 = 0x0036
	OperationSetDocumentAttributes           int16 = 0x0037
	OperationCancelJobs                      int16 = 0x0038
	OperationCancelMyJobs                    int16 = 0x0039
	OperationResubmitJob                     int16 = 0x003a
	OperationCloseJob                        int16 = 0x003b
	OperationIdentifyPrinter                 int16 = 0x003c
	OperationValidateDocument                int16 = 0x003d
	OperationAddDocumentImages               int16 = 0x003e
	OperationAcknowledgeDocument             int16 = 0x003f
	OperationAcknowledgeIdentifyPrinter      int16 = 0x0040
	OperationAcknowledgeJob                  int16 = 0x0041
	OperationFetchDocument                   int16 = 0x0042
	OperationFetchJob                        int16 = 0x0043
	OperationGetOutputDeviceAttributes       int16 = 0x0044
	OperationUpdateActiveJobs                int16 = 0x0045
	OperationDeregisterOutputDevice          int16 = 0x0046
	OperationUpdateDocumentStatus            int16 = 0x0047
	OperationUpdateJobStatus                 int16 = 0x0048
	OperationUpdateOutputDeviceAttributes    int16 = 0x0049
	OperationGetNextDocumentData             int16 = 0x004a
	OperationAllocatePrinterResources        int16 = 0x004b
	OperationCreatePrinter                   int16 = 0x004c
	OperationDeallocatePrinterResources      int16 = 0x004d
	OperationDeletePrinter                   int16 = 0x004e
	OperationGetPrinters                     int16 = 0x004f
	OperationShutdownOnePrinter              int16 = 0x0050
	OperationStartupOnePrinter               int16 = 0x0051
	OperationCancelResource                  int16 = 0x0052
	OperationCreateResource                  int16 = 0x0053
	OperationInstallResource                 int16 = 0x0054
	OperationSendResourceData                int16 = 0x0055
	OperationSetResourceAttributes           int16 = 0x0056
	OperationCreateResourceSubscriptions     int16 = 0x0057
	OperationCreateSystemSubscriptions       int16 = 0x0058
	OperationDisableAllPrinters              int16 = 0x0059
	OperationEnableAllPrinters               int16 = 0x005a
	OperationGetSystemAttributes             int16 = 0x005b
	OperationGetSystemSupportedValues        int16 = 0x005c
	OperationPauseAllPrinters                int16 = 0x005d
	OperationPauseAllPrintersAfterCurrentJob int16 = 0x005e
	OperationRegisterOutputDevice            int16 = 0x005f
	OperationRestartSystem                   int16 = 0x0060
	OperationResumeAllPrinters               int16 = 0x0061
	OperationSetSystemAttributes             int16 = 0x0062
	OperationShutdownAllPrinter              int16 = 0x0063
	OperationStartupAllPrinters              int16 = 0x0064
	OperationPrivate                         int16 = 0x4000
	OperationCupsGetDefault                  int16 = 0x4001
	OperationCupsGetPrinters                 int16 = 0x4002
	OperationCupsAddModifyPrinter            int16 = 0x4003
	OperationCupsDeletePrinter               int16 = 0x4004
	OperationCupsGetClasses                  int16 = 0x4005
	OperationCupsAddModifyClass              int16 = 0x4006
	OperationCupsDeleteClass                 int16 = 0x4007
	OperationCupsAcceptJobs                  int16 = 0x4008
	OperationCupsRejectJobs                  int16 = 0x4009
	OperationCupsSetDefault                  int16 = 0x400a
	OperationCupsGetDevices                  int16 = 0x400b
	OperationCupsGetPPDs                     int16 = 0x400c
	OperationCupsMoveJob                     int16 = 0x400d
	OperationCupsAuthenticateJob             int16 = 0x400e
	OperationCupsGetPpd                      int16 = 0x400f
	OperationCupsGetDocument                 int16 = 0x4027
	OperationCupsCreateLocalPrinter          int16 = 0x4028
)

type Tag int8

const (
	TagCupsInvalid       byte = -1
	TagZero              byte = 0x00
	TagOperation         byte = 0x01
	TagJob               byte = 0x02
	TagEnd               byte = 0x03
	TagPrinter           byte = 0x04
	TagUnsupportedGroup  byte = 0x05
	TagSubscription      byte = 0x06
	TagEventNotification byte = 0x07
	TagResource          byte = 0x08
	TagDocument          byte = 0x09
	TagSystem            byte = 0x0a
	TagUnsupportedValue  byte = 0x10
	TagDefault           byte = 0x11
	TagUnknown           byte = 0x12
	TagNoValue           byte = 0x13
	TagNotSettable       byte = 0x15
	TagDeleteAttr        byte = 0x16
	TagAdminDefine       byte = 0x17
	TagInteger           byte = 0x21
	TagBoolean           byte = 0x22
	TagEnum              byte = 0x23
	TagString            byte = 0x30
	TagDate              byte = 0x31
	TagResolution        byte = 0x32
	TagRange             byte = 0x33
	TagBeginCollection   byte = 0x34
	TagTextLang          byte = 0x35
	TagNameLang          byte = 0x36
	TagEndCollection     byte = 0x37
	TagText              byte = 0x41
	TagName              byte = 0x42
	TagReservedString    byte = 0x43
	TagKeyword           byte = 0x44
	TagUri               byte = 0x45
	TagUriScheme         byte = 0x46
	TagCharset           byte = 0x47
	TagLanguage          byte = 0x48
	TagMimeType          byte = 0x49
	TagMemberName        byte = 0x4a
	TagExtension         byte = 0x7f
	TagCupsMask          byte = 0x7fffffff
	TagCupsConst         byte = -0x7fffffff - 1
)

const (
	JobStatePending    uint8 = 0x03
	JobStateHeld       uint8 = 0x04
	JobStateProcessing uint8 = 0x05
	JobStateStopped    uint8 = 0x06
	JobStateCanceled   uint8 = 0x07
	JobStateAborted    uint8 = 0x08
	JobStateCompleted  uint8 = 0x09
)

const (
	DocumentStatePending    uint8 = 0x03
	DocumentStateProcessing uint8 = 0x05
	DocumentStateCanceled   uint8 = 0x07
	DocumentStateAborted    uint8 = 0x08
	DocumentStateCompleted  uint8 = 0x08
)

const (
	PrinterStateIdle       uint8 = 0x0003
	PrinterStateProcessing uint8 = 0x0004
	PrinterStateStopped    uint8 = 0x0005
)

const (
	JobStateFilterNotCompleted = "not-completed"
	JobStateFilterCompleted    = "completed"
	JobStateFilterAll          = "all"
)

type ErrorPolicy string

const (
	ErrorPolicyRetryJob        ErrorPolicy = "retry-job"
	ErrorPolicyAbortJob                    = "abort-job"
	ErrorPolicyRetryCurrentJob             = "retry-current-job"
	ErrorPolicyStopPrinter                 = "stop-printer"
)

const (
	CharsetLanguage      = "en-US"
	Charset              = "utf-8"
	ProtocolVersionMajor = uint8(2)
	ProtocolVersionMinor = uint8(0)

	DefaultJobPriority = 50

	MimeTypePostscript  = "application/postscript"
	MimeTypeOctetStream = "application/octet-stream"

	ContentTypeIPP = "application/ipp"
)

const (
	AttributeCopies                 = "copies"
	AttributeDocumentFormat         = "document-format"
	AttributeDocumentName           = "document-name"
	AttributeJobID                  = "job-id"
	AttributeJobName                = "job-name"
	AttributeJobPriority            = "job-priority"
	AttributeJobURI                 = "job-uri"
	AttributeLastDocument           = "last-document"
	AttributeMyJobs                 = "my-jobs"
	AttributePPDName                = "ppd-name"
	AttributePrinterIsShared        = "printer-is-shared"
	AttributePrinterURI             = "printer-uri"
	AttributePurgeJobs              = "purge-jobs"
	AttributeRequestedAttributes    = "requested-attributes"
	AttributeRequestingUserName     = "requesting-user-name"
	AttributeWhichJobs              = "which-jobs"
	AttributeFirstJobID             = "first-job-id"
	AttributeLimit                  = "limit"
	AttributeStatusMessage          = "status-message"
	AttributeCharset                = "attributes-charset"
	AttributeNaturalLanguage        = "attributes-natural-language"
	AttributeDeviceURI              = "device-uri"
	AttributeHoldJobUntil           = "job-hold-until"
	AttributePrinterErrorPolicy     = "printer-error-policy"
	AttributePrinterInfo            = "printer-info"
	AttributePrinterLocation        = "printer-location"
	AttributePrinterName            = "printer-name"
	AttributePrinterStateReason     = "printer-state-reason"
	AttributeJobPrinterURI          = "job-printer-uri"
	AttributeMemberURIs             = "member-uris"
	AttributeDocumentNumber         = "document-number"
	AttributeDocumentState          = "document-state"
	AttributeFinishings             = "finishings"
	AttributeJobHoldUntil           = "hold-job-until"
	AttributeJobSheets              = "job-sheets"
	AttributeJobState               = "job-state"
	AttributeJobStateReason         = "job-state-reason"
	AttributeMedia                  = "media"
	AttributeNumberUp               = "number-up"
	AttributeOrientationRequested   = "orientation-requested"
	AttributePrintQuality           = "print-quality"
	AttributePrinterIsAcceptingJobs = "printer-is-accepting-jobs"
	AttributePrinterResolution      = "printer-resolution"
	AttributePrinterState           = "printer-state"
	AttributeMemberNames            = "member-names"
	AttributePrinterType            = "printer-type"
	AttributePrinterMarkAndModel    = "printer-make-and-model"
	AttributePrinterStateMessage    = "printer-state-message"
	AttributePrinterUriSupported    = "printer-uri-supported"
	AttributeJobMediaProgress       = "job-media-progress"
	AttributeJobKilobyteOctets      = "job-k-octets"
	AttributeNumberOfDocuments      = "number-of-documents"
	AttributeJobOriginatingUserName = "job-originating-user-name"
)

var (
	DefaultClassAttributes   = []string{AttributePrinterName, AttributeMemberNames}
	DefaultPrinterAttributes = []string{AttributePrinterName, AttributePrinterType, AttributePrinterLocation, AttributePrinterInfo,
		AttributePrinterMarkAndModel, AttributePrinterState, AttributePrinterStateMessage, AttributePrinterStateReason,
		AttributePrinterUriSupported, AttributeDeviceURI, AttributePrinterIsShared}
	DefaultJobAttributes = []string{AttributeJobID, AttributeJobName, AttributePrinterURI, AttributeJobState, AttributeJobStateReason,
		AttributeJobHoldUntil, AttributeJobMediaProgress, AttributeJobKilobyteOctets, AttributeNumberOfDocuments, AttributeCopies,
		AttributeJobOriginatingUserName}

	AttributeTagMapping = map[string]byte{
		AttributeCharset:                TagCharset,
		AttributeNaturalLanguage:        TagLanguage,
		AttributeCopies:                 TagInteger,
		AttributeDeviceURI:              TagUri,
		AttributeDocumentFormat:         TagMimeType,
		AttributeDocumentName:           TagName,
		AttributeDocumentNumber:         TagInteger,
		AttributeDocumentState:          TagEnum,
		AttributeFinishings:             TagEnum,
		AttributeJobHoldUntil:           TagKeyword,
		AttributeHoldJobUntil:           TagKeyword,
		AttributeJobID:                  TagInteger,
		AttributeJobName:                TagName,
		AttributeJobPrinterURI:          TagUri,
		AttributeJobPriority:            TagInteger,
		AttributeJobSheets:              TagName,
		AttributeJobState:               TagEnum,
		AttributeJobStateReason:         TagKeyword,
		AttributeJobURI:                 TagUri,
		AttributeLastDocument:           TagBoolean,
		AttributeMedia:                  TagKeyword,
		AttributeMemberURIs:             TagUri,
		AttributeMyJobs:                 TagBoolean,
		AttributeNumberUp:               TagInteger,
		AttributeOrientationRequested:   TagEnum,
		AttributePPDName:                TagName,
		AttributePrintQuality:           TagEnum,
		AttributePrinterErrorPolicy:     TagName,
		AttributePrinterInfo:            TagText,
		AttributePrinterIsAcceptingJobs: TagBoolean,
		AttributePrinterIsShared:        TagBoolean,
		AttributePrinterLocation:        TagText,
		AttributePrinterResolution:      TagResolution,
		AttributePrinterState:           TagEnum,
		AttributePrinterStateReason:     TagKeyword,
		AttributePrinterURI:             TagUri,
		AttributePurgeJobs:              TagBoolean,
		AttributeRequestedAttributes:    TagKeyword,
		AttributeRequestingUserName:     TagName,
		AttributeWhichJobs:              TagKeyword,
		AttributeFirstJobID:             TagInteger,
		AttributeStatusMessage:          TagText,
	}
)
