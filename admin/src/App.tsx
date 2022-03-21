import React, {useEffect} from 'react';
import './App.css';
import {BrowserRouter, Route, Switch, Link} from 'react-router-dom';
import EventList from "./components/EventList";
import RegistrationList from "./components/RegistrationList";
import Login from "./components/Login";
import {useAppDispatch, useAppState} from "./AppContext";
import {Navbar, Nav} from "react-bootstrap";
import 'bootstrap/dist/css/bootstrap.min.css';
import jwtDecode, {JwtPayload} from "jwt-decode";


const App: React.FC = () => {
    const state = useAppState()
    const dispatch = useAppDispatch()

    useEffect(
        ()=>{
            if (!state.token) return
            let decodedToken = jwtDecode<JwtPayload>(state.token);

            if (!decodedToken) return

            let currentDate = new Date();

            // JWT exp is in seconds
            if (decodedToken.exp && decodedToken.exp * 1000 < currentDate.getTime()) {
                dispatch({type:"UNSET_TOKEN"})
            }
        }
    ,[state.token, dispatch])

    if(!state.token) {
        return <Login/>
    }


    return (
         <div className="wrapper">
           <BrowserRouter basename="/admin">
               <Navbar bg="light">
                   <Navbar.Brand>Leto 2021</Navbar.Brand>
                   <Nav style={{display:'flex', width:'100%'}}>
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
                       <Nav.Link style={{marginInlineStart:'auto'}}>
                            <button
                                onClick={() => dispatch({type:"UNSET_TOKEN"})}
                            >Odhlasit sa</button>
                       </Nav.Link>
                   </Nav>
               </Navbar>
             <Switch>
                <Route path="/events">
                  <EventList/>
                </Route>
                 <Route path="/registrations/:event">
                     <RegistrationList/>
                 </Route>
                <Route path="/registrations">
                  <RegistrationList/>
                </Route>

             </Switch>
           </BrowserRouter>
         </div>
    );
}

export default App;
