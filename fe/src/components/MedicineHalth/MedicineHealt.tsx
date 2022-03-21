import React from 'react';
import { Registration, ActionType } from '../../types/types';
import { IonGrid, IonRow, IonCol, IonItem, IonLabel, IonInput, IonRadioGroup, IonRadio } from '@ionic/react';

interface MedicineHealthProps {
    registration: Registration
    setValue: (action: ActionType, value: any) => void,
}

const MedicineHealth: React.FC<MedicineHealthProps> = (props) => {

    return (
        <IonGrid>
            <IonRow>
                <IonCol>
                    <h1>Lieky a zdravotný stav</h1>
                </IonCol>
            </IonRow>
            <IonRow>
                <IonCol>
                    <h4>Užíva vaše dieťa nejaké lieky ?</h4>
                    <IonRadioGroup value={props.registration.medicine.takes} onIonChange={e => props.setValue(ActionType.SET_MEDICINE, e.detail.value)}>
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
                props.registration.medicine.takes && 
                <IonRow>
                    <IonCol>
                        <h4>Prosím uveďte aké lieky užíva vaše dieťa</h4>
                        <IonItem>
                            <IonLabel position="floating"></IonLabel>
                            <IonInput value={props.registration.medicine.drugs} onIonChange={e => props.setValue(ActionType.SET_DRUGS, e.detail.value)} placeholder="zodac" required={true}></IonInput>
                        </IonItem>
                    </IonCol>
                </IonRow>
            }
            <IonRow>
                <IonCol>
                    <h4>Má vaše dieťa nejaké zdravotné ťažkosti alebo obmedzenia ?</h4>
                    <p>Alergie</p>
                    <IonRadioGroup value={props.registration.health.hasProblmes} onIonChange={e => props.setValue(ActionType.SET_HEALTH, e.detail.value)}>
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
                props.registration.health.hasProblmes && 
                <IonRow>
                    <IonCol>
                        <h4>Prosím uveďte aké zdravotné ťažkosti alebo obmedzenia má vaše dieťa</h4>
                        <IonItem>
                            <IonLabel position="floating"></IonLabel>
                            <IonInput value={props.registration.health.problems} onIonChange={e => props.setValue(ActionType.SET_PROBLEMS, e.detail.value)} required={true} placeholder="alergiu"></IonInput>
                        </IonItem>
                    </IonCol>
                </IonRow>
            }
        </IonGrid>
    );
};

export default MedicineHealth;
