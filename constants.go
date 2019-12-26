package ipp

type Status int8

const (
	StatusCupsInvalid                         Status = -1
	StatusOk                                         = 0x0000
	StatusOkIgnoredOrSubstituted                     = 0x0001
	StatusOkConflicting                              = 0x0002
	statusOkIgnoredSubscriptions                     = 0x0003
	statusOkIgnoredNotifications                     = 0x0004
	statusOkTooManyEvents                            = 0x0005
	statusOkButCancelSubscription                    = 0x0006
	statusOkEventsComplete                           = 0x0007
	statusRedirectionOtherSite                       = 0x0200
	statusCupsSeeOther                               = 0x0280
	statusErrorBadRequest                            = 0x0400
	StatusErrorForbidden                             = 0x0401
	StatusErrorNotAuthenticated                      = 0x0402
	StatusErrorNotAuthorized                         = 0x0403
	StatusErrorNotPossible                           = 0x0404
	StatusErrorTimeout                               = 0x0405
	StatusErrorNotFound                              = 0x0406
	StatusErrorGone                                  = 0x0407
	StatusErrorRequestEntity                         = 0x0408
	StatusErrorRequestValue                          = 0x0409
	StatusErrorDocumentFormatNotSupported            = 0x040a
	StatusErrorAttributesOrValues                    = 0x040b
	StatusErrorUriScheme                             = 0x040c
	StatusErrorCharset                               = 0x040d
	StatusErrorConflicting                           = 0x040e
	StatusErrorCompressionError                      = 0x040f
	StatusErrorDocumentFormatError                   = 0x0410
	StatusErrorDocumentAccess                        = 0x0411
	StatusErrorAttributesNotSettable                 = 0x0412
	StatusErrorIgnoredAllSubscriptions               = 0x0413
	StatusErrorTooManySubscriptions                  = 0x0414
	StatusErrorIgnoredAllNotifications               = 0x0415
	StatusErrorPrintSupportFileNotFound              = 0x0416
	StatusErrorDocumentPassword                      = 0x0417
	StatusErrorDocumentPermission                    = 0x0418
	StatusErrorDocumentSecurity                      = 0x0419
	StatusErrorDocumentUnprintable                   = 0x041a
	StatusErrorAccountInfoNeeded                     = 0x041b
	StatusErrorAccountClosed                         = 0x041c
	StatusErrorAccountLimitReached                   = 0x041d
	StatusErrorAccountAuthorizationFailed            = 0x041e
	StatusErrorNotFetchable                          = 0x041f
	StatusErrorCupsAccountInfoNeeded                 = 0x049C
	StatusErrorCupsAccountClosed                     = 0x049d
	StatusErrorCupsAccountLimitReached               = 0x049e
	StatusErrorCupsAccountAuthorizationFailed        = 0x049f
	StatusErrorInternal                              = 0x0500
	StatusErrorOperationNotSupported                 = 0x0501
	StatusErrorServiceUnavailable                    = 0x0502
	StatusErrorVersionNotSupported                   = 0x0503
	StatusErrorDevice                                = 0x0504
	StatusErrorTemporary                             = 0x0505
	StatusErrorNotAcceptingJobs                      = 0x0506
	StatusErrorBusy                                  = 0x0507
	StatusErrorJobCanceled                           = 0x0508
	StatusErrorMultipleJobsNotSupported              = 0x0509
	StatusErrorPrinterIsDeactivated                  = 0x050a
	StatusErrorTooManyJobs                           = 0x050b
	StatusErrorTooManyDocuments                      = 0x050c
	StatusErrorCupsAuthenticationCanceled            = 0x1000
	StatusErrorCupsPki                               = 0x1001
	StatusErrorCupsUpgradeRequired                   = 0x1002
)

type Operation int16

const (
	OperationCupsInvalid                     Operation = -0x0001
	OperationCupsNone                                  = 0x0000
	OperationPrintJob                                  = 0x0002
	OperationPrintUri                                  = 0x0003
	OperationValidateJob                               = 0x0004
	OperationCreateJob                                 = 0x0005
	OperationSendDocument                              = 0x0006
	OperationSendUri                                   = 0x0007
	OperationCancelJob                                 = 0x0008
	OperationGetJobAttributes                          = 0x0009
	OperationGetJobs                                   = 0x000a
	OperationGetPrinterAttributes                      = 0x000b
	OperationHoldJob                                   = 0x000c
	OperationReleaseJob                                = 0x000d
	OperationRestartJob                                = 0x000e
	OperationPausePrinter                              = 0x0010
	OperationResumePrinter                             = 0x0011
	OperationPurgeJobs                                 = 0x0012
	OperationSetPrinterAttributes                      = 0x0013
	OperationSetJobAttributes                          = 0x0014
	OperationGetPrinterSupportedValues                 = 0x0015
	OperationCreatePrinterSubscriptions                = 0x0016
	OperationCreateJobSubscriptions                    = 0x0017
	OperationGetSubscriptionAttributes                 = 0x0018
	OperationGetSubscriptions                          = 0x0019
	OperationRenewSubscription                         = 0x001a
	OperationCancelSubscription                        = 0x001b
	OperationGetNotifications                          = 0x001c
	OperationSendNotifications                         = 0x001d
	OperationGetResourceAttributes                     = 0x001e
	OperationGetResourceData                           = 0x001f
	OperationGetResources                              = 0x0020
	OperationGetPrintSupportFiles                      = 0x0021
	OperationEnablePrinter                             = 0x0022
	OperationDisablePrinter                            = 0x0023
	OperationPausePrinterAfterCurrentJob               = 0x0024
	OperationHoldNewJobs                               = 0x0025
	OperationReleaseHeldNewJobs                        = 0x0026
	OperationDeactivatePrinter                         = 0x0027
	OperationActivatePrinter                           = 0x0028
	OperationRestartPrinter                            = 0x0029
	OperationShutdownPrinter                           = 0x002a
	OperationStartupPrinter                            = 0x002b
	OperationReprocessJob                              = 0x002c
	OperationCancelCurrentJob                          = 0x002d
	OperationSuspendCurrentJob                         = 0x002e
	OperationResumeJob                                 = 0x002f
	OperationOperationPromoteJob                       = 0x0030
	OperationScheduleJobAfter                          = 0x0031
	OperationCancelDocument                            = 0x0033
	OperationGetDocumentAttributes                     = 0x0034
	OperationGetDocuments                              = 0x0035
	OperationDeleteDocument                            = 0x0036
	OperationSetDocumentAttributes                     = 0x0037
	OperationCancelJobs                                = 0x0038
	OperationCancelMyJobs                              = 0x0039
	OperationResubmitJob                               = 0x003a
	OperationCloseJob                                  = 0x003b
	OperationIdentifyPrinter                           = 0x003c
	OperationValidateDocument                          = 0x003d
	OperationAddDocumentImages                         = 0x003e
	OperationAcknowledgeDocument                       = 0x003f
	OperationAcknowledgeIdentifyPrinter                = 0x0040
	OperationAcknowledgeJob                            = 0x0041
	OperationFetchDocument                             = 0x0042
	OperationFetchJob                                  = 0x0043
	OperationGetOutputDeviceAttributes                 = 0x0044
	OperationUpdateActiveJobs                          = 0x0045
	OperationDeregisterOutputDevice                    = 0x0046
	OperationUpdateDocumentStatus                      = 0x0047
	OperationUpdateJobStatus                           = 0x0048
	OperationUpdateOutputDeviceAttributes              = 0x0049
	OperationGetNextDocumentData                       = 0x004a
	OperationAllocatePrinterResources                  = 0x004b
	OperationCreatePrinter                             = 0x004c
	OperationDeallocatePrinterResources                = 0x004d
	OperationDeletePrinter                             = 0x004e
	OperationGetPrinters                               = 0x004f
	OperationShutdownOnePrinter                        = 0x0050
	OperationStartupOnePrinter                         = 0x0051
	OperationCancelResource                            = 0x0052
	OperationCreateResource                            = 0x0053
	OperationInstallResource                           = 0x0054
	OperationSendResourceData                          = 0x0055
	OperationSetResourceAttributes                     = 0x0056
	OperationCreateResourceSubscriptions               = 0x0057
	OperationCreateSystemSubscriptions                 = 0x0058
	OperationDisableAllPrinters                        = 0x0059
	OperationEnableAllPrinters                         = 0x005a
	OperationGetSystemAttributes                       = 0x005b
	OperationGetSystemSupportedValues                  = 0x005c
	OperationPauseAllPrinters                          = 0x005d
	OperationPauseAllPrintersAfterCurrentJob           = 0x005e
	OperationRegisterOutputDevice                      = 0x005f
	OperationRestartSystem                             = 0x0060
	OperationResumeAllPrinters                         = 0x0061
	OperationSetSystemAttributes                       = 0x0062
	OperationShutdownAllPrinter                        = 0x0063
	OperationStartupAllPrinters                        = 0x0064
	OperationPrivate                                   = 0x4000
	OperationCupsGetDefault                            = 0x4001
	OperationCupsGetPrinters                           = 0x4002
	OperationCupsAddModifyPrinter                      = 0x4003
	OperationCupsDeletePrinter                         = 0x4004
	OperationCupsGetClasses                            = 0x4005
	OperationCupsAddModifyClass                        = 0x4006
	OperationCupsDeleteClass                           = 0x4007
	OperationCupsAcceptJobs                            = 0x4008
	OperationCupsRejectJobs                            = 0x4009
	OperationCupsSetDefault                            = 0x400a
	OperationCupsGetDevices                            = 0x400b
	OperationCupsGetPpds                               = 0x400c
	OperationCupsMoveJob                               = 0x400d
	OperationCupsAuthenticateJob                       = 0x400e
	OperationCupsGetPpd                                = 0x400f
	OperationCupsGetDocument                           = 0x4027
	OperationCupsCreateLocalPrinter                    = 0x4028
)

type Tag int8

const (
	TagCupsInvalid       Tag = -1
	TagZero                  = 0x00
	TagOperation             = 0x01
	TagJob                   = 0x02
	TagEnd                   = 0x03
	TagPrinter               = 0x04
	TagUnsupportedGroup      = 0x05
	TagSubscription          = 0x06
	TagEventNotification     = 0x07
	TagResource              = 0x08
	TagDocument              = 0x09
	TagSystem                = 0x0a
	TagUnsupportedValue      = 0x10
	TagDefault               = 0x11
	TagUnknown               = 0x12
	TagNoValue               = 0x013
	TagNotSettable           = 0x15
	TagDeleteAttr            = 0x16
	TagAdminDefine           = 0x17
	TagInteger               = 0x21
	TagBoolean               = 0x22
	TagEnum                  = 0x23
	TagString                = 0x30
	TagDate                  = 0x31
	TagResolution            = 0x32
	TagRange                 = 0x33
	TagBeginCollection       = 0x34
	TagTextLang              = 0x35
	TagNameLang              = 0x36
	TagEndCollection         = 0x37
	TagText                  = 0x41
	TagName                  = 0x42
	TagReservedString        = 0x43
	TagKeyword               = 0x44
	TagUri                   = 0x45
	TagUriScheme             = 0x46
	TagCharset               = 0x47
	TagLanguage              = 0x48
	TagMimeType              = 0x49
	TagMemberName            = 0x4a
	TagExtension             = 0x7f
	TagCupsMask              = 0x7fffffff
	TagCupsConst             = -0x7fffffff - 1
)

type JobState uint8

const (
	JobStatePending    JobState = 0x03
	JobStateHeld                = 0x04
	JobStateProcessing          = 0x05
	JobStateStopped             = 0x06
	JobStateCanceled            = 0x07
	JobStateAborted             = 0x08
	JobStateCompleted           = 0x09
)

type DocumentState uint8

const (
	DocumentStatePending    DocumentState = 0x03
	DocumentStateProcessing               = 0x05
	DocumentStateCanceled                 = 0x07
	DocumentStateAborted                  = 0x08
	DocumentStateCompleted                = 0x08
)

type PrinterState uint8

const (
	PrinterStateIdle       PrinterState = 0x0003
	PrinterStateProcessing              = 0x0004
	PrinterStateStopped                 = 0x0005
)

type JobStateFilter string

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
)

const (
	OperationAttributeCopies              string = "copies"
	OperationAttributeDocumentFormat             = "document-format"
	OperationAttributeDocumentName               = "document-name"
	OperationAttributeJobID                      = "job-id"
	OperationAttributeJobName                    = "job-name"
	OperationAttributeJobPriority                = "job-priority"
	OperationAttributeJobURI                     = "job-uri"
	OperationAttributeLastDocument               = "last-document"
	OperationAttributeMyJobs                     = "my-jobs"
	OperationAttributePPDName                    = "ppd-name"
	OperationAttributePrinterIsShared            = "printer-is-shared"
	OperationAttributePrinterURI                 = "printer-uri"
	OperationAttributePurgeJobs                  = "purge-jobs"
	OperationAttributeRequestedAttributes        = "requested-attributes"
	OperationAttributeRequestingUserName         = "requesting-user-name"
	OperationAttributeWhichJobs                  = "which-jobs"
	OperationAttributeFirstJobID                 = "first-job-id"
	OperationAttributeLimit                      = "limit"
)

const (
	PrinterAttributeDeviceURI          string = "device-uri"
	PrinterAttributeHoldJobUntil              = "job-hold-until"
	PrinterAttributePrinterErrorPolicy        = "printer-error-policy"
	PrinterAttributePrinterInfo               = "printer-info"
	PrinterAttributePrinterLocation           = "printer-location"
	PrinterAttributePrinterName               = "printer-name"
	PrinterAttributePrinterStateReason        = "printer-state-reason"
	PrinterAttributeJobPrinterURI             = "job-printer-uri"
	PrinterAttributeMemberURIs                = "member-uris"
)

var (
	DefaultClassAttributes   = []string{"printer-name", "member-names"}
	DefaultPrinterAttributes = []string{"printer-name", "printer-type", "printer-location", "printer-info",
		"printer-make-and-model", "printer-state", "printer-state-message", "printer-state-reason",
		"printer-uri-supported", "device-uri", "printer-is-shared"}
	DefaultJobAttributes = []string{"job-id", "job-name", "printer-uri", "job-state", "job-state-reasons",
		"job-hold-until", "job-media-progress", "job-k-octets", "number-of-documents", "copies",
		"job-originating-user-name"}

	AttributeTagMapping = map[string]Tag{
		"attributes-charset":          TagCharset,
		"attributes-natural-language": TagLanguage,
		"copies":                      TagInteger,
		"device-uri":                  TagUri,
		"document-format":             TagMimeType,
		"document-name":               TagName,
		"document-number":             TagInteger,
		"document-state":              TagEnum,
		"finishings":                  TagEnum,
		"hold-job-until":              TagKeyword,
		"job-hold-until":              TagKeyword,
		"job-id":                      TagInteger,
		"job-name":                    TagName,
		"job-printer-uri":             TagUri,
		"job-priority":                TagInteger,
		"job-sheets":                  TagName,
		"job-state":                   TagEnum,
		"job-state-reason":            TagKeyword,
		"job-uri":                     TagUri,
		"last-document":               TagBoolean,
		"media":                       TagKeyword,
		"member-uris":                 TagUri,
		"my-jobs":                     TagBoolean,
		"number-up":                   TagInteger,
		"orientation-requested":       TagEnum,
		"ppd-name":                    TagName,
		"print-quality":               TagEnum,
		"printer-error-policy":        TagName,
		"printer-info":                TagText,
		"printer-is-shared":           TagBoolean,
		"printer-location":            TagText,
		"printer-resolution":          TagResolution,
		"printer-state":               TagEnum,
		"printer-state-reason":        TagKeyword,
		"printer-uri":                 TagUri,
		"purge-jobs":                  TagBoolean,
		"requested-attributes":        TagKeyword,
		"requesting-user-name":        TagName,
		"which-jobs":                  TagKeyword,
		"first-job-id":                TagInteger,
	}
)
