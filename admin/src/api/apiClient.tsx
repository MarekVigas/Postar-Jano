import {AxiosInstance, AxiosRequestConfig} from 'axios'
import {Registrations} from "./registrations";
import {Events} from "./events";
import {User} from "./login";
import {getToken} from "./token";


class ApiClient {
    private readonly baseURL :string;
    private readonly axios :AxiosInstance;

    public events :Events;
    public registrations :Registrations;
    public user :User;

    constructor(axios: AxiosInstance, baseURL :string) {
        this.axios = axios
        this.baseURL = baseURL
        this.events = new Events(this)
        this.registrations = new Registrations(this)
        this.user = new User(this)
    }

    makeUrl(path: string) {
        return this.baseURL + path;
    }

    authHeader():AxiosRequestConfig {
        const token = getToken()
        if (!token)  return {}
        return {headers: {Authorization: 'Bearer ' + token}};
    }

    get<T = any>(relativeUrl: string): Promise<T> {
        const url = this.makeUrl(relativeUrl);
        const config = this.authHeader();

        return this.axios.get<T>(url, config)
            .then(res => res.data)
            .catch(err => {
                console.log(err);
                throw err;
            });
    }

    post<T = any>(relativeUrl: string, data: any): Promise<T> {
        const url = this.makeUrl(relativeUrl);
        const config = this.authHeader();

        return this.axios.post<T>(url,data, config)
            .then(res => res.data)
            .catch(err => {
                console.log(err);
                throw err;
            });
    }

    put<T = any>(relativeUrl: string, data: any): Promise<T> {
        const url = this.makeUrl(relativeUrl);
        const config = this.authHeader();

        return this.axios.put<T>(url,data, config)
            .then(res => res.data)
            .catch(err => {
                console.log(err);
                throw err;
            });
    }

    delete<T = any>(relativeUrl: string): Promise<T> {
        const url = this.makeUrl(relativeUrl);
        const config = this.authHeader();

        return this.axios.delete<T>(url, config)
            .then(res => res.data)
            .catch(err => {
                console.log(err);
                throw err;
            });
    }
}

export default ApiClient;