export enum UserActionStatus {
    FETCH_USERS_IN_PROGRESS = 'FETCH_USERS_IN_PROGRESS',
    FETCH_USERS_FAILED = 'FETCH_USERS_FAILED',
    FETCH_USERS_SUCCESSFUL = 'FETCH_USERS_SUCCESSFUL',
}

enum UserActions {
    FETCH_USERS = 'users/fetch-users',
}

export default UserActions