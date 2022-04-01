import React from 'react';
import './ChildInfo.scss';
import { Registration, Gender, ActionType, IEvent } from '../../types/types';
import DatePicker from 'react-datepicker';
import { IonGrid, IonRow, IonCol, IonItem, IonLabel, IonInput, IonRadioGroup, IonRadio, IonIcon, IonSelect, IonSelectOption } from '@ionic/react';
import "react-datepicker/dist/react-datepicker.css";

interface ChildInfoProps {
    registration: Registration,
    event: IEvent,
    setValue: (action: ActionType, value: any) => void,
}

const ChildInfo: React.FC<ChildInfoProps> = (props) => {
    const { child } = props.registration;

    const convertUTCToLocalDate = (date: any) => {
        if (!date) {
            return date
        }
        date = new Date(date)
        date = new Date(date.getUTCFullYear(), date.getUTCMonth(), date.getUTCDate())
        return date
    }

    const convertLocalToUTCDate = (date: any) => {
        if (!date) {
            return date
        }
        date = new Date(date)
        date = new Date(Date.UTC(date.getFullYear(), date.getMonth(), date.getDate()))
        return date
    }

    return (
        <IonGrid>
            <IonRow>
                <IonCol>
                    <h1>Informácie o dieťati</h1>
                </IonCol>
            </IonRow>
            <IonRow>
                <IonCol>
                    <h4>Meno</h4>
                    <IonItem>
                        <IonLabel position="floating"></IonLabel>
                        <IonInput value={child.name} onIonChange={e => props.setValue(ActionType.SET_CHILD_NAME, e.detail.value)} placeholder="Jožko" required={true}></IonInput>
                    </IonItem>
                </IonCol>
            </IonRow>
            <IonRow>
                <IonCol>
                    <h4>Priezvisko</h4>
                    <IonItem>
                        <IonLabel position="floating"></IonLabel>
                        <IonInput value={child.surname} onIonChange={e => props.setValue(ActionType.SET_CHILD_SURNAME, e.detail.value)} placeholder="Mrkvička" required={true}></IonInput>
                    </IonItem>
                </IonCol>
            </IonRow>
            <IonRow>
                <IonCol>
                    <h4>Pohlavie</h4>
                    <IonRadioGroup value={props.registration.child.gender} onIonChange={e => props.setValue(ActionType.SET_CHILD_GENDER, e.detail.value)}>
                        <IonItem>
                            <IonLabel>Chlapec</IonLabel>
                            <IonRadio slot="start" value={Gender.Male} />
                        </IonItem>
                        <IonItem>
                            <IonLabel>Dievča</IonLabel>
                            <IonRadio slot="start" value={Gender.Female} />
                        </IonItem>
                    </IonRadioGroup>
                </IonCol>
            </IonRow>
            <IonRow>
                <IonCol>
                    <h4>Dátum narodenia</h4>
                    <DatePicker
                        onChange={(date: any) => props.setValue(ActionType.SET_CHILD_BIRTH, convertLocalToUTCDate(date))}
                        selected={convertUTCToLocalDate(child.dateOfBirth)}
                        minDate={new Date(new Date().getFullYear() - props.event.max_age, 8, 15)}
                        maxDate={new Date(new Date().getFullYear() - props.event.min_age, 8, 15)}
                        dateFormat="dd/MM/yyyy"
                    />
                </IonCol>
            </IonRow>
            <IonRow>
                <IonCol>
                    <h4>Mesto / Obec trvalého bydliska</h4>
                    <IonItem>
                        <IonLabel></IonLabel>
                        <IonInput value={child.city} onIonChange={e => props.setValue(ActionType.SET_CHILD_CITY, e.detail.value)} required={true} placeholder="Banská Bystrica"></IonInput>
                    </IonItem>
                </IonCol>
            </IonRow>
            <IonRow>
                <IonCol>
                    <h4>Ukončený školský rok</h4>
                    <IonItem>
                    <IonSelect 
                        value={child.finishedSchoolYear} 
                        placeholder="Vybrať" 
                        onIonChange={e => props.setValue(ActionType.SET_CHILD_YEAR, e.detail.value)}
                        okText="Vybrať"
                        cancelText="Zrušiť"
                    >
                        <IonSelectOption value="1zs">1. ZŠ</IonSelectOption>
                        <IonSelectOption value="2zs">2. ZŠ</IonSelectOption>
                        <IonSelectOption value="3zs">3. ZŠ</IonSelectOption>
                        <IonSelectOption value="4zs">4. ZŠ</IonSelectOption>
                        <IonSelectOption value="5zs">5. ZŠ</IonSelectOption>
                        <IonSelectOption value="6zs">6. ZŠ</IonSelectOption>
                        <IonSelectOption value="7zs">7. ZŠ</IonSelectOption>
                        <IonSelectOption value="8zs">8. ZŠ</IonSelectOption>
                        <IonSelectOption value="9zs">9. ZŠ</IonSelectOption>
                        <IonSelectOption value="1ss">1. SŠ</IonSelectOption>
                        <IonSelectOption value="2ss">2. SŠ</IonSelectOption>
                        <IonSelectOption value="3ss">3. SŠ</IonSelectOption>
                        <IonSelectOption value="4ss">4. SŠ</IonSelectOption>
                    </IonSelect>
                    </IonItem>
                </IonCol>
            </IonRow>
            <IonRow>
                <IonCol>
                    <h4>Zúčastnilo sa vaše dieťa minuloročných akcií ?</h4>
                    <IonRadioGroup value={props.registration.child.attendedPreviousEvents} onIonChange={e => props.setValue(ActionType.SET_CHILD_ATTEND, e.detail.value)}>
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
        </IonGrid>
    );
};

export default ChildInfo;
