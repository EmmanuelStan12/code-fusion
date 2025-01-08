import { configureStore } from '@reduxjs/toolkit'
import authSlice from "../features/auth/auth.slice.ts";
import sessionSlice from "../features/code-session/session.slice.ts";
import dashboardSlice from "../features/dashboard/dashboard.slice.ts";
import usersSlice from "../features/users/users.slice.ts";

export const store = configureStore({
    reducer: {
        auth: authSlice,
        codeSession: sessionSlice,
        dashboard: dashboardSlice,
        users: usersSlice,
    },
})
export type RootState = ReturnType<typeof store.getState>

export type AppDispatch = typeof store.dispatch
