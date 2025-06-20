import type { ResourceSource } from 'solid-js';

interface ThingProps {
  custom_text: string;
  id: number;
  page: 'home' | 'moderator' | 'status';
  publication_date: string;
  publication_time: string;
  reloadList: ResourceSource<any>;
  tabIndex: string;
  thing_name: string;
  type: 'lost' | 'found';
}

export interface LostThingProps extends ThingProps {
  email?: string;
}

export interface FoundThingProps extends ThingProps {
  thing_location?: string;
}
