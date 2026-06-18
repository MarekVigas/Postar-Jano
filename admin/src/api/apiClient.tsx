import {Events} from "./events";
import {Registrations} from "./registrations";
import {User} from "./login";

class ApiClient {
    private readonly baseURL: string;

    public events: Events;
    public registrations: Registrations;
    public user: User;

    constructor(baseURL: string) {
        this.baseURL = baseURL;
        this.events = new Events(this);
        this.registrations = new Registrations(this);
        this.user = new User(this);
    }

    makeUrl(path: string): string {
        return this.baseURL + path;
    }

    private async request<T>(method: string, relativeUrl: string, body?: unknown): Promise<T> {
        const headers: Record<string, string> = {
            ...(body !== undefined ? {"Content-Type": "application/json"} : {}),
        };

        const res = await fetch(this.makeUrl(relativeUrl), {
            method,
            headers,
            credentials: 'include',
            body: body !== undefined ? JSON.stringify(body) : undefined,
        });

        if (!res.ok) {
            const err = new Error(`HTTP ${res.status}: ${res.statusText}`);
            console.log(err);
            throw err;
        }

        const text = await res.text();
        return (text ? JSON.parse(text) : undefined) as T;
    }

    get<T = unknown>(relativeUrl: string): Promise<T> {
        return this.request<T>("GET", relativeUrl);
    }

    post<T = unknown>(relativeUrl: string, data: unknown): Promise<T> {
        return this.request<T>("POST", relativeUrl, data);
    }

    put<T = unknown>(relativeUrl: string, data: unknown): Promise<T> {
        return this.request<T>("PUT", relativeUrl, data);
    }

    delete<T = unknown>(relativeUrl: string): Promise<T> {
        return this.request<T>("DELETE", relativeUrl);
    }
}

export default ApiClient;
