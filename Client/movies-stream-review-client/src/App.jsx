import { useState } from 'react'
import './App.css'
import Home from './components/home/Home'
import Header from './components/header/Header'
import Register from './components/register/Register'
import Login from './components/login/Login'
import {Route, Routes, useNavigate} from 'react-router-dom'
import RequiredAuth from './components/RequiredAuth'
import Layout from './components/Layout'
import Recommended from './components/recommended/Recommended'

function App() {

  return (
    <>
      <Header/>
      <Routes path="/" element={<Layout/>}>
        <Route path="/" element={<Home/>}></Route>
        <Route path="/users/register" element={<Register/>}></Route>
        <Route path="/users/login" element={<Login/>}></Route>
        <Route element={<RequiredAuth/>}>
          <Route path="/movies/recommended" element={<Recommended/>}></Route>
        </Route>
      </Routes>
    </>
  )
}

export default App
