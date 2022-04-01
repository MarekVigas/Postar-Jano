import React from 'react';
import "./Results.scss";
import load from "../load.png"
import { Registration, responseStatus } from '../../types/types';
import { IonGrid, IonRow, IonCol, IonButton } from '@ionic/react';

interface ResultsProps {
    registration: Registration,
    responseMsg: string | null,
    responseStatus: responseStatus | null
}

const Results: React.FC<ResultsProps> = (props) => {
    

    return (
        <IonGrid>
            {
                props.responseMsg == null && props.responseStatus == null && 
                <IonRow>
                    <IonCol></IonCol>
                    <IonCol>
                        <img src={load} alt="" className="rotate-scale-up" />
                    </IonCol>
                    <IonCol></IonCol>
                </IonRow>
            }
            {
                props.responseMsg != null && props.responseStatus != null && 
                <IonRow>
                    <IonCol></IonCol>
                    <IonCol size="8">
                        <div className={`${props.responseStatus} resultbox`}>
                            {props.responseMsg}
                        </div>
                    </IonCol>
                    <IonCol></IonCol>
                </IonRow>
            }
            {
                props.responseMsg != null && props.responseStatus === responseStatus.success && 
                <IonRow>
                    <IonCol></IonCol>
                    <IonCol size="4">
                        <IonButton onClick={() => {
                            window.location.replace("https://sbb.sk/leto");
                        }}>
                            Zobraziť ďalšie akcie
                        </IonButton>
                    </IonCol>
                    <IonCol></IonCol>
                </IonRow>
            }
        </IonGrid>
    );
};

export default Results;
