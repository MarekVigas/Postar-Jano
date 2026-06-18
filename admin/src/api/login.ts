import ApiClient from "./apiClient";

export interface IMeResponse {
    email: string;
}

export class User {
    constructor(protected client: ApiClient) {}

    public signIn(username: string, password: string): Promise<void> {
        return this.client.post<void>("/api/sign/in", {username, password})
    }

    public signOut(): Promise<void> {
        return this.client.post<void>("/api/sign/out", {})
    }

    public me(): Promise<IMeResponse | null> {
        return this.client.get<IMeResponse>("/api/me").catch(() => null)
    }
}
