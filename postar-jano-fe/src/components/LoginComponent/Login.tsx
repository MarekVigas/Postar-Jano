import React, { useState } from 'react'
import { IonGrid, IonRow, IonCol, IonInput, IonItem, IonLabel, IonButton, IonIcon } from '@ionic/react'
import { logInOutline } from 'ionicons/icons'

export default function Login() {
    const [email, setEmail] = useState<string>("")
    const [password, setPassword] = useState<string>("")

    const submitLogin = () => {

    }

    return (
        <IonGrid>
            <IonRow>
                <IonCol size="4"></IonCol>
                <IonCol>
                    <IonItem>
                        <IonLabel position="floating">
                            Email
                        </IonLabel>
                        <IonInput 
                            value={email} type="email" 
                            placeholder="john.doe@gmail.com"
                            onIonChange={e => setEmail(e.detail.value!)}
                        />
                    </IonItem>
                </IonCol>
                <IonCol size="4"></IonCol>
            </IonRow>
            <IonRow>
                <IonCol size="4"></IonCol>
                <IonCol>
                    <IonItem>
                        <IonLabel position="floating">
                            Password
                        </IonLabel>
                        <IonInput 
                            value={password} type="password" 
                            onIonChange={e => setPassword(e.detail.value!)}
                        />
                    </IonItem>
                </IonCol>
                <IonCol size="4"></IonCol>
            </IonRow>
            <IonRow>
                <IonCol size="4"></IonCol>
                <IonCol>
                    <IonButton onClick={() => submitLogin()}>
                        <IonIcon icon={logInOutline} /> Login
                    </IonButton>
                </IonCol>
                <IonCol size="4"></IonCol>
            </IonRow>
        </IonGrid>
    )
}
