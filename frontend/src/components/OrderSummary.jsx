import { formatPrice, serviceFee } from '../data/mockData.js'
import Button from './Button.jsx'

export default function OrderSummary({ event, quantity, buttonText, onNext }) {
  const subtotal = event.price * quantity
  const fee = serviceFee(event.price, quantity)
  return (
    <aside className="glass-card sticky top-24 rounded-3xl p-5">
      <h3 className="text-xl font-black">Resumen del pedido</h3>
      <div className="mt-5 flex gap-3">
        <img src={event.image} alt={event.title} className="h-20 w-20 rounded-2xl object-cover" />
        <div><p className="font-black">{event.title}</p><p className="text-sm text-gray-400">{event.date}</p><p className="text-sm text-gray-400">{event.venue}</p></div>
      </div>
      <div className="my-5 space-y-3 border-y border-ticket-border py-5 text-sm">
        <div className="flex justify-between text-gray-300"><span>General × {quantity}</span><span>{formatPrice(subtotal)}</span></div>
        <div className="flex justify-between text-gray-300"><span>Cargo por servicio</span><span>{formatPrice(fee)}</span></div>
      </div>
      <div className="mb-5 flex justify-between text-xl font-black"><span>Total</span><span className="text-violet-200">{formatPrice(subtotal + fee)}</span></div>
      {buttonText && <Button onClick={onNext} className="w-full">{buttonText}</Button>}
    </aside>
  )
}
