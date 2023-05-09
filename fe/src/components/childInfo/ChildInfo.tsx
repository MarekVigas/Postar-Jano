import React from 'react';
import { Registration, Gender, ActionType, IEvent } from '../../utils/types';
import { IonGrid, IonRow, IonCol, IonItem, IonLabel, IonInput, IonRadioGroup, IonRadio, IonIcon, IonSelect, IonSelectOption } from '@ionic/react';
import { Control, Controller, UseFormGetValues, UseFormRegister, UseFormSetValue } from 'react-hook-form';
import { RadioGroup } from '../from/radio';
import { Select } from '../from/select';

interface ChildInfoProps {
  register: UseFormRegister<Registration>,
  setValue: UseFormSetValue<Registration>,
  getValues: UseFormGetValues<Registration>
  control: Control<Registration>
}

const ChildInfo: React.FC<ChildInfoProps> = ({ register, control, setValue, getValues }) => {
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
            <IonInput
              {...register('child.name', { required: true })}
              placeholder='Jozef'
            />
          </IonItem>
        </IonCol>
      </IonRow>
      <IonRow>
        <IonCol>
          <h4>Priezvisko</h4>
          <IonItem>
            <IonInput
              {...register('child.surname')}
              placeholder='Mrkvicka'
            />
          </IonItem>
        </IonCol>
      </IonRow>
      <IonRow>
        <IonCol>
          <h4>Pohlavie</h4>
          <IonItem>
            <RadioGroup
              options={[{ value: Gender.Male, label: 'Chlapec' }, { value: Gender.Female, label: 'Dievča' }]}
              name="child.gender"
              register={register}
              setValue={setValue}
              getValues={getValues}
            />
          </IonItem>
        </IonCol>
      </IonRow>
      <IonRow>
        <IonCol>
          <h4>Dátum narodenia</h4>
          <IonItem>
            <IonInput {...register('child.dateOfBirthDay', { required: true })} label="Deň" type='number' min={1} max={31} labelPlacement="stacked" placeholder="23" size={4}></IonInput>
            <IonInput {...register('child.dateOfBirthMonth', { required: true })} label="Mesiac" type='number' min={1} max={12} labelPlacement="stacked" placeholder="5" size={4}></IonInput>
            <IonInput {...register('child.dateOfBirthYear', { required: true })} label="Rok" type='number' labelPlacement="stacked" placeholder="2022" size={8}></IonInput>
          </IonItem>
        </IonCol>
      </IonRow>
      <IonRow>
        <IonCol>
          <h4>Mesto / Obec trvalého bydliska</h4>
          <IonItem>
            <IonInput
              {...register('child.city', { required: true })}
              placeholder="Banská Bystrica"
            />
          </IonItem>
        </IonCol>
      </IonRow>
      <IonRow>
        <IonCol>
          <h4>Ukončený školský rok</h4>
          <IonItem>
            <Select
              options={[
                { value: '1zs', label: '1. ZŠ' },
                { value: '2zs', label: '2. ZŠ' },
                { value: '3zs', label: '3. ZŠ' },
                { value: '4zs', label: '4. ZŠ' },
                { value: '5zs', label: '5. ZŠ' },
                { value: '6zs', label: '6. ZŠ' },
                { value: '7zs', label: '7. ZŠ' },
                { value: '8zs', label: '8. ZŠ' },
                { value: '9zs', label: '9. ZŠ' },
                { value: '1ss', label: '1. SŠ' },
                { value: '2ss', label: '2. SŠ' },
                { value: '3ss', label: '3. SŠ' },
                { value: '4ss', label: '4. SŠ' }
              ]}
              name='child.finishedSchoolYear'
              setValue={setValue}
              control={control}
              placeholder="Vybrať"
              okText="Vybrať"
              cancelText="Zrušiť"
            />
          </IonItem>
        </IonCol>
      </IonRow>
      <IonRow>
        <IonCol>
          <h4>Zúčastnilo sa vaše dieťa minuloročných akcií ?</h4>
          <IonItem>
            <RadioGroup
              options={[{ value: false, label: 'Nie' }, { value: true, label: 'Áno' }]}
              name="child.attendedPreviousEvents"
              register={register}
              setValue={setValue}
              getValues={getValues}
            />
          </IonItem>
        </IonCol>
      </IonRow>
    </IonGrid>
  );
};

export default ChildInfo;
