import React from 'react';
import { Registration } from '../../utils/types';
import { IonGrid, IonRow, IonCol, IonItem, IonInput } from '@ionic/react';
import 'react-phone-number-input/style.css'
import PhoneInput from 'react-phone-number-input'
import "./ParentInfo.scss"
import { Controller, useFormContext } from 'react-hook-form';

const ParentInfo: React.FC = () => {
  const { register, setValue, watch, control } = useFormContext<Registration>();

  return (
    <IonGrid>
      <IonRow>
        <IonCol>
          <h1>Údaje o zákonnom zástupcovi</h1>
        </IonCol>
      </IonRow>
      <IonRow>
        <IonCol>
          <h4>Meno</h4>
          <IonItem>
            <IonInput
              {...register('parent.name')}
            />
          </IonItem>
        </IonCol>
      </IonRow>
      <IonRow>
        <IonCol>
          <h4>Priezvisko</h4>
          <IonItem>
            <IonInput
              {...register('parent.surname')}
            />
          </IonItem>
        </IonCol>
      </IonRow>
      <IonRow>
        <IonCol>
          <h4>Email</h4>
          <IonItem>
            <IonInput
              {...register('parent.email')}
              type='email'
            />
          </IonItem>
        </IonCol>
      </IonRow>
      <IonRow>
        <IonCol>
          <h4>Telefónne číslo</h4>
          <IonItem>
            <Controller
              render={() => (
                <PhoneInput
                  placeholder='0949 000 000'
                  inputRef={register}
                  autoComplete="phone_number"
                  defaultCountry="SK"
                  style={{
                    maxWidth: "25vw",
                  }}
                  useNationalFormatForDefaultCountryValue={true}
                  onChange={(value) => setValue('parent.phone', value || '')}
                />
              )}
              name="parent.phone"
              control={control}
              rules={{ required: true }}
            />
            <span style={{paddingLeft: '1vw'}}>{watch('parent.phone')}</span>
          </IonItem>
        </IonCol>
      </IonRow>
    </IonGrid>
  );
};

export default ParentInfo;
