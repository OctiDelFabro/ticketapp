import { Link } from 'react-router-dom'
import Button from './Button.jsx'
import Logo from './Logo.jsx'

export default function TopNavbar() {
  return (
    <header className="sticky top-0 z-40 border-b border-ticket-border/70 bg-ticket-bg/75 backdrop-blur-xl">
      <div className="container-page flex items-center gap-4 py-4">
        <Logo />
        <div className="hidden flex-1 md:block">
          <input className="input-dark" placeholder="Buscar artistas, eventos, venues..." />
        </div>
        <Button to="/login" className="hidden sm:inline-flex">Ingresar</Button>
        <Link to="/checkout" className="relative grid h-12 w-12 place-items-center rounded-2xl border border-ticket-border bg-ticket-card2 text-xl transition hover:border-ticket-purple2/60">
          🛒<span className="absolute -right-1 -top-1 grid h-5 w-5 place-items-center rounded-full bg-red-500 text-[10px] font-black">1</span>
        </Link>
      </div>
      <div className="container-page pb-4 md:hidden"><input className="input-dark" placeholder="Buscar artistas, eventos, venues..." /></div>
    </header>
  )
}
