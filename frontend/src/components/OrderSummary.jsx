import Button from './Button.jsx'

export default function OrderSummary({ event, quantity, buttonText, onNext, onRemove }) {
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
          <span>General x {quantity}</span>
          <span>Entrada</span>
        </div>
        <div className="flex justify-between text-gray-300">
          <span>Emisión backend</span>
          <span>Ticket real</span>
        </div>
      </div>
      <div className="mb-5 flex justify-between text-xl font-black">
        <span>Total</span>
        <span className="text-violet-200">General</span>
      </div>
      <Button onClick={onRemove} variant="danger" className="mb-3 w-full">Quitar del carrito</Button>
      {buttonText && <Button onClick={onNext} className="w-full">{buttonText}</Button>}
    </aside>
  )
}
