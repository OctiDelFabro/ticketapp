import { useState } from 'react'
import { Link } from 'react-router-dom'
import Button from '../components/Button.jsx'
import Logo from '../components/Logo.jsx'

export default function Login() {
  const [mode, setMode] = useState('login')
  const [showPassword, setShowPassword] = useState(false)
  const [message, setMessage] = useState('')
  const submit = (event) => { event.preventDefault(); setMessage('Funcionalidad mock por ahora') }

  return (
    <div className="app-shell bg-ticket-alt">
      <div className="container-page py-6"><Link to="/" className="font-bold text-gray-300 transition hover:text-white">← Volver</Link></div>
      <main className="container-page flex min-h-[78vh] items-center justify-center">
        <div className="w-full max-w-md">
          <Logo centered />
          <div className="glass-card mt-8 rounded-3xl p-6 sm:p-8">
            <div className="grid grid-cols-2 rounded-2xl bg-ticket-card2 p-1">
              {['login', 'register'].map((tab) => <button key={tab} onClick={() => { setMode(tab); setMessage('') }} className={`rounded-xl py-3 font-black transition ${mode === tab ? 'bg-ticket-purple text-white shadow-glow' : 'text-gray-400 hover:text-white'}`}>{tab === 'login' ? 'Iniciar sesión' : 'Crear cuenta'}</button>)}
            </div>
            <form onSubmit={submit} className="mt-6 space-y-4">
              {mode === 'register' && <div className="grid gap-4 sm:grid-cols-2"><input className="input-dark" placeholder="Nombre" /><input className="input-dark" placeholder="Apellido" /></div>}
              <input className="input-dark" type="email" placeholder="Email" />
              <div className="relative"><input className="input-dark pr-14" type={showPassword ? 'text' : 'password'} placeholder="Contraseña" /><button type="button" onClick={() => setShowPassword(!showPassword)} className="absolute right-4 top-1/2 -translate-y-1/2 text-lg">{showPassword ? '🙈' : '👁️'}</button></div>
              {mode === 'login' && <div className="text-right"><a href="#" className="text-sm font-bold text-violet-300 hover:text-violet-100">¿Olvidaste tu contraseña?</a></div>}
              <Button className="w-full" type="submit">{mode === 'login' ? 'Iniciar sesión' : 'Crear cuenta'}</Button>
            </form>
            {message && <p className="mt-4 rounded-2xl border border-ticket-purple2/40 bg-ticket-purple/15 p-3 text-center text-sm text-violet-100">{message}</p>}
            {mode === 'login' && <><div className="my-6 flex items-center gap-3 text-sm text-gray-500"><span className="h-px flex-1 bg-ticket-border" />o continuá con<span className="h-px flex-1 bg-ticket-border" /></div><div className="grid grid-cols-2 gap-3"><Button type="button" variant="secondary">G Google</Button><Button type="button" variant="secondary">f Facebook</Button></div></>}
          </div>
        </div>
      </main>
    </div>
  )
}
