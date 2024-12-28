import { configureStore } from '@reduxjs/toolkit'
import authSlice from "../features/auth/auth.slice.ts";
import authMiddleware from "../features/auth/auth.middleware.ts";

export const store = configureStore({
    reducer: {
        auth: authSlice,
    },
})
export type RootState = ReturnType<typeof store.getState>

export type AppDispatch = typeof store.dispatch
