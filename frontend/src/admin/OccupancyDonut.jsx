export default function OccupancyDonut({ value = 0 }) {
  const pct = Math.max(0, Math.min(100, Math.round(value)))
  const color = pct >= 85 ? '#ef4444' : pct >= 60 ? '#22c55e' : '#8b5cf6'
  return <div className="grid place-items-center"><div className="grid h-32 w-32 place-items-center rounded-full" style={{ background: `conic-gradient(${color} ${pct}%, #1f1f2e 0)` }}><div className="grid h-24 w-24 place-items-center rounded-full bg-ticket-card"><span className="text-2xl font-black">{pct}%</span></div></div></div>
}
