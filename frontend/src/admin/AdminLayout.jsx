import { NavLink, Outlet, useNavigate } from 'react-router-dom'
import Button from '../components/Button.jsx'
import Logo from '../components/Logo.jsx'
import { getUserInitials } from '../utils/admin.js'
import { clearAuthSession, getStoredUser } from '../utils/auth.js'

const navLink = ({ isActive }) => `flex items-center gap-3 rounded-2xl px-4 py-3 font-bold transition ${isActive ? 'bg-ticket-purple text-white shadow-glow' : 'text-gray-300 hover:bg-white/5 hover:text-white'}`

export default function AdminLayout() {
  const user = getStoredUser()
  const navigate = useNavigate()
  const logout = () => { clearAuthSession(); navigate('/login') }
  return <div className="min-h-screen bg-ticket-alt text-white lg:flex"><aside className="border-b border-ticket-border bg-ticket-card/80 p-5 lg:fixed lg:inset-y-0 lg:w-72 lg:border-b-0 lg:border-r"><Logo /><p className="mt-8 text-xs font-black uppercase tracking-[0.28em] text-ticket-purple2">Gestión</p><nav className="mt-4 grid gap-2"><NavLink className={navLink} to="/admin/eventos">🎫 Eventos</NavLink><NavLink className={navLink} to="/admin/reportes">📊 Reportes</NavLink></nav></aside><div className="flex min-h-screen flex-1 flex-col lg:pl-72"><header className="sticky top-0 z-20 border-b border-ticket-border bg-ticket-bg/90 backdrop-blur"><div className="flex flex-col gap-4 px-5 py-4 sm:flex-row sm:items-center sm:justify-between lg:px-8"><div><p className="text-xs font-black uppercase tracking-[0.32em] text-ticket-purple2">Modo administrador</p><h1 className="text-xl font-black">Panel Ticketar</h1></div><div className="flex items-center gap-3"><div className="flex h-11 w-11 items-center justify-center rounded-full border border-ticket-purple2/50 bg-ticket-purple/25 font-black">{getUserInitials(user)}</div><div className="hidden sm:block"><p className="text-sm font-bold">{user?.name || 'Administrador'}</p><p className="text-xs text-gray-400">{user?.email || 'Sin email'}</p></div><Button variant="secondary" className="px-4 py-2" onClick={logout}>Salir</Button></div></div></header><Outlet /></div></div>
}
