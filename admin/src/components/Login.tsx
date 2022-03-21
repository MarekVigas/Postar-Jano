import React, {useState} from 'react';
import {Button, Alert, Form} from "react-bootstrap";
import {useAPIClient, useAppDispatch} from "../AppContext";

const Login: React.FC = () => {
    const [username, setUserName] = useState<string>("")
    const [password, setPassword] = useState<string>("")

    const [errs, setErrs] = useState<string>()
    const apiClient = useAPIClient()
    const dispatch = useAppDispatch()

    const handleSubmit = (e :React.FormEvent) => {
        e.preventDefault();
        apiClient.user.signIn(username, password).then(
            (data) => {
                dispatch({type:"SET_TOKEN", token: data.token})
            }
        ).catch(
            (err) => {
                if (err.response !== undefined && err.response.status === 403) {
                    setErrs("Nespravne meno alebo heslo.")
                } else {
                    setErrs("Chyba prihlasenia, skuste to neskor.")
                }
            }
        )
    }

    return(
        <div className="login-wrapper">
            <h1>Prihláste sa</h1>
            {errs && <Alert variant="danger">{errs}</Alert>}
            <Form onSubmit={handleSubmit}>
                <Form.Group>
                    <Form.Label>Prihlasovacie meno</Form.Label>
                    <Form.Control
                        type="text"
                        onChange={e => setUserName(e.target.value)}
                    />
                </Form.Group>

                <Form.Group>
                    <Form.Label>Heslo</Form.Label>
                    <Form.Control
                        type="password"
                        onChange={e => setPassword(e.target.value)}
                    />
                </Form.Group>

                <Button type="submit" variant="primary">Prihlásiť sa</Button>
            </Form>
        </div>
    )
}
export default Login;