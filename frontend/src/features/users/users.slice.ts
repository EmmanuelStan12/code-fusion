import UserActions, {UserActionStatus} from "./users.actions.ts";
import {createAsyncThunk, createSlice} from "@reduxjs/toolkit";
import * as usersApi from "./users.api.ts";

export interface IUser {
    id: number
    firstName: string
    lastName: string
    email: string
    username: string
}

export interface UserState {
    status?: UserActionStatus
    data?: { users: IUser[] }
    message?: ''
    currentAction?: UserActions
}

const initialState: UserState = {}

const USERS_SLICE_NAME = 'USERS'

export const fetchUsers = createAsyncThunk(
    UserActions.FETCH_USERS,
    usersApi.fetchUsers,
)


const UsersSlice = createSlice({
    name: USERS_SLICE_NAME,
    initialState,
    reducers: {} as any,
    extraReducers: (builder) => {
        builder.addCase(fetchUsers.pending, (state: UserState, action) => {
            state.status = UserActionStatus.FETCH_USERS_IN_PROGRESS
            state.currentAction = UserActions.FETCH_USERS
            state.message = ''
        })

        builder.addCase(fetchUsers.fulfilled, (state: UserState, action) => {
            state.status = UserActionStatus.FETCH_USERS_SUCCESSFUL
            state.currentAction = UserActions.FETCH_USERS
            const payload = action.payload as object

            state.message = payload.message
            state.data = { users: payload.data }
        })

        builder.addCase(fetchUsers.rejected, (state: UserState, action) => {
            state.status = UserActionStatus.FETCH_USERS_FAILED
            state.currentAction = UserActions.FETCH_USERS
            state.message = (action.payload as any)?.message
        })
    }
})

export default UsersSlice.reducer