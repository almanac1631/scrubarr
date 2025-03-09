import {describe, expect, test} from "vitest";
import {formatFileSize} from "./fileSize.ts";

describe("format file size", () => {
    test("formats the file size on bytes", () => {
        const fileSize = 921;
        const expectedFileSizeString = "921 B";
        expect(formatFileSize(fileSize)).toBe(expectedFileSizeString);
    });
    test("formats the file size on kilobytes", () => {
        const fileSize = 49212;
        const expectedFileSizeString = "49.21 KB";
        expect(formatFileSize(fileSize)).toBe(expectedFileSizeString);
    });
    test("formats the file size on megabytes", () => {
        const fileSize = 9382912;
        const expectedFileSizeString = "9.38 MB";
        expect(formatFileSize(fileSize)).toBe(expectedFileSizeString);
    });
    test("formats the file size on gigabytes", () => {
        const fileSize = 3290802842;
        const expectedFileSizeString = "3.29 GB";
        expect(formatFileSize(fileSize)).toBe(expectedFileSizeString);
    });
    test("formats the file size on undefined", () => {
        const fileSize = undefined;
        const expectedFileSizeString = "0 B";
        expect(formatFileSize(fileSize)).toBe(expectedFileSizeString);
    });
});