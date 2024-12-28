import LocalStorage, {AUTH_TOKEN_KEY} from "../../services/storage.service.ts";

const authMiddleware = (store) => (next) => (action) => {
    const token = LocalStorage.get<string>(AUTH_TOKEN_KEY);
    if (token) {
        action.meta = {
            ...action.meta,
            token,
        };
    }
    return next(action)
};

export default authMiddleware;