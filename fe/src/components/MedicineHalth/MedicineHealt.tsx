import React from 'react';
import { IonGrid, IonRow, IonCol, IonItem, IonLabel, IonInput, IonRadioGroup, IonRadio } from '@ionic/react';
import { UseFormGetValues, UseFormRegister, UseFormSetValue, UseFormWatch } from 'react-hook-form';
import { Registration } from '../../utils/types';
import { RadioGroup } from '../from/radio';

interface MedicineHealthProps {
  register: UseFormRegister<Registration>
  setValue: UseFormSetValue<Registration>
  getValues: UseFormGetValues<Registration>
  watch: UseFormWatch<Registration>
}

const MedicineHealth: React.FC<MedicineHealthProps> = ({ register, setValue, watch, getValues }) => {

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
          <IonItem>
            <RadioGroup
              options={[{ value: false, label: 'Nie' }, { value: true, label: 'Áno' }]}
              name="medicine.takes"
              register={register}
              getValues={getValues}
              setValue={setValue}
            />
          </IonItem>
        </IonCol>
      </IonRow>
      {
        watch('medicine.takes') &&
        <IonRow>
          <IonCol>
            <h4>Prosím uveďte aké lieky užíva vaše dieťa</h4>
            <IonItem>
              <IonInput
                {...register('medicine.drugs')}
              />
            </IonItem>
          </IonCol>
        </IonRow>
      }
      <IonRow>
        <IonCol>
          <h4>Má vaše dieťa nejaké zdravotné ťažkosti alebo obmedzenia ?</h4>
          <p>Alergie / Intolerancie</p>
          <IonItem>
            <RadioGroup
              options={[{ value: false, label: 'Nie' }, { value: true, label: 'Áno' }]}
              name="health.hasProblmes"
              register={register}
              setValue={setValue}
              getValues={getValues}
            />
          </IonItem>
        </IonCol>
      </IonRow>
      {
        watch('health.hasProblmes') &&
        <IonRow>
          <IonCol>
            <h4>Prosím uveďte aké zdravotné ťažkosti alebo obmedzenia má vaše dieťa</h4>
            <IonItem>
              <IonInput
                {...register('health.problems')}
              />
            </IonItem>
          </IonCol>
        </IonRow>
      }
    </IonGrid>
  );
};

export default MedicineHealth;
