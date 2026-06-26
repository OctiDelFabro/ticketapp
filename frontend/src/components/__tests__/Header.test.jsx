import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { describe, expect, it, vi } from 'vitest'
import Header from '../Header.jsx'

const renderHeader = (props = {}) => render(
  <MemoryRouter initialEntries={["/"]}>
    <Header searchQuery="" setSearchQuery={vi.fn()} {...props} />
  </MemoryRouter>,
)

describe('Header', () => {
  it('cliente logueado ve Mis Entradas y opciones de cliente', () => {
    renderHeader({ isLoggedIn: true, isAdmin: false, cartCount: 2 })

    expect(screen.getByRole('link', { name: /mis entradas/i })).toHaveAttribute('href', '/mis-entradas')
    expect(screen.getByRole('link', { name: /ir al carrito/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /cerrar sesión/i })).toBeInTheDocument()
  })

  it('admin logueado ve Panel Admin', () => {
    renderHeader({ isLoggedIn: true, isAdmin: true })

    expect(screen.getByRole('link', { name: /panel admin/i })).toHaveAttribute('href', '/admin/eventos')
  })

  it('admin no ve carrito', () => {
    renderHeader({ isLoggedIn: true, isAdmin: true })

    expect(screen.queryByRole('link', { name: /ir al carrito/i })).not.toBeInTheDocument()
  })

  it('usuario no logueado ve opciones de login', () => {
    renderHeader({ isLoggedIn: false, isAdmin: false })

    expect(screen.getByRole('link', { name: /ingresar/i })).toHaveAttribute('href', '/login')
    expect(screen.queryByRole('button', { name: /cerrar sesión/i })).not.toBeInTheDocument()
  })
})
