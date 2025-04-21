import React, { useState } from 'react';
import './DaySelector.scss';
import { IEvent, Registration, Day, Stat } from '../../utils/types';
import { IonGrid, IonRow, IonCol, IonItem, IonLabel, IonCheckbox, IonProgressBar } from '@ionic/react';
import { useFormContext, UseFormRegister, UseFormSetValue, UseFormWatch } from 'react-hook-form';

interface ChildInfoProps {
    event: IEvent,
    stats: Stat[],
}

const DaySelector: React.FC<ChildInfoProps> = ({ event, stats }) => {
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

    const { setValue, watch } = useFormContext<Registration>()

    const selected = watch('days')

    return (
        <IonGrid>
            <IonRow>
                <IonCol>
                    <h1>Výber dní</h1>
                    <p>Je možnosť prihlásiť sa iba na niektoré dni alebo na celý čas.</p>
                </IonCol>
            </IonRow>
            <IonRow>
                <h4>Moje dieťa sa zúčastní týchto dní:</h4>
            </IonRow>
            <IonRow>
                <IonCol>
                    {event.days.map((day: Day, i) => (
                        <IonGrid className="dayGrid" key={i}>
                            <IonRow>
                                <IonCol>
                                    <IonItem lines="none">
                                        <IonLabel>
                                            {day.description}  Kapacita: {(parseFloat(`${(stats[i].boys_count+stats[i].girls_count)/stats[i].capacity}`)*100).toFixed(0)} %
                                        </IonLabel>
                                        <IonCheckbox
                                            slot="start"
                                            value={`${day.id}`} 
                                            checked={selected.includes(day.id)}
                                            disabled={(stats[i].boys_count+stats[i].girls_count)/stats[i].capacity === 1}
                                            onIonChange={e => {
                                                if (e.detail.checked) {
                                                    setValue('days', [...selected, day.id])
                                                } else {
                                                    setValue('days', [...selected.filter(id => id !== day.id)])
                                                }
                                            }}
                                        >
                                            
                                        </IonCheckbox>
                                    </IonItem>
                                </IonCol>
                            </IonRow>
                            <IonRow>
                                <IonCol>
                                    <IonProgressBar value={(stats[i].boys_count+stats[i].girls_count)/stats[i].capacity} color={capacitytoColor((stats[i].boys_count+stats[i].girls_count)/stats[i].capacity)}/>
                                </IonCol>
                            </IonRow>
                        </IonGrid>
                    ))}
                </IonCol>
            </IonRow>
        </IonGrid>
    );
};

export default DaySelector;
