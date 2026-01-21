import React from 'react'
import { Navigate, useLocation } from 'react-router'
import useUserStore from '@/store/useUserStore'

interface AuthGuardProps {
  children: React.ReactNode
}

const AuthGuard: React.FC<AuthGuardProps> = ({ children }) => {
  const token = useUserStore((state: any) => state.token)
  const location = useLocation()
  if (!token) {
    return <Navigate to="/login" state={{ from: location }} replace />
  }
  return <>{children}</>
}

export default AuthGuard