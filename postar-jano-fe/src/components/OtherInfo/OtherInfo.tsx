import React from 'react';
import { Registration } from '../../utils/types';
import { IonGrid, IonRow, IonCol, IonItem, IonInput } from '@ionic/react';
import { UseFormRegister } from 'react-hook-form';

interface OtherInfoProps {
  register: UseFormRegister<Registration>
}

const OtherInfo: React.FC<OtherInfoProps> = ({ register }) => {
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
