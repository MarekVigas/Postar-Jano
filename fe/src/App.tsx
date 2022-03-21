import React, { PureComponent } from 'react';
import Rollbar from "rollbar";
import { Redirect, Route } from 'react-router-dom';
import { IonApp, IonRouterOutlet } from '@ionic/react';
import { IonReactRouter } from '@ionic/react-router';
import Event from './pages/Event/Event';
import Login from './pages/Login/Login';

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

interface Props {}

interface State {}

class AppComponent extends PureComponent<Props, State> {
    protected rollbar: Rollbar;

    constructor(props: Props) {
        super(props);
        this.rollbar = new Rollbar({
            accessToken: '6f9a8edc3c4e4c828d763fa498d0952d',
            captureUncaught: true,
            captureUnhandledRejections: true,
            enabled: (process.env.NODE_ENV === 'production'),
        });
    }
    render() {
        return (
        <IonApp>
            <IonReactRouter>
            <IonRouterOutlet>
                <Route path="/login" component={Login} exact={true} />
                <Route path="/event/:id" component={Event} exact={true} />
                <Route render={() => <Redirect to="/login" />} />
            </IonRouterOutlet>
            </IonReactRouter>
        </IonApp>
        )
    }
}

export default AppComponent;
