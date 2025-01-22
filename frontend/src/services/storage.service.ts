export const AUTH_TOKEN_KEY = 'AUTH_TOKEN'

export default class LocalStorage {
    static get<T>(key: string): T {
        const data = localStorage.getItem(key)
        try {
            return JSON.parse(data) as T
        } catch (e) {
            return data as T
        }
    }

    static set<T>(key: string, value: T): T {
        try {
            if (typeof value === 'string') {
                localStorage.setItem(key, value)
                return
            }
            const data = JSON.stringify(value)
            localStorage.setItem(key, data)
        } catch (e) {
            localStorage.setItem(key, value)
        }
    }

    static del(key: string) {
        localStorage.removeItem(key)
    }

    static clear() {
        localStorage.clear()
    }
}