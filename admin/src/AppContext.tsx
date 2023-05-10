import React, { useContext, useReducer} from "react";
import ApiClient from "./api/apiClient";
import axios from "axios";
import {getToken, setToken} from "./api/token";

interface AppData {
    apiClient :ApiClient;
    state :AppState
    dispatch : (action :Action) => void;
}

const AppContext = React.createContext<AppData|undefined>(undefined)

const defaultApiHost = "https://leto-api.salezko.sk"

interface AppProps {
    host :string|undefined
    children :React.ReactNode
}

type AppState = {
    token :string|null;
}

type Action =
    | {
        type:'SET_TOKEN'; token:string
    }
    | {
        type:'UNSET_TOKEN';
    }
    | {
        type:'FETCH_REGISTRATION'
    }


function appReduce(state :AppState, action :Action) :AppState {
    switch (action.type) {
        case "SET_TOKEN":
            setToken(action.token)
            return {...state, token: action.token}
        case "FETCH_REGISTRATION":
            return {...state}
        case "UNSET_TOKEN":
            return {...state, token:null}
    }
}


function initState() :AppState {
    return {
        token: getToken(),
    }
}

function AppProvider({host, children} :AppProps) {
    const apiHost = host || defaultApiHost
    const apiClient = new ApiClient(axios.create(), apiHost)
    const [state, dispatch] = useReducer(appReduce, initState())

    return (
        <AppContext.Provider value={{apiClient, state, dispatch}}>
            {children}
        </AppContext.Provider>
    )
}

export function useSafeContext():AppData {
    const ctx = useContext(AppContext)
    if (!ctx) throw new Error("Context not initialized")
    return ctx
}

export function useAppState():AppState {
    const {state} = useSafeContext()
    return state
}

export function useAppDispatch() :(a :Action) => void {
    const {dispatch} = useSafeContext()
    return dispatch
}

export function useAPIClient() :ApiClient {
    const {apiClient} = useSafeContext()
    return apiClient
}

export {AppProvider, AppContext}