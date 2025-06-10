import type { ResourceSource } from 'solid-js';
export interface LostThingProps {
  syncList: ResourceSource<any>;
  tabIndex: string;
  id: number;
  publication_date: string;
  publication_time: string;
  thing_name: string;
  email: string;
  custom_text: string;
}
