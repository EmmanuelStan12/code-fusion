import {createAsyncThunk, createSlice} from '@reduxjs/toolkit'
import AuthActions, {AuthActionStatus} from "./auth.actions.ts";
import * as authApi from './auth.api.ts'
import LocalStorage, {AUTH_TOKEN_KEY} from "../../services/storage.service.ts";

export interface AuthState {
    status?: AuthActionStatus
    data?: { user: object }
    message?: ''
    currentAction?: AuthActions
}

const initialState: AuthState = {}

const AUTH_SLICE_NAME = 'AUTH'

export const login = createAsyncThunk(
    AuthActions.LOGIN,
    authApi.login
)

export const register = createAsyncThunk(
    AuthActions.REGISTER,
    authApi.register,
)

export const getAuthUser = createAsyncThunk(
    AuthActions.GET_AUTH_USER,
    authApi.getAuthUser,
)

const AuthSlice = createSlice({
    name: AUTH_SLICE_NAME,
    initialState,
    reducers: {} as any,
    extraReducers: (builder) => {
        builder.addCase(login.pending, (state: AuthState, action) => {
            state.status = AuthActionStatus.LOGIN_IN_PROGRESS
            state.currentAction = AuthActions.LOGIN
            state.message = ''
        })

        builder.addCase(login.fulfilled, (state: AuthState, action) => {
            state.status = AuthActionStatus.LOGIN_SUCCESSFUL
            state.currentAction = AuthActions.LOGIN
            const payload = action.payload as object

            const { token, user } = payload.data
            LocalStorage.set(AUTH_TOKEN_KEY, token)
            state.message = payload.message
            state.data = { user }
        })

        builder.addCase(login.rejected, (state: AuthState, action) => {
            state.status = AuthActionStatus.LOGIN_FAILED
            state.currentAction = AuthActions.LOGIN
            state.message = (action.payload as any)?.message
        })

        builder.addCase(register.pending, (state: AuthState, action) => {
            state.status = AuthActionStatus.REGISTER_IN_PROGRESS
            state.currentAction = AuthActions.REGISTER
            state.message = ''
        })

        builder.addCase(register.fulfilled, (state: AuthState, action) => {
            state.status = AuthActionStatus.REGISTER_SUCCESSFUL
            state.currentAction = AuthActions.REGISTER

            const payload = action.payload as object
            const { token, user } = payload.data
            LocalStorage.set(AUTH_TOKEN_KEY, token)
            state.data = { user }
            state.message = payload.message
        })

        builder.addCase(register.rejected, (state: AuthState, action) => {
            state.status = AuthActionStatus.REGISTER_FAILED
            state.currentAction = AuthActions.REGISTER
            state.message = (action.payload as any)?.message
        })

        builder.addCase(getAuthUser.pending, (state: AuthState, action) => {
            state.status = AuthActionStatus.AUTHENTICATION_IN_PROGRESS
            state.currentAction = AuthActions.GET_AUTH_USER
            state.message = ''
        })

        builder.addCase(getAuthUser.fulfilled, (state: AuthState, action) => {
            state.status = AuthActionStatus.AUTHENTICATION_SUCCESSFUL
            state.currentAction = AuthActions.GET_AUTH_USER
            const payload = action.payload as object

            state.data = { user: payload.data }
            state.message = payload.message
        })

        builder.addCase(getAuthUser.rejected, (state: AuthState, action) => {
            state.status = AuthActionStatus.AUTHENTICATION_FAILED
            state.currentAction = AuthActions.GET_AUTH_USER
            state.message = (action.payload as any)?.message
        })
    }
})

export default AuthSlice.reducer