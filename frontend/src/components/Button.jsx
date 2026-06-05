import { Link } from 'react-router-dom'

const styles = {
  primary: 'bg-ticket-purple text-white shadow-glow hover:bg-ticket-purple2',
  secondary: 'border border-ticket-border bg-ticket-card2 text-white hover:border-ticket-purple2/70 hover:bg-ticket-purple/10',
  danger: 'border border-red-500/40 bg-red-500/10 text-red-300 hover:bg-red-500/20',
}

export default function Button({ children, variant = 'primary', className = '', to, ...props }) {
  const classes = `inline-flex items-center justify-center gap-2 rounded-2xl px-5 py-3 font-bold transition ${styles[variant]} ${className}`
  if (to) return <Link to={to} className={classes}>{children}</Link>
  return <button className={classes} {...props}>{children}</button>
}
