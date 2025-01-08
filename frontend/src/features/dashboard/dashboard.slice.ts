import DashboardActions, {DashboardActionStatus} from "./dashboard.actions.ts";
import {createAsyncThunk, createSlice} from "@reduxjs/toolkit";
import * as dashboardApi from "./dashboard.api.ts";

export interface AnalyticsDTO {
    totalHours: number
    totalSessions: number
    totalLanguagesUsed: number
    totalCollaborators: number
}

export interface DashboardDTO {
    recentCollaborators: any[]
    recentSessions: any[]
    analytics: AnalyticsDTO
}

export interface DashboardState {
    status?: DashboardActionStatus
    data?: DashboardDTO
    message?: ''
    currentAction?: DashboardActions
}

const initialState: DashboardState = {}

const DASHBOARD_SLICE_NAME = 'DASHBOARD'

export const fetchDashboard = createAsyncThunk(
    DashboardActions.FETCH_DASHBOARD,
    dashboardApi.fetchDashboardAnalytics,
)


const DashboardsSlice = createSlice({
    name: DASHBOARD_SLICE_NAME,
    initialState,
    reducers: {} as any,
    extraReducers: (builder) => {
        builder.addCase(fetchDashboard.pending, (state: DashboardState, action) => {
            state.status = DashboardActionStatus.FETCH_DASHBOARD_IN_PROGRESS
            state.currentAction = DashboardActions.FETCH_DASHBOARD
            state.message = ''
        })

        builder.addCase(fetchDashboard.fulfilled, (state: DashboardState, action) => {
            state.status = DashboardActionStatus.FETCH_DASHBOARD_SUCCESSFUL
            state.currentAction = DashboardActions.FETCH_DASHBOARD
            const payload = action.payload as object

            state.message = payload.message
            state.data = payload.data
        })

        builder.addCase(fetchDashboard.rejected, (state: DashboardState, action) => {
            state.status = DashboardActionStatus.FETCH_DASHBOARD_FAILED
            state.currentAction = DashboardActions.FETCH_DASHBOARD
            state.message = (action.payload as any)?.message
        })
    }
})

export default DashboardsSlice.reducer