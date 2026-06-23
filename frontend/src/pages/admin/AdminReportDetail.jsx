import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import AdminStatCard from '../../admin/AdminStatCard.jsx'
import OccupancyDonut from '../../admin/OccupancyDonut.jsx'
import { money, normalizeStatsEvent } from '../../admin/adminData.js'
import Badge from '../../components/Badge.jsx'
import Button from '../../components/Button.jsx'
import { getAdminEventReport } from '../../services/api.js'
import { formatEventDate, formatEventTime } from '../../utils/events.js'

const reportErrorMessage = (error) => {
  if (error?.status === 404) return 'Reporte no encontrado.'
  if (error?.status === 401 || error?.status === 403) return 'No tenés permisos de administrador para ver este reporte.'
  return error?.message || 'No pudimos cargar el reporte.'
}

const formatPercent = (value) => `${new Intl.NumberFormat('es-AR', { maximumFractionDigits: 2 }).format(Number(value ?? 0))}%`

export default function AdminReportDetail() {
  const { id } = useParams()
  const [event, setEvent] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    let isMounted = true
    getAdminEventReport(id)
      .then((report) => {
        if (isMounted) setEvent(normalizeStatsEvent(report))
      })
      .catch((e) => {
        if (isMounted) setError(reportErrorMessage(e))
      })
      .finally(() => {
        if (isMounted) setLoading(false)
      })
    return () => { isMounted = false }
  }, [id])

  if (loading) return <main className="px-5 py-8 lg:px-8"><div className="rounded-3xl border border-ticket-border bg-ticket-card p-10 text-center text-gray-400">Cargando reporte...</div></main>
  if (!event) return <main className="px-5 py-8 lg:px-8"><Link to="/admin/reportes" className="font-bold text-gray-300 hover:text-white">← Volver a reportes</Link><p className="mt-6 rounded-2xl border border-red-500/30 bg-red-500/10 p-4 text-red-200">{error || 'Reporte no encontrado.'}</p></main>

  const occupancy = event.occupancy_rate_percent

  return <main className="px-5 py-8 lg:px-8">
    <Link to="/admin/reportes" className="font-bold text-gray-300 hover:text-white">← Volver a reportes</Link>
    <section className="mt-6 rounded-[2rem] border border-ticket-purple2/30 bg-ticket-card p-6 shadow-glow">
      <div className="flex flex-col gap-5 lg:flex-row lg:items-start lg:justify-between">
        <div>
          <div className="flex flex-wrap gap-3">
            <Badge>{event.category}</Badge>
            <Badge tone={event.active ? 'purple' : 'red'}>{event.active ? 'Activo' : 'Inactivo'}</Badge>
            <Badge tone={occupancy >= 85 ? 'red' : 'purple'}>{formatPercent(occupancy)} ocupación</Badge>
          </div>
          <h1 className="mt-4 text-4xl font-black">{event.title}</h1>
          {event.description && <p className="mt-3 max-w-3xl text-gray-300">{event.description}</p>}
          <p className="mt-3 text-gray-300">📅 {formatEventDate(event.start_date)} · ⏰ {formatEventTime(event.start_date)} · 📍 {event.location}</p>
          <p className="mt-2 text-sm text-gray-500">Duración: {event.duration_minutes ?? 0} minutos · ID evento: {event.event_id}</p>
        </div>
        <Button variant="secondary" onClick={() => window.print()}>Exportar reporte</Button>
      </div>
    </section>

    <div className="mt-6 grid gap-4 md:grid-cols-2 xl:grid-cols-4">
      <AdminStatCard label="Entradas activas" value={event.active_tickets} />
      <AdminStatCard label="Entradas canceladas" value={event.cancelled_tickets} />
      <AdminStatCard label="Tickets totales" value={event.total_tickets} />
      <AdminStatCard label="Cupo disponible" value={event.available_capacity} helper={`Capacidad total ${event.capacity}`} />
      <AdminStatCard label="Ocupación" value={formatPercent(occupancy)} />
      <AdminStatCard label="Ingresos estimados" value={money.format(event.estimated_revenue)} helper="Calculados por el backend" />
      <AdminStatCard label="Precio entrada" value={money.format(event.price)} />
    </div>

    <div className="mt-6 grid gap-6 xl:grid-cols-[1fr_320px]">
      <section className="rounded-3xl border border-ticket-border bg-ticket-card p-5">
        <h2 className="text-xl font-black">Ventas y compradores</h2>
        <div className="mt-6 rounded-2xl border border-dashed border-ticket-border p-10 text-center text-gray-400">El backend actual todavía no expone compradores ni ventas por semana.</div>
      </section>
      <section className="rounded-3xl border border-ticket-border bg-ticket-card p-5">
        <h2 className="text-xl font-black">Ocupación</h2>
        <div className="mt-6"><OccupancyDonut value={occupancy} /></div>
      </section>
    </div>
  </main>
}
