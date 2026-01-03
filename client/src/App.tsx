import { BrowserRouter, Routes, Route } from "react-router-dom"
import LoginPage from "./pages/LoginPage"
import FeedPage from "./pages/FeedPage"
import CreatePage from "./pages/CreatePage"

export default function App() {
  return (

    <BrowserRouter >
      <Routes>
        <Route path='/' element={<LoginPage />} />
        <Route path='/create' element={<CreatePage />} />
        <Route path='/:email' element={<FeedPage />} />
      </Routes>
    </BrowserRouter>
  )
}