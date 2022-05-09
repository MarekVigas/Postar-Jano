import React, { useState } from 'react';
import { Registration, ActionType } from '../../types/types';
import { IonGrid, IonRow, IonCol, IonItem, IonLabel, IonInput, IonRadioGroup, IonRadio } from '@ionic/react';

interface OtherInfoProps {
    registration: Registration
    setValue: (action: ActionType, value: any) => void,
}

const OtherInfo: React.FC<OtherInfoProps> = (props) => {
    const [attended, setAttended] = useState<boolean>(false);

    return (
        <IonGrid>
            <IonRow>
                <IonCol>
                    <h1>Ostatné</h1>
                </IonCol>
            </IonRow>
            <IonRow>
                <IonCol>
                    <h4>Navštevovalo vaše dieťa nejaký krúžok, stretko v stredisku počas roka?</h4>
                    <IonRadioGroup value={attended} onIonChange={e => setAttended(e.detail.value)}>
                        <IonItem>
                            <IonLabel>Áno</IonLabel>
                            <IonRadio slot="start" value={true} />
                        </IonItem>
                        <IonItem>
                            <IonLabel>Nie</IonLabel>
                            <IonRadio slot="start" value={false} />
                        </IonItem>
                    </IonRadioGroup>
                </IonCol>
            </IonRow>
            {
                attended && 
                <IonRow>
                    <IonCol>
                        <h4>Prosím uveďte ktoré</h4>
                        <IonItem>
                            <IonLabel position="floating"></IonLabel>
                            <IonInput value={props.registration.memberShip.attendedActivities} onIonChange={e => props.setValue(ActionType.SET_ATTENDED_ACTIVITIES, e.detail.value)} required={true}></IonInput>
                        </IonItem>
                    </IonCol>
                </IonRow>
            }
            <IonRow>
                <IonCol>
                    <h4>Poznámky a pripomienky</h4>
                    <p>Chceli by ste niečo dodať k prihláške alebo povzbudiť organizačný tím ? Tu je na to priestor.</p>
                    <IonItem>
                            <IonLabel position="floating"></IonLabel>
                            <IonInput value={props.registration.notes} onIonChange={e => props.setValue(ActionType.SET_NOTES, e.detail.value)} required={true}></IonInput>
                        </IonItem>
                </IonCol>
            </IonRow>
        </IonGrid>
    );
};

export default OtherInfo;
