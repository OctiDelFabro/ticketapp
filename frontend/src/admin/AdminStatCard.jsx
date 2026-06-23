export default function AdminStatCard({ label, value, helper }) {
  return <div className="rounded-3xl border border-ticket-purple2/20 bg-ticket-card p-5 shadow-soft"><p className="text-sm font-bold text-gray-400">{label}</p><p className="mt-2 text-3xl font-black">{value}</p>{helper && <p className="mt-2 text-sm text-gray-500">{helper}</p>}</div>
}
