import React, {useEffect, useState} from 'react'

export default function useStorage<T>(key :string, init:T):[T, React.Dispatch<React.SetStateAction<T>>] {
    const [val, setVal] = useState<T>(() => {
        const stored = sessionStorage.getItem(key)
        if (stored) {
            return JSON.parse(stored)
        }
        return init
    })

    useEffect(() => {
        if (val != null) {
            sessionStorage.setItem(key, JSON.stringify(val))
        } else {
            sessionStorage.removeItem(key)
        }

    },[key, val])
    return [val, setVal]
}