import Button from './Button.jsx'
import { formatPrice, normalizePrice } from '../utils/formatters.js'

export default function OrderSummary({ event, quantity = 1, mode = 'purchase', buttonText, onNext, onRemove, disabled = false }) {
  const eventPrice = normalizePrice(event?.price)
  const isGift = mode === 'gift'
  const displayQuantity = isGift ? 1 : quantity
  const total = eventPrice * displayQuantity
  const eventPath = event ? `/evento/${event.id}` : '/'

  if (!event) {
    return (
      <aside className="glass-card sticky top-24 rounded-3xl p-5">
        <h3 className="text-xl font-black">Tu carrito está vacío</h3>
        <p className="mt-3 text-sm leading-6 text-gray-400">
          Agregá una entrada desde el detalle del evento para continuar.
        </p>
        <div className="mt-6 grid gap-3">
          <Button to={eventPath} className="w-full">Volver al evento</Button>
          <Button to="/" variant="secondary" className="w-full">Ir al inicio</Button>
        </div>
      </aside>
    )
  }

  return (
    <aside className="glass-card sticky top-24 rounded-3xl p-5">
      <h3 className="text-xl font-black">Resumen del pedido</h3>
      <div className="mt-5 flex gap-3">
        <img src={event.image} alt={event.title} className="h-20 w-20 rounded-2xl object-cover" />
        <div>
          <p className="font-black">{event.title}</p>
          <p className="text-sm text-gray-400">{event.date}</p>
          <p className="text-sm text-gray-400">{event.venue}</p>
        </div>
      </div>
      <div className="my-5 space-y-3 border-y border-ticket-border py-5 text-sm">
        <div className="flex justify-between text-gray-300">
          <span>{isGift ? 'Regalo para otro usuario' : `Entrada general x ${displayQuantity}`}</span>
          <span>{formatPrice(eventPrice)}</span>
        </div>
        {isGift && <div className="flex justify-between text-gray-300"><span>Cantidad</span><span>1</span></div>}
      </div>
      <div className="mb-5 flex justify-between text-xl font-black">
        <span>Total</span>
        <span className="text-violet-200">{formatPrice(total)}</span>
      </div>
      <Button onClick={onRemove} variant="danger" className="mb-3 w-full" disabled={disabled}>Quitar del carrito</Button>
      {buttonText && <Button onClick={onNext} className="w-full" disabled={disabled}>{buttonText}</Button>}
    </aside>
  )
}
