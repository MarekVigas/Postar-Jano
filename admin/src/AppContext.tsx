import React, {useContext, useReducer} from "react";
import ApiClient from "./api/apiClient";

interface AppData {
    apiClient: ApiClient;
    state: AppState;
    dispatch: (action: Action) => void;
}

const AppContext = React.createContext<AppData | undefined>(undefined)

const defaultApiHost = ""

interface AppProps {
    host: string | undefined;
    children: React.ReactNode;
}

type AppState = {
    isAuthenticated: boolean;
    isLoading: boolean;
}

type Action =
    | { type: 'SIGN_IN' }
    | { type: 'SIGN_OUT' }
    | { type: 'INIT_DONE'; isAuthenticated: boolean }

function appReduce(state: AppState, action: Action): AppState {
    switch (action.type) {
        case "SIGN_IN":
            return {...state, isAuthenticated: true}
        case "SIGN_OUT":
            return {...state, isAuthenticated: false}
        case "INIT_DONE":
            return {isAuthenticated: action.isAuthenticated, isLoading: false}
    }
}

function initState(): AppState {
    return {isAuthenticated: false, isLoading: true}
}

function AppProvider({host, children}: AppProps) {
    const apiHost = host || defaultApiHost
    const apiClient = new ApiClient(apiHost)
    const [state, dispatch] = useReducer(appReduce, initState())

    return (
        <AppContext.Provider value={{apiClient, state, dispatch}}>
            {children}
        </AppContext.Provider>
    )
}

export function useSafeContext(): AppData {
    const ctx = useContext(AppContext)
    if (!ctx) throw new Error("Context not initialized")
    return ctx
}

export function useAppState(): AppState {
    const {state} = useSafeContext()
    return state
}

export function useAppDispatch(): (a: Action) => void {
    const {dispatch} = useSafeContext()
    return dispatch
}

export function useAPIClient(): ApiClient {
    const {apiClient} = useSafeContext()
    return apiClient
}

export {AppProvider, AppContext}
