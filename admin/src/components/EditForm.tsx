import React from "react";
import {IExtendedRegistration} from "../api/registrations";
import {Button, Modal, Form} from "react-bootstrap";

interface Props {
    show :boolean;
    reg :IExtendedRegistration |null;
    handleChange :(mutator :(prev :IExtendedRegistration) => IExtendedRegistration) => void;
    handleSubmit :() => void;
    handleClose :() => void;
}

const EditForm:React.FC<Props> = ({show, reg:r,handleClose, handleChange, handleSubmit}) => {
    const onClose = () => {
        handleClose()
    }

    if (r == null) return null

    const days = r.days.length === 1 ?
            <span> {r.days[0].description}</span> :
            <ul>
                {r.days.map((d) => <li key={d.id}>{d.description}</li>)}
            </ul>

    const onChange = (mutator :((prev :IExtendedRegistration) => IExtendedRegistration) ) => {
       return handleChange(mutator)
    }

    const renderDate = () => {
        const d = new Date(r.date_of_birth)
        return `${d.getDate()}/${d.getMonth()+1}/${d.getFullYear()}`
    }

    const renderForm = () => {
       const renderMaybeNumber = (val :number|null):number|undefined =>  val || undefined

       return <Form>
           <Form.Group>
               <Form.Label>Suma</Form.Label>
               <Form.Control
                   type="number"
                   onChange={
                       (e :React.ChangeEvent<HTMLInputElement>) => {
                            const val = parseInt(e.target.value) || null
                            if (val == null) return
                            onChange((prev :IExtendedRegistration) => {
                                prev.amount = val
                                return prev
                            })
                       }
                   }
                   value={r.amount}
               />
           </Form.Group>
           <Form.Group>
               <Form.Label>Zaplatené</Form.Label>
           <Form.Control
               type="number"
               onChange={
                   (e :React.ChangeEvent<HTMLInputElement>) => {
                       const val = parseInt(e.target.value) || null
                       onChange((prev :IExtendedRegistration) => {
                           prev.payed = val
                           return prev
                       })
                   }
               }
               value={renderMaybeNumber(r.payed)}
           />
           </Form.Group>
           <Form.Group>
               <Form.Label>Poznamka admin</Form.Label>
           <Form.Control
               type="textarea"
               onChange={
                   (e :React.ChangeEvent<HTMLInputElement>) => {
                       const val = e.target.value
                       onChange((prev :IExtendedRegistration) => {
                           prev.admin_note = val
                           return prev
                       })
                   }
               }
               value={r.admin_note}
           />
           </Form.Group>
        </Form>
    }

    return (<>
        <Modal show={show} onHide={onClose}>
            <Modal.Header closeButton>
                <Modal.Title>{r.name} {r.surname} {renderDate()}</Modal.Title>
            </Modal.Header>

            <Modal.Body>
                <span>{r.title}</span>
                {days}
                <br/>
                {r !== null && renderForm()}
                <Button variant="secondary" onClick={onClose}>Zrušiť</Button>
                <Button variant="primary" onClick={() => {
                    handleSubmit()
                }}>Uloziť zmeny</Button>
            </Modal.Body>

            <Modal.Footer>
                <div className="editForm-readonly-wrapper">
                    <span>Rodič:</span> <span>{r.parent_name} {r.parent_surname}</span>
                    <span>Email:</span> <span>{r.email}</span>
                    <span>Telefon:</span> <span>{r.phone}</span>
                    <span>Poznamka:</span> <span>{r.notes}</span>
                </div>
            </Modal.Footer>
        </Modal>
    </>)
}

export default EditForm