import {createAsyncThunk, createSlice} from '@reduxjs/toolkit'
import CodeSessionActions, {CodeSessionActionStatus} from "./session.actions.ts";
import * as sessionApi from './session.api.ts'

export interface CodeSessionState {
    status?: CodeSessionActionStatus
    data?: { currentSession: object }
    message?: ''
    currentAction?: CodeSessionActions
}

const initialState: CodeSessionState = {}

const CODE_SESSION_SLICE_NAME = 'CODE_SESSION'

export const createCodeSession = createAsyncThunk(
    CodeSessionActions.CREATE_SESSION,
    sessionApi.createSession
)

export const fetchCodeSession = createAsyncThunk(
    CodeSessionActions.FETCH_CODE_SESSION,
    sessionApi.fetchCodeSession,
)

const CodeSessionSlice = createSlice({
    name: CODE_SESSION_SLICE_NAME,
    initialState,
    reducers: {} as any,
    extraReducers: (builder) => {
        builder.addCase(createCodeSession.pending, (state: CodeSessionState, action) => {
            state.status = CodeSessionActionStatus.CREATE_SESSION_IN_PROGRESS
            state.currentAction = CodeSessionActions.CREATE_SESSION
            state.message = ''
        })

        builder.addCase(createCodeSession.fulfilled, (state: CodeSessionState, action) => {
            state.status = CodeSessionActionStatus.CREATE_SESSION_SUCCESSFUL
            state.currentAction = CodeSessionActions.CREATE_SESSION
            const payload = action.payload as object

            state.message = payload.message
            state.data = { currentSession: payload.data }
        })

        builder.addCase(createCodeSession.rejected, (state: CodeSessionState, action) => {
            state.status = CodeSessionActionStatus.CREATE_SESSION_FAILED
            state.currentAction = CodeSessionActions.CREATE_SESSION
            state.message = (action.payload as any)?.message
        })

        builder.addCase(fetchCodeSession.pending, (state: CodeSessionState, action) => {
            state.status = CodeSessionActionStatus.FETCH_CODE_SESSION_IN_PROGRESS
            state.currentAction = CodeSessionActions.FETCH_CODE_SESSION
            state.message = ''
        })

        builder.addCase(fetchCodeSession.fulfilled, (state: CodeSessionState, action) => {
            state.status = CodeSessionActionStatus.FETCH_CODE_SESSION_SUCCESSFUL
            state.currentAction = CodeSessionActions.FETCH_CODE_SESSION
            const payload = action.payload as object

            state.message = payload.message
            state.data = { currentSession: payload.data }
        })

        builder.addCase(fetchCodeSession.rejected, (state: CodeSessionState, action) => {
            state.status = CodeSessionActionStatus.FETCH_CODE_SESSION_FAILED
            state.currentAction = CodeSessionActions.FETCH_CODE_SESSION
            state.message = (action.payload as any)?.message
        })
    }
})

export default CodeSessionSlice.reducer