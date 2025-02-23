import {Configuration, DefaultApi} from "../api";

let apiClient: DefaultApi | null = null;

export const basePath = "api";

export function initializeApiClient(jwt: string) {
    apiClient = new DefaultApi(new Configuration({
        accessToken: jwt,
    }), basePath);
}

export function getApiClient(): DefaultApi {
    if (apiClient === null) {
        throw new Error('API Client is not initialized yet');
    }
    return apiClient;
}
