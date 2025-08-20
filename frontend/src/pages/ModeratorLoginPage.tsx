import { JSX } from 'solid-js';

import ModeratorLoginForm from '../ui/ModeratorLoginForm';
import ModeratorAuthPage from '../ui/ModeratorAuthPage/ModeratorAuthPage';

const ModeratorLoginPage = (): JSX.Element => {
  return (
    <ModeratorAuthPage>
      <ModeratorLoginForm />
    </ModeratorAuthPage>
  );
};

export default ModeratorLoginPage;
