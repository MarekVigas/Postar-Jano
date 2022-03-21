import React, { useEffect, useState} from 'react';
import {useAPIClient} from "../AppContext";
import {IEvent, IStat} from "../api/events";
import EventEntry from "./EventEntry";


const EventList:React.FC = () :JSX.Element => {
    const apiClient = useAPIClient()
    const [events, setEvents] = useState<IEvent[]>([])
    const [stats, setStats] = useState<IStat[]>([])

    useEffect(() => {
        apiClient.events.list().then((events) => setEvents(events))
        apiClient.events.stats().then((stats) => setStats(stats))
    },[apiClient])

    const statsForEvent = (eventID :number) => stats.filter(s => s.event_id === eventID)

    return (
            <div>
                {events.map(e => <EventEntry key={e.id} event={e} stats={statsForEvent(e.id)} />)}
            </div>
    )
}
export default EventList;