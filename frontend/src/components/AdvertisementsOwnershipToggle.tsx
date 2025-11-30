import { JSX } from 'solid-js';
import type { Setter } from 'solid-js';

import { AdvertisementsOwnership } from '../utils/consts';
import TabsToggle from '../ui/TabsToggle/TabsToggle';

const AdvertisementsOwnershipToggle = (props: {
  setter: Setter<AdvertisementsOwnership>;
  fullscreen?: boolean;
}): JSX.Element => (
  <TabsToggle
    tabs={[
      AdvertisementsOwnership.not_my,
      AdvertisementsOwnership.my,
      AdvertisementsOwnership.all,
    ]}
    tabsNames={['Не мои', 'Мои', 'Все']}
    setter={props.setter}
    tabsHTMLElementId='things_list_toggle'
    fullscreen={props.fullscreen}
  />
);

export default AdvertisementsOwnershipToggle;
