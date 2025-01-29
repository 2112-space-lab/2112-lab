



export type TileSatelliteMapping = {
    MappingID: string;
    NoradID: string;
    TileID: string;
    TileCenterLat: string;
    TileCenterLon: string;
    TileZoomLevel: number;
    IntersectionLongitude: number;
    IntersectionLatitude: number;
    Intersection: {
        Longitude: number;
        Latitude: number;
    };
};

export interface TileMapping {
    mappings: any[];
    totalItems: number;
}

export interface Tile {
    Quadkey: string;
    ZoomLevel: number;
    CenterLat: number;
    CenterLon: number;
    SpatialIndex?: string;
    NbFaces: number;
    Radius: number;
    BoundariesJSON?: string;
    ID: string;
}
