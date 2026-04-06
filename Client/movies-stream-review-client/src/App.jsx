import { useState } from 'react'
import './App.css'
import Home from './components/home/Home'
import Header from './components/header/Header'
import Register from './components/register/Register'
import Login from './components/login/Login'
import RequiredAuth from './components/RequiredAuth'
import Layout from './components/Layout'
import Recommended from './components/recommended/Recommended'
import Review from './components/review/Review'
import {Route, Routes, useNavigate} from 'react-router-dom'

function App() {
  const navigate = useNavigate();

  const updateMovieReview = (imdb_id) => {
    navigate(`/movies/updatereview/${imdb_id}`);
  }

  return (
    <>
      <Header/>
      <Routes path="/" element={<Layout/>}>
        <Route path="/" element={<Home updateMovieReview={updateMovieReview}/>}></Route>
        <Route path="/users/register" element={<Register/>}></Route>
        <Route path="/users/login" element={<Login/>}></Route>
        <Route element={<RequiredAuth/>}>
          <Route path="/movies/recommended" element={<Recommended/>}></Route>
          <Route path="/movies/updatereview/:imdb_id" element={<Review/>}></Route>
        </Route>
      </Routes>
    </>
  )
}

export default App;
