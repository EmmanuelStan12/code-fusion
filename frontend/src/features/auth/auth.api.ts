import {API_BASE_URL} from "../../config/env.config.ts";
import LocalStorage, {AUTH_TOKEN_KEY} from "../../services/storage.service.ts";

export interface LoginDTO {
    email: string
    password: string
}

export interface RegisterDTO {
    email: string
    password: string
    firstName: string
    lastName: string
    username: string
}

export async function login(data: LoginDTO, { rejectWithValue }) {
    try {
        const response = await fetch(`${API_BASE_URL}/api/v1/login`, {
            method: 'POST',
            body: JSON.stringify(data),
        })
        if (!response.ok) {
            const error = await response.json()
            return rejectWithValue(error)
        }
        return response.json()
    } catch (e) {
        return rejectWithValue({
            message: e.message || 'Something went wrong'
        })
    }
}

export async function register(data: RegisterDTO, { rejectWithValue }) {
    try {
        const response = await fetch(`${API_BASE_URL}/api/v1/register`, {
            method: 'POST',
            body: JSON.stringify(data),
        })
        if (!response.ok) {
            const error = await response.json()
            return rejectWithValue(error)
        }
        const result = await response.json()

        return result
    } catch (e) {
        return rejectWithValue({
            message: e.message || 'Something went wrong'
        })
    }
}

export async function getAuthUser(_, { rejectWithValue }) {
    const token = LocalStorage.get(AUTH_TOKEN_KEY)
    try {
        const response = await fetch(`${API_BASE_URL}/api/v1/users/me`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
        })
        if (!response.ok) {
            const error = await response.json()
            return rejectWithValue(error)
        }
        const result = await response.json()

        return result
    } catch (e) {
        return rejectWithValue({
            message: e.message || 'Something went wrong'
        })
    }
}
