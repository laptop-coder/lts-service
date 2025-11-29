import { JSX } from 'solid-js';
import type { Setter } from 'solid-js';

import { ThingType } from '../utils/consts';
import TabsToggle from '../ui/TabsToggle/TabsToggle';

const ThingsTypeToggle = (props: {
  setter: Setter<ThingType>;
  fullscreen?: boolean;
}): JSX.Element => (
  <TabsToggle
    tabs={[ThingType.lost, ThingType.found]}
    tabsNames={['Потеряно', 'Найдено']}
    setter={props.setter}
    tabsHTMLElementId='things_list_toggle'
    fullscreen={props.fullscreen}
  />
);

export default ThingsTypeToggle;
