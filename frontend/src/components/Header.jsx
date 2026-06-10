import { Link, useLocation } from 'react-router-dom'
import Button from './Button.jsx'
import Logo from './Logo.jsx'

export default function Header({ cartCount = 0, isLoggedIn = false, searchQuery, setSearchQuery, onLogout }) {
  const { pathname } = useLocation()
  const showSearch = pathname === '/'

  return (
    <header className="sticky top-0 z-40 border-b border-ticket-border/70 bg-ticket-bg/80 backdrop-blur-xl">
      <div className="container-page flex items-center gap-4 py-4">
        <Logo />
        <div className="hidden flex-1 md:block">
          {showSearch && (
            <input
              className="input-dark"
              onChange={(event) => setSearchQuery(event.target.value)}
              placeholder="Buscar artistas, eventos, venues..."
              value={searchQuery}
            />
          )}
        </div>
        <div className="hidden items-center gap-3 sm:flex">
          <Button to={isLoggedIn ? '/mis-entradas' : '/login'}>
            {isLoggedIn ? 'Mis Entradas' : 'Ingresar'}
          </Button>
          {isLoggedIn && <Button onClick={onLogout} type="button" variant="secondary">Cerrar sesión</Button>}
        </div>
        <Link to="/checkout" className="relative grid h-12 w-12 place-items-center rounded-2xl border border-ticket-border bg-ticket-card2 text-xl transition hover:border-ticket-purple2/60" aria-label="Ir al carrito">
          🛒
          <span className="absolute -right-1 -top-1 grid h-5 w-5 place-items-center rounded-full bg-red-500 text-[10px] font-black">
            {cartCount}
          </span>
        </Link>
      </div>
      {showSearch && (
        <div className="container-page pb-4 md:hidden">
          <input
            className="input-dark"
            onChange={(event) => setSearchQuery(event.target.value)}
            placeholder="Buscar artistas, eventos, venues..."
            value={searchQuery}
          />
        </div>
      )}
    </header>
  )
}
