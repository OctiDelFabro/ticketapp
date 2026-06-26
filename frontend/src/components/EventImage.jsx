import { useEffect, useState } from 'react'
import { fallbackEventImage, getEventImage } from '../utils/events.js'

export default function EventImage({ event, src, alt = 'Imagen del evento', className = '', ...props }) {
  const imageSrc = src || getEventImage(event)
  const [currentSrc, setCurrentSrc] = useState(imageSrc)

  useEffect(() => {
    setCurrentSrc(imageSrc)
  }, [imageSrc])

  return (
    <img
      src={currentSrc || fallbackEventImage}
      alt={alt}
      className={className}
      onError={() => {
        if (currentSrc !== fallbackEventImage) setCurrentSrc(fallbackEventImage)
      }}
      {...props}
    />
  )
}
