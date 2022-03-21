import ApiClient from "./apiClient";

export interface ISignInResponse {
    token :string;
}

export class User {
    constructor(protected client: ApiClient) {}

    public signIn( username: string, password :string) :Promise<ISignInResponse>{
        return this.client.post<ISignInResponse>("/api/sign/in", {username, password})
    }
}