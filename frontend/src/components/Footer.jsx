import { Link } from 'react-router-dom'
import Logo from './Logo.jsx'

export default function Footer() {
  return (
    <footer className="border-t border-[#1f1f2e] bg-[#111118] pb-28 pt-10 text-white sm:pb-24">
      <div className="container-page flex flex-col gap-8 md:flex-row md:items-start md:justify-between">
        <div className="max-w-md">
          <Logo />
          <p className="mt-4 leading-7 text-gray-400">
            La forma más simple de encontrar y comprar entradas para tus eventos favoritos.
          </p>
        </div>
        <nav className="grid gap-3 text-sm font-bold text-gray-400 sm:grid-cols-3 sm:gap-8">
          <Link to="/" className="transition hover:text-white">Eventos</Link>
          <Link to="/mis-entradas" className="transition hover:text-white">Ayuda</Link>
          <Link to="/login" className="transition hover:text-white">Contacto</Link>
        </nav>
      </div>
      <div className="container-page mt-8 border-t border-ticket-border pt-6 text-sm text-gray-400">
        © 2026 TicketApp. Todos los derechos reservados.
      </div>
    </footer>
  )
}
