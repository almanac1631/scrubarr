import {describe, expect, test} from "vitest";
import {Retriever} from "../api";
import {sortRetrieverList} from "./retriever-sorting.ts";

describe("sort retriever list", () => {
    test("sort by category", () => {
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
    test("sort by category and software name", () => {
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
    test("sort by category, software name and name", () => {
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