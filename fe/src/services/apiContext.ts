import * as React from "react";
import {ApiClient} from "./apiClient";

export const ApiClientContext = React.createContext<ApiClient|null>(null);