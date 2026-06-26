import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import AdminStatCard from '../../admin/AdminStatCard.jsx'
import { getId, money, normalizeStatsEvents } from '../../admin/adminData.js'
import Badge from '../../components/Badge.jsx'
import { getAdminEventStats, getAdminStatsSummary } from '../../services/api.js'

const adminStatsErrorMessage = (error) => {
  if (error?.status === 401 || error?.status === 403) return 'No tenés permisos de administrador para ver reportes.'
  return error?.message || 'No pudimos cargar los reportes.'
}

const formatPercent = (value) => `${new Intl.NumberFormat('es-AR', { maximumFractionDigits: 2 }).format(Number(value ?? 0))}%`

export default function AdminReports() {
  const [summary, setSummary] = useState(null)
  const [events, setEvents] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    let isMounted = true
    Promise.all([getAdminStatsSummary(), getAdminEventStats()])
      .then(([summaryData, eventStats]) => {
        if (!isMounted) return
        setSummary(summaryData)
        setEvents(normalizeStatsEvents(eventStats))
      })
      .catch((e) => {
        if (isMounted) setError(adminStatsErrorMessage(e))
      })
      .finally(() => {
        if (isMounted) setLoading(false)
      })
    return () => { isMounted = false }
  }, [])

  return <main className="px-5 py-8 lg:px-8">
    <p className="text-sm font-black uppercase tracking-[0.28em] text-ticket-purple2">Reportes</p>
    <h1 className="mt-2 text-4xl font-black">Estadísticas admin</h1>
    <p className="mt-3 max-w-3xl text-gray-400">Métricas de usuarios, eventos, entradas activas y recaudación estimada.</p>

    {error && <p className="mt-4 rounded-2xl border border-red-500/30 bg-red-500/10 p-4 text-red-200">{error}</p>}

    {loading ? <div className="mt-8 rounded-3xl border border-ticket-border bg-ticket-card p-8 text-gray-400">Cargando reportes...</div> : !error && <>
      {summary && <section className="mt-8 grid gap-4 sm:grid-cols-2 xl:grid-cols-5">
        <AdminStatCard label="Usuarios totales" value={summary.total_users ?? 0} helper={`${summary.admin_users ?? 0} admin · ${summary.client_users ?? 0} clientes`} />
        <AdminStatCard label="Eventos activos" value={summary.active_events ?? 0} helper={`${summary.total_events ?? 0} eventos totales`} />
        <AdminStatCard label="Entradas activas" value={summary.active_tickets ?? 0} helper={`${summary.cancelled_tickets ?? 0} canceladas`} />
        <AdminStatCard label="Ocupación general" value={formatPercent(summary.occupancy_rate_percent)} helper={`${summary.available_capacity ?? 0} cupos disponibles`} />
        <AdminStatCard label="Ingresos estimados" value={money.format(Number(summary.estimated_revenue ?? 0))} helper="Tickets ACTIVE por precio actual" />
      </section>}

      <section className="mt-8">
        <h2 className="text-2xl font-black">Eventos</h2>
        <div className="mt-4 grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          {events.length ? events.map((event) => <Link key={getId(event)} to={`/admin/reportes/${event.event_id}`} className="rounded-3xl border border-ticket-border bg-ticket-card p-5 transition hover:border-ticket-purple2/60 hover:bg-ticket-card2">
            <div className="flex flex-wrap items-center gap-2">
              <Badge>{event.category}</Badge>
              <Badge tone={event.active ? 'purple' : 'red'}>{event.active ? 'Activo' : 'Inactivo'}</Badge>
            </div>
            <h3 className="mt-4 text-2xl font-black">{event.title}</h3>
            <p className="mt-2 text-gray-400">📍 {event.location}</p>
            <dl className="mt-5 grid grid-cols-2 gap-3 text-sm">
              <div><dt className="text-gray-500">Entradas activas</dt><dd className="font-black">{event.active_tickets}</dd></div>
              <div><dt className="text-gray-500">Cupo disponible</dt><dd className="font-black">{event.available_capacity}</dd></div>
              <div><dt className="text-gray-500">Ocupación</dt><dd className="font-black">{formatPercent(event.occupancy_rate_percent)}</dd></div>
              <div><dt className="text-gray-500">Ingresos estimados</dt><dd className="font-black">{money.format(event.estimated_revenue)}</dd></div>
            </dl>
            <p className="mt-5 text-sm font-bold text-ticket-purple2">Ver métricas →</p>
          </Link>) : <div className="rounded-3xl border border-ticket-border bg-ticket-card p-8 text-gray-400">No hay estadísticas de eventos disponibles.</div>}
        </div>
      </section>
    </>}
  </main>
}
