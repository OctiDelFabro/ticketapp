export default function Badge({ children, tone = 'purple', className = '' }) {
  const tones = {
    purple: 'border-ticket-purple2/40 bg-ticket-purple/20 text-violet-200',
    red: 'border-red-500/40 bg-red-500/15 text-red-200',
    amber: 'border-amber-500/40 bg-amber-500/15 text-amber-200',
    dark: 'border-ticket-border bg-ticket-card2 text-gray-300',
  }
  return <span className={`inline-flex rounded-full border px-3 py-1 text-xs font-black uppercase tracking-[0.16em] ${tones[tone]} ${className}`}>{children}</span>
}
