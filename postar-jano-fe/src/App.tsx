import { IonApp, IonRouterOutlet, IonSplitPane, setupIonicReact } from '@ionic/react';
import { IonReactRouter } from '@ionic/react-router';
import { Redirect, Route } from 'react-router-dom';
import axios from 'axios';

/* Core CSS required for Ionic components to work properly */
import '@ionic/react/css/core.css';

/* Basic CSS for apps built with Ionic */
import '@ionic/react/css/normalize.css';
import '@ionic/react/css/structure.css';
import '@ionic/react/css/typography.css';

/* Optional CSS utils that can be commented out */
import '@ionic/react/css/padding.css';
import '@ionic/react/css/float-elements.css';
import '@ionic/react/css/text-alignment.css';
import '@ionic/react/css/text-transformation.css';
import '@ionic/react/css/flex-utils.css';
import '@ionic/react/css/display.css';

/* Theme variables */
import './theme/variables.css';
import EventPage from './pages/Event/EventPage';
import { SWRConfig } from 'swr';
import { SWRFetcher } from './utils/api';

setupIonicReact();

axios.defaults.baseURL = `${import.meta.env.VITE_API_HOST}/api`

const App: React.FC = () => {
  return (
    <IonApp>
      <SWRConfig value={{fetcher: SWRFetcher}}>
        <IonReactRouter>
          <IonSplitPane contentId="main" >
            <IonRouterOutlet id="main">
              <Route path="/" exact={true}>
                <Redirect to="/event/1" />
              </Route>
              <Route path="/event/:id" exact={true} >
                <EventPage />
              </Route>
            </IonRouterOutlet>
          </IonSplitPane>
        </IonReactRouter>
      </SWRConfig>
    </IonApp>
  );
};

export default App;
