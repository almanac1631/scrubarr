export function formatFileSize(sizeInBytes: number | undefined): string {
    if (sizeInBytes === undefined) {
        return "0 B";
    }
    if (sizeInBytes < 1000) {
        return `${sizeInBytes} B`;
    }
    const sizeInKilobytes = sizeInBytes / 1000;
    if (sizeInKilobytes < 1000) {
        return `${sizeInKilobytes.toFixed(2)} KB`;
    }
    const sizeInMegabytes = sizeInKilobytes / 1000;
    if (sizeInMegabytes < 1000) {
        return `${sizeInMegabytes.toFixed(2)} MB`;
    }
    const sizeInGigabytes = sizeInMegabytes / 1000;
    return `${sizeInGigabytes.toFixed(2)} GB`;
}
