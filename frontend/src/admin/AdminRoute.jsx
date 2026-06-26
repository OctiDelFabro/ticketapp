import { Navigate } from 'react-router-dom'
import Button from '../components/Button.jsx'
import { getStoredUser, isAuthenticated } from '../utils/auth.js'
import { isAdminUser } from '../utils/admin.js'

export default function AdminRoute({ children }) {
  if (!isAuthenticated()) return <Navigate to="/login" replace state={{ message: 'Iniciá sesión con una cuenta administradora para continuar.' }} />
  if (!isAdminUser(getStoredUser())) {
    return <main className="container-page flex min-h-[70vh] items-center justify-center"><div className="glass-card max-w-lg rounded-3xl p-8 text-center"><p className="text-sm font-black uppercase tracking-[0.3em] text-red-300">Acceso denegado</p><h1 className="mt-3 text-3xl font-black">Esta sección es solo para administradores.</h1><p className="mt-3 text-gray-400">Tu sesión no tiene permisos ADMIN. No se otorgan permisos si el token o usuario guardado no informa un rol administrador.</p><Button to="/" className="mt-6">Volver al home</Button></div></main>
  }
  return children
}
