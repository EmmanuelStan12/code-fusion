import {API_BASE_URL} from "../../config/env.config.ts";
import {LoginDTO} from "../auth/auth.api.ts";
import LocalStorage, {AUTH_TOKEN_KEY} from "../../services/storage.service.ts";

export async function fetchDashboardAnalytics(_, { rejectWithValue }) {
    const token = LocalStorage.get(AUTH_TOKEN_KEY)
    try {
        const response = await fetch(`${API_BASE_URL}/api/v1/analytics`, {
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

