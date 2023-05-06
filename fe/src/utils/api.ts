import axios from "axios";
type FetcherProps = {
  url: string,
  args: any,
}

export const SWRFetcher = async ({ url, args }: FetcherProps) => {
  return axios(url, {
    params: args
  })
  .then((r) => r.data)
};

// register(id: number, registration: Registration): Promise<RegistrationRespone> {
//   return this.client.post<RegistrationRespone>(`/registrations/${id}`, registration);
// }

// stats(id: number): Promise<Stat[]> {
//   return this.client.get<Stat[]>(`/stats/${id}`);
// }

// get(id: number): Promise<IEvent> {
//   return this.client.get<IEvent>(`/events/${id}`);
// }