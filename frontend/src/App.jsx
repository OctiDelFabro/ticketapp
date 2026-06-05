import { Route, Routes } from 'react-router-dom'
import BottomNav from './components/BottomNav.jsx'
import Checkout from './pages/Checkout.jsx'
import EventDetail from './pages/EventDetail.jsx'
import Home from './pages/Home.jsx'
import Login from './pages/Login.jsx'
import MyTickets from './pages/MyTickets.jsx'

export default function App() {
  return (
    <>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/evento/:id" element={<EventDetail />} />
        <Route path="/login" element={<Login />} />
        <Route path="/checkout" element={<Checkout />} />
        <Route path="/mis-entradas" element={<MyTickets />} />
      </Routes>
      <BottomNav />
    </>
  )
}
