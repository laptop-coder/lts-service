import { JSX } from 'solid-js';
import type { Setter } from 'solid-js';

import { NoticesVerification } from '../utils/consts';
import TabsToggle from '../ui/TabsToggle/TabsToggle';

const NoticesVerificationToggle = (props: {
  setter: Setter<NoticesVerification>;
  fullscreen?: boolean;
}): JSX.Element => (
  <TabsToggle
    tabs={[
      NoticesVerification.not_verified,
      NoticesVerification.approved,
      NoticesVerification.rejected,
    ]}
    tabsNames={['Не проверено', 'Одобрено', 'Отклонено']}
    setter={props.setter}
    tabsHTMLElementId='notices_verification_toggle'
    fullscreen={props.fullscreen}
  />
);

export default NoticesVerificationToggle;
