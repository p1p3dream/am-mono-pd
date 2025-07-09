import { jest } from '@jest/globals';

export interface MockMapInstance {
  getCenter: jest.Mock;
  getZoom: jest.Mock;
  getLayer: jest.Mock;
  remove: jest.Mock;
  flyTo: jest.Mock;
  on: jest.Mock;
  off: jest.Mock;
  addLayer: jest.Mock;
  addSource: jest.Mock;
  setLayoutProperty: jest.Mock;
  options: any;
}

export interface MockMapOptions {
  dragPan: boolean;
  scrollZoom: boolean;
  doubleClickZoom: boolean;
  touchZoomRotate: boolean;
}

const mockMapInstance = {
  getCenter: jest.fn().mockReturnValue({ lng: 0, lat: 0 }),
  getZoom: jest.fn().mockReturnValue(10),
  getLayer: jest.fn().mockReturnValue(true),
  remove: jest.fn(),
  flyTo: jest.fn(),
  on: jest.fn(),
  off: jest.fn(),
  addLayer: jest.fn(),
  addSource: jest.fn(),
  setLayoutProperty: jest.fn(),
  options: {},
};

mockMapInstance.flyTo.mockReturnValue(mockMapInstance);
mockMapInstance.on.mockReturnValue(mockMapInstance);
mockMapInstance.off.mockReturnValue(mockMapInstance);
mockMapInstance.addLayer.mockReturnValue(mockMapInstance);
mockMapInstance.addSource.mockReturnValue(mockMapInstance);
mockMapInstance.setLayoutProperty.mockReturnValue(mockMapInstance);

const mockMapClass = jest.fn().mockImplementation((options) => {
  const instance = Object.create(mockMapInstance);
  instance.options = options;
  return instance;
});

jest.doMock('mapbox-gl', () => ({
  Map: mockMapClass,
  NavigationControl: jest.fn(),
  Marker: jest.fn().mockImplementation(() => ({
    setLngLat: jest.fn().mockReturnThis(),
    addTo: jest.fn().mockReturnThis(),
    remove: jest.fn(),
  })),
  LngLat: jest.fn(),
  LngLatBounds: jest.fn(),
}));
