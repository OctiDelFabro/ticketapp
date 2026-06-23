import { useEffect, useState } from 'react'
import { Route, Routes, useLocation, useNavigate } from 'react-router-dom'
import BottomNav from './components/BottomNav.jsx'
import Footer from './components/Footer.jsx'
import Header from './components/Header.jsx'
import Checkout from './pages/Checkout.jsx'
import EventDetail from './pages/EventDetail.jsx'
import Home from './pages/Home.jsx'
import Login from './pages/Login.jsx'
import MyTickets from './pages/MyTickets.jsx'
import AdminLayout from './admin/AdminLayout.jsx'
import AdminRoute from './admin/AdminRoute.jsx'
import AdminEvents from './pages/admin/AdminEvents.jsx'
import AdminReports from './pages/admin/AdminReports.jsx'
import AdminReportDetail from './pages/admin/AdminReportDetail.jsx'
import { clearAuthSession, isAuthenticated } from './utils/auth.js'
import { getStoredCart, saveStoredCart } from './utils/cart.js'

export default function App() {
  const [searchQuery, setSearchQuery] = useState('')
  const [cartItem, setCartItem] = useState(getStoredCart)
  const [isLoggedIn, setIsLoggedIn] = useState(isAuthenticated)
  const navigate = useNavigate()
  const location = useLocation()
  const isAdminPath = location.pathname.startsWith('/admin')

  const handleLogout = () => {
    clearAuthSession()
    setIsLoggedIn(false)
    setCartItem(null)
    navigate('/login')
  }

  useEffect(() => {
    saveStoredCart(cartItem)
  }, [cartItem])

  return (
    <div className="app-shell">
      {!isAdminPath && <Header
        cartCount={cartItem?.quantity ?? 0}
        isLoggedIn={isLoggedIn}
        searchQuery={searchQuery}
        setSearchQuery={setSearchQuery}
        onLogout={handleLogout}
      />}
      <Routes>
        <Route path="/" element={<Home searchQuery={searchQuery} />} />
        <Route path="/evento/:id" element={<EventDetail onAddToCart={setCartItem} />} />
        <Route path="/login" element={<Login setIsLoggedIn={setIsLoggedIn} />} />
        <Route path="/checkout" element={<Checkout cartItem={cartItem} setCartItem={setCartItem} />} />
        <Route path="/mis-entradas" element={<MyTickets isLoggedIn={isLoggedIn} />} />
        <Route path="/admin" element={<AdminRoute><AdminLayout /></AdminRoute>}>
          <Route path="eventos" element={<AdminEvents />} />
          <Route path="eventos/nuevo" element={<AdminEvents />} />
          <Route path="eventos/:id/editar" element={<AdminEvents />} />
          <Route path="reportes" element={<AdminReports />} />
          <Route path="reportes/:id" element={<AdminReportDetail />} />
        </Route>
      </Routes>
      {!isAdminPath && <Footer />}
      {!isAdminPath && <BottomNav isLoggedIn={isLoggedIn} onLogout={handleLogout} />}
    </div>
  )
}
