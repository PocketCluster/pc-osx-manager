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

    ctx.hostDeviceSerial            = findSerialNumber()

    ctx.processorCount              = findSystemProcessorCount()
    ctx.activeProcessorCount        = findSystemActiveProcessorCount()
    ctx.physicalMemorySize          = findSystemPhysicalMemorySize()

    return nil
}
