import React, { useEffect, useState } from 'react';
import "./Stepper.scss"
import { Registration, IEvent, Stat, RegistrationRespone, PromoResponse, Gender } from '../../utils/types';
import { FormProvider, useForm } from "react-hook-form";
import { IonButton, IonCol, IonGrid, IonIcon, IonItem, IonProgressBar, IonRow, useIonRouter, useIonToast } from '@ionic/react';
import { arrowBackOutline, arrowForwardOutline } from 'ionicons/icons';
import IntroInfo from '../IntroInfo/IntroInfo';
import MedicineHealth from '../MedicineHalth/MedicineHealt';
import ChildInfo from '../childInfo/ChildInfo';
import ParentInfo from '../ParentInfo/ParentInfo';
import OtherInfo from '../OtherInfo/OtherInfo';
import DaySelector from '../DaySelector/DaySelector';
import axios from 'axios'
import { useStateMachine } from 'little-state-machine';
import useSWR from 'swr';
import { yupResolver } from "@hookform/resolvers/yup";
import * as yup from "yup";
import { sk } from 'yup-locales';

yup.setLocale(sk)

const schema = yup.object().shape({
  child: yup.object().shape({
    name: yup.string().min(3).required().label('Meno dieťaťa'),
    surname: yup.string().min(3).required().label('Priezvisko dieťaťa'),
    city: yup.string().min(3).required().label('Mesto'),
    gender: yup.string().required().label('Pohlavie'),
    finishedSchoolYear: yup.string().required().label('Ukončený školský rok'),
    attendedPreviousEvents: yup.boolean().required().label('Zúčastnil sa minuloročných akcií'),
    dateOfBirthDay: yup.number().min(1).max(31).required().label('Deň narodenia'),
    dateOfBirthMonth: yup.number().min(1).max(12).required().label('Mesiac narodenia'),
    dateOfBirthYear: yup.number().min(1990).max(2018).required().label('Rok narodenia')
  }),
  medicine: yup.object().shape({
    takes: yup.boolean().required().label('Lieky'),
    drugs: yup.string().optional()
  }),
  health: yup.object().shape({
    hasProblmes: yup.boolean().required().label('Zdravotné ťažkosti'),
    problems: yup.string().optional()
  }),
  parent: yup.object().shape({
    name: yup.string().min(3).required().label('Meno zákonneho zástupcu'),
    surname: yup.string().min(3).required().label('Priezvisko zákonneho zástupcu'),
    email: yup.string().email().required().label('Email'),
    phone: yup.string().required().label('Telefónne číslo')
  }),
  days: yup.array().of(yup.number().required()).min(1).required().label('Dni akcie'),
  promo_code: yup.string().nullable().label('Promo kód'),
  notes: yup.string().nullable().label('Poznámka'),
  memberShip: yup.object().shape({
    attendedActivities: yup.string().required().label('Zúčastnené aktivity')
  }).required().label('Členstvo')
});

interface StepperProps {
  event: IEvent,
  stats: Stat[]
}

enum ActivePage {
  Intro,
  ChildInfo,
  MedicineHealth,
  DaySelector,
  ParentInfo,
  OtherInfo,
  Results
}

const eventFull = (stats: Stat[]) => {
  let total = 0;
  let sum = 0;

  if (stats) {
    for (const stat of stats) {
      sum += stat.boys_count + stat.girls_count;
      total += stat.capacity;
    }
    return (total === sum)
  } else {
    return true
  }
}


const Stepper: React.FC<StepperProps> = ({ event, stats }) => {
  const methods = useForm({
    resolver: yupResolver(schema),
    mode: 'onBlur'
  })

  const { setValue, handleSubmit, formState: { errors, isValid } } = methods;

  const { child: childErrors, medicine: medicineErrors, health: healthErrors, parent: parentErrors, days: daysErrors } = errors

  const childErrorMessages = [
    childErrors?.name,
    childErrors?.surname,
    childErrors?.dateOfBirthDay,
    childErrors?.dateOfBirthMonth,
    childErrors?.dateOfBirthYear,
    childErrors?.city,
    childErrors?.attendedPreviousEvents

  ].map(error => error?.message).filter(e => e)

  const MedicineHealthErrorMessages = [
    medicineErrors?.takes,
    healthErrors?.hasProblmes
  ].map(error => error?.message).filter(e => e)

  const daysErrorMessages = [
    daysErrors

  ].map(error => error?.message).filter(e => e)

  const parentErrorMessages = [
    parentErrors?.name,
    parentErrors?.surname,
    parentErrors?.email,
    parentErrors?.phone
  ].map(error => error?.message).filter(e => e)

  const router = useIonRouter()

  const [pageCount, setPageCount] = useState(5)
  const [page, setPage] = useState(0)
  const [activePage, setActivePage] = useState(0)
  const [canGoBack, setCanGoBack] = useState(true)
  const [isFull, setIsFull] = useState(false)
  const { state: { promo } } = useStateMachine()

  const [present] = useIonToast()

  const { data: promoValidation } = useSWR<PromoResponse>(promo ? {
    url: `promo_codes/validate`, config: {
      method: 'POST',
      data: {
        promo_code: promo
      }
    }
  } : null)

  useEffect(() => {
    if (event.days.length > 1) {
      setPageCount(6)
    } else {
      setValue('days', [event.days[0].id])
    }
    setIsFull(eventFull(stats))

  }, [event, stats])

  useEffect(() => {
    let pages: ActivePage[] = [];
    if (event.days.length > 1) {
      pages = [
        ActivePage.Intro,
        ActivePage.ChildInfo,
        ActivePage.MedicineHealth,
        ActivePage.DaySelector,
        ActivePage.ParentInfo,
        ActivePage.OtherInfo,
        ActivePage.Results,
      ]
    } else {
      pages = [
        ActivePage.Intro,
        ActivePage.ChildInfo,
        ActivePage.MedicineHealth,
        ActivePage.ParentInfo,
        ActivePage.OtherInfo,
        ActivePage.Results,
      ]
    }
    setActivePage(pages[page])
  }, [page, pageCount, event])

  return (
    <FormProvider {...methods}>
      <form onSubmit={handleSubmit((data) => {
        const { child, days, parent, notes, medicine, health, memberShip } = data
        
        const padded_month = `${child.dateOfBirthMonth}`.padStart(2, '0')
        const padded_day = `${child.dateOfBirthDay}`.padStart(2, '0')
        
        const registration: Registration = {
          child: {
            ...child,
            gender: child.gender == Gender.Male ? Gender.Male : Gender.Female,
            dateOfBirth: `${child.dateOfBirthYear}-${padded_month}-${padded_day}T05:00:00.00Z`
          },
          days: event.days.length === 1 ? [event.days[0].id] : days,
          parent,
          promo_code: promo ?? null,
          notes: notes ?? "",
          medicine: {
            takes: medicine.takes,
            drugs: medicine.drugs ?? "",
          },
          health: {
            hasProblmes: health.hasProblmes,
            problems: health.problems ?? "",
          },
          memberShip: {
            attendedActivities: memberShip.attendedActivities
          }
        }

        setPage(page + 1)

        if (promo) {
          registration.promo_code = promo
        }

        const dataString = JSON.stringify(registration, null, 2)
        console.log(dataString)

        axios.post<RegistrationRespone>(`/registrations/${event.id}`, registration)
          .then((response) => response.data)
          .then((res: RegistrationRespone) => {
            if (res.success) {
              present({
                duration: 2000,
                message: "Prihláška bola úspešne spracovaná",
                color: "success",
                position: "top"
              })
              setCanGoBack(false)
              router.push('/events')
            } else if (res.registeredIDs) {
              if (res.registeredIDs.length !== registration.days.length) {
                const notRegistred = registration.days.filter(d => !res.registeredIDs?.includes(d))
                let msg = 'Nepodarilo sa prihlásiť na tieto termíny: '
                for (const dayId of notRegistred) {
                  const day = event.days.filter(d => d.id === dayId)[0];
                  msg += `${day.description} `
                }
                setCanGoBack(true)
                present({
                  duration: 2000,
                  message: msg,
                  color: "danger",
                  position: "top"
                })
              }
            } else {
              present({
                duration: 2000,
                message: "Došlo k neznámej chybe",
                color: "danger",
                position: "top"
              })
              router.push('/events')
            }
          })
          .catch(async err => {
            present({
              duration: 2000,
              message: "Došlo k neznámej chybe",
              color: "danger",
              position: "top"
            })
            setCanGoBack(true)
            console.log(err)
          })
      })} >
        <IonGrid>
          <IonRow>
            <IonCol size="2"></IonCol>
            <IonCol>
              <IonProgressBar value={page / (pageCount - 1)}></IonProgressBar>
            </IonCol>
            <IonCol size="2"></IonCol>
          </IonRow>
          {/* <IonRow><pre>{JSON.stringify(watch(), null, 2)}</pre></IonRow> */}
          {
            event && stats &&
            <IonRow>
              <IonCol size="2"></IonCol>
              <IonCol>
                {
                  activePage === ActivePage.Intro && <IntroInfo event={event} stats={stats} />
                }
                {
                  activePage === ActivePage.ChildInfo && <ChildInfo/>
                }
                {
                  activePage === ActivePage.MedicineHealth && <MedicineHealth/>
                }
                {
                  activePage === ActivePage.DaySelector && <DaySelector event={event} stats={stats} />
                }
                {
                  activePage === ActivePage.ParentInfo && <ParentInfo />
                }
                {
                  activePage === ActivePage.OtherInfo && <OtherInfo />
                }
              </IonCol>
              <IonCol size="2"></IonCol>
            </IonRow>
          }
          {
            !isValid &&
            <IonRow>
              <IonCol></IonCol>
              <IonCol>
                {
                  activePage === ActivePage.ChildInfo && childErrorMessages.map((error) => (
                    <IonItem color='danger'>
                      {error}
                    </IonItem>
                  ))
                }
                {
                  activePage === ActivePage.MedicineHealth && MedicineHealthErrorMessages.map((error) => (
                    <IonItem color='danger'>
                      {error}
                    </IonItem>
                  ))
                }
                {
                  activePage === ActivePage.DaySelector && daysErrorMessages.map((error) => (
                    <IonItem color='danger'>
                      {error}
                    </IonItem>
                  ))
                }
                {
                  activePage === ActivePage.ParentInfo && parentErrorMessages.map((error) => (
                    <IonItem color='danger'>
                      {error}
                    </IonItem>
                  ))
                }
                {
                  activePage === ActivePage.OtherInfo && [...childErrorMessages, ...MedicineHealthErrorMessages, ...daysErrorMessages, ...parentErrorMessages].map((error) => (
                    <IonItem color='danger'>
                      {error}
                    </IonItem>
                  ))
                }
              </IonCol>
              <IonCol></IonCol>
            </IonRow>
          }
          <IonRow>
            <IonCol></IonCol>
            <IonCol size="3">
              <div className="previous">
                {
                  page > 0 && canGoBack &&
                  <IonButton expand="full" shape="round" onClick={() => setPage(page - 1)}>
                    <IonIcon icon={arrowBackOutline} />
                    Späť
                  </IonButton>
                }
              </div>
            </IonCol>
            <IonCol size="3">
              <div className="next">
                {
                  page < pageCount - 1 && !isFull && (event.active || (event.promo_registration && promoValidation?.status == "ok" && promoValidation.available_registrations > 0)) &&
                  <IonButton expand="full" shape="round" onClick={async () => {
                    if (page < pageCount) {
                      setPage(page + 1)
                    }
                  }}>
                    Ďalej
                    <IonIcon icon={arrowForwardOutline} />
                  </IonButton>
                }
                {
                  activePage == ActivePage.OtherInfo &&
                  <IonButton
                    expand="full"
                    shape="round"
                    color="success"
                    disabled={!isValid}
                    type='submit'
                  >
                    Odoslať
                  </IonButton>
                }
              </div>
            </IonCol>
            <IonCol></IonCol>
          </IonRow>
        </IonGrid>
      </form>
    </FormProvider>
  )
};

export default Stepper;
