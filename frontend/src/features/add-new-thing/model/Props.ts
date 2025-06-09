import type { ResourceSource, Setter } from 'solid-js';
export interface Props {
  syncLostThingsList: ResourceSource<any>;
  syncFoundThingsList: ResourceSource<any>;
  setAddNewThing: Setter<boolean>;
}
