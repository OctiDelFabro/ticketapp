import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import Button from '../components/Button.jsx'
import CheckoutProgress from '../components/CheckoutProgress.jsx'
import OrderSummary from '../components/OrderSummary.jsx'
import { formatPrice, getEventById, serviceFee } from '../data/mockData.js'

const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
const phonePattern = /^\+?[\d\s]{8,20}$/

export default function Checkout({ cartItem, setCartItem }) {
  const [step, setStep] = useState(1)
  const [payment, setPayment] = useState('Crédito')
  const [customer, setCustomer] = useState({ firstName: '', lastName: '', email: '', phone: '', dni: '' })
  const [card, setCard] = useState({ number: '', expiry: '', cvv: '', holder: '' })
  const [errors, setErrors] = useState({})
  const navigate = useNavigate()

  const event = cartItem?.event ?? getEventById(cartItem?.eventId)
  const quantity = cartItem?.quantity ?? 1
  const total = event ? event.price * quantity + serviceFee(event.price, quantity) : 0

  const eventPath = event ? `/evento/${event.id}` : '/evento/arctic-monkeys'

  const updateCustomer = (field, value) => {
    setCustomer((current) => ({ ...current, [field]: value }))
    setErrors((current) => ({ ...current, [field]: '' }))
  }

  const updateCard = (field, value) => {
    setCard((current) => ({ ...current, [field]: value }))
    setErrors((current) => ({ ...current, [field]: '' }))
  }

  const validateCustomer = () => {
    const nextErrors = {}
    if (!customer.firstName.trim()) nextErrors.firstName = 'Ingresá tu nombre.'
    if (!customer.lastName.trim()) nextErrors.lastName = 'Ingresá tu apellido.'
    if (!customer.email.trim() || !emailPattern.test(customer.email)) nextErrors.email = 'Ingresá un email válido.'
    if (!phonePattern.test(customer.phone.trim())) nextErrors.phone = 'Ingresá un teléfono válido.'
    if (!/^\d{7,8}$/.test(customer.dni.trim())) nextErrors.dni = 'El DNI debe contener solo números y tener entre 7 y 8 dígitos.'
    setErrors(nextErrors)
    return Object.keys(nextErrors).length === 0
  }

  const validatePayment = () => {
    const nextErrors = {}
    if (!/^\d{13,19}$/.test(card.number.trim())) nextErrors.number = 'El número de tarjeta debe contener entre 13 y 19 números.'
    if (!/^(0[1-9]|1[0-2])\/\d{2}$/.test(card.expiry.trim())) nextErrors.expiry = 'Ingresá el vencimiento con formato MM/AA.'
    if (!/^\d{3,4}$/.test(card.cvv.trim())) nextErrors.cvv = 'El CVV debe contener 3 o 4 números.'
    if (!card.holder.trim()) nextErrors.holder = 'Ingresá el titular de la tarjeta.'
    setErrors(nextErrors)
    return Object.keys(nextErrors).length === 0
  }

  const next = () => {
    if (!event) return
    if (step === 2 && !validateCustomer()) return
    if (step === 3 && !validatePayment()) return
    setErrors({})
    setStep((value) => Math.min(4, value + 1))
  }

  const goBack = () => {
    if (step === 1) {
      navigate(eventPath)
      return
    }
    setErrors({})
    setStep((value) => Math.max(1, value - 1))
  }

  const updateQuantity = (nextQuantity) => {
    setCartItem((item) => item ? { ...item, quantity: nextQuantity } : item)
  }

  const removeCartItem = () => {
    setCartItem(null)
    setStep(1)
  }

  return (
    <div className="bg-ticket-alt">
      <main className="container-page py-8">
        <div className="mb-6 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <button onClick={goBack} className="w-fit font-bold text-gray-300 hover:text-white">
            {step === 1 ? '← Volver al evento' : '← Volver'}
          </button>
          <h1 className="text-xl font-black">{event ? `${event.title} - TicketApp Checkout` : 'TicketApp Checkout'}</h1>
        </div>
        <CheckoutProgress currentStep={step} />
        {!event ? (
          <div className="mt-8 grid gap-8 lg:grid-cols-[1fr_380px]">
            <section className="glass-card rounded-3xl p-6 sm:p-8">
              <h2 className="text-3xl font-black">Tu carrito está vacío</h2>
              <p className="mt-3 max-w-xl text-gray-400">
                Agregá una entrada desde el detalle del evento para continuar.
              </p>
              <div className="mt-6 flex flex-wrap gap-3">
                <Button to={eventPath}>Volver al evento</Button>
                <Button to="/" variant="secondary">Ir al inicio</Button>
              </div>
            </section>
            <OrderSummary event={null} />
          </div>
        ) : step < 4 ? (
          <div className="mt-8 grid gap-8 lg:grid-cols-[1fr_380px]">
            <section className="space-y-6">
              {step === 1 && (
                <>
                  <div className="glass-card rounded-3xl p-5">
                    <div className="flex flex-col gap-5 sm:flex-row">
                      <img src={event.image} alt={event.title} className="h-44 rounded-2xl object-cover sm:w-60" />
                      <div>
                        <BadgeLine text="Evento seleccionado" />
                        <h2 className="mt-2 text-3xl font-black">{event.title}</h2>
                        <p className="mt-3 text-gray-300">{event.date} - {event.time}</p>
                        <p className="text-gray-400">{event.venue}</p>
                      </div>
                    </div>
                  </div>
                  <div className="glass-card rounded-3xl p-6">
                    <h2 className="text-2xl font-black">Tipo de entrada</h2>
                    <div className="mt-4 rounded-2xl border border-ticket-purple2 bg-ticket-purple/15 p-4">
                      <div className="flex justify-between">
                        <span className="font-black">General</span>
                        <span className="font-black text-violet-200">{formatPrice(event.price)}</span>
                      </div>
                    </div>
                    <h3 className="mt-6 font-black">Cantidad</h3>
                    <div className="mt-3 flex flex-wrap gap-3">
                      {[1, 2, 3, 4, 5, 6].map((n) => (
                        <button
                          key={n}
                          onClick={() => updateQuantity(n)}
                          className={`h-12 w-12 rounded-2xl border font-black ${quantity === n ? 'border-ticket-purple bg-ticket-purple shadow-glow' : 'border-ticket-border bg-ticket-card2 text-gray-300'}`}
                        >
                          {n}
                        </button>
                      ))}
                    </div>
                  </div>
                </>
              )}
              {step === 2 && (
                <div className="glass-card rounded-3xl p-6">
                  <h2 className="text-2xl font-black">Tus datos</h2>
                  <form className="mt-5 grid gap-4">
                    <div className="grid gap-4 sm:grid-cols-2">
                      <Field error={errors.firstName}><input className={`input-dark ${errors.firstName ? 'border-red-500' : ''}`} onChange={(event) => updateCustomer('firstName', event.target.value)} placeholder="Nombre" value={customer.firstName} /></Field>
                      <Field error={errors.lastName}><input className={`input-dark ${errors.lastName ? 'border-red-500' : ''}`} onChange={(event) => updateCustomer('lastName', event.target.value)} placeholder="Apellido" value={customer.lastName} /></Field>
                    </div>
                    <Field error={errors.email}><input className={`input-dark ${errors.email ? 'border-red-500' : ''}`} onChange={(event) => updateCustomer('email', event.target.value)} type="email" placeholder="Email" value={customer.email} /></Field>
                    <Field error={errors.phone}><input className={`input-dark ${errors.phone ? 'border-red-500' : ''}`} onChange={(event) => updateCustomer('phone', event.target.value)} placeholder="Teléfono" value={customer.phone} /></Field>
                    <Field error={errors.dni}><input className={`input-dark ${errors.dni ? 'border-red-500' : ''}`} inputMode="numeric" onChange={(event) => updateCustomer('dni', event.target.value.replace(/\D/g, ''))} placeholder="DNI" value={customer.dni} /></Field>
                  </form>
                </div>
              )}
              {step === 3 && (
                <div className="glass-card rounded-3xl p-6">
                  <h2 className="text-2xl font-black">Pago</h2>
                  <div className="mt-5 flex flex-wrap gap-3">
                    {['Crédito', 'Débito', 'MercadoPago'].map((item) => (
                      <button
                        key={item}
                        onClick={() => setPayment(item)}
                        className={`rounded-2xl border px-5 py-3 font-black ${payment === item ? 'border-ticket-purple bg-ticket-purple' : 'border-ticket-border bg-ticket-card2 text-gray-300'}`}
                      >
                        {item}
                      </button>
                    ))}
                  </div>
                  <form className="mt-6 grid gap-4">
                    <Field error={errors.number}><input className={`input-dark ${errors.number ? 'border-red-500' : ''}`} inputMode="numeric" onChange={(event) => updateCard('number', event.target.value.replace(/\D/g, ''))} placeholder="Número de tarjeta" value={card.number} /></Field>
                    <div className="grid gap-4 sm:grid-cols-2">
                      <Field error={errors.expiry}><input className={`input-dark ${errors.expiry ? 'border-red-500' : ''}`} onChange={(event) => updateCard('expiry', event.target.value)} placeholder="Vencimiento MM/AA" value={card.expiry} /></Field>
                      <Field error={errors.cvv}><input className={`input-dark ${errors.cvv ? 'border-red-500' : ''}`} inputMode="numeric" onChange={(event) => updateCard('cvv', event.target.value.replace(/\D/g, ''))} placeholder="CVV" value={card.cvv} /></Field>
                    </div>
                    <Field error={errors.holder}><input className={`input-dark ${errors.holder ? 'border-red-500' : ''}`} onChange={(event) => updateCard('holder', event.target.value)} placeholder="Titular de la tarjeta" value={card.holder} /></Field>
                  </form>
                </div>
              )}
            </section>
            <OrderSummary
              event={event}
              quantity={quantity}
              buttonText={step === 3 ? 'Finalizar compra' : 'Continuar'}
              onNext={next}
              onRemove={removeCartItem}
            />
          </div>
        ) : (
          <section className="mx-auto mt-10 max-w-2xl text-center">
            <div className="mx-auto grid h-24 w-24 place-items-center rounded-full bg-ticket-purple text-5xl font-black shadow-glow">✓</div>
            <h2 className="mt-6 text-4xl font-black">¡Compra confirmada!</h2>
            <p className="mt-3 text-gray-400">Recibirás la confirmación de tu compra en tu email.</p>
            <div className="glass-card mt-8 rounded-3xl p-6 text-left">
              <div className="mt-6 grid gap-3 text-sm text-gray-300 sm:grid-cols-2">
                <p>Artista / evento: <b className="text-white">{event.title}</b></p>
                <p>Fecha: <b className="text-white">{event.date}</b></p>
                <p>Lugar: <b className="text-white">{event.venue}</b></p>
                <p>Cantidad de entradas: <b className="text-white">{quantity}</b></p>
                <p>Tipo de entrada: <b className="text-white">General</b></p>
                <p className="sm:col-span-2">Total pagado: <b className="text-violet-200">{formatPrice(total)}</b></p>
              </div>
            </div>
            <div className="mt-6 flex flex-wrap justify-center gap-4">
              <Button onClick={() => navigate('/mis-entradas')}>Ir a Mis Entradas</Button>
              <Button onClick={() => navigate('/')} variant="secondary">Volver al inicio</Button>
            </div>
          </section>
        )}
      </main>
    </div>
  )
}

function BadgeLine({ text }) {
  return <p className="text-xs font-black uppercase tracking-[0.25em] text-ticket-purple2">{text}</p>
}

function Field({ children, error }) {
  return (
    <label className="block">
      {children}
      {error && <span className="mt-1 block text-sm text-red-400">{error}</span>}
    </label>
  )
}
