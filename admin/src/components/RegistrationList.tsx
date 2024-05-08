import React, {useEffect, useState} from 'react';
import {IExtendedRegistration} from "../api/registrations";
import RegistrationEntry from "./RegistrationEntry";
import ViewFilter, {IViewFields} from "./ViewFilter";
import {useParams} from 'react-router-dom'
import {useAPIClient} from "../AppContext";
import useStorage from "../hooks/useStorage";
import EditForm from "./EditForm";
import {Table} from 'react-bootstrap'
import CopyButton from "./CopyButton";

const RegistrationList:React.FC = () :JSX.Element => {
    const [filter, setFilter] = useState<string>("")
    const [registrations, setRegistrations] = useState<IExtendedRegistration[]>([])
    const [fields, setFields] = useStorage<IViewFields>("VIEW_FILTERS",{
        "id": { display: "ID", show: true},
        "name": { display: "Meno", show: true},
        "surname": { display: "Priezvisko", show: true},
        "gender": {display: "Pohlavie", show: true},
        "date_of_birth": {display: "Datum narodenia", show: false},
        "finished_school": {display: "Dokonceny rocnik", show:false},
        "attended_previous": {display: "Ucast na akciach", show:false},
        "attended_activities": {display: "Navsetvovane aktivity", show:false},

        "city" : {display: "Mesto", show:false},
        "pills" :{display: "Lieky", show:false},
        "problems" :{display: "Zdravotne problemy", show:false},
        "notes": {display: "Poznamka", show:false},

        "parent_name" :{display: "Meno rodica", show:false},
        "parent_surname" :{display: "Priezvisko rodica", show:false},
        "email" :{display: "email", show:false},
        "phone" :{display: "telefon", show:false},

        "amount": {display: "Suma", show:false},
        "payed": {display: "Zaplatene", show:false},
        "discount": {display: "Zlava", show:false},
        "specific_symbol": {display: "Špecifický symbol", show:false},

        "title": { display: "Akcia", show: true},
        "days": { display: "Dni", show: true},

        "buttons": {display: "Upravy", show:true},
    })
    const [expandViewFilter, setExpandViewFilter] = useState(false)
    const [showEdit, setShowEdit] = useState(false)
    const [editedRegistration, setEditedRegistration] = useState<IExtendedRegistration|null>(null)
    const apiClient = useAPIClient()
    const {event} = useParams<{event:string}>()

    useEffect(
        () => {
            apiClient.registrations.list().then(
                (data) => setRegistrations(data)
            )
        }
    ,[apiClient])

    const displayedRegistrations = ():IExtendedRegistration[] => {
        const eventID = parseInt(event || "")

        const byEvent = registrations.filter(r => (!eventID) || r.eventID === eventID)

        if (filter === "") {
            return byEvent
        }

        const lowerFilter = filter.toLowerCase()
        return byEvent.filter(r =>
            r.name.toLowerCase().includes(lowerFilter) ||
            r.surname.toLowerCase().includes(lowerFilter) ||
            r.email.toLowerCase().includes(lowerFilter)
        )
    }

    const renderViewFilter = () => {
        if (expandViewFilter) {
            return (
                <>
                    <button onClick={()=>setExpandViewFilter(false)}>Skryť</button>
                    <ViewFilter fields={fields} setFields={setFields}/>
                </>
            )
        }
        return (
            <>
                <button onClick={()=>setExpandViewFilter(true)}>Zobrazené stĺpce</button>
            </>
        )
    }

    const handleUpdate = () => {
        if (editedRegistration === null) return
        apiClient.registrations.update(editedRegistration).then(
            () => {
                setRegistrations((prevState) => {
                    return prevState.map((r) => {
                        if (r.id !== editedRegistration.id) return r
                        return {...r, ...editedRegistration}
                    })
                })
                setShowEdit(false)
            }
        )
    }

    const handleEdit = (mutator:(prev :IExtendedRegistration) => IExtendedRegistration) => {
        setEditedRegistration((p :IExtendedRegistration|null) => {
            if (p == null) return p
            const copyVal :IExtendedRegistration = {...p}

            const newVal = mutator(copyVal)
            return newVal
        })
    }


    const isShown = (name :string): boolean => {
        const entry = fields[name]
        if (!entry) return true
        return entry.show
    }

    return (
        <div>
        <h2>Zoznam prihlasenych</h2>
            <input
                value={filter}
                placeholder='filter...'
                onChange={event => setFilter(event.target.value)}
            />
            {editedRegistration!= null && <EditForm
                show={showEdit}
                reg={{...editedRegistration}}
                handleChange={handleEdit}
                handleSubmit={handleUpdate}
                handleClose={() => {setShowEdit(false)}}
            />}
            {renderViewFilter()}
            <CopyButton selector="table"/>
            <Table>
                <thead>
                {isShown("id") && <th>ID</th>}
                {isShown("name") && <th>Meno</th>}
                {isShown("surname") && <th>Priezvisko</th>}
                {isShown("date_of_birth") && <th>Datum narodenia</th>}
                {isShown("gender") && <th>Pohlavie</th>}
                {isShown("finished_school") && <th>Ukonceny rocnik</th>}
                {isShown("attended_previous") && <th>Ucast na predchadzajucich akciach</th>}
                {isShown("attended_activities") && <th>Navsetvovane aktivity</th>}

                {isShown("city") && <th>Mesto</th>}
                {isShown("pills") && <th>Lieky</th>}
                {isShown("problems") && <th>Zdravotne problemy</th>}
                {isShown("notes") && <th>Poznamky</th>}

                {isShown("title") && <th>Nazov akcie</th>}
                {isShown("days") && <th>Dni</th>}

                {isShown("parent_name") && <th>Meno rodica</th>}
                {isShown("parent_surname") && <th>Priezvisko rodica</th>}
                {isShown("email") && <th>Email</th>}
                {isShown("phone") && <th>Telefon</th>}


                {isShown("specific_symbol") && <th>Specificky symbol</th>}
                {isShown("amount") && <th>Suma</th>}
                {isShown("payed") && <th>Zaplatene</th>}
                {isShown("discount") && <th>Zlava</th>}

                {isShown("buttons") && <th>Upravy </th>}
                </thead>
                <tbody>
                {displayedRegistrations().map(
                    r =>
                        <RegistrationEntry
                            key={r.id}
                            fields={fields}
                            registration={r}
                            deleteRegByID={(id:number) => {
                                apiClient.registrations.delete(id).then(()=>{
                                    setRegistrations((prev => prev.filter(r => r.id !== id)))
                                })
                            }}
                            editRegByID={() => {
                                setEditedRegistration({...r})
                                setShowEdit(true)
                            }}
                        />
                )}
                </tbody>
            </Table>
        </div>
    )
}
export default RegistrationList;