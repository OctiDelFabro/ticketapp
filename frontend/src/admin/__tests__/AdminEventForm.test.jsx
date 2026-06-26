import React from 'react'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import AdminEventForm from '../AdminEventForm.jsx'

const validForm = {
  title: 'Festival Test',
  description: 'Una descripción válida para el evento',
  category: 'Música',
  location: 'Estadio Test',
  startDate: '2026-08-15T21:00',
  duration: '90',
  capacity: '100',
}

const setup = () => {
  const onSubmit = vi.fn()
  const onClose = vi.fn()
  render(<AdminEventForm onSubmit={onSubmit} onClose={onClose} loading={false} />)
  return { onSubmit, onClose }
}

const fillRequiredFields = async (user) => {
  await user.type(screen.getByLabelText(/título/i), validForm.title)
  await user.selectOptions(screen.getByLabelText(/categoría/i), validForm.category)
  await user.type(screen.getByLabelText(/lugar/i), validForm.location)
  await user.type(screen.getByLabelText(/fecha y hora/i), validForm.startDate)
  await user.clear(screen.getByLabelText(/duración/i))
  await user.type(screen.getByLabelText(/duración/i), validForm.duration)
  await user.clear(screen.getByLabelText(/cupo/i))
  await user.type(screen.getByLabelText(/cupo/i), validForm.capacity)
  await user.type(screen.getByLabelText(/descripción/i), validForm.description)
}

describe('AdminEventForm', () => {
  it('muestra el campo Imagen URL y su texto de ayuda', () => {
    setup()

    expect(screen.getByLabelText(/imagen url/i)).toBeInTheDocument()
    expect(screen.getByText(/pegá una url directa de imagen/i)).toBeInTheDocument()
  })

  it('muestra preview cuando se ingresa una URL válida', async () => {
    const user = userEvent.setup()
    setup()

    await user.type(screen.getByLabelText(/imagen url/i), 'https://example.com/evento.jpg')

    expect(screen.getByRole('img', { name: /preview de imagen del evento/i })).toHaveAttribute('src', 'https://example.com/evento.jpg')
  })

  it('muestra error si se ingresa una URL inválida', async () => {
    const user = userEvent.setup()
    setup()

    await fillRequiredFields(user)
    await user.type(screen.getByLabelText(/imagen url/i), 'imagen cualquiera')
    await user.click(screen.getByRole('button', { name: /guardar/i }))

    expect(screen.getByText(/usá una url http:\/\//i)).toBeInTheDocument()
  })

  it('permite guardar cuando los campos requeridos son válidos', async () => {
    const user = userEvent.setup()
    const { onSubmit } = setup()

    await fillRequiredFields(user)
    await user.type(screen.getByLabelText(/imagen url/i), '/events/default.jpg')
    await user.click(screen.getByRole('button', { name: /guardar/i }))

    expect(onSubmit).toHaveBeenCalledWith(expect.objectContaining({
      title: validForm.title,
      description: validForm.description,
      category: validForm.category,
      location: validForm.location,
      image_url: '/events/default.jpg',
      duration_minutes: Number(validForm.duration),
      capacity: Number(validForm.capacity),
    }))
  })
})
