export enum DashboardActionStatus {
    FETCH_DASHBOARD_IN_PROGRESS = 'FETCH_DASHBOARD_IN_PROGRESS',
    FETCH_DASHBOARD_FAILED = 'FETCH_DASHBOARD_FAILED',
    FETCH_DASHBOARD_SUCCESSFUL = 'FETCH_DASHBOARD_SUCCESSFUL',
}

enum DashboardActions {
    FETCH_DASHBOARD = 'dashboard/fetch-dashboard',
}

export default DashboardActions