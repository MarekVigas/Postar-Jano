import React from 'react';
import { arrowForwardOutline, arrowBackOutline } from 'ionicons/icons'
import IntroInfo from "../IntroInfo/IntroInfo";
import "./Stepper.scss"
import { Registration, IEvent, ActionType, Stat, RegistrationRespone, responseStatus } from '../../types/types';
import { IonIcon, IonProgressBar, IonButton, IonContent, IonGrid, IonRow, IonCol } from '@ionic/react';
import ChildInfo from '../childInfo/ChildInfo';
import DaySelector from '../DaySelector/DaySelector';
import ParentInfo from '../ParentInfo/ParentInfo';
import MedicineHealth from '../MedicineHalth/MedicineHealt';
import { toastController } from '@ionic/core'
import OtherInfo from '../OtherInfo/OtherInfo';
import Results from '../Result/Results';
import { ApiClientContext } from '../../services/apiContext';
import { isPossiblePhoneNumber } from 'react-phone-number-input'

interface StepperProps {
    event: IEvent,
    stats: Stat[]
}

interface StepperState {
    registration: Registration;
    event: IEvent | null,
    stats: Stat[] | null,
    page: number,
    pageCount: number,
    valid: boolean,
    canGoBack: boolean,
    responseMsg: string | null,
    responseStatus: responseStatus | null
}

const defaultState: StepperState = {
    registration: {
        child: {
            name: "",
            surname: "",
            gender: null,
            city: "",
            dateOfBirth: new Date(),
            finishedSchoolYear: null,
            attendedPreviousEvents: null
        },
        days: [],
        medicine: {
            takes: null,
            drugs: ""
        },
        health: {
            hasProblmes: null,
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
    },
    stats: null,
    event: null,
    page: 0,
    pageCount: 5,
    valid: true,
    canGoBack: true,
    responseMsg: null,
    responseStatus: null
}

class Stepper extends React.Component<StepperProps, StepperState> {
    static contextType = ApiClientContext;
    context!: React.ContextType<typeof ApiClientContext>;
    state: StepperState;
    
    constructor(props: StepperProps) {
        super(props);
        if (this.props.event.days.length > 1) {
            defaultState.pageCount += 1
        }
        this.state = defaultState;
        this.state.event = props.event;
        this.state.stats = props.stats;
    }

    protected removeCapslock = (str: string) => {
        return str.charAt(0).toUpperCase() + str.slice(1).toLowerCase();
    }

    protected setValueHandler = (type: ActionType, value: any) => {
        const state = {...this.state};
        switch (type) {
            case ActionType.SET_CHILD_NAME:
                state.registration.child.name = this.removeCapslock(value);
                break;
            case ActionType.SET_CHILD_SURNAME:
                state.registration.child.surname = this.removeCapslock(value);
                break;
            case ActionType.SET_CHILD_GENDER:
                state.registration.child.gender = value;
                break;
            case ActionType.SET_CHILD_BIRTH:
                state.registration.child.dateOfBirth = value;
                break;
            case ActionType.SET_CHILD_CITY:
                state.registration.child.city = value;
                break;
            case ActionType.SET_CHILD_YEAR:
                state.registration.child.finishedSchoolYear = value;
                break;
            case ActionType.SET_CHILD_ATTEND:
                state.registration.child.attendedPreviousEvents = value;
                break;
            case ActionType.SET_DAYS:
                state.registration.days = value;
                break;
            case ActionType.SET_PARNET_NAME:
                state.registration.parent.name = this.removeCapslock(value);
                break;
            case ActionType.SET_PARENT_SURNAME:
                state.registration.parent.surname = this.removeCapslock(value);
                break;
            case ActionType.SET_PARENT_PHONE:
                state.registration.parent.phone = value;
                break;
            case ActionType.SET_PARENT_EMAIL:
                state.registration.parent.email = value;
                break;
            case ActionType.SET_MEDICINE:
                state.registration.medicine.takes = value;
                break;
            case ActionType.SET_HEALTH:
                state.registration.health.hasProblmes = value;
                break;
            case ActionType.SET_DRUGS:
                state.registration.medicine.drugs = value;
                break;
            case ActionType.SET_PROBLEMS:
                state.registration.health.problems = value;
                break;
            case ActionType.SET_ATTENDED_ACTIVITIES:
                state.registration.memberShip.attendedActivities = value;
                break;
            case ActionType.SET_NOTES:
                state.registration.notes = value;
                break;
            default:
                break;
        }
        this.setState({...state})
    }

    protected validate = async () => {
        const problems: string[] = [];
        let page = this.state.page;
        const { child, medicine, health, parent } = this.state.registration

        const validEmail = (email: string): boolean => {
            const re = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
            return re.test(String(email).toLowerCase());
        }

        if (page == 1) {
            if (child.name.length < 2){
                problems.push("Vyplňte meno dieťaťa.");
            }
            if (child.surname.length < 2){
                problems.push("Vyplňte priezvisko dieťaťa.");
            }
            if (child.gender == null){
                problems.push("Vyplňte pohlavie dieťaťa.");
            }
            if (child.city.length <= 2){
                problems.push("Vyplňte mesto/obec trvalého bydliska, nie skratku.");
            }
            if (child.finishedSchoolYear == null){
                problems.push("Vyplňte ukončený školský rok.");
            }
            if (child.attendedPreviousEvents == null){
                problems.push("Vyplňte účasť na minuloročných akciách.");
            }
            if (child.finishedSchoolYear == null){
                problems.push("Vyplňte ukončený školský rok.");
            }
        }

        if (page == 2) {
            if (medicine.takes === null){
                problems.push("Vyplňte či vaše dieťa berie lieky.");
            }
            if (medicine.takes === true && medicine.drugs.length < 2){
                problems.push("Vyplňte aké lieky berie vaše dieťa.");
            }
            if (health.hasProblmes === null){
                problems.push("Vyplňte zdravotný stav.");
            }
            if (health.hasProblmes === true && health.problems.length < 2){
                problems.push("Vyplňte zdravotné problémy.");
            }
        }

        if (this.state.event!.days.length > 1) {
            page -= 1; 
        }

        if (page == 3) {
            if (parent.name.length < 2){
                problems.push("Vyplňte meno rodiča.")
            }
            if (parent.surname.length < 2){
                problems.push("Vyplňte priezvisko rodiča.")
            }
            if (!validEmail(parent.email)){
                problems.push("Nesprávny formát emailu.")
            }
            if (!isPossiblePhoneNumber(parent.phone)){
                problems.push("Nesprávny formát tel. čísla.")
            }
        }
        
        if (problems.length > 0) {
            const sendToast = await toastController.create({
                duration: 3500,
                message: problems[0],
                color: "danger",
                position: "top"
            })
            sendToast.present()
            return false;
        }
        return true;
    }

    protected eventFull = (): boolean => {
        let total = 0;
        let sum = 0;

        if (this.state.stats) {
            for (const stat of this.state.stats) {
                sum += stat.boys_count + stat.girls_count;
                total += stat.capacity;   
            }
            return (total === sum)
        } else {
            return true
        }
    }

    protected handleSubmit = async () => {
        if (!this.context || !this.state.event) {
            return;
        }
        this.setState({
            ...this.state,
            page: this.state.page +1,
            responseMsg: null,
            responseStatus: null
        })
        const sendToast = await toastController.create({
            duration: 2000,
            message: "Prihláška bola odoslaná na spracovanie",
            color: "secondary",
            position: "bottom"
        })
        sendToast.present()

        const registration = {...this.state.registration};

        if (this.state.event && this.state.event.days.length === 1) {
            registration.days = []
            registration.days.push(this.state.event.days[0].id)
        }
        
        this.context.event.register(this.state.event.id, registration)
        .then(async(res)  => {
            if(res.success){
                const processToast = await toastController.create({
                    duration: 2000,
                    message: "Prihláška bola úspešne spracovaná",
                    color: "success",
                    position: "bottom"
                })
                processToast.present()
                this.setState({
                    ...this.state,
                    canGoBack: false,
                    responseMsg: "Prihláška bola úspešne spracovaná",
                    responseStatus: responseStatus.success
                })
            } else if (res.registeredIDs) {
                if (res.registeredIDs.length !== this.state.registration.days.length) {
                    const notRegistred = this.state.registration.days.filter(d => !res.registeredIDs?.includes(d))
                    let msg = 'Nepodarilo sa prihlásiť na tieto termíny: '
                    for (const dayId of notRegistred) {
                        const day = this.props.event.days.filter(d => d.id === dayId)[0];
                        msg += `${day.description} `
                    }
                    this.setState({
                        ...this.state,
                        canGoBack: true,
                        responseMsg: msg,
                        responseStatus: responseStatus.fail
                    })
                }
            } else {
                this.setState({
                    ...this.state,
                    canGoBack: true,
                    responseMsg: "Došlo k neznámej chybe",
                    responseStatus: responseStatus.fail
                })
            }
        })
        .catch(async err => {
            const errorToast = await toastController.create({
                duration: 2000,
                message: err,
                color: "danger",
                position: "bottom"
            })
            errorToast.present()
            this.setState({
                ...this.state,
                canGoBack: true,
                responseMsg: "Došlo k chybe pri odosielaní prihlášky",
                responseStatus: responseStatus.fail
            })
            console.log(err)
        })
    }

    componentDidUpdate(prevProps: StepperProps) {
        if(JSON.stringify(prevProps.stats) !== JSON.stringify(this.props.stats)) {
            this.setState({
                stats: this.props.stats
            });
        }
    }

    render(): React.ReactNode {
        return (
            <IonContent>
            <IonGrid>
                <IonRow>
                    <IonCol size="2"></IonCol>
                    <IonCol>
                        <IonProgressBar value={this.state.page/(this.state.pageCount-1)}></IonProgressBar>
                    </IonCol>
                    <IonCol size="2"></IonCol>
                </IonRow>
                {
                    this.state.event && 
                    <IonRow>
                        <IonCol size="2"></IonCol>
                        <IonCol>
                            {
                                this.state.page === 0 && this.state.stats && <IntroInfo event={this.state.event} stats={this.state.stats} />
                            }
                            {
                                this.state.page === 1 && <ChildInfo 
                                    event={this.props.event}
                                    registration={this.state.registration}
                                    setValue={(t,v) => this.setValueHandler(t, v)}
                                />
                            }
                            {
                                this.state.page === 2 && <MedicineHealth
                                    registration={this.state.registration}
                                    setValue={(t,v) => this.setValueHandler(t, v)}
                                />
                            }
                            {
                                this.state.page === 3 && this.state.event.days.length > 1 && this.state.stats && <DaySelector
                                    stats={this.state.stats}
                                    event={this.state.event}
                                    registration={this.state.registration}
                                    setValue={(t,v) => this.setValueHandler(t, v)}
                                />
                            }
                            {
                                this.state.page === 3 && this.state.event.days.length === 1 && <ParentInfo
                                    registration={this.state.registration}
                                    setValue={(t,v) => this.setValueHandler(t, v)}
                                />
                            }
                            {
                                this.state.page === 4 && this.state.event.days.length > 1 && <ParentInfo
                                    registration={this.state.registration}
                                    setValue={(t,v) => this.setValueHandler(t, v)}
                                />
                            }
                            {
                                this.state.page === 4 && this.state.event.days.length === 1 && <OtherInfo
                                    registration={this.state.registration}
                                    setValue={(t,v) => this.setValueHandler(t, v)}
                                />
                            }
                            {
                                this.state.page === 5 && this.state.event.days.length > 1 && <OtherInfo
                                    registration={this.state.registration}
                                    setValue={(t,v) => this.setValueHandler(t, v)}
                                />
                            }
                            {
                                this.state.page === 5 && this.state.event.days.length === 1 && <Results
                                    registration={this.state.registration}
                                    responseMsg={this.state.responseMsg}
                                    responseStatus={this.state.responseStatus}
                                />
                            }
                            {
                                this.state.page === 6 && this.state.event.days.length > 1 && <Results
                                    registration={this.state.registration}
                                    responseMsg={this.state.responseMsg}
                                    responseStatus={this.state.responseStatus}
                                />
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
                                this.state.page > 0 && this.state.canGoBack && 
                                <IonButton expand="full" shape="round" onClick={() => {
                                    if (this.state.page > 0) {
                                        this.setState({...this.state, page: this.state.page - 1})
                                    }
                                }}>
                                    <IonIcon icon={arrowBackOutline}/>
                                    Späť
                                </IonButton>
                            }
                        </div>
                    </IonCol>
                    <IonCol size="3">
                        <div className="next">
                            {
                                this.state.page < this.state.pageCount -1 && !this.eventFull() && this.state.event?.active &&
                                <IonButton expand="full" shape="round" onClick={async () => {
                                    if (this.state.page < this.state.pageCount && await this.validate()) {
                                        this.setState({...this.state, page: this.state.page + 1})
                                    }
                                }}>
                                    Ďalej
                                    <IonIcon icon={arrowForwardOutline}/>
                                </IonButton>
                            }
                            {
                                this.state.page === this.state.pageCount -1 &&
                                <IonButton 
                                    expand="full" 
                                    shape="round"
                                    color="success"
                                    disabled={!this.state.valid}
                                    onClick={this.handleSubmit}
                                >
                                    Odoslať
                                </IonButton>
                            }
                        </div>
                    </IonCol>
                    <IonCol></IonCol>
                </IonRow>
            </IonGrid>
        </IonContent>
        );
    }
};

export default Stepper;
