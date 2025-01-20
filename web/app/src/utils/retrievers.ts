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

export interface RetrieverCategory {
    displayName: string,
    name: string,
    logoFilename: string
}

const retrieverCategoryMapping = {
    "torrent_client": "Torrent Clients",
    "folder": "Folders",
    "arr_app": "*arr apps",
}

export function getCategoriesFromRetrieverList(retrieverList: Retriever[]): RetrieverCategory[] {
    return retrieverList.reduce((categoryList, retriever) => {
        if (categoryList.filter((category) => category.name === retriever.category).length > 0) {
            return categoryList;
        }
        const categoryName = retriever.category;
        console.log(retrieverCategoryMapping);
        categoryList.push({
            displayName: retrieverCategoryMapping[categoryName],
            name: categoryName,
            logoFilename: `category/${categoryName.replace("_", "-")}-logo.svg`
        });
        return categoryList;
    }, [] as RetrieverCategory[])
}