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
import StreamMovie from './components/stream/StreamMovie'
import {Route, Routes, useNavigate} from 'react-router-dom'
import axiosClient from './api/axiosConfig'
import useAuth from './hooks/useAuth'

function App() {
  const navigate = useNavigate();
  const {auth, setAuth} = useAuth();

  const updateMovieReview = (imdb_id) => {
    navigate(`/movies/updatereview/${imdb_id}`);
  }

  const handleLogout = async () => {
    try {
        const response = await axiosClient.post("/users/logout", {user_id: auth.user_id});
        console.log(response.data);
        setAuth(null);
        // localStorage.removeItem('user');
        console.log('User logged out');
    } catch(error) {
        console.error('Error logging out:', error)
    }
  }

  return (
    <>
      <Header handleLogout={handleLogout} />
      <Routes path="/" element={<Layout/>}>
        <Route path="/" element={<Home updateMovieReview={updateMovieReview}/>}></Route>
        <Route path="/users/register" element={<Register/>}></Route>
        <Route path="/users/login" element={<Login/>}></Route>
        <Route element={<RequiredAuth/>}>
          <Route path="/movies/recommended" element={<Recommended/>}></Route>
          <Route path="/movies/updatereview/:imdb_id" element={<Review/>}></Route>
          <Route path="/movies/stream/:yt_id" element={<StreamMovie/>}></Route>
        </Route>
      </Routes>
    </>
  )
}

export default App;
