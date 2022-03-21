const TOKEN_STORAGE_KEY='leto.sbb.token'

export function getToken() {
   return sessionStorage.getItem(TOKEN_STORAGE_KEY)
}

export function setToken(token :string) {
    return sessionStorage.setItem(TOKEN_STORAGE_KEY, token)
}