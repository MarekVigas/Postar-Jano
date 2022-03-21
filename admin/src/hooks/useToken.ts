import React, {useEffect, useState} from 'react'

const TOKEN_STORAGE_KEY='leto.sbb.token'

export default function useToken():[string|null, React.Dispatch<React.SetStateAction<string|null>>] {
    const [token, setToken] = useState<string|null>(() => {
        return sessionStorage.getItem(TOKEN_STORAGE_KEY)
    })

    useEffect(() => {
        if (token != null) {
            sessionStorage.setItem(TOKEN_STORAGE_KEY, token)
        } else {
            sessionStorage.removeItem(TOKEN_STORAGE_KEY)
        }

    },[token])
    return [token, setToken]
}