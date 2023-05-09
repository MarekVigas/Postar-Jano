import { IonButton, IonButtons, IonCard, IonCardContent, IonCardHeader, IonCardTitle, IonContent, IonHeader, IonItem, IonLabel, IonList, IonLoading, IonMenuButton, IonPage, IonThumbnail, IonTitle, IonToolbar } from '@ionic/react';
import useSWR from 'swr';
import { IEvent, PromoResponse } from '../utils/types';
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
  const { actions, state } = useStateMachine({ updatePromoCode })

  const { data: promoValidation } = useSWR<PromoResponse>(state.promo ? {
    url: `promo_codes/validate`, config: {
      method: 'POST',
      data: {
        promo_code: state.promo
      }
    }
  } : null)

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
          state.promo && promoValidation &&
          <IonItem>
            Váš promo kód: {state.promo} <br />
            Zostávajúci počet použití: {promoValidation?.available_registrations}
          </IonItem>
        }
        {
          events &&
          <IonCard>
            <IonCardHeader>
              <IonCardTitle>Akcie</IonCardTitle>
            </IonCardHeader>
            <IonCardContent>
              <IonList>
                {events.filter(event => event.title.toLowerCase() !== 'test').map((event) => (
                  <IonItem href={`/event/${event.id}`} key={event.id}>
                    <IonThumbnail slot="start">
                      <img alt={`${event.title} photo`} src={event.photo} />
                    </IonThumbnail>
                    <IonLabel>{event.title}</IonLabel>
                    <IonButton slot='end'>Prihlásiť</IonButton>
                  </IonItem>
                ))}
              </IonList>
            </IonCardContent>
          </IonCard>
        }
      </IonContent>
    </IonPage>
  );
};

export default EventsPage;
