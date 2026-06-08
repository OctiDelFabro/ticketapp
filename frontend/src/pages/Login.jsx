import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import Button from '../components/Button.jsx'
import Logo from '../components/Logo.jsx'
import { login, register } from '../services/api.js'
import { saveAuthSession } from '../utils/auth.js'

const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

export default function Login({ setIsLoggedIn }) {
  const [mode, setMode] = useState('login')
  const [showPassword, setShowPassword] = useState(false)
  const [form, setForm] = useState({ firstName: '', lastName: '', email: '', password: '' })
  const [errors, setErrors] = useState({})
  const [submitError, setSubmitError] = useState('')
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()

  const updateField = (field, value) => {
    setForm((current) => ({ ...current, [field]: value }))
    setErrors((current) => ({ ...current, [field]: '' }))
    setSubmitError('')
  }

  const validate = () => {
    const nextErrors = {}
    if (mode === 'register' && !form.firstName.trim()) nextErrors.firstName = 'Ingresá tu nombre.'
    if (mode === 'register' && !form.lastName.trim()) nextErrors.lastName = 'Ingresá tu apellido.'
    if (!form.email.trim() || !emailPattern.test(form.email)) nextErrors.email = 'Ingresá un email válido.'
    if (!form.password.trim()) nextErrors.password = 'Ingresá tu contraseña.'
    setErrors(nextErrors)
    return Object.keys(nextErrors).length === 0
  }

  const submit = async (event) => {
    event.preventDefault()
    if (!validate()) return

    setLoading(true)
    setSubmitError('')
    try {
      const payload = mode === 'register'
        ? { name: `${form.firstName} ${form.lastName}`.trim(), email: form.email, password: form.password }
        : { email: form.email, password: form.password }
      const response = mode === 'register' ? await register(payload) : await login(payload)
      saveAuthSession(response)
      setIsLoggedIn(true)
      navigate('/')
    } catch (err) {
      setSubmitError(err.message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="bg-ticket-alt">
      <div className="container-page py-6"><button onClick={() => navigate('/')} className="font-bold text-gray-300 transition hover:text-white">← Volver</button></div>
      <main className="container-page flex min-h-[78vh] items-center justify-center">
        <div className="w-full max-w-md">
          <Logo centered />
          <div className="glass-card mt-8 rounded-3xl p-6 sm:p-8">
            <div className="grid grid-cols-2 rounded-2xl bg-ticket-card2 p-1">
              {['login', 'register'].map((tab) => <button key={tab} type="button" onClick={() => { setMode(tab); setErrors({}); setSubmitError('') }} className={`rounded-xl py-3 font-black transition ${mode === tab ? 'bg-ticket-purple text-white shadow-glow' : 'text-gray-400 hover:text-white'}`}>{tab === 'login' ? 'Iniciar sesión' : 'Crear cuenta'}</button>)}
            </div>
            <form onSubmit={submit} className="mt-6 space-y-4">
              {mode === 'register' && (
                <div className="grid gap-4 sm:grid-cols-2">
                  <Field error={errors.firstName}><input className={`input-dark ${errors.firstName ? 'border-red-500' : ''}`} onChange={(event) => updateField('firstName', event.target.value)} placeholder="Nombre" value={form.firstName} /></Field>
                  <Field error={errors.lastName}><input className={`input-dark ${errors.lastName ? 'border-red-500' : ''}`} onChange={(event) => updateField('lastName', event.target.value)} placeholder="Apellido" value={form.lastName} /></Field>
                </div>
              )}
              <Field error={errors.email}><input className={`input-dark ${errors.email ? 'border-red-500' : ''}`} onChange={(event) => updateField('email', event.target.value)} type="email" placeholder="Email" value={form.email} /></Field>
              <Field error={errors.password}>
                <div className="relative"><input className={`input-dark pr-14 ${errors.password ? 'border-red-500' : ''}`} onChange={(event) => updateField('password', event.target.value)} type={showPassword ? 'text' : 'password'} placeholder="Contraseña" value={form.password} /><button type="button" onClick={() => setShowPassword(!showPassword)} className="absolute right-4 top-1/2 -translate-y-1/2 text-lg">{showPassword ? '🙈' : '👁️'}</button></div>
              </Field>
              {mode === 'login' && <div className="text-right"><button type="button" className="text-sm font-bold text-violet-300 hover:text-violet-100">¿Olvidaste tu contraseña?</button></div>}
              {submitError && <div className="rounded-2xl border border-red-500/40 bg-red-500/10 p-3 text-sm text-red-200">{submitError}</div>}
              <Button className="w-full" type="submit" disabled={loading}>{loading ? 'Procesando...' : mode === 'login' ? 'Iniciar sesión' : 'Crear cuenta'}</Button>
            </form>
            {mode === 'login' && <><div className="my-6 flex items-center gap-3 text-sm text-gray-500"><span className="h-px flex-1 bg-ticket-border" />o continuá con<span className="h-px flex-1 bg-ticket-border" /></div><div className="grid grid-cols-2 gap-3"><Button type="button" variant="secondary" disabled>G Google</Button><Button type="button" variant="secondary" disabled>f Facebook</Button></div></>}
          </div>
        </div>
      </main>
    </div>
  )
}

function Field({ children, error }) {
  return (
    <label className="block">
      {children}
      {error && <span className="mt-1 block text-sm text-red-400">{error}</span>}
    </label>
  )
}
