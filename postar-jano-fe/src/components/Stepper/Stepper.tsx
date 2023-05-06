import React, { useEffect, useState } from 'react';
import "./Stepper.scss"
import { Registration, IEvent, Stat, responseStatus } from '../../utils/types';
import { useForm, Controller } from "react-hook-form";
import { IonButton, IonCol, IonGrid, IonIcon, IonProgressBar, IonRow } from '@ionic/react';
import { arrowBackOutline, arrowForwardOutline } from 'ionicons/icons';
import IntroInfo from '../IntroInfo/IntroInfo';
import MedicineHealth from '../MedicineHalth/MedicineHealt';
import ChildInfo from '../childInfo/ChildInfo';

interface StepperProps {
  event: IEvent,
  stats: Stat[]
}

const initialValues = {
  child: {
    name: "",
    surname: "",
    gender: null,
    city: "",
    dateOfBirth: null,
    finishedSchoolYear: null,
    attendedPreviousEvents: null
  },
  days: [],
  medicine: {
    takes: false,
    drugs: ""
  },
  health: {
    hasProblmes: false,
    problems: ""
  },
  parent: {
    name: "",
    surname: "",
    email: "",
    phone: "",
  },
  memberShip: {
    attendedActivities: ""
  },
  notes: ""
}

const validate = () => {
  return true
}

const eventFull = (stats) => {
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
  const { control, watch, setValue, register, handleSubmit, errors, formState } = useForm<Registration>({
    defaultValues: initialValues
  })

  useEffect(() => {
    if (event.days.length > 1) {
      setPageCount(pageCount + 1)
    }
    setIsFull(eventFull(stats))

  }, [event, stats])

  const [pageCount, setPageCount] = useState(5)
  const [page, setPage] = useState(0)
  const [canGoBack, setCanGoBack] = useState(true)
  const [isFull, setIsFull] = useState(false)

  const onSubmit = data => {
    alert(JSON.stringify(data, null, 2));
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} >
      <IonGrid>
        <IonRow>
          <IonCol size="2"></IonCol>
          <IonCol>
            <IonProgressBar value={page / (pageCount - 1)}></IonProgressBar>
          </IonCol>
          <IonCol size="2"></IonCol>
        </IonRow>
        {
          event && stats &&
          <IonRow>
            <IonCol size="2"></IonCol>
            <IonCol>
              {
                page === 0 && <IntroInfo event={event} stats={stats} />
              }
              {
                  page === 1 && <ChildInfo register={register} control={control} setValue={setValue}/>
              }
              {
                page === 2 && <MedicineHealth register={register} setValue={setValue} watch={watch} />
              }
            </IonCol>
            <IonCol size="2"></IonCol>
          </IonRow>
        }
        <IonRow>
          <IonCol></IonCol>
          <IonCol size="3">
            <div className="previous">
              {
                page > 0 && canGoBack &&
                <IonButton expand="full" shape="round" onClick={() => {
                  if (page > 0) {
                    setPage(page - 1)
                  }
                }}>
                  <IonIcon icon={arrowBackOutline} />
                  Späť
                </IonButton>
              }
            </div>
          </IonCol>
          <IonCol size="3">
            <div className="next">
              {
                page < pageCount - 1 && !isFull && event.active &&
                <IonButton expand="full" shape="round" onClick={async () => {
                  if (page < pageCount && validate()) {
                    setPage(page + 1)
                  }
                }}>
                  Ďalej
                  <IonIcon icon={arrowForwardOutline} />
                </IonButton>
              }
              {
                page === pageCount - 1 &&
                <IonButton
                  expand="full"
                  shape="round"
                  color="success"
                  // disabled={!this.state.valid}
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
  )
};

export default Stepper;
