import React, { useEffect } from 'react';
import { Registration } from '../../utils/types';
import { IonGrid, IonRow, IonCol, IonItem, IonInput } from '@ionic/react';
import { UseFormRegister, UseFormTrigger } from 'react-hook-form';

interface OtherInfoProps {
  register: UseFormRegister<Registration>
  trigger: UseFormTrigger<Registration>
}

const OtherInfo: React.FC<OtherInfoProps> = ({ register, trigger }) => {
  useEffect(() => {
    trigger()
  }, [])

  return (
    <IonGrid>
      <IonRow>
        <IonCol>
          <h1>Ostatné</h1>
        </IonCol>
      </IonRow>
      <IonRow>
        <IonCol>
          <h4>Poznámky a pripomienky</h4>
          <p>Chceli by ste niečo dodať k prihláške alebo povzbudiť organizačný tím ? Tu je na to priestor.</p>
          <IonItem>
            <IonInput
              {...register('notes')}
            />
          </IonItem>
        </IonCol>
      </IonRow>
    </IonGrid>
  );
};

export default OtherInfo;
