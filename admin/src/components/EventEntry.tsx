import React from 'react';
import {IEvent, IStat} from "../api/events";
import {Link} from "react-router-dom";
import Limit from "./Limit";

interface Props {
    event :IEvent;
    stats :IStat[];
}

const EventEntry:React.FC<Props> = (props) :JSX.Element => {
    const {
        event:e,
        stats,
    } = props

    const statForDay = (id :number) => {
        const filtered = stats.filter(s => s.day_id === id)
        return filtered[0]
    }

    return (
        <div className="event-entry-wrapper">
            <h5>{e.title}</h5>
            {e.days.map(d => (
                <div key={d.id} >
                    <span>{d.description}</span>
                    {statForDay(d.id) && <Limit day={d} stat={statForDay(d.id)}/>}
                </div>
            ))}
            <Link to={`/registrations/${e.id}`}>
                Účastníci
            </Link>
        </div>
    )
}

export default EventEntry;
