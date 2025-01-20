import {getPageList} from "./pagination.ts";
import {describe, expect, test} from "vitest";

describe("calculate pages", () => {
    test("calculates the amount of pages correctly", () => {
        const pageList = getPageList(100, 10, 1);
        expect(pageList.length).toBe(10);
    });
    test("calculates the amount of pages correctly for an odd amount of items", () => {
        const pageList = getPageList(81, 10, 1);
        expect(pageList.length).toBe(9);
    });
    test("calculates the amount of pages correctly for no items", () => {
        const pageList = getPageList(0, 0, 1);
        expect(pageList.length).toBe(0);
    });
    test("calculates the page numbers correctly", () => {
        const pageList = getPageList(100, 10, 1);
        for (let i = 0; i < 10; i++) {
            expect(pageList[i].pageNumber).toBe(i + 1);
        }
    });
    test("sets the correct page as selected", () => {
        const pageList = getPageList(81, 10, 6);
        expect(pageList[5]).toStrictEqual({
            pageNumber: 6,
            isSelected: true,
        })
    });
    test("includes placeholder correctly on page 1", () => {
        const pageList = getPageList(81, 10, 1, 6);
        expect(pageList).toStrictEqual([
            {pageNumber: 1, isSelected: true},
            {pageNumber: 2, isSelected: false},
            {pageNumber: null, isSelected: false},
            {pageNumber: 9, isSelected: false},
        ])
    });
    test("does not set placeholder when the active page is 5", () => {
        const pageList = getPageList(81, 10, 5, 6);
        expect(pageList).toStrictEqual([
            {pageNumber: 1, isSelected: false},
            {pageNumber: null, isSelected: false},
            {pageNumber: 4, isSelected: false},
            {pageNumber: 5, isSelected: true},
            {pageNumber: 6, isSelected: false},
            {pageNumber: null, isSelected: false},
            {pageNumber: 9, isSelected: false},
        ])
    });
});