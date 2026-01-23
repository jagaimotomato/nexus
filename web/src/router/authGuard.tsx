import React from 'react'
import { Navigate, useLocation } from 'react-router-dom'
import useUserStore from '@/store/useUserStore'

interface AuthGuardProps {
  children: React.ReactNode
}

const AuthGuard: React.FC<AuthGuardProps> = ({ children }) => {
  const token = useUserStore((state) => state.token)
  const location = useLocation()
  if (!token) {
    return <Navigate to="/login" state={{ from: location }} replace />
  }
  return <>{children}</>
}

export default AuthGuard