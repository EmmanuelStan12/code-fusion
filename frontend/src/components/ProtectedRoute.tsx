import {useLocation, Outlet, Navigate} from "react-router";
import {useAppSelector} from "../redux/hooks.ts";
import {useSelector} from "react-redux";
import {useEffect, useMemo, useState} from "react";

export default function ProtectedRoute() {
    const state = useAppSelector((state) => state.auth);
    const location = useLocation()

    if (!state.data?.user) {
        return <Navigate to="/login" state={{ from: location }} />;
    }
    return <Outlet />;
};