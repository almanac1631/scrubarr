import {Retriever} from "../api";

const compareAttributeList = [
    (retriever: Retriever) => retriever.category,
    (retriever: Retriever) => retriever.softwareName,
    (retriever: Retriever) => retriever.name,
];

export function sortRetrieverList(retrievers: Retriever[]) {
    return retrievers.sort((retrieverA, retrieverB) => {
        for (const compareAttributeFunction of compareAttributeList) {
            const compareResult = compareAttributeFunction(retrieverA).localeCompare(compareAttributeFunction(retrieverB));
            if (compareResult === 0) {
                continue
            }
            return compareResult;
        }
        return 0;
    });
}