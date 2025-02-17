import {initializeApiClient} from "../utils/api.ts";

const localStorageKey = "authToken";

export function isAuthenticated(): boolean {
    const jwtStr = localStorage.getItem(localStorageKey);
    if (!jwtStr) {
        return false;
    }
    return isTokenStillValidClaimsBased(jwtStr);
}

export function initializeAuthToken(jwtStr: string): void {
    localStorage.setItem(localStorageKey, jwtStr);
    initializeApiClient(jwtStr);
}

export function isTokenStillValidClaimsBased(jwtStr: string): boolean {
    const jwtStrSplit = jwtStr.split(".");
    if (jwtStrSplit.length !== 3) {
        return false;
    }
    try {
        const claims = JSON.parse(atob(jwtStrSplit[1]));
        return claims.exp > (Date.now() / 1000)
    } catch (SyntaxError) {
        return false;
    }
}