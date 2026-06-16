import React, {useState} from "react";
import {IExtendedRegistration} from "../api/registrations";
import {Button, Form, Modal} from "react-bootstrap";

interface Props {
    show: boolean;
    reg: IExtendedRegistration | null;
    handleSubmit: (email: string) => void;
    handleClose: () => void;
}

const ResendConfirmationModal: React.FC<Props> = ({show, reg: r, handleClose, handleSubmit}) => {
    const [email, setEmail] = useState<string>("")

    if (r == null) return null

    const onShow = () => setEmail(r.email)

    return (
        <Modal show={show} onHide={handleClose} onShow={onShow}>
            <Modal.Header closeButton>
                <Modal.Title>Preposlať potvrdzovací email — {r.name} {r.surname}</Modal.Title>
            </Modal.Header>
            <Modal.Body>
                <p><strong>{r.title}</strong></p>
                <Form>
                    <Form.Group>
                        <Form.Label>Email</Form.Label>
                        <Form.Control
                            type="email"
                            value={email}
                            onChange={(e: React.ChangeEvent<HTMLInputElement>) => setEmail(e.target.value)}
                        />
                    </Form.Group>
                </Form>
            </Modal.Body>
            <Modal.Footer>
                <Button variant="secondary" onClick={handleClose}>Zrušiť</Button>
                <Button variant="primary" onClick={() => handleSubmit(email)}>Preposlať</Button>
            </Modal.Footer>
        </Modal>
    )
}

export default ResendConfirmationModal
