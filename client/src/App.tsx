import { BrowserRouter, Routes, Route } from "react-router-dom"
import LoginPage from "./pages/LoginPage"
import FeedPage from "./pages/FeedPage"

export default function App() {
  return (

    <BrowserRouter >
      <Routes>
        <Route path='/' element={<LoginPage />} />
        <Route path='/:email' element={<FeedPage />} />
      </Routes>
    </BrowserRouter>
  )
}