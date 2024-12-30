import {API_BASE_URL} from "../../config/env.config.ts";
import {LoginDTO} from "../auth/auth.api.ts";
import LocalStorage, {AUTH_TOKEN_KEY} from "../../services/storage.service.ts";

export interface CreateSessionDTO {
    language: string
    memoryLimit: number
    timeout: number
    title: string
}

export async function createSession(data: CreateSessionDTO, { rejectWithValue }) {
    const token = LocalStorage.get(AUTH_TOKEN_KEY)
    try {
        const response = await fetch(`${API_BASE_URL}/api/v1/sessions/create`, {
            method: 'POST',
            body: JSON.stringify(data),
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
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

export async function fetchCodeSession(sessionId: string, { rejectWithValue }) {
    const token = LocalStorage.get(AUTH_TOKEN_KEY)
    try {
        const response = await fetch(`${API_BASE_URL}/api/v1/sessions/${sessionId}`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
            },
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

