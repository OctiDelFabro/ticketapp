import { describe, expect, it } from 'vitest'
import { isSupportedEventImageURL } from '../events.js'

describe('isSupportedEventImageURL', () => {
  it('acepta URLs http://', () => {
    expect(isSupportedEventImageURL('http://example.com/image.jpg')).toBe(true)
  })

  it('acepta URLs https://', () => {
    expect(isSupportedEventImageURL('https://example.com/image.jpg')).toBe(true)
  })

  it('acepta rutas locales que empiezan con /', () => {
    expect(isSupportedEventImageURL('/events/default.jpg')).toBe(true)
  })

  it('acepta vacío', () => {
    expect(isSupportedEventImageURL('')).toBe(true)
  })

  it('rechaza texto inválido', () => {
    expect(isSupportedEventImageURL('imagen cualquiera')).toBe(false)
  })
})
