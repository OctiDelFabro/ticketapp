const styles = {
  success: {
    icon: '✓',
    classes: 'border-emerald-500/40 bg-emerald-500/10 text-emerald-100',
    iconClasses: 'bg-emerald-500/20 text-emerald-200',
  },
  error: {
    icon: '!',
    classes: 'border-red-500/40 bg-red-500/10 text-red-100',
    iconClasses: 'bg-red-500/20 text-red-200',
  },
  info: {
    icon: 'i',
    classes: 'border-ticket-purple2/40 bg-ticket-purple/10 text-violet-100',
    iconClasses: 'bg-ticket-purple/25 text-violet-200',
  },
}

export default function AlertMessage({ type = 'info', message, onClose }) {
  if (!message) return null

  const style = styles[type] ?? styles.info

  return (
    <div className={`flex items-start gap-3 rounded-2xl border p-4 text-sm font-bold shadow-soft ${style.classes}`} role={type === 'error' ? 'alert' : 'status'}>
      <span className={`grid h-7 w-7 shrink-0 place-items-center rounded-full text-sm font-black ${style.iconClasses}`}>{style.icon}</span>
      <p className="flex-1 leading-6">{message}</p>
      {onClose && (
        <button aria-label="Cerrar mensaje" className="rounded-full px-2 text-lg leading-none opacity-70 transition hover:opacity-100" onClick={onClose} type="button">
          ×
        </button>
      )}
    </div>
  )
}
