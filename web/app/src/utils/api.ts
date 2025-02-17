import {Configuration, DefaultApi} from "../api";

let apiClient: DefaultApi | null = null;

export function initializeApiClient(jwt: string) {
    apiClient = new DefaultApi(new Configuration({
        apiKey: jwt,
    }), "/api");
}

export function getApiClient(): DefaultApi {
    if (apiClient === null) {
        throw new Error('API Client is not initialized yet');
    }
    return apiClient;
}
