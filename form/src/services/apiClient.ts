import { AxiosInstance } from "axios";
import { IEvent, Registration, RegistrationRespone, Stat } from "../types/types"

export class ApiClient {

    protected token: string = "";
    public owner: Owner;
    public event: Event;

    constructor(public axios: AxiosInstance, protected baseUrl: string) {
        this.token = "";
        this.owner = new Owner(this);
        this.event = new Event(this);
    }

    setToken(token: string) {
        this.token = token;
    }

    getToken() {
        return this.token;
    }

    makeUrl(url: string) {
        return `${this.baseUrl}/api${url}`;
    }

    getRequestConfig() {
        if(!this.token) {
            return;
        }

        return {
            headers: {
                Authorization: "Bearer " + this.token
            }
        }
    }

    get<T = any>(relativeUrl: string): Promise<T> {
        const url = this.makeUrl(relativeUrl);
        const config = this.getRequestConfig();

        return this.axios.get<T>(url, config)
            .then(res => res.data)
            .catch(err => {
                console.log(err);
                throw err;
            });
    }

    post<T = any>(relativeUrl: string, data: any): Promise<T> {
        const url = this.makeUrl(relativeUrl);
        const config = this.getRequestConfig();

        return this.axios.post<T>(url,data, config)
            .then(res => res.data)
            .catch(err => {
                console.log(err);
                throw err;
            });
    }

    put<T = any>(relativeUrl: string, data: any): Promise<T> {
        const url = this.makeUrl(relativeUrl);
        const config = this.getRequestConfig();

        return this.axios.put<T>(url,data, config)
            .then(res => res.data)
            .catch(err => {
                console.log(err);
                throw err;
            });
    }
}

interface LoginData {
    userName: string,
    pass: string
}

class Owner {
    constructor(protected client: ApiClient) {}

    logIn(data: LoginData): Promise<{token: string}> {
        return this.client.post<{token: string}>("/signIn/", data);
    }
}

class Event {
    constructor(protected client: ApiClient) {}

    register(id: number, registration: Registration): Promise<RegistrationRespone> {
        return this.client.post<RegistrationRespone>(`/registrations/${id}`, registration);
    }
    
    stats(id: number): Promise<Stat[]> {
        return this.client.get<Stat[]>(`/stats/${id}`);
    }

    get(id: number): Promise<IEvent> {
        return this.client.get<IEvent>(`/events/${id}`);
    }
}
