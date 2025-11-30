import { JSX } from 'solid-js';
import type { Setter } from 'solid-js';

import { UserProfileSection } from '../utils/consts';
import TabsToggle from '../ui/TabsToggle/TabsToggle';

const UserProfileSectionToggle = (props: {
  setter: Setter<UserProfileSection>;
  fullscreen?: boolean;
}): JSX.Element => (
  <TabsToggle
    tabs={[UserProfileSection.advertisements, UserProfileSection.settings]}
    tabsNames={['Мои объявления', 'Настройки']}
    setter={props.setter}
    tabsHTMLElementId='user_profile_sections_toggle'
    fullscreen={props.fullscreen}
  />
);

export default UserProfileSectionToggle;
