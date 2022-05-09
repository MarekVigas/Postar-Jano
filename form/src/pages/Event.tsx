import { IonContent, IonHeader, IonPage, IonTitle, IonToolbar, IonLoading } from '@ionic/react';
import React from 'react';
import { RouteComponentProps } from 'react-router-dom';
import Stepper from '../components/Stepper/Stepper';
import { IEvent, Stat } from '../types/types';
import { ApiClientContext } from '../services/apiContext';

interface Props extends RouteComponentProps<{
    id: string;
}> { }

interface State {
    event: IEvent | null,
    stats: Stat[] | null
}

class Event extends React.Component<Props, State> {
    static contextType = ApiClientContext;
    context!: React.ContextType<typeof ApiClientContext>;
    timer: NodeJS.Timeout | null;

    constructor(props: Props) {
        super(props);
        this.state = {
            event: null,
            stats: null
        };
        this.timer = null;
    }

    componentDidMount() {
        if (!this.context) {
            return;
        }

        this.context.event.get(parseInt(this.props.match.params.id, 10))
        .then(event => {
            this.setState({
                ...this.state,
                event,
            })
        })
        .catch(err => {
            throw err            
        })
        this.getStats();
        this.timer = setInterval(() => this.getStats(), 5000)
    }

    componentWillUnmount() {
        clearInterval(this.timer!);
        this.timer = null;
    }

    protected getStats = () => {
        if (!this.context) {
            return;
        }

        this.context.event.stats(parseInt(this.props.match.params.id, 10))
        .then(stats => {
            this.setState({
                ...this.state,
                stats,
            })
        })
        .catch(err => {
            throw err            
        })
    }

    render() {
        return (
            <IonPage>
                <IonHeader>
                    <IonToolbar>
                        <IonTitle>{this.state.event?.title.toUpperCase()}</IonTitle>
                    </IonToolbar>
                </IonHeader>
                <IonContent>
                    <IonHeader collapse="condense">
                        <IonToolbar>
                            <IonTitle size="large">{this.state.event?.title.toUpperCase()}</IonTitle>
                        </IonToolbar>
                    </IonHeader>
                    {
                        this.state.event && this.state.stats && <Stepper event={this.state.event} stats={this.state.stats} />
                    }
                    {
                        this.state.event == null &&
                        <IonLoading
                            isOpen={true}
                            message={'Načítavam akciu...'}
                        />
                    }
                </IonContent>
            </IonPage>
        )
    }
};

export default Event;
