import axios from "axios";
import ApiClient from "./apiClient";

export interface IExtendedRegistration {
    id :number;
    name :string;
    surname :string;
    gender :string;
    date_of_birth: string;
    finished_school: string;
    attended_previous: boolean;
    attended_activities: string|null;

    city :string;
    pills :string|null;
    problems :string|null;
    notes: string;

    parent_name :string;
    parent_surname :string;
    email :string;
    phone :string;

    amount :number;
    payed :number|null;
    discount :number|null
    specific_symbol :string;

    eventID :number;
    title: string;
    days: IDay[];

    admin_note :string;
}

export interface IRegistration {
    id :number;
}

export interface IDay {
    id: number;
    description: string;
}

export class Registrations {
    constructor(protected client: ApiClient) {}

    public list () :Promise<IExtendedRegistration[]> {
        return this.client.get<IExtendedRegistration[]>("/api/registrations")
    }

    public update(reg :IExtendedRegistration) :Promise<number> {
        return new Promise<number>((resolve, reject) => {
            this.client.put(`/api/registrations/${reg.id}`, reg)
                .then(() => resolve(reg.id)).catch(err => reject(err))
        })
    }

    public delete(id :number) :Promise<IRegistration> {
        return this.client.delete<IRegistration>(`/api/registrations/${id}`)
    }
}

export const updateRegistration = (apiHost :string, token:string|null, reg :IExtendedRegistration) :Promise<number> => {
    return new Promise<number>((resolve, reject) => {
        axios.put(
            `${apiHost}/api/registrations/${reg.id}`,
            reg,
            {
                headers: {
                    Authorization: `Bearer ${token}`,
                }
            }
        ).then(() => resolve(reg.id)).catch(err => reject(err))
    })
}