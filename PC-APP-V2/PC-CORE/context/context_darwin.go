// +build darwin
package context

func (ctx *hostContext) RefreshStatus() error {

    ctx.cocoaHomePath               = findCocoaHomeDirectory()
    ctx.posixHomePath               = findPosixHomeDirectory()
    ctx.fullUserName                = findFullUserName()
    ctx.loginUserName               = findLoginUserName()
    ctx.userTempPath                = findUserTemporaryDirectory()

    ctx.applicationSupportPath      = findApplicationSupportDirectory()
    ctx.applicationDocumentPath     = findApplicationDocumentsDirectoru()
    ctx.applicationTempPath         = findApplicationTemporaryDirectory()
    ctx.applicationLibCachePath     = findApplicationLibraryCacheDirectory()
    ctx.applicationResourcePath     = findApplicationResourceDirectory()
    ctx.applicationExecutablePath   = findApplicationExecutableDirectory()

    ctx.processorCount              = findSystemProcessorCount()
    ctx.activeProcessorCount        = findSystemActiveProcessorCount()
    ctx.physicalMemorySize          = findSystemPhysicalMemorySize()

    ctx.hostDeviceSerial            = findSerialNumber()
    ctx.currentLanguageCode         = findCurrentLanguageCode()
    ctx.currentCountryCode          = findCurrentCountryCode()

    return nil
}
