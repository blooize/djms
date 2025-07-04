export interface Event {
    id: number;
    name: string;
    club_id: number;
}

export interface DancerSlots {
    club_id: number;
    event_id: number;
    dancerslots: Array<DancerSlot>;
}

export interface DancerSlot {
    date: number;
    dancers : Array<Dancer>;
}

export interface Dancer {
    id: number;
    name: string;
}

export interface TalentSlots {
    event_id: number;
    club_id: number;
    talentslots: Array<TalentSlot>;
}

export interface TalentSlot {
    date: number;
    talents: Array<Talent>;
}

export interface Talent {
    id: number;
    name: string;
}

export interface Club {
    id: number;
    name: string;
}