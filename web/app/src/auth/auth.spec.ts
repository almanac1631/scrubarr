import {describe, expect, test} from "vitest";
import {isTokenStillValidClaimsBased} from "./auth.ts";

describe("check token validity using claims", () => {
    test("checks valid token claims", async () => {
        const payload = {
            iat: Math.floor(Date.now() / 1000),
            exp: Math.floor(Date.now() / 1000) + 3600,
        };
        const encodedPayload = btoa(JSON.stringify(payload));
        let jwtStr = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.${encodedPayload}.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c`;//gitleaks:allow
        expect(isTokenStillValidClaimsBased(jwtStr)).toBe(true);
    });
    test("checks expired token claims", async () => {
        let jwtStr = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3Mzk3MTI0MzMsImV4cCI6MTczOTc1MjQzM30.1LOyxfn7rwwFyrKeVWnAVSt9GOohlZvMjUpykXGpzpI`; //gitleaks:allow
        expect(isTokenStillValidClaimsBased(jwtStr)).toBe(false);
    });
    test("checks malformed token", async () => {
        let jwtStr = `someinvalidtoken`;
        expect(isTokenStillValidClaimsBased(jwtStr)).toBe(false);
    });
    test("checks malformed jwt token", async () => {
        let jwtStr = `some.invalid.token`;
        expect(isTokenStillValidClaimsBased(jwtStr)).toBe(false);
    });
    test("checks malformed jwt token with a valid json", async () => {
        const encodedPayload = btoa(JSON.stringify({}));
        let jwtStr = `some.${encodedPayload}.token`;
        expect(isTokenStillValidClaimsBased(jwtStr)).toBe(false);
    });
})
