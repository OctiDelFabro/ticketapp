import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { getId, normalizeEvents } from '../../admin/adminData.js'
import Badge from '../../components/Badge.jsx'
import { getAdminEvents } from '../../services/api.js'
import { formatEventDate } from '../../utils/events.js'

export default function AdminReports() {
  const [events, setEvents] = useState([]), [loading, setLoading] = useState(true), [error, setError] = useState('')
  useEffect(() => { getAdminEvents().then((data)=>setEvents(normalizeEvents(data))).catch((e)=>setError(e.message)).finally(()=>setLoading(false)) }, [])
  return <main className="px-5 py-8 lg:px-8"><p className="text-sm font-black uppercase tracking-[0.28em] text-ticket-purple2">Reportes</p><h1 className="mt-2 text-4xl font-black">Elegí un evento</h1>{error && <p className="mt-4 text-red-300">{error}</p>}<div className="mt-8 grid gap-4 md:grid-cols-2 xl:grid-cols-3">{loading ? <div className="rounded-3xl border border-ticket-border bg-ticket-card p-8 text-gray-400">Cargando eventos...</div> : events.length ? events.map(e => <Link key={getId(e)} to={`/admin/reportes/${getId(e)}`} className="rounded-3xl border border-ticket-border bg-ticket-card p-5 transition hover:border-ticket-purple2/60 hover:bg-ticket-card2"><Badge>{e.category}</Badge><h2 className="mt-4 text-2xl font-black">{e.title}</h2><p className="mt-2 text-gray-400">{formatEventDate(e.start_date)} · {e.location}</p><p className="mt-4 text-sm font-bold text-ticket-purple2">Ver métricas →</p></Link>) : <div className="rounded-3xl border border-ticket-border bg-ticket-card p-8 text-gray-400">No hay eventos disponibles.</div>}</div></main>
}
