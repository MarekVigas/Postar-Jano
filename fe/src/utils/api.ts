import axios from "axios";
import { AxiosRequestConfig } from 'axios'
type FetcherProps = {
  url: string,
  config: AxiosRequestConfig,
}

export const SWRFetcher = async ({ url, config }: FetcherProps) => {
  return axios(url, config)
  .then((r) => r.data)
};
