import React from 'react'
import { render, screen } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import AdminRoute from '../AdminRoute.jsx'
import { getStoredUser, isAuthenticated } from '../../utils/auth.js'
import { isAdminUser } from '../../utils/admin.js'

vi.mock('../../utils/auth.js', () => ({
  getStoredUser: vi.fn(),
  isAuthenticated: vi.fn(),
}))

vi.mock('../../utils/admin.js', () => ({
  isAdminUser: vi.fn(),
}))

const renderRoute = () => render(
  <MemoryRouter initialEntries={["/admin/eventos"]}>
    <Routes>
      <Route path="/login" element={<p>Login requerido</p>} />
      <Route path="/admin/eventos" element={<AdminRoute><p>Contenido admin</p></AdminRoute>} />
    </Routes>
  </MemoryRouter>,
)

describe('AdminRoute', () => {
  beforeEach(() => vi.clearAllMocks())

  it('redirige al login si no hay usuario autenticado', () => {
    isAuthenticated.mockReturnValue(false)

    renderRoute()

    expect(screen.getByText('Login requerido')).toBeInTheDocument()
    expect(screen.queryByText('Contenido admin')).not.toBeInTheDocument()
  })

  it('bloquea el acceso si hay usuario cliente', () => {
    const user = { role: 'client' }
    isAuthenticated.mockReturnValue(true)
    getStoredUser.mockReturnValue(user)
    isAdminUser.mockReturnValue(false)

    renderRoute()

    expect(screen.getByText(/acceso denegado/i)).toBeInTheDocument()
    expect(screen.queryByText('Contenido admin')).not.toBeInTheDocument()
  })

  it('permite renderizar el contenido protegido si hay usuario admin', () => {
    const user = { role: 'admin' }
    isAuthenticated.mockReturnValue(true)
    getStoredUser.mockReturnValue(user)
    isAdminUser.mockReturnValue(true)

    renderRoute()

    expect(screen.getByText('Contenido admin')).toBeInTheDocument()
  })
})
