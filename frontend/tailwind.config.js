/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{js,jsx}'],
  theme: {
    extend: {
      fontFamily: {
        sans: ['DM Sans', 'ui-sans-serif', 'system-ui'],
      },
      colors: {
        ticket: {
          bg: '#0a0a0f',
          alt: '#090911',
          card: '#111118',
          card2: '#151522',
          border: '#1f1f2e',
          purple: '#7c3aed',
          purple2: '#8b5cf6',
        },
      },
      boxShadow: {
        glow: '0 24px 80px rgba(124, 58, 237, 0.22)',
        soft: '0 20px 60px rgba(0,0,0,0.35)',
      },
    },
  },
  plugins: [],
}
