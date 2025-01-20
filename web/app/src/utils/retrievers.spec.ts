import {describe, expect, test} from "vitest";
import {Retriever} from "../api";
import {getCategoriesFromRetrieverList, RetrieverCategory, sortRetrieverList} from "./retrievers.ts";

describe("sort retriever list", () => {
    test("sorts by category", () => {
        const retrievers: Retriever[] = [
            {id: "1", name: "main", category: "torrent_client", softwareName: "deluge"},
            {id: "2", name: "main", category: "arr_app", softwareName: "sonarr"},
            {id: "3", name: "main", category: "folder", softwareName: "folder"},
        ];
        sortRetrieverList(retrievers);
        const expectedRetrievers: Retriever[] = [
            {id: "2", name: "main", category: "arr_app", softwareName: "sonarr"},
            {id: "3", name: "main", category: "folder", softwareName: "folder"},
            {id: "1", name: "main", category: "torrent_client", softwareName: "deluge"},
        ];
        expect(retrievers).toStrictEqual(expectedRetrievers);
    });
    test("sorts by category and software name", () => {
        const retrievers: Retriever[] = [
            {id: "2", name: "main", category: "torrent_client", softwareName: "sonarr"},
            {id: "1", name: "main", category: "torrent_client", softwareName: "deluge"},
            {id: "3", name: "main", category: "folder", softwareName: "folder"},
        ];
        sortRetrieverList(retrievers);
        const expectedRetrievers: Retriever[] = [
            {id: "3", name: "main", category: "folder", softwareName: "folder"},
            {id: "1", name: "main", category: "torrent_client", softwareName: "deluge"},
            {id: "2", name: "main", category: "torrent_client", softwareName: "sonarr"},
        ];
        expect(retrievers).toStrictEqual(expectedRetrievers);
    });
    test("sorts by category, software name and name", () => {
        const retrievers: Retriever[] = [
            {id: "5", name: "main2", category: "torrent_client", softwareName: "sonarr"},
            {id: "4", name: "main3", category: "torrent_client", softwareName: "sonarr"},
            {id: "2", name: "main1", category: "torrent_client", softwareName: "sonarr"},
            {id: "1", name: "main", category: "torrent_client", softwareName: "deluge"},
            {id: "3", name: "main", category: "folder", softwareName: "folder"},
        ];
        sortRetrieverList(retrievers);
        const expectedRetrievers: Retriever[] = [
            {id: "3", name: "main", category: "folder", softwareName: "folder"},
            {id: "1", name: "main", category: "torrent_client", softwareName: "deluge"},
            {id: "2", name: "main1", category: "torrent_client", softwareName: "sonarr"},
            {id: "5", name: "main2", category: "torrent_client", softwareName: "sonarr"},
            {id: "4", name: "main3", category: "torrent_client", softwareName: "sonarr"},
        ];
        expect(retrievers).toStrictEqual(expectedRetrievers);
    });
});

describe("get retriever categories from retriever list", () => {
    test("gets retriever categories from retriever list", () => {
        const retrievers: Retriever[] = [
            {id: "2", name: "main", category: "torrent_client", softwareName: "rtorrent"},
            {id: "1", name: "main", category: "torrent_client", softwareName: "deluge"},
            {id: "3", name: "main", category: "folder", softwareName: "folder"},
            {id: "4", name: "main", category: "arr_app", softwareName: "sonarr"},
        ];
        const actualCategories = getCategoriesFromRetrieverList(retrievers);
        const expectedCategories: RetrieverCategory[] = [
            {displayName: "Torrent Clients", name: "torrent_client", logoFilename: "category/torrent-client-logo.svg"},
            {displayName: "Folders", name: "folder", logoFilename: "category/folder-logo.svg"},
            {displayName: "*arr apps", name: "arr_app", logoFilename: "category/arr-app-logo.svg"},
        ]
        expect(actualCategories).toStrictEqual(expectedCategories);
    })
});