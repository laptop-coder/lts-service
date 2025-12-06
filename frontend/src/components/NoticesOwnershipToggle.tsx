import { JSX } from 'solid-js';
import type { Setter } from 'solid-js';

import { NoticesOwnership } from '../utils/consts';
import TabsToggle from '../ui/TabsToggle/TabsToggle';

const NoticesOwnershipToggle = (props: {
  setter: Setter<NoticesOwnership>;
  fullscreen?: boolean;
}): JSX.Element => (
  <TabsToggle
    tabs={[NoticesOwnership.not_my, NoticesOwnership.my, NoticesOwnership.all]}
    tabsNames={['Не мои', 'Мои', 'Все']}
    setter={props.setter}
    tabsHTMLElementId='things_list_toggle'
    fullscreen={props.fullscreen}
  />
);

export default NoticesOwnershipToggle;
