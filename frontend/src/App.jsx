import { useEffect, useState } from 'react'
import { Route, Routes } from 'react-router-dom'
import BottomNav from './components/BottomNav.jsx'
import Footer from './components/Footer.jsx'
import Header from './components/Header.jsx'
import Checkout from './pages/Checkout.jsx'
import EventDetail from './pages/EventDetail.jsx'
import Home from './pages/Home.jsx'
import Login from './pages/Login.jsx'
import MyTickets from './pages/MyTickets.jsx'
import { isAuthenticated } from './utils/auth.js'
import { getStoredCart, saveStoredCart } from './utils/cart.js'

export default function App() {
  const [searchQuery, setSearchQuery] = useState('')
  const [cartItem, setCartItem] = useState(getStoredCart)
  const [isLoggedIn, setIsLoggedIn] = useState(isAuthenticated)

  useEffect(() => {
    saveStoredCart(cartItem)
  }, [cartItem])

  return (
    <div className="app-shell">
      <Header
        cartCount={cartItem?.quantity ?? 0}
        isLoggedIn={isLoggedIn}
        searchQuery={searchQuery}
        setSearchQuery={setSearchQuery}
      />
      <Routes>
        <Route path="/" element={<Home searchQuery={searchQuery} />} />
        <Route path="/evento/:id" element={<EventDetail onAddToCart={setCartItem} />} />
        <Route path="/login" element={<Login setIsLoggedIn={setIsLoggedIn} />} />
        <Route path="/checkout" element={<Checkout cartItem={cartItem} setCartItem={setCartItem} />} />
        <Route path="/mis-entradas" element={<MyTickets isLoggedIn={isLoggedIn} />} />
      </Routes>
      <Footer />
      <BottomNav />
    </div>
  )
}
