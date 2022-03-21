import React from "react";
import {IExtendedRegistration} from "../api/registrations";
import {IViewFields} from "./ViewFilter";

interface Props {
    registration : IExtendedRegistration;
    fields : IViewFields;
    deleteRegByID :(id:number)=>void;
    editRegByID :()=>void;
}

const RegistrationEntry :React.FC<Props> = ({fields, registration, deleteRegByID, editRegByID}) => {
    const {
        id,
        name,
        surname,
        date_of_birth,
        finished_school,
        attended_previous,
        attended_activities,
        gender,

        city,
        pills,
        problems,
        notes,

        parent_name,
        parent_surname,
        email,
        phone,

        title,
        days,

        amount,
        payed,
        discount,
    } = registration

    const isShown = (name :string): boolean => {
        const entry = fields[name]
        if (!entry) return true
        return entry.show
    }

    const renderDate = () => {
        const d = new Date(date_of_birth)
        return `${d.getDate()}/${d.getUTCMonth()+1}/${d.getFullYear()}`
    }

    return (
        <tr>
            {isShown("id") && <td>{id}</td>}
            {isShown("name") && <td>{name}</td>}
            {isShown("surname") && <td>{surname}</td>}
            {isShown("date_of_birth") && <td>{renderDate()}</td>}
            {isShown("gender") && <td>{gender === "male"? "CH": "D"}</td>}
            {isShown("finished_school") && <td>{finished_school}</td>}
            {isShown("attended_previous") && <td>{attended_previous}</td>}
            {isShown("attended_activities") && <td>{attended_activities}</td>}

            {isShown("city") && <td>{city}</td>}
            {isShown("pills") && <td>{pills}</td>}
            {isShown("problems") && <td>{problems}</td>}
            {isShown("notes") && <td>{notes}</td>}

            {isShown("title") && <td>{title}</td>}
            {isShown("days") &&
            <td>{days.map((d) => (
                    <span key={d.id}>
                   {d.description}
                </span>
                )
            )}</td>
            }

            {isShown("parent_name") && <td>{parent_name}</td>}
            {isShown("parent_surname") && <td>{parent_surname}</td>}
            {isShown("email") && <td>{email}</td>}
            {isShown("phone") && <td>{phone}</td>}


            {isShown("amount") && <td>{amount}€</td>}
            {isShown("payed") && <td>{payed}€</td>}
            {isShown("discount") && <td>{discount}€</td>}

            {isShown("buttons") && <td>
            <button onClick={() => {
                editRegByID()
            }}>Upravit</button>
            <button onClick={(e) => {
                if (window.confirm(`Naozaj chcete vymazat registraciu ${name} ${surname} (ID ${id})?`)) {
                    console.log("Deleting ", id)
                    deleteRegByID(id)
                }
            }}>&times;</button>
            </td>}
        </tr>
    )
}

export default RegistrationEntry;