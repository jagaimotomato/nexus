import { type JSX } from "react";
import { Navigate, useLocation } from "react-router-dom";
import useUserStore from "@/store/useUserStore";

interface AuthGuardProps {
  children: JSX.Element;
}

const AuthGuard = ({ children }: AuthGuardProps) => {
  const token = useUserStore((state) => state.token);
  const location = useLocation();
  if (!token) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }
  return children;
};
export default AuthGuard;
