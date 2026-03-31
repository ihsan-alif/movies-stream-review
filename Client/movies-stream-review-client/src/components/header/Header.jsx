import {useState} from 'react'
import Button from 'react-bootstrap/Button'
import Container from 'react-bootstrap/Container'
import Nav from 'react-bootstrap/Nav'
import Navbar from 'react-bootstrap/Navbar'
import {useNavigate, NavLink, Link} from 'react-router-dom'

const Header = () => {
    const navigate = useNavigate();
    const [auth, setAuth] = useState(false);

    return (
        <Navbar className='shadow-sm' bg="dark" variant="dark" stick='top' expand="lg">
            <Container>
                <Navbar.Brand>Movies Stream</Navbar.Brand>
                <Navbar.Toggle aria-controls='main-navbar-nav'/>
                <Navbar.Collapse>
                    <Nav className="me-auto">
                        <Nav.Link as={NavLink} to="/">
                            Home
                        </Nav.Link>
                        <Nav.Link as={NavLink} to="/movies/recommended">
                            Recommended
                        </Nav.Link>
                    </Nav>
                    <Nav className="ms-auto align-items-center">
                        {auth ? (
                            <>
                                <span>
                                    Hello, <strong>Name</strong>
                                </span>
                                <Button variant="outline-light" size="sm">
                                    Logout
                                </Button>
                            </>
                        ):(
                            <>
                                <Button className="me-2" variant="outline-info" size="sm" onClick={() => navigate("/users/login")}>
                                    Login
                                </Button>
                                <Button variant="info" size="sm" onClick={() => navigate("/users/register")}>
                                    Register
                                </Button>
                            </>
                        )} 
                    </Nav>
                </Navbar.Collapse>
            </Container>
        </Navbar>
    )
}

export default Header;