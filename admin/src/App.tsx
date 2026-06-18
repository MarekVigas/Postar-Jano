import React, {useEffect} from 'react';
import './App.css';
import {BrowserRouter, Route, Routes, Link, Navigate} from 'react-router-dom';
import EventList from "./components/EventList";
import RegistrationList from "./components/RegistrationList";
import Login from "./components/Login";
import {useAPIClient, useAppDispatch, useAppState} from "./AppContext";
import {Navbar, Nav} from "react-bootstrap";
import 'bootstrap/dist/css/bootstrap.min.css';


const AppContent: React.FC = () => {
    const state = useAppState()
    const dispatch = useAppDispatch()
    const apiClient = useAPIClient()

    useEffect(() => {
        apiClient.user.me().then(result => {
            dispatch({type: 'INIT_DONE', isAuthenticated: result !== null})
        })
    }, []) // eslint-disable-line react-hooks/exhaustive-deps

    if (state.isLoading) {
        return null
    }

    if (!state.isAuthenticated) {
        return <Login/>
    }

    const handleSignOut = () => {
        apiClient.user.signOut().finally(() => {
            dispatch({type: 'SIGN_OUT'})
        })
    }

    return (
        <div className="wrapper">
            <Navbar bg="light">
                <Navbar.Brand>Leto 2024</Navbar.Brand>
                <Nav style={{display: 'flex', width: '100%'}}>
                    <Nav.Link>
                        <Link to="/events">
                            Akcie
                        </Link>
                    </Nav.Link>
                    <Nav.Link>
                        <Link to="/registrations">
                            Prihlaseni
                        </Link>
                    </Nav.Link>
                    <Nav.Link style={{marginInlineStart: 'auto'}}>
                        <button onClick={handleSignOut}>Odhlasit sa</button>
                    </Nav.Link>
                </Nav>
            </Navbar>
            <Routes>
                <Route path="/events" element={<EventList/>}/>
                <Route path="/registrations/:event" element={<RegistrationList/>}/>
                <Route path="/registrations" element={<RegistrationList/>}/>
                <Route path="*" element={<Navigate to="/events" replace/>}/>
            </Routes>
        </div>
    );
}

const App: React.FC = () => {
    return (
        <BrowserRouter basename="/admin">
            <AppContent/>
        </BrowserRouter>
    );
}

export default App;
