import { IonButtons, IonContent, IonHeader, IonItem, IonLoading, IonMenuButton, IonPage, IonTitle, IonToolbar } from '@ionic/react';
import useSWR from 'swr';
import { IEvent, Stat } from '../utils/types';

const EventsPage: React.FC = () => {
  const { data: events, isLoading } = useSWR<IEvent[]>({ url: `/events` })

  return (
    <IonPage>
      <IonHeader>
        <IonToolbar>
          <IonButtons slot="start">
            <IonMenuButton />
          </IonButtons>
          <IonTitle>Ivo Akcie</IonTitle>
        </IonToolbar>
      </IonHeader>

      <IonContent fullscreen>
        <IonHeader collapse="condense">
          <IonToolbar>
            <IonTitle size="large">Ivo Akcie</IonTitle>
          </IonToolbar>
        </IonHeader>
        <IonLoading
          isOpen={isLoading}
          title='Načítavam...'
        />
        {
          events && events.map((event) => (
            <IonItem>
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
