import { useEffect, useMemo, useState } from 'react'
import Button from '../components/Button.jsx'
import EventImage from '../components/EventImage.jsx'
import { eventCategories } from '../constants/eventCategories.js'
import { isSupportedEventImageURL } from '../utils/events.js'
import { toInputDate } from './adminData.js'

const initial = { title: '', description: '', image_url: '', category: '', location: '', start_date: '', duration_minutes: 120, capacity: 1, active: true, price: '' }

export default function AdminEventForm({ event, onClose, onSubmit, loading }) {
  const [form, setForm] = useState(initial)
  const [errors, setErrors] = useState({})
  useEffect(() => { setForm(event ? { ...initial, ...event, start_date: toInputDate(event.start_date), price: event.price ?? event.ticket_price ?? '' } : initial) }, [event])
  const set = (k, v) => { setForm((f) => ({ ...f, [k]: v })); setErrors((e) => ({ ...e, [k]: '' })) }
  const imageURL = useMemo(() => form.image_url.trim(), [form.image_url])
  const canPreviewImage = isSupportedEventImageURL(imageURL)
  const submit = (e) => {
    e.preventDefault(); const next = {}
    ;['title', 'description', 'category', 'location', 'start_date'].forEach((k) => { if (!String(form[k] ?? '').trim()) next[k] = 'Campo requerido.' })
    if (!isSupportedEventImageURL(form.image_url)) next.image_url = 'Usá una URL http://, https:// o una ruta local que empiece con /.'
    if (Number(form.duration_minutes) <= 0) next.duration_minutes = 'Debe ser mayor a 0.'
    if (Number(form.capacity) <= 0) next.capacity = 'Debe ser mayor a 0.'
    if (form.price !== '' && Number(form.price) < 0) next.price = 'No puede ser negativo.'
    if (Number.isNaN(new Date(form.start_date).getTime())) next.start_date = 'Fecha inválida.'
    setErrors(next); if (Object.keys(next).length) return
    const payload = { title: form.title.trim(), description: form.description.trim(), image_url: form.image_url.trim(), category: form.category.trim(), location: form.location.trim(), start_date: new Date(form.start_date).toISOString(), duration_minutes: Number(form.duration_minutes), capacity: Number(form.capacity), active: Boolean(form.active) }
    if (form.price !== '') payload.price = Number(form.price)
    onSubmit(payload)
  }
  return <div className="fixed inset-0 z-50 overflow-y-auto bg-black/70 p-4 backdrop-blur"><form onSubmit={submit} className="mx-auto my-8 max-w-3xl rounded-3xl border border-ticket-purple2/30 bg-ticket-card p-6 shadow-glow"><div className="flex items-start justify-between gap-4"><div><p className="text-xs font-black uppercase tracking-[0.28em] text-ticket-purple2">{event ? 'Editar evento' : 'Nuevo evento'}</p><h2 className="mt-1 text-3xl font-black">Datos del evento</h2></div><button type="button" onClick={onClose} className="text-2xl text-gray-400 hover:text-white">×</button></div><div className="mt-6 grid gap-4 sm:grid-cols-2"><Field label="Título" error={errors.title}><input className="input-dark" value={form.title} onChange={(e)=>set('title', e.target.value)} /></Field><Field label="Categoría" error={errors.category}><select className="input-dark" value={form.category} onChange={(e)=>set('category', e.target.value)}><option value="">Seleccionar categoría</option>{eventCategories.map((item) => <option key={item} value={item}>{item}</option>)}</select></Field><Field label="Lugar" error={errors.location}><input className="input-dark" value={form.location} onChange={(e)=>set('location', e.target.value)} /></Field><Field label="Imagen URL" error={errors.image_url} help="Pegá una URL directa de imagen (http:// o https://) o una ruta local existente, por ejemplo /events/default.jpg."><input className="input-dark" placeholder="https://sitio.com/imagen.jpg" value={form.image_url} onChange={(e)=>set('image_url', e.target.value)} /></Field><Field label="Fecha y hora" error={errors.start_date}><input type="datetime-local" className="input-dark" value={form.start_date} onChange={(e)=>set('start_date', e.target.value)} /></Field><Field label="Duración (min)" error={errors.duration_minutes}><input type="number" min="1" className="input-dark" value={form.duration_minutes} onChange={(e)=>set('duration_minutes', e.target.value)} /></Field><Field label="Cupo" error={errors.capacity}><input type="number" min="1" className="input-dark" value={form.capacity} onChange={(e)=>set('capacity', e.target.value)} /></Field><Field label="Precio opcional" error={errors.price}><input type="number" min="0" className="input-dark" value={form.price} onChange={(e)=>set('price', e.target.value)} /></Field><label className="flex items-center gap-3 rounded-2xl border border-ticket-border bg-ticket-card2 px-4 py-3"><input type="checkbox" checked={form.active} onChange={(e)=>set('active', e.target.checked)} /> Activo</label><div className="rounded-2xl border border-ticket-border bg-ticket-card2 p-3"><p className="mb-2 text-sm font-bold text-gray-300">Preview de imagen</p>{canPreviewImage ? <EventImage src={imageURL} alt={form.title || 'Preview de imagen del evento'} className="h-36 w-full rounded-xl object-cover" /> : <div className="grid h-36 place-items-center rounded-xl border border-dashed border-ticket-border text-center text-sm text-gray-400">Ingresá una URL http://, https:// o una ruta local válida.</div>}<p className="mt-2 text-xs text-gray-500">Si la imagen no carga, se muestra la imagen fallback automáticamente.</p></div><Field className="sm:col-span-2" label="Descripción" error={errors.description}><textarea rows="4" className="input-dark" value={form.description} onChange={(e)=>set('description', e.target.value)} /></Field></div><div className="mt-6 flex justify-end gap-3"><Button type="button" variant="secondary" onClick={onClose}>Cancelar</Button><Button disabled={loading}>{loading ? 'Guardando...' : 'Guardar'}</Button></div></form></div>
}
function Field({ label, error, help, children, className='' }) { return <label className={`block ${className}`}><span className="mb-2 block text-sm font-bold text-gray-300">{label}</span>{children}{help && <span className="mt-1 block text-xs text-gray-500">{help}</span>}{error && <span className="mt-1 block text-sm text-red-400">{error}</span>}</label> }
