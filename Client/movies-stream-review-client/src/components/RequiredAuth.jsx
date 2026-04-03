import { useLocation, Navigate, Outlet } from "react-router-dom";
import useAuth from "../hooks/useAuth";

const RequiredAuth = () => {
    const {auth} = useAuth();
    const location = useLocation();

    return auth ? (
        <Outlet />
    ) : (
        <Navigate to='/users/login' state={{from:location}} replace />
    )
}
export default RequiredAuth;