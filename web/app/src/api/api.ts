/* tslint:disable */
/* eslint-disable */
/**
 * scrubarr API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 0.0.1
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


import type { Configuration } from './configuration';
import type { AxiosPromise, AxiosInstance, RawAxiosRequestConfig } from 'axios';
import globalAxios from 'axios';
// Some imports not used depending on template conditions
// @ts-ignore
import { DUMMY_BASE_URL, assertParamExists, setApiKeyToObject, setBasicAuthToObject, setBearerAuthToObject, setOAuthToObject, setSearchParams, serializeDataIfNeeded, toPathString, createRequestFunction } from './common';
import type { RequestArgs } from './base';
// @ts-ignore
import { BASE_PATH, COLLECTION_FORMATS, BaseAPI, RequiredError, operationServerMap } from './base';

/**
 * 
 * @export
 * @interface EntryMapping
 */
export interface EntryMapping {
    /**
     * The name of this entry.
     * @type {string}
     * @memberof EntryMapping
     */
    'name': string;
    /**
     * The date and time this entry was added.
     * @type {string}
     * @memberof EntryMapping
     */
    'dateAdded': string;
    /**
     * The size of this entry in bytes.
     * @type {number}
     * @memberof EntryMapping
     */
    'size': number;
    /**
     * 
     * @type {Array<EntryMappingRetrieverFindingsInner>}
     * @memberof EntryMapping
     */
    'retrieverFindings': Array<EntryMappingRetrieverFindingsInner>;
}
/**
 * 
 * @export
 * @interface EntryMappingRetrieverFindingsInner
 */
export interface EntryMappingRetrieverFindingsInner {
    /**
     * Id used to identify retrievers.
     * @type {string}
     * @memberof EntryMappingRetrieverFindingsInner
     */
    'id': string;
}
/**
 * 
 * @export
 * @interface ErrorResponseBody
 */
export interface ErrorResponseBody {
    /**
     * 
     * @type {string}
     * @memberof ErrorResponseBody
     */
    'error': string;
    /**
     * 
     * @type {string}
     * @memberof ErrorResponseBody
     */
    'detail': string;
}
/**
 * 
 * @export
 * @interface GetEntryMappings200Response
 */
export interface GetEntryMappings200Response {
    /**
     * 
     * @type {Array<EntryMapping>}
     * @memberof GetEntryMappings200Response
     */
    'entries': Array<EntryMapping>;
    /**
     * The total amount of entries that could be returned for the provided filter.
     * @type {number}
     * @memberof GetEntryMappings200Response
     */
    'totalAmount': number;
}
/**
 * 
 * @export
 * @interface GetRetrievers200Response
 */
export interface GetRetrievers200Response {
    /**
     * 
     * @type {Array<Retriever>}
     * @memberof GetRetrievers200Response
     */
    'retrievers': Array<Retriever>;
}
/**
 * 
 * @export
 * @interface Login200Response
 */
export interface Login200Response {
    /**
     * 
     * @type {string}
     * @memberof Login200Response
     */
    'message': string;
    /**
     * 
     * @type {string}
     * @memberof Login200Response
     */
    'token': string;
}
/**
 * 
 * @export
 * @interface LoginRequestBody
 */
export interface LoginRequestBody {
    /**
     * The username to login with.
     * @type {string}
     * @memberof LoginRequestBody
     */
    'username': string;
    /**
     * The password to login with.
     * @type {string}
     * @memberof LoginRequestBody
     */
    'password': string;
}
/**
 * 
 * @export
 * @interface RefreshEntryMappings200Response
 */
export interface RefreshEntryMappings200Response {
    /**
     * The status message to display.
     * @type {string}
     * @memberof RefreshEntryMappings200Response
     */
    'message': string;
}
/**
 * 
 * @export
 * @interface Retriever
 */
export interface Retriever {
    /**
     * Id used to identify retrievers.
     * @type {string}
     * @memberof Retriever
     */
    'id': string;
    /**
     * The category this retriever belongs to.
     * @type {string}
     * @memberof Retriever
     */
    'category': RetrieverCategoryEnum;
    /**
     * The name of the retriever\'s software.
     * @type {string}
     * @memberof Retriever
     */
    'softwareName': RetrieverSoftwareNameEnum;
    /**
     * The provided name used to differentiate between multiple instances of the same software retrievers.
     * @type {string}
     * @memberof Retriever
     */
    'name': string;
}

export const RetrieverCategoryEnum = {
    TorrentClient: 'torrent_client',
    Folder: 'folder',
    ArrApp: 'arr_app'
} as const;

export type RetrieverCategoryEnum = typeof RetrieverCategoryEnum[keyof typeof RetrieverCategoryEnum];
export const RetrieverSoftwareNameEnum = {
    Deluge: 'deluge',
    Rtorrent: 'rtorrent',
    Folder: 'folder',
    Sonarr: 'sonarr',
    Radarr: 'radarr'
} as const;

export type RetrieverSoftwareNameEnum = typeof RetrieverSoftwareNameEnum[keyof typeof RetrieverSoftwareNameEnum];


/**
 * DefaultApi - axios parameter creator
 * @export
 */
export const DefaultApiAxiosParamCreator = function (configuration?: Configuration) {
    return {
        /**
         * 
         * @summary Get a list of entry mappings.
         * @param {number} page The page number to display.
         * @param {number} pageSize The amount of items to display per each page.
         * @param {GetEntryMappingsFilterEnum} [filter] The filter to apply before returning the entries.
         * @param {GetEntryMappingsSortByEnum} [sortBy] The criteria to sort the entries by.
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getEntryMappings: async (page: number, pageSize: number, filter?: GetEntryMappingsFilterEnum, sortBy?: GetEntryMappingsSortByEnum, options: RawAxiosRequestConfig = {}): Promise<RequestArgs> => {
            // verify required parameter 'page' is not null or undefined
            assertParamExists('getEntryMappings', 'page', page)
            // verify required parameter 'pageSize' is not null or undefined
            assertParamExists('getEntryMappings', 'pageSize', pageSize)
            const localVarPath = `/entry-mappings`;
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, DUMMY_BASE_URL);
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }

            const localVarRequestOptions = { method: 'GET', ...baseOptions, ...options};
            const localVarHeaderParameter = {} as any;
            const localVarQueryParameter = {} as any;

            // authentication BearerAuth required
            // http bearer authentication required
            await setBearerAuthToObject(localVarHeaderParameter, configuration)

            if (page !== undefined) {
                localVarQueryParameter['page'] = page;
            }

            if (pageSize !== undefined) {
                localVarQueryParameter['pageSize'] = pageSize;
            }

            if (filter !== undefined) {
                localVarQueryParameter['filter'] = filter;
            }

            if (sortBy !== undefined) {
                localVarQueryParameter['sortBy'] = sortBy;
            }


    
            setSearchParams(localVarUrlObj, localVarQueryParameter);
            let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
            localVarRequestOptions.headers = {...localVarHeaderParameter, ...headersFromBaseOptions, ...options.headers};

            return {
                url: toPathString(localVarUrlObj),
                options: localVarRequestOptions,
            };
        },
        /**
         * 
         * @summary Get a list of retrievers.
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getRetrievers: async (options: RawAxiosRequestConfig = {}): Promise<RequestArgs> => {
            const localVarPath = `/retrievers`;
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, DUMMY_BASE_URL);
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }

            const localVarRequestOptions = { method: 'GET', ...baseOptions, ...options};
            const localVarHeaderParameter = {} as any;
            const localVarQueryParameter = {} as any;

            // authentication BearerAuth required
            // http bearer authentication required
            await setBearerAuthToObject(localVarHeaderParameter, configuration)


    
            setSearchParams(localVarUrlObj, localVarQueryParameter);
            let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
            localVarRequestOptions.headers = {...localVarHeaderParameter, ...headersFromBaseOptions, ...options.headers};

            return {
                url: toPathString(localVarUrlObj),
                options: localVarRequestOptions,
            };
        },
        /**
         * 
         * @summary Login to the application using the provided credentials.
         * @param {LoginRequestBody} loginRequestBody 
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        login: async (loginRequestBody: LoginRequestBody, options: RawAxiosRequestConfig = {}): Promise<RequestArgs> => {
            // verify required parameter 'loginRequestBody' is not null or undefined
            assertParamExists('login', 'loginRequestBody', loginRequestBody)
            const localVarPath = `/login`;
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, DUMMY_BASE_URL);
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }

            const localVarRequestOptions = { method: 'POST', ...baseOptions, ...options};
            const localVarHeaderParameter = {} as any;
            const localVarQueryParameter = {} as any;


    
            localVarHeaderParameter['Content-Type'] = 'application/json';

            setSearchParams(localVarUrlObj, localVarQueryParameter);
            let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
            localVarRequestOptions.headers = {...localVarHeaderParameter, ...headersFromBaseOptions, ...options.headers};
            localVarRequestOptions.data = serializeDataIfNeeded(loginRequestBody, localVarRequestOptions, configuration)

            return {
                url: toPathString(localVarUrlObj),
                options: localVarRequestOptions,
            };
        },
        /**
         * 
         * @summary Trigger a refresh of the entry mappings.
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        refreshEntryMappings: async (options: RawAxiosRequestConfig = {}): Promise<RequestArgs> => {
            const localVarPath = `/entry-mappings`;
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, DUMMY_BASE_URL);
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }

            const localVarRequestOptions = { method: 'POST', ...baseOptions, ...options};
            const localVarHeaderParameter = {} as any;
            const localVarQueryParameter = {} as any;

            // authentication BearerAuth required
            // http bearer authentication required
            await setBearerAuthToObject(localVarHeaderParameter, configuration)


    
            setSearchParams(localVarUrlObj, localVarQueryParameter);
            let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
            localVarRequestOptions.headers = {...localVarHeaderParameter, ...headersFromBaseOptions, ...options.headers};

            return {
                url: toPathString(localVarUrlObj),
                options: localVarRequestOptions,
            };
        },
    }
};

/**
 * DefaultApi - functional programming interface
 * @export
 */
export const DefaultApiFp = function(configuration?: Configuration) {
    const localVarAxiosParamCreator = DefaultApiAxiosParamCreator(configuration)
    return {
        /**
         * 
         * @summary Get a list of entry mappings.
         * @param {number} page The page number to display.
         * @param {number} pageSize The amount of items to display per each page.
         * @param {GetEntryMappingsFilterEnum} [filter] The filter to apply before returning the entries.
         * @param {GetEntryMappingsSortByEnum} [sortBy] The criteria to sort the entries by.
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async getEntryMappings(page: number, pageSize: number, filter?: GetEntryMappingsFilterEnum, sortBy?: GetEntryMappingsSortByEnum, options?: RawAxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<GetEntryMappings200Response>> {
            const localVarAxiosArgs = await localVarAxiosParamCreator.getEntryMappings(page, pageSize, filter, sortBy, options);
            const localVarOperationServerIndex = configuration?.serverIndex ?? 0;
            const localVarOperationServerBasePath = operationServerMap['DefaultApi.getEntryMappings']?.[localVarOperationServerIndex]?.url;
            return (axios, basePath) => createRequestFunction(localVarAxiosArgs, globalAxios, BASE_PATH, configuration)(axios, localVarOperationServerBasePath || basePath);
        },
        /**
         * 
         * @summary Get a list of retrievers.
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async getRetrievers(options?: RawAxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<GetRetrievers200Response>> {
            const localVarAxiosArgs = await localVarAxiosParamCreator.getRetrievers(options);
            const localVarOperationServerIndex = configuration?.serverIndex ?? 0;
            const localVarOperationServerBasePath = operationServerMap['DefaultApi.getRetrievers']?.[localVarOperationServerIndex]?.url;
            return (axios, basePath) => createRequestFunction(localVarAxiosArgs, globalAxios, BASE_PATH, configuration)(axios, localVarOperationServerBasePath || basePath);
        },
        /**
         * 
         * @summary Login to the application using the provided credentials.
         * @param {LoginRequestBody} loginRequestBody 
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async login(loginRequestBody: LoginRequestBody, options?: RawAxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<Login200Response>> {
            const localVarAxiosArgs = await localVarAxiosParamCreator.login(loginRequestBody, options);
            const localVarOperationServerIndex = configuration?.serverIndex ?? 0;
            const localVarOperationServerBasePath = operationServerMap['DefaultApi.login']?.[localVarOperationServerIndex]?.url;
            return (axios, basePath) => createRequestFunction(localVarAxiosArgs, globalAxios, BASE_PATH, configuration)(axios, localVarOperationServerBasePath || basePath);
        },
        /**
         * 
         * @summary Trigger a refresh of the entry mappings.
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async refreshEntryMappings(options?: RawAxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<RefreshEntryMappings200Response>> {
            const localVarAxiosArgs = await localVarAxiosParamCreator.refreshEntryMappings(options);
            const localVarOperationServerIndex = configuration?.serverIndex ?? 0;
            const localVarOperationServerBasePath = operationServerMap['DefaultApi.refreshEntryMappings']?.[localVarOperationServerIndex]?.url;
            return (axios, basePath) => createRequestFunction(localVarAxiosArgs, globalAxios, BASE_PATH, configuration)(axios, localVarOperationServerBasePath || basePath);
        },
    }
};

/**
 * DefaultApi - factory interface
 * @export
 */
export const DefaultApiFactory = function (configuration?: Configuration, basePath?: string, axios?: AxiosInstance) {
    const localVarFp = DefaultApiFp(configuration)
    return {
        /**
         * 
         * @summary Get a list of entry mappings.
         * @param {number} page The page number to display.
         * @param {number} pageSize The amount of items to display per each page.
         * @param {GetEntryMappingsFilterEnum} [filter] The filter to apply before returning the entries.
         * @param {GetEntryMappingsSortByEnum} [sortBy] The criteria to sort the entries by.
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getEntryMappings(page: number, pageSize: number, filter?: GetEntryMappingsFilterEnum, sortBy?: GetEntryMappingsSortByEnum, options?: RawAxiosRequestConfig): AxiosPromise<GetEntryMappings200Response> {
            return localVarFp.getEntryMappings(page, pageSize, filter, sortBy, options).then((request) => request(axios, basePath));
        },
        /**
         * 
         * @summary Get a list of retrievers.
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getRetrievers(options?: RawAxiosRequestConfig): AxiosPromise<GetRetrievers200Response> {
            return localVarFp.getRetrievers(options).then((request) => request(axios, basePath));
        },
        /**
         * 
         * @summary Login to the application using the provided credentials.
         * @param {LoginRequestBody} loginRequestBody 
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        login(loginRequestBody: LoginRequestBody, options?: RawAxiosRequestConfig): AxiosPromise<Login200Response> {
            return localVarFp.login(loginRequestBody, options).then((request) => request(axios, basePath));
        },
        /**
         * 
         * @summary Trigger a refresh of the entry mappings.
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        refreshEntryMappings(options?: RawAxiosRequestConfig): AxiosPromise<RefreshEntryMappings200Response> {
            return localVarFp.refreshEntryMappings(options).then((request) => request(axios, basePath));
        },
    };
};

/**
 * DefaultApi - object-oriented interface
 * @export
 * @class DefaultApi
 * @extends {BaseAPI}
 */
export class DefaultApi extends BaseAPI {
    /**
     * 
     * @summary Get a list of entry mappings.
     * @param {number} page The page number to display.
     * @param {number} pageSize The amount of items to display per each page.
     * @param {GetEntryMappingsFilterEnum} [filter] The filter to apply before returning the entries.
     * @param {GetEntryMappingsSortByEnum} [sortBy] The criteria to sort the entries by.
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof DefaultApi
     */
    public getEntryMappings(page: number, pageSize: number, filter?: GetEntryMappingsFilterEnum, sortBy?: GetEntryMappingsSortByEnum, options?: RawAxiosRequestConfig) {
        return DefaultApiFp(this.configuration).getEntryMappings(page, pageSize, filter, sortBy, options).then((request) => request(this.axios, this.basePath));
    }

    /**
     * 
     * @summary Get a list of retrievers.
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof DefaultApi
     */
    public getRetrievers(options?: RawAxiosRequestConfig) {
        return DefaultApiFp(this.configuration).getRetrievers(options).then((request) => request(this.axios, this.basePath));
    }

    /**
     * 
     * @summary Login to the application using the provided credentials.
     * @param {LoginRequestBody} loginRequestBody 
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof DefaultApi
     */
    public login(loginRequestBody: LoginRequestBody, options?: RawAxiosRequestConfig) {
        return DefaultApiFp(this.configuration).login(loginRequestBody, options).then((request) => request(this.axios, this.basePath));
    }

    /**
     * 
     * @summary Trigger a refresh of the entry mappings.
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof DefaultApi
     */
    public refreshEntryMappings(options?: RawAxiosRequestConfig) {
        return DefaultApiFp(this.configuration).refreshEntryMappings(options).then((request) => request(this.axios, this.basePath));
    }
}

/**
 * @export
 */
export const GetEntryMappingsFilterEnum = {
    IncompleteEntries: 'incomplete_entries',
    CompleteEntries: 'complete_entries'
} as const;
export type GetEntryMappingsFilterEnum = typeof GetEntryMappingsFilterEnum[keyof typeof GetEntryMappingsFilterEnum];
/**
 * @export
 */
export const GetEntryMappingsSortByEnum = {
    DateAddedAsc: 'date_added_asc',
    DateAddedDesc: 'date_added_desc',
    SizeAsc: 'size_asc',
    SizeDesc: 'size_desc'
} as const;
export type GetEntryMappingsSortByEnum = typeof GetEntryMappingsSortByEnum[keyof typeof GetEntryMappingsSortByEnum];


