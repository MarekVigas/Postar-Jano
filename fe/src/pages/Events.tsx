import { IonButtons, IonContent, IonHeader, IonItem, IonLoading, IonMenuButton, IonPage, IonTitle, IonToolbar } from '@ionic/react';
import useSWR from 'swr';
import { IEvent } from '../utils/types';
import queryString from 'query-string'
import { useLocation } from 'react-router-dom'
import { useEffect } from 'react';
import { useStateMachine } from 'little-state-machine'

const updatePromoCode = (state: any, payload: string) => {
  return {
    ...state,
    promo: payload
  }
}

const EventsPage: React.FC = () => {
  const { data: events, isLoading } = useSWR<IEvent[]>({ url: `/events` })
  const { search } = useLocation()
  const { promo } = queryString.parse(search)
  const { actions } = useStateMachine({ updatePromoCode })

  useEffect(() => {
    if (promo) {
      console.log('set promo')
      actions.updatePromoCode(promo as string)
    }
  }, [promo])

  return (
    <IonPage>
      <IonHeader>
        <IonToolbar>
          <IonButtons slot="start">
            <IonMenuButton />
          </IonButtons>
          <IonTitle>Salezko</IonTitle>
        </IonToolbar>
      </IonHeader>

      <IonContent fullscreen>
        <IonHeader collapse="condense">
          <IonToolbar>
            <IonTitle size="large">Salezko</IonTitle>
          </IonToolbar>
        </IonHeader>
        <IonLoading
          isOpen={isLoading}
          title='Načítavam...'
        />
        {
          events && events.map((event, i) => (
            <IonItem key={i} href={`/event/${event.id}`}>
              {event.title}
            </IonItem>
          )
          )
        }
      </IonContent>
    </IonPage>
  );
};

export default EventsPage;
