import React, { useState } from 'react';
import './DaySelector.scss';
import { IEvent, Registration, ActionType, Day, Stat } from '../../utils/types';
import { IonGrid, IonRow, IonCol, IonItem, IonLabel, IonCheckbox, IonProgressBar } from '@ionic/react';

interface ChildInfoProps {
    event: IEvent,
    registration: Registration,
    stats: Stat[],
    setValue: (action: ActionType, value: any) => void,
}

const DaySelector: React.FC<ChildInfoProps> = (props) => {
    const { days } = props.event;
    const { stats } = props;
    const [selected, setSelected] = useState<string[]>([...props.registration.days.map(d => `${d}`)]);

    const handleSetSelected = (days: string[]) => {
        setSelected(days);
        props.setValue(ActionType.SET_DAYS, days.map(d => parseInt(d, 10)));
    }

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
                    {days.map((day: Day, i) => (
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
                                            checked={selected.includes(`${day.id}`)}
                                            disabled={(stats[i].boys_count+stats[i].girls_count)/stats[i].capacity === 1}
                                            onIonChange={e => {
                                                if (e.detail.checked) {
                                                    handleSetSelected([...selected, `${day.id}`])
                                                } else {
                                                    handleSetSelected([...selected.filter(d => d !== `${day.id}`)])
                                                }
                                            }}
                                        />
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
