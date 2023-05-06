import { IonButtons, IonContent, IonHeader, IonLoading, IonMenuButton, IonPage, IonTitle, IonToolbar } from '@ionic/react';
import { useParams } from 'react-router';
import useSWR from 'swr';
import { IEvent, Stat } from '../utils/types';
import Stepper from '../components/Stepper/Stepper'

const EventPage: React.FC = () => {

  const { id } = useParams<{ id: string }>();
  const { data: event, isLoading } = useSWR<IEvent>({ url: `/events/${id}` })
  const { data: stats, isLoading: isStatsLoading } = useSWR<Stat[]>({ url: `/stats/${id}` }, { refreshInterval: 5000 })

  return (
    <IonPage>
      <IonHeader>
        <IonToolbar>
          <IonButtons slot="start">
            <IonMenuButton />
          </IonButtons>
          <IonTitle>{event?.title}</IonTitle>
        </IonToolbar>
      </IonHeader>

      <IonContent fullscreen>
        <IonHeader collapse="condense">
          <IonToolbar>
            <IonTitle size="large">{event?.title}</IonTitle>
          </IonToolbar>
        </IonHeader>
        <IonLoading
          isOpen={isLoading || isStatsLoading}
          title='Načítavam...'
        />
        {
          event && stats && <Stepper event={event} stats={stats} />
        }
      </IonContent>
    </IonPage>
  );
};

export default EventPage;
