// +build darwin
package context

func (ctx *hostContext) RefreshStatus() error {
    ctx.Lock()
    defer ctx.Unlock()

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
    ctx.applicationBundleVersion    = findApplicationBundleVersionString()
    ctx.applicationBundleExpiration = findApplicationBundleExpirationString()

    ctx.processorCount              = findSystemProcessorCount()
    ctx.activeProcessorCount        = findSystemActiveProcessorCount()
    ctx.physicalMemorySize          = findSystemPhysicalMemorySize()
    ctx.physicalCoreCount           = findSystemPhysicalCoreCount()

    // (2017-05-21) This will be replaced with randomly generated 16 chars
    //ctx.clusterPublicName         = findSerialNumber()
    ctx.currentLanguageCode         = findCurrentLanguageCode()
    ctx.currentCountryCode          = findCurrentCountryCode()

    return nil
}
