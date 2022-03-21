import ApiClient from "./apiClient";

export interface IEvent {
    id :number;
    title :string;
    date_from :string;
    date_to :string;
    description :string;
    location :string;
    owner_id :number;
    days :IDay[];
}

export interface IDay {
    id :number;
    description :string;
    price :number;
    capacity: number;
    limit_boys: number|null;
    limit_girls: number|null;
}

export interface IStat {
    day_id: number;
    event_id: number;
    capacity: number;
    limit_boys: number|null;
    limit_girls: number|null;
    boys_count: number;
    girls_count: number;
}

export class Events {
    constructor(protected client: ApiClient) {}

    public list () :Promise<IEvent[]> {
        return this.client.get<IEvent[]>("/api/events")
    }

    public stats (): Promise<IStat[]> {
        return this.client.get("/api/stats")
    }
}
