import React from 'react';
import './IntroInfo.scss';
import { IEvent, Stat } from '../../types/types';
import { calendarOutline, locateOutline, personOutline, logoEuro, manOutline, womanOutline, timeOutline, informationCircleOutline } from 'ionicons/icons';
import { IonIcon, IonItem, IonAvatar, IonLabel, IonCol, IonGrid, IonRow, IonProgressBar } from '@ionic/react';

interface IntroInfoProps {
    event: IEvent,
    stats: Stat[]
}

const IntroInfo: React.FC<IntroInfoProps> = (props) => {
    const { event, stats } = props;

    const capacitytoColor = (capacity: number): string => {
        let color = "primary";

        if (capacity >= 0.5) {
            if (capacity < 0.75) {
                color = "warning" 
            } else if (capacity >= 0.75) {
                color = "danger"
            }
        }
        return color
    }

    const calculatePercentage = (a: number, b: number): string => {
        return `${parseFloat(`${(a/b)*100}`).toFixed(0)} %`
    }

    return (
        <IonGrid>
            <IonRow>
                <IonCol>
                    <h1>{event.title.toUpperCase()}</h1>
                </IonCol>
            </IonRow>
            {
                event.photo && 
                <IonRow>
                    <IonCol></IonCol>
                    <IonCol>
                        <img className="eventPhoto" src={event.photo} alt=""/>
                    </IonCol>
                    <IonCol></IonCol>
                </IonRow>
            }
            <IonItem>
                    <IonCol slot="start" size="0.9" className="infIcon" >
                        <IonIcon icon={calendarOutline} class="bigIcon" />
                    </IonCol>
                    <IonCol size="4">
                        <h5>Dátum:</h5>
                    </IonCol>
                    <IonCol>
                        <span className="value">                    
                            {event.date_from} - {event.date_to}
                        </span>
                    </IonCol>
            </IonItem>
            <IonItem>
                <IonCol slot="start" size="0.9" className="infIcon" >
                    <IonIcon icon={locateOutline} class="bigIcon" />
                </IonCol>
                <IonCol size="4">
                    <h5>Miesto:</h5>
                </IonCol>
                <IonCol>
                    <span className="value">
                        {event.location}
                    </span>
                </IonCol>
            </IonItem>
            {
                event.time && 
                <IonItem>
                    <IonCol slot="start" size="0.9" className="infIcon" >
                        <IonIcon icon={timeOutline} class="bigIcon" />
                    </IonCol>
                    <IonCol size="4">
                        <h5>Čas:</h5>
                    </IonCol>
                    <IonCol>
                        <span className="value">
                            {event.time}
                        </span>
                    </IonCol>
                </IonItem>
            }
            <IonItem>
                <IonCol slot="start" size="0.9" className="infIcon" >
                    <IonIcon icon={personOutline} class="bigIcon" />
                </IonCol>
                <IonCol size="4">
                    <h5>Vek:</h5>
                </IonCol>
                <IonCol>
                    <span className="value">
                        {event.description}
                    </span>
                </IonCol>
            </IonItem>
            <IonItem>
                <IonCol slot="start" size="0.9" className="infIcon" >
                    <IonIcon icon={logoEuro} class="bigIcon" />
                </IonCol>
                <IonCol size="4">
                    <h5>Cena:</h5>
                </IonCol>
                <IonCol>
                    <span className="value">
                        {event.price}
                    </span>
                </IonCol>
            </IonItem>
            {
                event.info && 
                <IonItem>
                    <IonCol slot="start" size="0.9" className="infIcon" >
                        <IonIcon icon={informationCircleOutline} class="bigIcon" />
                    </IonCol>
                    <IonCol size="4">
                        <h5>Ďalšie informácie:</h5>
                    </IonCol>
                    <IonCol>
                        <span className="value">
                            {event.info}
                        </span>
                    </IonCol>
                </IonItem>
            }
            <IonRow>
                <h5>Kapacita: </h5>
            </IonRow>
            {
                stats[0].limit_boys != null && stats.map((stat, index) => (
                    <IonItem key={index}>
                        <IonCol slot="start" size="3">
                            {event.days[index].description}
                        </IonCol>
                        <IonCol size="1">
                            <IonIcon icon={manOutline} className="bigIcon" />
                        </IonCol>
                        <IonCol size="1" className="capacityPercentage">
                            {calculatePercentage(stat.boys_count, stat.limit_boys!)}
                        </IonCol>
                        <IonCol>
                            <IonProgressBar 
                                value={stat.boys_count/stat.limit_boys!} 
                                color={capacitytoColor((stat.boys_count+stat.girls_count)/stat.capacity)}/>
                        </IonCol>
                        <IonCol size="1">
                            <IonIcon icon={womanOutline} className="bigIcon" />
                        </IonCol>
                        <IonCol size="1" className="capacityPercentage">
                            {calculatePercentage(stat.girls_count, stat.limit_girls!)}
                        </IonCol>
                        <IonCol>
                            <IonProgressBar 
                                value={stat.girls_count/stat.limit_girls!} 
                                color={capacitytoColor((stat.boys_count+stat.girls_count)/stat.capacity)}/>
                        </IonCol>
                    </IonItem>
                ))
            }
            {
                stats[0].limit_boys == null && stats.map((stat, index) => (
                    <IonItem key={index}>
                        <IonCol slot="start" size="3">
                            {event.days[index].description}
                        </IonCol>
                        <IonCol size="1" className="capacityPercentage">
                            {calculatePercentage(stat.boys_count+stat.girls_count, stat.capacity)}
                        </IonCol>
                        <IonCol>
                            <IonProgressBar 
                                value={(stat.boys_count+stat.girls_count)/stat.capacity} 
                                color={capacitytoColor((stat.boys_count+stat.girls_count)/stat.capacity)}/>
                        </IonCol>
                </IonItem>
                ))
            }
            <IonRow>
                <IonCol>
                    <h5>Za akciu zodpovedá:</h5>
                </IonCol>
            </IonRow>
            <IonRow>
                <IonCol>
                <IonItem>
                    <IonAvatar slot="start">
                        <img src={event.owner.photo} alt={`${event.owner.name} ${event.owner.surname}`} />
                    </IonAvatar>
                    <IonLabel>{event.owner.name} {event.owner.surname} <br/>{event.owner.email} <br/> {event.owner.phone}</IonLabel>
                </IonItem>
                </IonCol>
            </IonRow>
        </IonGrid>
    );
};

export default IntroInfo;
