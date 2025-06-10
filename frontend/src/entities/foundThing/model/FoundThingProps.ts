import type { ResourceSource } from 'solid-js';

export interface FoundThingProps {
  syncList: ResourceSource<any>;
  tabIndex: string;
  id: number;
  publication_date: string;
  publication_time: string;
  thing_name: string;
  thing_location: string;
  custom_text: string;
}
