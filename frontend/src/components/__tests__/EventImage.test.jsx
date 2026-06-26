import { fireEvent, render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import EventImage from '../EventImage.jsx'
import { fallbackEventImage } from '../../utils/events.js'

describe('EventImage', () => {
  it('renderiza la imagen cuando recibe una URL válida', () => {
    render(<EventImage src="https://example.com/event.jpg" alt="Evento principal" />)

    expect(screen.getByRole('img', { name: /evento principal/i })).toHaveAttribute('src', 'https://example.com/event.jpg')
  })

  it('cambia a fallback cuando ocurre error de carga', () => {
    render(<EventImage src="https://example.com/broken.jpg" alt="Evento con error" />)

    fireEvent.error(screen.getByRole('img', { name: /evento con error/i }))

    expect(screen.getByRole('img', { name: /evento con error/i })).toHaveAttribute('src', fallbackEventImage)
  })

  it('usa fallback si no recibe imagen', () => {
    render(<EventImage alt="Evento sin imagen" />)

    expect(screen.getByRole('img', { name: /evento sin imagen/i })).toHaveAttribute('src', fallbackEventImage)
  })
})
