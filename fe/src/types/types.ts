export interface Event {
    id: number;
    title: string;
    description: string;
    photo: string;
    date_from: string;
    date_to: string;
    location: string;
    min_age: number;
    max_age: number;
    price: number;
    owner: Owner;
    info: string;
    days: Day[],
    time: string
    stats: Stat[] | null
}

export interface Stat {
    event_id: number,
    day_id: number,
    capacity: number,
    limit_boys: number | null,
    limit_girls: number | null,
    boys_count: number,
    girls_count: number,
}
export interface Day {
    id: number;
    capacity: number;
    limit_boys: number | null,
    limit_girls: number | null,
    description: string;
}

interface Owner {
    name: string;
    surname: string;
    photo: string;
    gender: Gender;
    email: string;
    phone: string;
}

export interface Registration {
    child: Child;
    days: number[]; // ids of days
    medicine: Medicine;
    health: Health;
    parent: Parent;
    memberShip: MemberShip; 
    notes: string; // long string
}

interface MemberShip {
    attendedActivities: string;
}

interface Medicine {
    takes: boolean | null;
    drugs: string;
}

interface Health {
    hasProblmes: boolean | null;
    problems: string;
}

interface Child {
    name: string;
    surname: string;
    gender: Gender | null;
    city: string;
    dateOfBirth: Date;
    finishedSchoolYear: string | null;
    attendedPreviousEvents: boolean | null;
}

interface Parent {
    name: string;
    surname: string;
    email: string;
    phone: string;
}

export enum Gender {
    Male = "male",
    Female = "female"
}

export enum ActionType {
    SET_CHILD_NAME,
    SET_CHILD_SURNAME,
    SET_CHILD_GENDER,
    SET_CHILD_BIRTH,
    SET_CHILD_CITY,
    SET_CHILD_YEAR,
    SET_CHILD_ATTEND,
    SET_DAYS,
    SET_PARNET_NAME,
    SET_PARENT_SURNAME,
    SET_PARENT_PHONE,
    SET_PARENT_EMAIL,
    SET_MEDICINE,
    SET_HEALTH,
    SET_PROBLEMS,
    SET_DRUGS,
    SET_ATTENDED_ACTIVITIES,
    SET_NOTES
}

export interface RegistrationRespone {
    success: boolean,
    days?: number[],
    token: string
}


export enum responseStatus {
    success =  "success",
    fail = "fail"
}