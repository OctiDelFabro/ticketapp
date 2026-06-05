import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import Button from '../components/Button.jsx'
import CheckoutProgress from '../components/CheckoutProgress.jsx'
import OrderSummary from '../components/OrderSummary.jsx'
import { events, formatPrice, serviceFee } from '../data/mockData.js'

export default function Checkout() {
  const event = events[0]
  const [step, setStep] = useState(1)
  const [quantity, setQuantity] = useState(2)
  const [payment, setPayment] = useState('Crédito')
  const navigate = useNavigate()
  const total = event.price * quantity + serviceFee(event.price, quantity)
  const next = () => setStep((value) => Math.min(4, value + 1))

  return (
    <div className="app-shell bg-ticket-alt">
      <header className="border-b border-ticket-border bg-ticket-bg/70 backdrop-blur-xl"><div className="container-page flex flex-col gap-3 py-5 sm:flex-row sm:items-center sm:justify-between"><Link to="/evento/arctic-monkeys" className="font-bold text-gray-300 hover:text-white">← Volver al evento</Link><h1 className="text-xl font-black">Arctic Monkeys — Ticketar Checkout</h1></div></header>
      <main className="container-page py-8">
        <CheckoutProgress currentStep={step} />
        {step < 4 ? <div className="mt-8 grid gap-8 lg:grid-cols-[1fr_380px]">
          <section className="space-y-6">
            {step === 1 && <><div className="glass-card rounded-3xl p-5"><div className="flex flex-col gap-5 sm:flex-row"><img src={event.image} alt={event.title} className="h-44 rounded-2xl object-cover sm:w-60" /><div><BadgeLine text="Evento seleccionado" /><h2 className="mt-2 text-3xl font-black">{event.title}</h2><p className="mt-3 text-gray-300">{event.date} · {event.time}</p><p className="text-gray-400">{event.venue}</p></div></div></div><div className="glass-card rounded-3xl p-6"><h2 className="text-2xl font-black">Tipo de entrada</h2><div className="mt-4 rounded-2xl border border-ticket-purple2 bg-ticket-purple/15 p-4"><div className="flex justify-between"><span className="font-black">General</span><span className="font-black text-violet-200">{formatPrice(event.price)}</span></div></div><h3 className="mt-6 font-black">Cantidad</h3><div className="mt-3 flex flex-wrap gap-3">{[1,2,3,4,5,6].map((n) => <button key={n} onClick={() => setQuantity(n)} className={`h-12 w-12 rounded-2xl border font-black ${quantity === n ? 'border-ticket-purple bg-ticket-purple shadow-glow' : 'border-ticket-border bg-ticket-card2 text-gray-300'}`}>{n}</button>)}</div></div></>}
            {step === 2 && <div className="glass-card rounded-3xl p-6"><h2 className="text-2xl font-black">Tus datos</h2><form className="mt-5 grid gap-4"><div className="grid gap-4 sm:grid-cols-2"><input className="input-dark" placeholder="Nombre" /><input className="input-dark" placeholder="Apellido" /></div><input className="input-dark" type="email" placeholder="Email" /><input className="input-dark" placeholder="Teléfono" /><input className="input-dark" placeholder="DNI" /></form></div>}
            {step === 3 && <div className="glass-card rounded-3xl p-6"><h2 className="text-2xl font-black">Pago</h2><div className="mt-5 flex flex-wrap gap-3">{['Crédito','Débito','MercadoPago'].map((item) => <button key={item} onClick={() => setPayment(item)} className={`rounded-2xl border px-5 py-3 font-black ${payment === item ? 'border-ticket-purple bg-ticket-purple' : 'border-ticket-border bg-ticket-card2 text-gray-300'}`}>{item}</button>)}</div><form className="mt-6 grid gap-4"><input className="input-dark" placeholder="Número de tarjeta" /><div className="grid gap-4 sm:grid-cols-2"><input className="input-dark" placeholder="Vencimiento" /><input className="input-dark" placeholder="CVV" /></div><input className="input-dark" placeholder="Titular de la tarjeta" /></form></div>}
          </section>
          <OrderSummary event={event} quantity={quantity} buttonText={step === 3 ? 'Finalizar compra' : 'Continuar →'} onNext={next} />
        </div> : <section className="mx-auto mt-10 max-w-2xl text-center"><div className="mx-auto grid h-24 w-24 place-items-center rounded-full bg-ticket-purple text-5xl shadow-glow">✓</div><h2 className="mt-6 text-4xl font-black">¡Compra confirmada!</h2><p className="mt-3 text-gray-400">Recibirás tus entradas en tu email en los próximos minutos.</p><div className="glass-card mt-8 rounded-3xl p-6 text-left"><div className="mx-auto grid h-36 w-36 place-items-center rounded-2xl border border-ticket-border bg-white text-5xl text-ticket-bg">▦</div><p className="mt-5 text-center font-black tracking-widest">TICKETAR-2026-AM-0047819</p><div className="mt-6 grid gap-3 text-sm text-gray-300 sm:grid-cols-2"><p>Evento: <b className="text-white">{event.title}</b></p><p>Fecha: <b className="text-white">{event.date}</b></p><p>Lugar: <b className="text-white">{event.venue}</b></p><p>Tipo: <b className="text-white">General × {quantity}</b></p><p className="sm:col-span-2">Total pagado: <b className="text-violet-200">{formatPrice(total)}</b></p></div></div><div className="mt-6 flex flex-wrap justify-center gap-4"><Button onClick={() => navigate('/mis-entradas')}>Ir a Mis Entradas</Button><Button variant="secondary">Descargar PDF</Button></div></section>}
      </main>
    </div>
  )
}

function BadgeLine({ text }) {
  return <p className="text-xs font-black uppercase tracking-[0.25em] text-ticket-purple2">{text}</p>
}
