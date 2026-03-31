import { useState, useEffect } from "react";
import axiosClient from "../../api/axiosConfig";
import Movies from '../movies/Movies';

const Home = () => {
    const [movies, setMovies] = useState([]);
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState();

    useEffect(() => {
        const fetchMovies = async () => {
            setLoading(true)
            setMessage("")
        try {
            const response = await axiosClient.get('/movies');

            // === DEBUGGING PENTING ===
                console.log("Full response:", response);
                console.log("response.data:", response.data);
                console.log("Tipe data:", typeof response.data);
                console.log("Apakah array?", Array.isArray(response.data));

            setMovies(response.data);
            if (response.data.length === 0) {
                setMessage('There are currently no movies available')
            }
        } catch(error) {
            console.error('Error fetching movies', error)
            setMessage('Error fetching movies')
        } finally {
            setLoading(false)
        }
    }
    fetchMovies();
    }, [])

    return (
        <>
            {loading ? (
                <h2>Loading...</h2>
            ): (
                <Movies movies ={movies} message ={message}/>
            )}
        </>
    );
};

export default Home;