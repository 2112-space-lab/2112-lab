export interface TLE {
    ID: string;
    SpaceID: string;
    Line1: string;
    Line2: string;
    Epoch: string; // ISO 8601 formatted timestamp
}

export interface Satellite {
    ID: string;
    CreatedAt: string; // ISO 8601 formatted timestamp
    UpdatedAt: string; // ISO 8601 formatted timestamp
    Name: string;
    SpaceID: string;
    Type: string;
    LaunchDate?: string | null; // Nullable ISO 8601 formatted date
    DecayDate?: string | null; // Nullable ISO 8601 formatted date
    IntlDesignator: string;
    Owner: string;
    ObjectType: string;
    Period?: number | null;
    Inclination?: number | null;
    Apogee?: number | null;
    Perigee?: number | null;
    RCS?: number | null;
    TleUpdatedAt?: string | null; // Nullable ISO 8601 formatted timestamp
    Altitude?: number | null;
}

export interface SatelliteInfo {
    Satellite: Satellite;
    TLEs: TLE[];
}

export interface OrbitDataItem {
    latitude: number;
    longitude: number;
    altitude: number;
    time: string;
}