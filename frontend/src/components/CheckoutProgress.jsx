const steps = ['Entradas', 'Tus datos', 'Pago', 'Confirmación']

export default function CheckoutProgress({ currentStep }) {
  return (
    <div className="glass-card rounded-3xl p-4 sm:p-6">
      <div className="flex items-center">
        {steps.map((step, index) => {
          const number = index + 1
          const complete = number < currentStep
          const active = number === currentStep
          return (
            <div key={step} className="flex flex-1 items-center last:flex-none">
              <div className="flex flex-col items-center gap-2">
                <div className={`grid h-10 w-10 place-items-center rounded-full border text-sm font-black ${complete || active ? 'border-ticket-purple2 bg-ticket-purple text-white shadow-glow' : 'border-ticket-border bg-ticket-card2 text-gray-500'}`}>{complete ? '✓' : number}</div>
                <span className={`hidden text-xs font-bold sm:block ${active ? 'text-white' : 'text-gray-500'}`}>{step}</span>
              </div>
              {index < steps.length - 1 && <div className={`mx-2 h-px flex-1 sm:mx-4 ${complete ? 'bg-ticket-purple2' : 'bg-ticket-border'}`} />}
            </div>
          )
        })}
      </div>
    </div>
  )
}
