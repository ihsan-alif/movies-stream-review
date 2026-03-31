import { useState } from 'react'
import './App.css'
import Home from './components/home/Home'
import Header from './components/header/Header'

function App() {
  const [count, setCount] = useState(0)

  return (
    <>
      <Header/>
      <Home/>
    </>
  )
}

export default App
